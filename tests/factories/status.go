package factories

import (
	"context"
	"encoding/json"

	"github.com/bluele/factory-go/factory"
	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/app/models/meta"
	"github.com/mises-id/sns/lib/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var StatusFactory = factory.NewFactory(
	&models.Status{},
).SeqInt("ID", func(n int) (interface{}, error) {
	return primitive.NewObjectID(), nil
}).Attr("ParentID", func(args factory.Args) (interface{}, error) {
	return primitive.NilObjectID, nil
}).Attr("OriginID", func(args factory.Args) (interface{}, error) {
	return primitive.NilObjectID, nil
}).Attr("UID", func(args factory.Args) (interface{}, error) {
	return uint64(0), nil
}).Attr("StatusType", func(args factory.Args) (interface{}, error) {
	return enum.TextStatus, nil
}).Attr("Meta", func(args factory.Args) (interface{}, error) {
	return json.RawMessage{}, nil
}).Attr("Content", func(args factory.Args) (interface{}, error) {
	return "test status", nil
}).Attr("CommentsCount", func(args factory.Args) (interface{}, error) {
	return uint64(0), nil
}).Attr("LikesCount", func(args factory.Args) (interface{}, error) {
	return uint64(0), nil
}).Attr("ForwardsCount", func(args factory.Args) (interface{}, error) {
	return uint64(0), nil
}).OnCreate(func(args factory.Args) error {
	_, err := db.DB().Collection("statuses").InsertOne(context.Background(), args.Instance())
	return err
})

func InitStatuses(args ...*models.Status) {
	for _, arg := range args {
		StatusFactory.MustCreateWithOption(map[string]interface{}{
			"ID":            arg.ID,
			"ParentID":      arg.ParentID,
			"OriginID":      arg.OriginID,
			"UID":           arg.UID,
			"StatusType":    arg.StatusType,
			"Meta":          arg.Meta,
			"Content":       arg.Content,
			"CommentsCount": arg.CommentsCount,
			"LikesCount":    arg.LikesCount,
			"ForwardsCount": arg.ForwardsCount,
		})
	}
}

func InitDefaultStatuses() []*models.Status {
	statuses := make([]*models.Status, 0)
	userIDs := []uint64{1001, 1002}
	linkMeta := &meta.LinkMeta{
		Title:        "link title",
		Host:         "www.test.com",
		AttachmentID: uint64(1),
		Link:         "http://www.test.com/articles/1",
	}
	linkMetaData, _ := json.Marshal(linkMeta)
	// text status
	for _, uid := range userIDs {
		item := StatusFactory.MustCreateWithOption(map[string]interface{}{
			"UID":        uid,
			"StatusType": enum.TextStatus,
			"Meta":       json.RawMessage{},
			"Content":    "test text status",
		}).(*models.Status)
		statuses = append(statuses, item)
	}
	// link status
	for _, uid := range userIDs {
		item := StatusFactory.MustCreateWithOption(map[string]interface{}{
			"UID":        uid,
			"StatusType": enum.LinkStatus,
			"Meta":       linkMetaData,
			"Content":    "test link status",
		}).(*models.Status)
		statuses = append(statuses, item)
	}
	// forward status
	for _, uid := range userIDs {
		item1 := StatusFactory.MustCreateWithOption(map[string]interface{}{
			"UID":        uid,
			"OriginID":   statuses[0].ID,
			"ParentID":   statuses[0].ID,
			"StatusType": enum.TextStatus,
			"Meta":       linkMetaData,
			"Content":    "test forward status",
		}).(*models.Status)
		item2 := StatusFactory.MustCreateWithOption(map[string]interface{}{
			"UID":        uid,
			"OriginID":   statuses[0].ID,
			"ParentID":   item1.ID,
			"StatusType": enum.TextStatus,
			"Meta":       linkMetaData,
			"Content":    "test forward status",
		}).(*models.Status)
		statuses = append(statuses, item1, item2)
	}
	return statuses
}
