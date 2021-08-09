package v1

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/app/models/meta"
	svc "github.com/mises-id/sns/app/services/status"
	"github.com/mises-id/sns/lib/codes"
	"github.com/mises-id/sns/lib/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListUserStatusParams struct {
	*pagination.PageQuickParams
}

type CreateStatusParams struct {
	StatusType string             `json:"status_type"`
	ParentID   primitive.ObjectID `json:"parent_status_id"`
	Content    string             `json:"content"`
	Meta       json.RawMessage    `json:"meta"`
}

type LinkMeta struct {
	Title         string `json:"title"`
	Host          string `json:"host"`
	Link          string `json:"link"`
	AttachmentID  uint64 `json:"attachment_id"`
	AttachmentURL string `json:"attachment_url"`
}

type StatusResp struct {
	ID            string      `json:"id"`
	User          *UserResp   `json:"user"`
	Content       string      `json:"content"`
	FromType      string      `json:"from_type"`
	StatusType    string      `json:"status_type"`
	ParentStatus  *StatusResp `json:"parent_status"`
	OriginStatus  *StatusResp `json:"origin_status"`
	CommentsCount uint64      `json:"comments_count"`
	LikesCount    uint64      `json:"likes_count"`
	ForwardsCount uint64      `json:"forwards_count"`
	LinkMeta      *LinkMeta   `json:"link_meta"`
	CreatedAt     time.Time   `json:"created_at"`
}

// list user status
func ListUserStatus(c echo.Context) error {
	uidParam := c.Param("uid")
	uid, err := strconv.ParseUint(uidParam, 10, 64)
	if err != nil {
		return codes.ErrInvalidArgument.Newf("invalid uid %s", uidParam)
	}
	params := &ListUserStatusParams{}
	if err = c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	statuses, page, err := svc.ListStatus(c.Request().Context(), &svc.ListStatusParams{
		PageQuickParams: params.PageQuickParams,
		UID:             uid,
	})
	if err != nil {
		return err
	}
	resp := batchBuildStatusResp(statuses)
	return rest.BuildSuccessRespWithPagination(c, resp, page.BuildJSONResult())
}

func Timeline(c echo.Context) error {
	uid := c.Get("CurrentUser").(*models.User).UID
	params := &ListUserStatusParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	statuses, page, err := svc.UserTimeline(c.Request().Context(), uid, params.PageQuickParams)
	if err != nil {
		return err
	}
	resp := batchBuildStatusResp(statuses)
	return rest.BuildSuccessRespWithPagination(c, resp, page.BuildJSONResult())
}

func RecommendStatus(c echo.Context) error {
	iuser := c.Get("CurrentUser")
	var uid uint64
	if iuser != nil {
		uid = iuser.(*models.User).UID
	}

	params := &ListUserStatusParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid query params")
	}
	statuses, page, err := svc.RecommendStatus(c.Request().Context(), uid, params.PageQuickParams)
	if err != nil {
		return err
	}
	resp := batchBuildStatusResp(statuses)
	return rest.BuildSuccessRespWithPagination(c, resp, page.BuildJSONResult())
}

func CreateStatus(c echo.Context) error {
	params := &CreateStatusParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid status params")
	}
	uid := c.Get("CurrentUser").(*models.User).UID
	fromType := enum.FromPost
	if !params.ParentID.IsZero() {
		fromType = enum.FromForward
	}
	status, err := svc.CreateStatus(c.Request().Context(), uid, &svc.CreateStatusParams{
		StatusType: params.StatusType,
		Content:    params.Content,
		ParentID:   params.ParentID,
		Meta:       params.Meta,
		FromType:   fromType,
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

func LikeStatus(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return codes.ErrInvalidArgument.New("invalid status id")
	}
	uid := c.Get("CurrentUser").(*models.User).UID
	_, err = svc.LikeStatus(c.Request().Context(), uid, id)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func UnlikeStatus(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return codes.ErrInvalidArgument.New("invalid status id")
	}
	uid := c.Get("CurrentUser").(*models.User).UID
	err = svc.UnlikeStatus(c.Request().Context(), uid, id)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, nil)
}

func batchBuildStatusResp(statuses []*models.Status) []*StatusResp {
	result := make([]*StatusResp, len(statuses))
	for i, status := range statuses {
		result[i] = buildStatusResp(status)
	}
	return result
}

func buildStatusResp(status *models.Status) *StatusResp {
	if status == nil {
		return nil
	}
	resp := &StatusResp{
		ID:            status.ID.Hex(),
		User:          buildUserResp(status.User),
		Content:       status.Content,
		FromType:      status.FromType.String(),
		StatusType:    status.StatusType.String(),
		ParentStatus:  buildStatusResp(status.ParentStatus),
		OriginStatus:  buildStatusResp(status.OriginStatus),
		CommentsCount: status.CommentsCount,
		LikesCount:    status.LikesCount,
		ForwardsCount: status.ForwardsCount,
		CreatedAt:     status.CreatedAt,
	}
	return resp
}

func buildLinkMeta(meta *meta.LinkMeta) *LinkMeta {
	if meta == nil {
		return nil
	}
	return &LinkMeta{
		Title:         meta.Title,
		Host:          meta.Host,
		Link:          meta.Link,
		AttachmentID:  meta.AttachmentID,
		AttachmentURL: meta.AttachmentURL,
	}
}
