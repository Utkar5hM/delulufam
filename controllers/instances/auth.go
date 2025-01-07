package instances

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/doug-martin/goqu/v9"
	"github.com/labstack/echo/v4"
)

type Scope struct {
	Username string `json:"username"`
}

func (h *instanceHandler) deviceAuthorization(c echo.Context) error {
	clientId := c.FormValue("client_id")
	scope := c.FormValue("scope")
	clientIP := c.RealIP()

	sql, _, _ := goqu.From("instances").Where(goqu.C("client_id").Eq(clientId)).Select("id").ToSQL()

	row := h.DB.QueryRow(context.Background(), sql)
	var instanceId int
	err := row.Scan(&instanceId)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": "Invalid client_id",
		})
	}
	var scopeData Scope
	err = json.Unmarshal([]byte(scope), &scopeData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "Invalid scope",
			"message": "Scope must be a valid JSON string",
		})
	}
	return c.JSON(200, echo.Map{
		"clientId": clientId,
		"username": scopeData.Username,
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
