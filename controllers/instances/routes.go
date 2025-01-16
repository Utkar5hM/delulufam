package instances

import (
	"github.com/Utkar5hM/delulufam/controllers/authentication"
	"github.com/Utkar5hM/delulufam/utils/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func UseSubroute(g *echo.Group, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	instanceGroup := g.Group("")
	instanceGroup.Use(authentication.IsLoggedIn(cfg.JWT_SECRET))
	useInstanceRoutes(instanceGroup, db, cfg)
}

func UseOAuthServerSubroute(g *echo.Group, db *pgxpool.Pool, rdb *redis.Client, cfg *config.Config) {
	h := &instanceHandler{config.Handler{DB: db, Config: cfg, RDB: rdb}}
	g.POST("/device_authorization", h.deviceAuthorization)
	g.POST("/token", h.token)
	verify := g.Group("/verify")
	verify.Use(authentication.IsLoggedIn(cfg.JWT_SECRET))
	verify.GET("", h.VerificationPage)
	verify.POST("", h.VerifyUserCode)
}

func useInstanceRoutes(g *echo.Group, db *pgxpool.Pool, cfg *config.Config) {
	h := &instanceHandler{config.Handler{DB: db, Config: cfg}}
	g.POST("", h.createInstance)
	g.Use(h.isAdminOrCreatorMiddleware)
	g.POST("/host_user/:id", h.addInstanceHostUser)
	g.DELETE("/host_user/:id", h.deleteInstanceHostUser)
	g.PUT("/status/:id", h.setStatusInstance)
	g.POST("/access/:id", h.addUserInstanceAccess)
	g.DELETE("/access/:id", h.deleteUserInstanceAccess)
	g.DELETE("/:id", h.deleteInstance)
}
