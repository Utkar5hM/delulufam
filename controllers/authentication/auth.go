package authentication

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	Id       int    `json:"id"`
}

type JwtCustomClaims struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Id       int    `json:"id"`
	jwt.RegisteredClaims
}

func IsLoggedIn(jwt_secret string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey:  []byte(jwt_secret),
		TokenLookup: "header:Authorization:Bearer ,cookie:jwt",
	})
}
