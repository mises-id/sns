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
)

type CreateStatusParams struct {
	StatusType string
	ParentID   primitive.ObjectID
	OriginID   primitive.ObjectID
	Content    string
	Meta       json.RawMessage
}

func ListUserStatus(ctx context.Context, uid uint64, pageParams *pagination.PageQuickParams) ([]*models.Status, pagination.Pagination, error) {
	// check user exsist
	_, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, nil, err
	}
	return models.ListUserStatus(ctx, uid, pageParams)
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
	var originStatus, parentStatus *models.Status
	if !params.OriginID.IsZero() {
		originStatus, err = models.FindStatus(ctx, params.OriginID)
		if err != nil {
			return nil, err
		}
		if err = originStatus.IncStatusCounter(ctx, "forwards_count"); err != nil {
			return nil, err
		}
	}
	if !params.ParentID.IsZero() {
		parentStatus, err = models.FindStatus(ctx, params.OriginID)
		if err != nil {
			return nil, err
		}
		if err = parentStatus.IncStatusCounter(ctx, "forwards_count"); err != nil {
			return nil, err
		}
	}
	return models.CreateStatus(ctx, &models.CreateStatusParams{
		UID:        uid,
		StatusType: statusType,
		Content:    params.Content,
		OriginID:   params.OriginID,
		ParentID:   params.ParentID,
		MetaData:   metaData,
	})
}

func DeleteStatus(ctx context.Context, uid uint64, id primitive.ObjectID) error {
	status, err := models.FindStatus(ctx, id)
	if err != nil {
		return err
	}
	if status.UID != uid {
		return codes.ErrNotFound
	}
	return models.DeleteStatus(ctx, id)
}
