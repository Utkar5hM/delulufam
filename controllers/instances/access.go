package instances

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/labstack/echo/v4"
)

func (h *instanceHandler) addUserInstanceAccess(c echo.Context) error {
	instanceId := c.Param("id")
	username := c.FormValue("username")
	instanceHostUsername := c.FormValue("host_username")
	sql, _, _ := goqu.From("users").Where(goqu.C("username").Eq(username)).Select("id").ToSQL()
	row := h.DB.QueryRow(context.Background(), sql)
	var userID string
	err := row.Scan(&userID)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to add access for user to Instance.",
			"message": err.Error(),
			"status":  "error",
		})
	}
	sql, _, _ = goqu.Insert("instance_users").Rows(
		goqu.Record{
			"instance_id":            instanceId,
			"user_id":                userID,
			"instance_host_username": instanceHostUsername,
		},
	).ToSQL()
	_, err = h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to add access for user to Instance.",
			"message": err.Error(),
			"status":  "error",
		})
	}
	return c.JSON(200, echo.Map{
		"status": "success",
	})
}

func (h *instanceHandler) deleteUserInstanceAccess(c echo.Context) error {
	instanceId := c.Param("id")
	username := c.FormValue("username")
	hostUsername := c.FormValue("host_username")
	sql, _, _ := goqu.From("users").Where(goqu.C("username").Eq(username)).Select("id").ToSQL()
	row := h.DB.QueryRow(context.Background(), sql)
	var userID string
	err := row.Scan(&userID)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to delete instance access to the user.",
			"message": err.Error(),
			"status":  "error",
		})
	}
	sql, _, _ = goqu.From("instance_users").Where(goqu.C("instance_id").Eq(instanceId), goqu.C("user_id").Eq(userID), goqu.C("instance_host_username").Eq(hostUsername)).Delete().ToSQL()
	_, err = h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to delete instance access to the user.",
			"message": err.Error(),
			"status":  "error",
		})
	}
	return c.JSON(200, echo.Map{
		"status": "success",
	})
}
