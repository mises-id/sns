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
	pagination.QuickPagination
	RelationType string `query:"relation_type"`
}

type FollowParams struct {
	ToUserID uint64 `json:"to_user_id" query:"to_user_id"`
}

type FriendshipResp struct {
	User         *UserResp `json:"user"`
	RelationType string    `json:"relation_type"`
	CreatedAt    time.Time `json:"created_at"`
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
	relationType, err := enum.RelationTypeFromString(params.RelationType)
	if err != nil {
		relationType = enum.Fan
	}
	follows, page, err := followSVC.ListFriendship(c.Request().Context(), uid, relationType, &params.QuickPagination)
	if err != nil {
		return err
	}
	resp := batchBuildFriendshipResp(relationType, follows)
	return rest.BuildSuccessRespWithPagination(c, resp, page.BuildJSONResult())
}

func Follow(c echo.Context) error {
	uid := c.Get("CurrentUser").(*models.User).UID
	params := &FollowParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument
	}
	_, err := followSVC.Follow(c.Request().Context(), uid, params.ToUserID)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func Unfollow(c echo.Context) error {
	uid := c.Get("CurrentUser").(*models.User).UID
	params := &FollowParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument
	}
	err := followSVC.Unfollow(c.Request().Context(), uid, params.ToUserID)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func batchBuildFriendshipResp(relationType enum.RelationType, friendships []*models.Follow) []*FriendshipResp {
	resp := make([]*FriendshipResp, len(friendships))
	for i, friendship := range friendships {
		user := friendship.ToUser
		currentRelationType := enum.Following
		if relationType == enum.Fan {
			user = friendship.FromUser
			currentRelationType = enum.Fan
		}
		if friendship.IsFriend {
			currentRelationType = enum.Friend
		}
		resp[i] = &FriendshipResp{
			RelationType: currentRelationType.String(),
			CreatedAt:    friendship.CreatedAt,
			User:         buildUserResp(user),
		}
	}
	return resp
}
