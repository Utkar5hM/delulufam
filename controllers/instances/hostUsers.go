package instances

import (
	"context"
	"net/http"

	"github.com/doug-martin/goqu/v9"
	"github.com/labstack/echo/v4"
)

func (h *instanceHandler) addInstanceHostUser(c echo.Context) error {
	id := c.Param("id")
	username := c.FormValue("username")
	sql, _, _ := goqu.Insert("instance_host_users").Rows(
		goqu.Record{
			"instance_id": id,
			"username":    username,
		},
	).ToSQL()
	_, err := h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": "Failed to add host user: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully added host user",
	})

}

func (h *instanceHandler) deleteInstanceHostUser(c echo.Context) error {
	id := c.Param("id")
	username := c.FormValue("username")
	sql, _, _ := goqu.From("instance_host_users").Where(goqu.C("instance_id").Eq(id), goqu.C("username").Eq(username)).Delete().ToSQL()
	_, err := h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": "Failed to delete host user: " + err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status": "success",
	})
}
