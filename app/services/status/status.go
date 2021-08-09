package status

import (
	"context"
	"encoding/json"

	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/app/models/meta"
	"github.com/mises-id/sns/lib/codes"
	"github.com/mises-id/sns/lib/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateStatusParams struct {
	StatusType string
	ParentID   primitive.ObjectID
	Content    string
	Meta       json.RawMessage
	FromType   enum.FromType
}

type ListStatusParams struct {
	*pagination.PageQuickParams
	UID      uint64
	ParentID primitive.ObjectID
	FromType string
}

func ListStatus(ctx context.Context, params *ListStatusParams) ([]*models.Status, pagination.Pagination, error) {
	var fromType *enum.FromTypeFilter
	if params.FromType != "" {
		tp, err := enum.FromTypeFromString(params.FromType)
		if err != nil {
			return nil, nil, err
		}
		fromType = &enum.FromTypeFilter{FromType: tp}
	}
	return models.ListStatus(ctx, []uint64{params.UID}, params.ParentID, fromType, params.PageQuickParams)
}

func UserTimeline(ctx context.Context, uid uint64, pageParams *pagination.PageQuickParams) ([]*models.Status, pagination.Pagination, error) {
	friendIDs, err := models.ListFollowingUserIDs(ctx, uid)
	if err != nil {

	}
	return models.ListStatus(ctx, friendIDs, primitive.NilObjectID, nil, pageParams)
}

func RecommendStatus(ctx context.Context, uid uint64, pageParams *pagination.PageQuickParams) ([]*models.Status, pagination.Pagination, error) {
	// check user exsist
	return models.ListStatus(ctx, nil, primitive.NilObjectID, nil, pageParams)
}

func CreateStatus(ctx context.Context, uid uint64, params *CreateStatusParams) (*models.Status, error) {
	statusType, err := enum.StatusTypeFromString(params.StatusType)
	if err != nil {
		return nil, err
	}
	metaData, err := meta.BuildStatusMeta(statusType, params.Meta)
	if err != nil {
		return nil, err
	}
	return models.CreateStatus(ctx, &models.CreateStatusParams{
		UID:        uid,
		StatusType: statusType,
		Content:    params.Content,
		ParentID:   params.ParentID,
		FromType:   params.FromType,
		MetaData:   metaData,
	})
}

func LikeStatus(ctx context.Context, uid uint64, statusID primitive.ObjectID) (*models.Like, error) {
	status, err := models.FindStatus(ctx, statusID)
	if err != nil {
		return nil, err
	}
	like, err := models.FindLike(ctx, uid, statusID, enum.LikeStatus)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == nil {
		return like, nil
	}
	like, err = models.CreateLike(ctx, uid, statusID, enum.LikeStatus)
	if err != nil {
		return nil, err
	}
	return like, status.IncStatusCounter(ctx, "likes_count")
}

func UnlikeStatus(ctx context.Context, uid uint64, statusID primitive.ObjectID) error {
	like, err := models.FindLike(ctx, uid, statusID, enum.LikeStatus)
	if err != nil {
		return err
	}
	status, err := models.FindStatus(ctx, statusID)
	if err != nil {
		return err
	}
	if err = models.DeleteLike(ctx, like.ID); err != nil {
		return err
	}
	return status.IncStatusCounter(ctx, "likes_count")
}

func DeleteStatus(ctx context.Context, uid uint64, id primitive.ObjectID) error {
	status, err := models.FindStatus(ctx, id)
	if err != nil {
		return err
	}
	if status.UID != uid {
		return codes.ErrForbidden
	}
	return models.DeleteStatus(ctx, id)
}
