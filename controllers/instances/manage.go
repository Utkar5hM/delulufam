package instances

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (h *instanceHandler) createInstance(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*auth.JwtCustomClaims)
	// username := claims.Username
	// message := fmt.Sprintf("Welcome %s!\nYour Role: %s", username, claims.Role)
	return c.JSON(http.StatusOK, echo.Map{
		"clientId": user,
	})
}

func (h *instanceHandler) deleteInstance(c echo.Context) error {
	return c.JSON(200, echo.Map{
		"clientId": "test",
	})
}
