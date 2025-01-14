package instances

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	userCodeLength = 8
	charset        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateUserCode() (string, error) {
	code := make([]byte, userCodeLength)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}
	// Optionally, insert a hyphen in the middle for readability
	return string(code[:4]) + "-" + string(code[4:]), nil
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
			"error":   "Invalid client_id",
			"message": "Instance with the specified client_id does not exist",
		})
	}

	deviceCode := uuid.New().String()
	userCode, err := generateUserCode()

	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to generate user code",
			"message": err.Error(),
		})
	}

	expiration := time.Now().Add(10 * time.Minute).Unix()

	data := map[string]interface{}{
		"client_id":   clientId,
		"scope":       scope,
		"user_code":   userCode,
		"device_code": deviceCode,
		"expires_at":  expiration,
		"approved":    false,
		"clientIP":    clientIP,
	}
	_, err = h.RDB.HSet(context.Background(), deviceCode, data).Result()
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to store device code",
			"message": err.Error(),
		})
	}
	_, err = h.RDB.HSet(context.Background(), userCode, data).Result()
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to store user code",
			"message": err.Error(),
		})
	}
	_, err = h.RDB.Expire(context.Background(), deviceCode, 10*time.Minute).Result()
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to set expiration for device code",
			"message": err.Error(),
		})
	}
	_, err = h.RDB.Expire(context.Background(), userCode, 10*time.Minute).Result()
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":   "Failed to set expiration for device code",
			"message": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"device_code":               deviceCode,
		"user_code":                 userCode,
		"verification_uri":          c.Request().Host + "/instances/auth/device",
		"verification_uri_complete": c.Request().Host + "/instances/auth/device/verify?user_code=" + userCode,
		"expires_in":                (10 * time.Minute) / time.Second,
		"interval":                  5,
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

func (h *instanceHandler) VerifyUserCodeGET(c echo.Context) error {
	return nil
}
