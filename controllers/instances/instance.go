package instances

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Utkar5hM/delulufam/controllers/authentication"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (h *instanceHandler) createInstance(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)
	name := c.FormValue("name")
	description := c.FormValue("description")
	host_address := c.FormValue("host_address")
	sql, _, _ := goqu.Insert("instances").Rows(
		goqu.Record{
			"name":         name,
			"description":  description,
			"host_address": host_address,
			"created_by":   claims.Id,
		},
	).ToSQL()
	fmt.Println(sql)
	_, err := h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": "Failed to create instance: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully created instance",
	})
}

func (h *instanceHandler) deleteInstance(c echo.Context) error {
	id := c.Param("id")
	sql, _, _ := goqu.From("instances").Where(goqu.C("id").Eq(id)).Delete().ToSQL()
	_, err := h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error": "Failed to delete instance: " + err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status": "success",
	})
}

func (h *instanceHandler) setStatusInstance(c echo.Context) error {
	id := c.Param("id")
	status := c.FormValue("status")
	sql, _, _ := goqu.From("instances").Where(goqu.C("id").Eq(id)).Update().Set(goqu.Record{"status": status}).ToSQL()
	_, err := h.DB.Exec(context.Background(), sql)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   fmt.Sprintf("Failed to set instance status to %s.", status),
			"details": err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"status": "success",
	})
}

func (h *instanceHandler) isAdminOrCreatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*authentication.JwtCustomClaims)
		if claims.Role != "admin" {
			id := c.Param("id")
			sql, _, _ := goqu.From("instances").Where(goqu.C("id").Eq(id)).Select("created_by").ToSQL()
			row := h.DB.QueryRow(context.Background(), sql)
			var instance_created_by int
			err := row.Scan(&instance_created_by)
			if err != nil {
				return c.JSON(400, echo.Map{
					"message": "Failed to Authroize your action.",
					"error":   err.Error(),
				})
			}
			if claims.Id != instance_created_by {
				return c.JSON(403, echo.Map{
					"message": "You are not authorized to disable this instance",
				})
			}
		}
		return next(c)
	}
}
