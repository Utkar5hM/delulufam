package auth

import (
	"net/http"

	"github.com/Utkar5hM/delulufam/utils/config"
	"github.com/Utkar5hM/delulufam/utils/render"
	"github.com/Utkar5hM/delulufam/views"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func UseSubroute(g *echo.Group, db *pgxpool.Pool, cfg *config.Config) {
	h := &Handler{DB: db, config: cfg}
	g.POST("/oauth/google/login", h.GoogleLogin)
	g.GET("/oauth/google/callback", h.GoogleCallback)
	g.GET("/login", func(c echo.Context) error {
		return render.Render(c, http.StatusOK, views.Login())
	})
}
