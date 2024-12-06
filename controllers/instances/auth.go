package instances

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/labstack/echo/v4"
)

func (h *instanceHandler) deviceAuthorization(c echo.Context) error {
	clientId := c.FormValue("client_id")
	username := c.FormValue("username")
	clientIP := c.RealIP()

	sql, _, _ := goqu.From("instances").Where(goqu.C("client_id").Eq(clientId)).Select("client_id").ToSQL()

	row := h.DB.QueryRow(context.Background(), sql)
	var instanceId int
	err := row.Scan(&instanceId)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": "Invalid client_id",
		})
	}
	return c.JSON(200, echo.Map{
		"clientId": clientId,
		"username": username,
		"clientIP": clientIP,
	})
}

func (h *instanceHandler) token(c echo.Context) error {
	return nil
}

func (h *instanceHandler) VerificationPage(c echo.Context) error {
	return nil
}

func (h *instanceHandler) VerifyUserCode(c echo.Context) error {
	return nil
}
