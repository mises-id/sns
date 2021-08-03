package factories

import (
	"context"

	"github.com/bluele/factory-go/factory"
	"github.com/google/uuid"
	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/lib/db"
)

var UserFactory = factory.NewFactory(
	&models.User{},
).SeqInt("UID", func(n int) (interface{}, error) {
	return uint64(n), nil
}).Attr("Username", func(args factory.Args) (interface{}, error) {
	return "", nil
}).Attr("Gender", func(args factory.Args) (interface{}, error) {
	return enum.GenderFemale, nil
}).Attr("Misesid", func(args factory.Args) (interface{}, error) {
	return uuid.New().String(), nil
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
		UserFactory.MustCreateWithOption(map[string]interface{}{
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
