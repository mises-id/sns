package middleware

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns/lib/codes"
	log "github.com/sirupsen/logrus"
)

var ErrorResponseMiddleware = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			if _, ok := err.(*echo.HTTPError); ok {
				return err
			}

			code, ok := err.(codes.Code)
			if !ok {
				log.WithFields(map[string]interface{}{
					"RequestID": c.Response().Header().Get(echo.HeaderXRequestID),
				}).Error("Unkown Error:", err)
				code = codes.ErrInternal
			}

			return c.JSON(code.HTTPStatus, code)
		}
		return nil
	}
}
