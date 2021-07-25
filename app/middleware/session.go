package middleware

import (
	"strings"

	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/services/session"
	"github.com/mises-id/sns/lib/codes"
)

var (
	validAuthMethods = []string{
		"Bearer",
	}
)

var SetCurrentUserMiddleware = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")
		if authorization != "" {
			strs := strings.Split(authorization, " ")
			if err := validateAuthToken(strs); err != nil {
				return err
			}

			user, err := session.Auth(c.Request().Context(), strs[1])
			if err != nil {
				return err
			}
			c.Set("CurrentUser", user)
		}

		return next(c)
	}
}

var RequireCurrentUserMiddleware = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("CurrentUser").(*models.User)
		if !ok || user == nil {
			return codes.ErrUnauthorized
		}
		return next(c)
	}
}

func validateAuthToken(strs []string) error {
	if len(strs) != 2 {
		return codes.ErrInvalidAuth
	}
	authMethod, authToken := strs[0], strs[1]
	if len(authToken) > 1000 {
		return codes.ErrInvalidAuthToken
	}
	for _, m := range validAuthMethods {
		if m == authMethod {
			return nil
		}
	}
	return codes.ErrInvalidAuthMethod
}
