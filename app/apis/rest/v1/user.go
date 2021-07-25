package v1

import (
	"strconv"

	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
	svc "github.com/mises-id/sns/app/services/user"
	"github.com/mises-id/sns/lib/codes"
)

type SigninParams struct {
	Provider  string `json:"provider"`
	UserAuthz *struct {
		Misesid  string `json:"misesid"`
		AuthCode string `json:"auth_code"`
	} `json:"user_authz"`
}

type AvatarResp struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type UserResp struct {
	UID      uint64      `json:"uid"`
	Username string      `json:"username"`
	Misesid  string      `json:"misesid"`
	Gender   string      `json:"gender"`
	Mobile   string      `json:"mobile"`
	Email    string      `json:"email"`
	Address  string      `json:"address"`
	Avatar   *AvatarResp `json:"avatar"`
}

func Signin(c echo.Context) error {
	params := &SigninParams{}
	if err := c.Bind(params); err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func FindUser(c echo.Context) error {
	uidParam := c.Param("uid")
	uid, err := strconv.ParseUint(uidParam, 10, 64)
	if err != nil {
		return codes.ErrInvalidArgument.Newf("invalid uid %s", uidParam)
	}
	user, err := svc.FindUser(c.Request().Context(), uid)
	if err != nil {
		return err
	}
	resp := &UserResp{
		UID:      user.UID,
		Username: user.Username,
		Misesid:  user.Misesid,
		Gender:   user.Gender.String(),
		Mobile:   user.Mobile,
		Email:    user.Email,
		Address:  user.Address,
	}
	if user.Avatar != nil {
		resp.Avatar = &AvatarResp{
			// TODO support multiple sizes avatar
			Small:  user.Avatar.FileUrl(),
			Medium: user.Avatar.FileUrl(),
			Large:  user.Avatar.FileUrl(),
		}
	}
	return rest.BuildSuccessResp(c, resp)
}
