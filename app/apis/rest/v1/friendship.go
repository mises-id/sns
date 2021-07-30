package v1

import (
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	followSVC "github.com/mises-id/sns/app/services/follow"
	"github.com/mises-id/sns/lib/codes"
	"github.com/mises-id/sns/lib/pagination"
)

type ListFriendshipParams struct {
	pagination.TraditionalParams
	Relate string `query:"relate"`
}

type FriendshipResp struct {
	User      *UserResp `json:"user"`
	Relate    string    `json:"relate"`
	CreatedAt time.Time `json:"created_at"`
}

func ListFriendship(c echo.Context) error {
	uidParam := c.Param("uid")
	uid, err := strconv.ParseUint(uidParam, 10, 64)
	if err != nil {
		return codes.ErrInvalidArgument.Newf("invalid uid %s", uidParam)
	}
	params := &ListFriendshipParams{}
	if err := c.Bind(params); err != nil {
		return err
	}
	relate, err := enum.RelationTypeFromString(params.Relate)
	if err != nil {
		relate = enum.Fan
	}
	follows, page, err := followSVC.ListFriendship(c.Request().Context(), uid, relate, &params.TraditionalParams)
	if err != nil {
		return err
	}
	resp := batchBuildFriendshipResp(relate, follows)
	return rest.BuildSuccessRespWithPagination(c, resp, page.BuildJSONResult())
}

func Follow(c echo.Context) error {
	uid := c.Get("CurrentUser").(*models.User).UID
	uidParam := c.QueryParam("to_user")
	toUserID, err := strconv.ParseUint(uidParam, 10, 64)
	if err != nil {
		return codes.ErrInvalidArgument.Newf("invalid to user id %s", uidParam)
	}
	_, err = followSVC.Follow(c.Request().Context(), uid, toUserID)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func Unfollow(c echo.Context) error {
	uid := c.Get("CurrentUser").(*models.User).UID
	uidParam := c.QueryParam("to_user")
	toUserID, err := strconv.ParseUint(uidParam, 10, 64)
	if err != nil {
		return codes.ErrInvalidArgument.Newf("invalid to user id %s", uidParam)
	}
	err = followSVC.Unfollow(c.Request().Context(), uid, toUserID)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func batchBuildFriendshipResp(relate enum.RelationType, friendships []*models.Follow) []*FriendshipResp {
	resp := make([]*FriendshipResp, len(friendships))
	for i, friendship := range friendships {
		user := friendship.ToUser
		currentRelationType := enum.Following
		if relate == enum.Fan {
			user = friendship.FromUser
			currentRelationType = enum.Fan
		}
		if friendship.IsFriend {
			currentRelationType = enum.Friend
		}
		resp[i] = &FriendshipResp{
			Relate:    currentRelationType.String(),
			CreatedAt: friendship.CreatedAt,
			User:      buildUser(user),
		}
	}
	return resp
}
