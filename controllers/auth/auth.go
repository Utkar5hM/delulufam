package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Utkar5hM/delulufam/utils/config"
	"github.com/Utkar5hM/delulufam/utils/render"
	"github.com/Utkar5hM/delulufam/views"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	DB     *pgxpool.Pool
	config *config.Config
}

type UserRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtCustomClaims struct {
	Username string `json:"username"`
	Admin    string `json:"admin"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    string `json:"admin"`
	Role     string `json:"role"`
}

func (h *Handler) RegisterPost(c echo.Context) error {
	user := new(UserRegister)
	if err := c.Bind(user); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = h.DB.Exec(context.Background(),
		"INSERT INTO users (username, password, email, phone, role, admin) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Username, string(hashedPassword), user.Email, user.Phone, "user", false)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

func (h *Handler) LoginPost(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	fetched_user := &User{}
	var isAdmin bool
	sql, _, _ := goqu.From("users").Where(goqu.C("username").Eq(username)).Select("username", "password", "admin", "role").ToSQL()
	row := h.DB.QueryRow(context.Background(), sql)
	if err := row.Scan(&fetched_user.Username, &fetched_user.Password, &isAdmin, &fetched_user.Role); err != nil {
		log.Fatalf("Error fetching row: %v", err)
	}
	if isAdmin {
		fetched_user.Admin = "true"
	} else {
		fetched_user.Admin = "false"
	}
	if err := bcrypt.CompareHashAndPassword([]byte(fetched_user.Password), []byte(password)); err != nil {
		return echo.ErrUnauthorized
	}

	claims := &JwtCustomClaims{
		fetched_user.Username,
		fetched_user.Admin,
		fetched_user.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(h.config.JWT_SECRET))
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})

}

func UseSubroute(g *echo.Group, db *pgxpool.Pool, cfg *config.Config) {
	h := &Handler{DB: db, config: cfg}
	g.POST("/register", h.RegisterPost)
	g.POST("/login", h.LoginPost)
	g.GET("/login", func(c echo.Context) error {
		return render.Render(c, http.StatusOK, views.Login())
	})
}

func IsLoggedIn(jwt_secret string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey:  []byte(jwt_secret),
		TokenLookup: "header:Authorization,cookie:jwt",
	})
}
