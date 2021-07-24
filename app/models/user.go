package models

import (
	"context"

	"github.com/mises-id/sns/lib/db"
	"go.mongodb.org/mongo-driver/bson"
)

type Gender uint8

const (
	Male Gender = iota
	Female
)

type User struct {
	UID      uint64 `bson:"_id"`
	Username string `bson:"username,omitempty"`
	Misesid  string `bson:"misesid,omitempty"`
	Gender   Gender `bson:"gender,misesid"`
	Mobile   string `bson:"mobile,omitempty"`
	Email    string `bson:"email,omitempty"`
	Address  string `bson:"address,omitempty"`
}

func (u *User) BeforeCreate(ctx context.Context) error {
	var err error
	u.UID, err = getNextSeq(ctx, "userid")
	return err
}

func CreateUser(ctx context.Context) (*User, error) {
	user := &User{
		Username: "",
	}
	err := user.BeforeCreate(ctx)
	if err != nil {
		return nil, err
	}
	_, err = db.DB().Collection("users").InsertOne(ctx, user)
	return user, err
}

func FindUser(ctx context.Context, uid uint64) (*User, error) {
	user := &User{}
	result := db.DB().Collection("users").FindOne(ctx, &bson.M{
		"_id": uid,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	return user, result.Decode(user)
}
