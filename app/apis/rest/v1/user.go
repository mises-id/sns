package v1

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
)

type SigninParams struct {
	Provider  string `json:"provider"`
	UserAuthz *struct {
		Misesid  string `json:"misesid"`
		AuthCode string `json:"auth_code"`
	} `json:"user_authz"`
}

func Signin(c echo.Context) error {
	params := &SigninParams{}
	if err := c.Bind(params); err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}
