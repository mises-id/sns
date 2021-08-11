package v1

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	svc "github.com/mises-id/sns/app/services/status"
	"github.com/mises-id/sns/lib/codes"
	"github.com/mises-id/sns/lib/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCommentParams struct {
	CommentableID   primitive.ObjectID `json:"commentable_id"`
	CommentableType string             `json:"commentable_type"`
	Content         string             `json:"content"`
}

type ListCommentParams struct {
	pagination.PageQuickParams
	CommentableID   primitive.ObjectID `query:"commentable_id"`
	CommentableType string             `query:"commentable_type"`
}

func ListComment(c echo.Context) error {
	params := &ListCommentParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}

	statuses, page, err := svc.ListStatus(c.Request().Context(), &svc.ListStatusParams{
		PageQuickParams: &params.PageQuickParams,
		ParentID:        params.CommentableID,
		FromType:        enum.FromComment.String(),
	})
	if err != nil {
		return err
	}
	resp := batchBuildStatusResp(statuses)
	return rest.BuildSuccessRespWithPagination(c, resp, page.BuildJSONResult())
}

func CreateComment(c echo.Context) error {
	params := &CreateCommentParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid comment params")
	}
	uid := c.Get("CurrentUser").(*models.User).UID
	status, err := svc.CreateStatus(c.Request().Context(), uid, &svc.CreateStatusParams{
		StatusType: enum.TextStatus.String(),
		Content:    params.Content,
		ParentID:   params.CommentableID,
		FromType:   enum.FromComment,
	})
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, buildStatusResp(status))
}
