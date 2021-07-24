package attachment

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/mises-id/sns/app/models"
)

func CreateAttachment(ctx context.Context, fileType, filename string, file multipart.File) (*models.Attachment, error) {
	filenames := strings.Split(filename, "/")
	tp, err := convertFileType(fileType)
	if err != nil {
		return nil, err
	}
	return models.CreateAttachment(ctx, tp, filenames[len(filenames)-1], file)
}

func convertFileType(tp string) (models.FileType, error) {
	switch tp {
	default:
		return models.ImageFile, errors.New("invalid file type")
	case "image":
		return models.ImageFile, nil
	case "video":
		return models.VideoFile, nil
	}
}
