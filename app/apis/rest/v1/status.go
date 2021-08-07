package v1

import (
	"encoding/json"

	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
	"github.com/mises-id/sns/app/models"
	svc "github.com/mises-id/sns/app/services/status"
	"github.com/mises-id/sns/lib/codes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateSsatusParams struct {
	StatusType string             `json:"status_type"`
	ParentID   primitive.ObjectID `json:"parent_id"`
	OriginID   primitive.ObjectID `json:"origin_id"`
	Content    string             `json:"content"`
	Meta       json.RawMessage    `json:"meta"`
}

type StatusResp struct {
	ID string
}

func CreateStatus(c echo.Context) error {
	params := &CreateSsatusParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid status params")
	}
	uid := c.Get("CurrentUser").(*models.User).UID
	status, err := svc.CreateStatus(c.Request().Context(), uid, &svc.CreateStatusParams{
		StatusType: params.StatusType,
		Content:    params.Content,
		OriginID:   params.OriginID,
		ParentID:   params.ParentID,
		Meta:       params.Meta,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, buildStatusResp(status))
}

func DeleteStatus(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return codes.ErrInvalidArgument.New("invalid status id")
	}
	uid := c.Get("CurrentUser").(*models.User).UID
	err = svc.DeleteStatus(c.Request().Context(), uid, id)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func buildStatusResp(status *models.Status) *StatusResp {
	return &StatusResp{
		ID: status.ID.String(),
	}
}
