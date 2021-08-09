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
	CurrentUID uint64
	UID        uint64
	ParentID   primitive.ObjectID
	FromType   string
}

func GetStatus(ctx context.Context, currentUID uint64, id primitive.ObjectID) (*models.Status, error) {
	status, err := models.FindStatus(ctx, id)
	if err != nil {
		return nil, err
	}
	return status, batchSetIsLiked(ctx, currentUID, status)
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
	statues, page, err := models.ListStatus(ctx, []uint64{params.UID}, params.ParentID, fromType, params.PageQuickParams)
	if err != nil {
		return nil, nil, err
	}
	return statues, page, batchSetIsLiked(ctx, params.CurrentUID, statues...)
}

func UserTimeline(ctx context.Context, uid uint64, pageParams *pagination.PageQuickParams) ([]*models.Status, pagination.Pagination, error) {
	friendIDs, err := models.ListFollowingUserIDs(ctx, uid)
	if err != nil {
		return nil, nil, err
	}
	if len(friendIDs) == 0 {
		return []*models.Status{}, &pagination.QuickPagination{
			Limit: pageParams.Limit,
		}, nil
	}

	statues, page, err := models.ListStatus(ctx, friendIDs, primitive.NilObjectID, nil, pageParams)
	if err != nil {
		return nil, nil, err
	}
	return statues, page, batchSetIsLiked(ctx, uid, statues...)
}

func RecommendStatus(ctx context.Context, uid uint64, pageParams *pagination.PageQuickParams) ([]*models.Status, pagination.Pagination, error) {
	statues, page, err := models.ListStatus(ctx, nil, primitive.NilObjectID, nil, pageParams)
	if err != nil {
		return nil, nil, err
	}
	return statues, page, batchSetIsLiked(ctx, uid, statues...)
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
	return status.IncStatusCounter(ctx, "likes_count", -1)
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

func batchSetIsLiked(ctx context.Context, uid uint64, statuses ...*models.Status) error {
	if uid == 0 {
		return nil
	}
	statusIDs := make([]primitive.ObjectID, len(statuses))
	for i, status := range statuses {
		statusIDs[i] = status.ID
	}
	likeMap, err := models.GetStatusLikeMap(ctx, uid, statusIDs)
	if err != nil {
		return err
	}
	for _, status := range statuses {
		status.IsLiked = likeMap[status.ID] != nil
	}
	return nil
}
