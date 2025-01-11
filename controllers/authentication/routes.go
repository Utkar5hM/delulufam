package authentication

import (
	"net/http"

	"github.com/Utkar5hM/delulufam/utils/config"
	"github.com/Utkar5hM/delulufam/utils/render"
	"github.com/Utkar5hM/delulufam/views"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func UseSubroute(g *echo.Group, db *pgxpool.Pool, cfg *config.Config) {
	h := &authHandler{config.Handler{DB: db, Config: cfg}}
	g.POST("/oauth/google/login", h.GoogleLogin)
	g.GET("/oauth/google/callback", h.GoogleCallback)
	g.GET("/login", func(c echo.Context) error {
		return render.Render(c, http.StatusOK, views.Login())
	})
}

func IsAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*JwtCustomClaims)
		if user.Role != "admin" {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}

func IsAdmin(c echo.Context) bool {
	user := c.Get("user").(*JwtCustomClaims)
	if user.Role != "admin" {
		return true
	}
	return false
}
