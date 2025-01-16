package instances

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/Utkar5hM/delulufam/controllers/authentication"
	"github.com/doug-martin/goqu/v9"
	"github.com/golang-jwt/jwt/v5"
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
		"status":      "pending",
		"clientIP":    clientIP,
		"instance_id": instanceId,
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
		"verification_uri_complete": c.Request().Host + "/instances/auth/device?user_code=" + userCode,
		"expires_in":                (10 * time.Minute) / time.Second,
		"interval":                  5,
	})
}

func (h *instanceHandler) token(c echo.Context) error {
	grantType := c.FormValue("grant_type")
	clientId := c.FormValue("client_id")
	deviceCode := c.FormValue("device_code")
	if grantType != "urn:ietf:params:oauth:grant-type:device_code" {
		return c.JSON(400, echo.Map{
			"error":             "unsupported_grant_type",
			"error_description": "Only 'urn:ietf:params:oauth:grant-type:device_code' is supported",
		})
	}
	value, err := h.RDB.HGetAll(context.Background(), deviceCode).Result()
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":             "Invalid device_code",
			"error_description": "Device code does not exist or has expired",
		})
	}
	if value["client_id"] != clientId {
		return c.JSON(400, echo.Map{
			"error":             "invalid_client",
			"error_description": "Client ID does not match the device code",
		})
	}
	if value["status"] == "pending" {
		return c.JSON(400, echo.Map{
			"error":             "authorization_pending",
			"error_description": "User has not yet completed the authorization",
		})
	}
	if value["status"] == "denied" {
		return c.JSON(400, echo.Map{
			"error":             "access_denied",
			"error_description": "User has denied the authorization",
		})
	}
	expiresAt, err := strconv.ParseInt(value["expires_at"], 10, 64)
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":             "invalid_token",
			"error_description": "Invalid expiration time format",
		})
	}
	if expiresAt < time.Now().Unix() {
		return c.JSON(400, echo.Map{
			"error":             "expired_token",
			"error_description": "Device code has expired",
		})
	}
	if value["status"] == "approved" {
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":   value["client_id"],
			"scope": value["scope"],
		})
		token, err := accessToken.SignedString([]byte(h.Config.JWT_SECRET))

		if err != nil {
			return c.JSON(400, echo.Map{
				"error":             "server_error",
				"error_description": "Failed to generate access token",
			})
		}
		c.Response().Header().Add("Pragma", "no-cache")
		c.Response().Header().Add("Cache-Control", "no-store")
		return c.JSON(200, echo.Map{
			"access_token": token,
			"token_type":   "N/A",
			"expires_in":   0,
			"scope":        value["scope"],
		})
	}

	return c.JSON(400, echo.Map{
		"error":             "invalid_request",
		"error_description": "Invalid request",
	})
}

func (h *instanceHandler) VerificationPage(c echo.Context) error {
	return nil
}

func (h *instanceHandler) VerifyUserCode(c echo.Context) error {
	userCode := c.FormValue("user_code")
	value, err := h.RDB.HGetAll(context.Background(), userCode).Result()
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":             "Invalid user_code",
			"error_description": "User code does not exist or has expired",
		})
	}
	if value["status"] != "pending" {
		return c.JSON(400, echo.Map{
			"error":             "invalid_request",
			"error_description": "User code has already been verified or denied",
		})
	}
	scopes := strings.Split(value["scope"], ",")
	firstScope := scopes[0]
	if strings.Split(firstScope, ":")[0] != "user" {
		return c.JSON(400, echo.Map{
			"error":             "invalid_scope",
			"error_description": "Invalid scope",
		})
	}
	hostUsername := strings.Split(firstScope, ":")[1]
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*authentication.JwtCustomClaims)
	userId := claims.Id
	instanceId := value["instance_id"]
	sql, _, err := goqu.From("instance_users").
		Where(
			goqu.C("instance_id").Eq(instanceId),
			goqu.C("user_id").Eq(userId),
			goqu.C("instance_host_username").Eq(hostUsername),
		).Select(goqu.COUNT("*")).ToSQL()
	fmt.Println(sql)
	row := h.DB.QueryRow(context.Background(), sql)
	var count int
	err = row.Scan(&count)
	if err != nil {
		_, err = h.RDB.HSet(context.Background(), userCode, "status", "denied").Result()
		if err != nil {
			return c.JSON(400, echo.Map{
				"error":             "Failed to verify user code",
				"error_description": err.Error(),
			})
		}
		return c.JSON(400, echo.Map{
			"error":             "Failed to verify user code",
			"error_description": err.Error(),
		})
	}
	if count == 0 {
		return c.JSON(400, echo.Map{
			"error":             "access_denied",
			"error_description": "User does not have access to the specified instance",
		})
	}
	_, err = h.RDB.HSet(context.Background(), userCode, "status", "approved").Result()
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":             "Failed to verify user code",
			"error_description": err.Error(),
		})
	}
	_, err = h.RDB.HSet(context.Background(), value["device_code"], "status", "approved").Result()
	if err != nil {
		return c.JSON(400, echo.Map{
			"error":             "Failed to verify user code",
			"error_description": err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"message": "User code verified successfully",
		"status":  "success",
	})
}
