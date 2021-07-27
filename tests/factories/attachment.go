package factories

import (
	"context"
	"time"

	"github.com/bluele/factory-go/factory"
	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/lib/db"
)

var attachmentFactory = factory.NewFactory(
	&models.Attachment{},
).Attr("ID", func(args factory.Args) (interface{}, error) {
	return uint64(0), nil
}).Attr("Filename", func(args factory.Args) (interface{}, error) {
	return "", nil
}).Attr("FileType", func(args factory.Args) (interface{}, error) {
	return enum.ImageFile, nil
}).Attr("CreatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
}).Attr("UpdatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
}).OnCreate(func(args factory.Args) error {
	_, err := db.DB().Collection("attachments").InsertOne(context.Background(), args.Instance())
	return err
})

func InitAttachments(args ...*models.Attachment) {
	for _, arg := range args {
		attachmentFactory.MustCreateWithOption(map[string]interface{}{
			"ID":        arg.ID,
			"Filename":  arg.Filename,
			"FileType":  arg.FileType,
			"CreatedAt": arg.CreatedAt,
			"UpdatedAt": arg.UpdatedAt,
		})
	}
}
