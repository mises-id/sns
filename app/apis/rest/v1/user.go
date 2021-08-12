package v1

import (
	"strconv"

	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	sessionSVC "github.com/mises-id/sns/app/services/session"
	svc "github.com/mises-id/sns/app/services/user"
	"github.com/mises-id/sns/lib/codes"
)

type SignInParams struct {
	Provider  string `json:"provider"`
	UserAuthz *struct {
		Auth string `json:"auth"`
	} `json:"user_authz"`
}

type AvatarResp struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type UserResp struct {
	UID        uint64      `json:"uid"`
	Username   string      `json:"username"`
	Misesid    string      `json:"misesid"`
	Gender     string      `json:"gender"`
	Mobile     string      `json:"mobile"`
	Email      string      `json:"email"`
	Address    string      `json:"address"`
	Avatar     *AvatarResp `json:"avatar"`
	IsFollowed bool        `json:"is_followed"`
}

func SignIn(c echo.Context) error {
	params := &SignInParams{}
	if err := c.Bind(params); err != nil {
		return err
	}
	token, err := sessionSVC.SignIn(c.Request().Context(), params.UserAuthz.Auth)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, echo.Map{
		"token": token,
	})
}

func MyProfile(c echo.Context) error {
	uid := c.Get("CurrentUser").(*models.User).UID
	user, err := svc.FindUser(c.Request().Context(), uid)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, buildUserResp(user))
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
	return rest.BuildSuccessResp(c, buildUserResp(user))
}

type UserProfileParams struct {
	Gender  string `json:"gender"`
	Mobile  string `json:"mobile"`
	Eamil   string `json:"email"`
	Address string `json:"address"`
}

type UserNameParams struct {
	Username string `json:"username"`
}

type UserAvatarParams struct {
	AttachmentID uint64 `json:"attachment_id"`
}

type UserUpdateParams struct {
	By       string             `json:"by"`
	Profile  *UserProfileParams `json:"profile"`
	Username *UserNameParams    `json:"username"`
	Avatar   *UserAvatarParams  `json:"avatar"`
}

func UpdateUser(c echo.Context) error {
	uid := c.Get("CurrentUser").(*models.User).UID
	params := &UserUpdateParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument
	}
	var user *models.User
	var err error
	switch params.By {
	default:
		return codes.ErrInvalidArgument
	case "profile":
		gender, err := enum.GenderFromString(params.Profile.Gender)
		if err != nil {
			return codes.ErrInvalidArgument
		}
		user, err = svc.UpdateUserProfile(c.Request().Context(), uid, &svc.UserProfileParams{
			Gender:  gender,
			Mobile:  params.Profile.Mobile,
			Email:   params.Profile.Eamil,
			Address: params.Profile.Address,
		})
	case "avatar":
		user, err = svc.UpdateUserAvatar(c.Request().Context(), uid, params.Avatar.AttachmentID)
	case "username":
		user, err = svc.UpdateUsername(c.Request().Context(), uid, params.Username.Username)
	}
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, buildUserResp(user))
}

func buildUserResp(user *models.User) *UserResp {
	if user == nil {
		return nil
	}
	resp := &UserResp{
		UID:        user.UID,
		Username:   user.Username,
		Misesid:    user.Misesid,
		Gender:     user.Gender.String(),
		Mobile:     user.Mobile,
		Email:      user.Email,
		Address:    user.Address,
		IsFollowed: user.IsFollowed,
	}
	if user.Avatar != nil {
		resp.Avatar = &AvatarResp{
			// TODO support multiple sizes avatar
			Small:  user.Avatar.FileUrl(),
			Medium: user.Avatar.FileUrl(),
			Large:  user.Avatar.FileUrl(),
		}
	}
	return resp
}
