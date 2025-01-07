package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Utkar5hM/delulufam/controllers/auth"
	"github.com/Utkar5hM/delulufam/controllers/instances"
	"github.com/Utkar5hM/delulufam/utils/config"
	"github.com/Utkar5hM/delulufam/utils/render"
	"github.com/Utkar5hM/delulufam/views"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	username := claims.Username
	message := fmt.Sprintf("Welcome %s!\nYour Role: %s", username, claims.Role)
	return c.String(http.StatusOK, message)
	// return c.String(http.StatusOK, "Welcome!")
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Connection to Database first. :)
	dbpool, err := pgxpool.New(context.Background(), cfg.DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.Secure())

	e.GET("/", func(c echo.Context) error {
		return render.Render(c, http.StatusOK, views.Index())
	})

	e.GET("/view", func(c echo.Context) error {
		return render.Render(c, http.StatusOK, views.ViewPlaylist())
	})
	authGroup := e.Group("/users")
	auth.UseSubroute(authGroup, dbpool, cfg)

	instanceGroup := e.Group("/instances")
	instanceGroup.Use(auth.IsLoggedIn(cfg.JWT_SECRET))
	instances.UseSubroute(instanceGroup, dbpool, cfg)
	// instanceGroup.Use(auth.IsLoggedIn(cfg.JWT_SECRET))

	e.Static("/static/", "static")
	// Restricted group
	r := e.Group("/restricted")

	// Configure middleware with the custom claims type
	r.Use(auth.IsLoggedIn(cfg.JWT_SECRET))

	// r.Use(echojwt.WithConfig(config))
	r.GET("", restricted)
	e.Logger.Fatal(e.Start(":4000"))

}
