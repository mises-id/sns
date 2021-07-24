package v1

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
	svc "github.com/mises-id/sns/app/services/attachment"
)

type AttachmentResp struct {
	ID       uint64
	Filename string
	FileType string
}

func Upload(c echo.Context) error {
	fileType := c.FormValue("file_type")
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	attachment, err := svc.CreateAttachment(c.Request().Context(), fileType, file.Filename, src)
	if err != nil {
		return err
	}
	return rest.BuildSuccessResp(c, &AttachmentResp{
		ID:       attachment.ID,
		Filename: attachment.Filename,
		FileType: attachment.FileType.String(),
	})
}
