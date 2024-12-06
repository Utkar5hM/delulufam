package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Utkar5hM/delulufam/utils/config"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type authHandler struct {
	config.Handler
}

func (h *authHandler) GoogleLogin(c echo.Context) error {
	url := h.Config.GoogleLoginConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *authHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	token, err := h.Config.GoogleLoginConfig.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	client := h.Config.GoogleLoginConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email    string `json:"email"`
		Username string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return err
	}

	// Check if user exists in the database
	sql_fetched_user := &User{}
	sql, _, _ := goqu.From("users").Where(goqu.C("email").Eq(userInfo.Email)).Select("username", "email", "role").ToSQL()
	row := h.DB.QueryRow(context.Background(), sql)
	err = row.Scan(&sql_fetched_user.Username, &sql_fetched_user.Email, &sql_fetched_user.Role)
	if err != nil {
		// User does not exist, create a new user
		_, err = h.DB.Exec(context.Background(),
			"INSERT INTO users (username, email, role) VALUES ($1, $2, $3)",
			userInfo.Username, userInfo.Email, "user")
		if err != nil {
			return err
		}
		sql_fetched_user.Username = userInfo.Username
		// sql_fetched_user.Email = userInfo.Email
		sql_fetched_user.Role = "user"
	}

	// Generate JWT token
	claims := &JwtCustomClaims{
		sql_fetched_user.Username,
		sql_fetched_user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.Config.JWT_SECRET))
	if err != nil {
		return err
	}

	// Set JWT token as cookie
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})
}
