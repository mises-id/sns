package factories

import (
	"context"

	"github.com/bluele/factory-go/factory"
	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/lib/db"
)

var userFactory = factory.NewFactory(
	&models.User{},
).Attr("UID", func(args factory.Args) (interface{}, error) {
	return uint64(0), nil
}).Attr("Username", func(args factory.Args) (interface{}, error) {
	return "", nil
}).Attr("Misesid", func(args factory.Args) (interface{}, error) {
	return "", nil
}).Attr("Gender", func(args factory.Args) (interface{}, error) {
	return enum.GenderOther, nil
}).Attr("Mobile", func(args factory.Args) (interface{}, error) {
	return "", nil
}).Attr("Email", func(args factory.Args) (interface{}, error) {
	return "", nil
}).Attr("Address", func(args factory.Args) (interface{}, error) {
	return "", nil
}).Attr("AvatarID", func(args factory.Args) (interface{}, error) {
	return uint64(0), nil
}).OnCreate(func(args factory.Args) error {
	_, err := db.DB().Collection("users").InsertOne(context.Background(), args.Instance())
	return err
})

func InitUsers(args ...*models.User) {
	for _, arg := range args {
		userFactory.MustCreateWithOption(map[string]interface{}{
			"UID":      arg.UID,
			"Username": arg.Username,
			"Misesid":  arg.Misesid,
			"Gender":   arg.Gender,
			"Mobile":   arg.Mobile,
			"Email":    arg.Email,
			"Address":  arg.Address,
			"AvatarID": arg.AvatarID,
		})
	}
}
