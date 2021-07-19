package rest

import (
	"net/http"

	"github.com/labstack/echo"
)

// BuildSuccessResp return a success response with payload
func BuildSuccessResp(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code": 0,
		"data": data,
	})
}

// BuildSuccessResp return a success response with payload
func BuildSuccessRespWithPagination(c echo.Context, data interface{}, pagination interface{}) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code":       0,
		"data":       data,
		"pagination": pagination,
	})
}

// Probe for k8s liveness
func Probe(c echo.Context) error {
	return BuildSuccessResp(c, nil)
}
