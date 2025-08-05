package middlewares

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"go_api/cmd/common"
	"go_api/cmd/internal/models"
	"strings"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

type AppMiddleware struct {
	Logger echo.Logger
	DB     *gorm.DB
}

func (appMiddleware *AppMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if strings.HasPrefix(authHeader, "Bearer") == false {
			return common.SendFailedValidationResponse(c, "Unauthorized", "Please provide a Bearer token")
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := common.ParseJWTAccessToken(accessToken)
		if err != nil {
			spew.Dump(66666)
			return common.SendFailedValidationResponse(c, "Unauthorized", err.Error())
		}

		if common.IsClaimExpired(claims) {
			return common.SendFailedValidationResponse(c, "Unauthorized", "Token has expired")
		}

		var user models.User
		result := appMiddleware.DB.First(&user, claims.ID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return common.SendFailedValidationResponse(c, "Unauthorized", "Invalid user")
		}

		if result.Error != nil {
			return common.SendFailedValidationResponse(c, "Unauthorized", "Invalid access token")
		}
		c.Set("user", user)
		return next(c)
	}
}
