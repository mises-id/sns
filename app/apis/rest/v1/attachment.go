package v1

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
	svc "github.com/mises-id/sns/app/services/attachment"
	"github.com/mises-id/sns/lib/codes"
)

type UploadParams struct {
	FileType string `form:"file_type"`
}

type AttachmentResp struct {
	ID       uint64 `json:"id"`
	Filename string `json:"filename"`
	FileType string `json:"file_type"`
	Url      string `json:"url"`
}

func Upload(c echo.Context) error {
	params := &UploadParams{}
	if err := c.Bind(params); err != nil {
		return codes.ErrInvalidArgument.New("invalid upload params")
	}
	file, err := c.FormFile("file")
	if err != nil {
		return codes.ErrInvalidArgument.New("receive file failed")
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	attachment, err := svc.CreateAttachment(c.Request().Context(), params.FileType, file.Filename, src)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, &AttachmentResp{
		ID:       attachment.ID,
		Filename: attachment.Filename,
		FileType: attachment.FileType.String(),
		Url:      attachment.FileUrl(),
	})
}
