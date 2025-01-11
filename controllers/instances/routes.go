package instances

import (
	"github.com/Utkar5hM/delulufam/utils/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func UseSubroute(g *echo.Group, db *pgxpool.Pool, cfg *config.Config) {
	authGroup := g.Group("/auth")
	useAuthSubroute(authGroup, db, cfg)
	useInstanceRoutes(g, db, cfg)
}

func useAuthSubroute(g *echo.Group, db *pgxpool.Pool, cfg *config.Config) {
	h := &instanceHandler{config.Handler{DB: db, Config: cfg}}
	g.POST("/device_authorization", h.deviceAuthorization)
	g.POST("/token", h.token)
	g.GET("/device", h.VerificationPage)
	g.POST("/device/verify", h.VerifyUserCode)
}

func useInstanceRoutes(g *echo.Group, db *pgxpool.Pool, cfg *config.Config) {
	h := &instanceHandler{config.Handler{DB: db, Config: cfg}}
	g.POST("", h.createInstance)
	controlGroup := g.Group("")
	controlGroup.Use(h.isAdminOrCreatorMiddleware)
	controlGroup.PUT("/status/:id", h.setStatusInstance)
	controlGroup.DELETE("/:id", h.deleteInstance)
}
