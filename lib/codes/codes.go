package codes

import (
	"fmt"
	"net/http"
)

// Code is a error of error code
type Code struct {
	HTTPStatus int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"message"`
}

// Error return error message
func (e Code) Error() string {
	return e.Msg
}

// Equal return weather it is the same code
func (e Code) Equal(err error) bool {
	sourceErr, ok := err.(Code)
	if !ok {
		return false
	}

	return sourceErr.Code == e.Code
}

// New will modify the msg of Code
func (e Code) New(msg string) Code {
	return Code{
		HTTPStatus: e.HTTPStatus,
		Code:       e.Code,
		Msg:        msg,
	}
}

// Newf will modify the msg of Code with params
func (e Code) Newf(msg string, args ...interface{}) Code {
	return Code{
		HTTPStatus: e.HTTPStatus,
		Code:       e.Code,
		Msg:        fmt.Sprintf(msg, args...),
	}
}

const (
	SuccessCode         = 0
	InvalidArgumentCode = 400000
	UnauthorizedCode    = 401000
	ForbiddenCode       = 403000
	NotFoundCode        = 404000
	InternalCode        = 500000
	UnimplementedCode   = 500001
)

var (
	Success            = Code{HTTPStatus: http.StatusOK, Code: SuccessCode, Msg: "success"}
	ErrInvalidArgument = Code{HTTPStatus: http.StatusBadRequest, Code: InvalidArgumentCode, Msg: "invalid params"}
	ErrUnauthorized    = Code{HTTPStatus: http.StatusUnauthorized, Code: UnauthorizedCode, Msg: "unauthorized"}
	ErrForbidden       = Code{HTTPStatus: http.StatusForbidden, Code: ForbiddenCode, Msg: "forbidden"}
	ErrNotFound        = Code{HTTPStatus: http.StatusNotFound, Code: NotFoundCode, Msg: "not found"}
	ErrInternal        = Code{HTTPStatus: http.StatusInternalServerError, Code: InternalCode, Msg: "internal error"}
	ErrUnimplemented   = Code{HTTPStatus: http.StatusInternalServerError, Code: InternalCode, Msg: "internal error"}
)
