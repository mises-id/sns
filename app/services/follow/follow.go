package follow

import (
	"context"

	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/lib/pagination"
	"go.mongodb.org/mongo-driver/mongo"
)

func ListFriendship(ctx context.Context, uid uint64, relationType enum.RelationType, pageParams *pagination.TraditionalParams) ([]*models.Follow, pagination.Pagination, error) {
	// check user exsit
	_, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, nil, err
	}
	return models.ListFollow(ctx, uid, relationType, pageParams)
}

func Follow(ctx context.Context, uid, focusUserID uint64) (*models.Follow, error) {
	isFriend := false
	follow, err := models.GetFollow(ctx, uid, focusUserID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	fansFollow, err := models.GetFollow(ctx, focusUserID, uid)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == nil {
		isFriend = true
		if !fansFollow.IsFriend {
			if err = fansFollow.SetFriend(ctx, true); err != nil {
				return nil, err
			}
		}
	}
	if follow != nil {
		return follow, follow.SetFriend(ctx, isFriend)
	}
	return models.CreateFollow(ctx, uid, focusUserID, isFriend)
}

func Unfollow(ctx context.Context, fromUID, toUID uint64) error {
	_, err := models.GetFollow(ctx, fromUID, toUID)
	if err != nil {
		return nil
	}
	fansFollow, err := models.GetFollow(ctx, toUID, fromUID)
	if err == nil && fansFollow.IsFriend {
		if err = fansFollow.SetFriend(ctx, false); err != nil {
			return err
		}
	} else if err != mongo.ErrNoDocuments {
		return err
	}
	return models.DeleteFollow(ctx, fromUID, toUID)
}
