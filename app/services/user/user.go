package user

import (
	"context"

	"github.com/mises-id/sns/app/models"
)

func FindUser(ctx context.Context, uid uint64) (*models.User, error) {
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user, preloadAvatar(ctx, user)
}

func preloadAvatar(ctx context.Context, users ...*models.User) error {
	avatarIDs := make([]uint64, len(users))
	for i, user := range users {
		avatarIDs[i] = user.AvatarID
	}
	attachmentMap, err := models.FindAttachmentMap(ctx, avatarIDs)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.Avatar = attachmentMap[user.AvatarID]
	}
	return nil
}
