package models

import (
	"context"
	"time"

	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	UID            uint64      `bson:"_id"`
	Username       string      `bson:"username,omitempty"`
	Misesid        string      `bson:"misesid,omitempty"`
	Gender         enum.Gender `bson:"gender,misesid"`
	Mobile         string      `bson:"mobile,omitempty"`
	Email          string      `bson:"email,omitempty"`
	Address        string      `bson:"address,omitempty"`
	AvatarID       uint64      `bson:"avatar_id,omitempty"`
	FollowingCount int64       `bson:"following_count,omitempty"`
	FansCount      int64       `bson:"fans_count,omitempty"`
	CreatedAt      time.Time   `bson:"created_at,omitempty"`
	UpdatedAt      time.Time   `bson:"updated_at,omitempty"`
	Avatar         *Attachment `bson:"-"`
}

func (u *User) BeforeCreate(ctx context.Context) error {
	var err error
	u.UID, err = getNextSeq(ctx, "userid")
	if err != nil {
		return err
	}
	u.CreatedAt = time.Now()
	return u.BeforeUpdate(ctx)
}

func (u *User) BeforeUpdate(ctx context.Context) error {
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) IncFollowingCount(ctx context.Context) error {
	return db.DB().Collection("users").FindOneAndUpdate(ctx, bson.M{"_id": u.UID},
		bson.D{{
			Key: "$inc",
			Value: bson.D{{
				Key:   "following_count",
				Value: 1,
			}}},
		}).Err()
}

func (u *User) IncFansCount(ctx context.Context) error {
	return db.DB().Collection("users").FindOneAndUpdate(ctx, bson.M{"_id": u.UID},
		bson.D{{
			Key: "$inc",
			Value: bson.D{{
				Key:   "fans_count",
				Value: 1,
			}}},
		}).Err()
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

func FindOrCreateUserByMisesid(ctx context.Context, misesid string) (*User, error) {
	user := &User{}
	result := db.DB().Collection("users").FindOne(ctx, &bson.M{
		"misesid": misesid,
	})
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return createMisesUser(ctx, misesid)
	}
	if err != nil {
		return nil, err
	}
	return user, result.Decode(user)
}

func UpdateUserProfile(ctx context.Context, user *User) error {
	_, err := db.DB().Collection("users").UpdateOne(ctx, &bson.M{
		"_id": user.UID,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"gender":     user.Gender,
			"mobile":     user.Mobile,
			"email":      user.Email,
			"address":    user.Address,
			"updated_at": time.Now(),
		}}})
	return err
}

func UpdateUsername(ctx context.Context, user *User) error {
	_, err := db.DB().Collection("users").UpdateOne(ctx, &bson.M{
		"_id": user.UID,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"username":   user.Username,
			"updated_at": time.Now(),
		}}})
	return err
}

func UpdateUserAvatar(ctx context.Context, user *User) error {
	_, err := db.DB().Collection("users").UpdateOne(ctx, &bson.M{
		"_id": user.UID,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"avatar_id":  user.AvatarID,
			"updated_at": time.Now(),
		}}})
	return err
}

func createMisesUser(ctx context.Context, misesid string) (*User, error) {
	user := &User{
		Misesid: misesid,
	}
	err := user.BeforeCreate(ctx)
	if err != nil {
		return nil, err
	}
	_, err = db.DB().Collection("users").InsertOne(ctx, user)
	return user, err
}

func PreloadUserAvatar(ctx context.Context, users ...*User) error {
	avatarIds := make([]uint64, 0)
	for _, user := range users {
		if user.AvatarID != 0 {
			avatarIds = append(avatarIds, user.AvatarID)
		}
	}
	attachments := make([]*Attachment, 0)
	err := db.ODM(ctx).Where(bson.M{"_id": bson.M{"$in": avatarIds}}).Find(&attachments).Error
	if err != nil {
		return err
	}
	avatarMap := make(map[uint64]*Attachment)
	for _, attachment := range attachments {
		avatarMap[attachment.ID] = attachment
	}
	for _, user := range users {
		user.Avatar = avatarMap[user.AvatarID]
	}
	return nil
}
