package models

import (
	"context"
	"time"

	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Like struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty"`
	UID        uint64              `bson:"uid,omitempty"`
	TargetID   primitive.ObjectID  `bson:"target_id,omitempty"`
	TargetType enum.LikeTargetType `bson:"target_type,omitempty"`
	DeletedAt  time.Time           `bson:"deleted_at,omitempty"`
	CreatedAt  time.Time           `bson:"created_at,omitempty"`
	UpdatedAt  time.Time           `bson:"updated_at,omitempty"`
}

func CreateLike(ctx context.Context, uid uint64, targetID primitive.ObjectID, targetType enum.LikeTargetType) (*Like, error) {
	like := &Like{
		UID:        uid,
		TargetID:   targetID,
		TargetType: targetType,
	}
	return like, db.ODM(ctx).Create(like).Error
}

func DeleteLike(ctx context.Context, id primitive.ObjectID) error {
	return db.DB().Collection("counters").FindOneAndUpdate(ctx, bson.M{
		"_id":        id,
		"deleted_at": nil,
	}, bson.M{"deleted_at": time.Now()}).Err()
}

func FindLike(ctx context.Context, uid uint64, targetID primitive.ObjectID, targetType enum.LikeTargetType) (*Like, error) {
	like := &Like{}
	return like, db.ODM(ctx).Where(bson.M{
		"uid":         uid,
		"target_id":   targetID,
		"target_type": targetType,
		"deleted_at":  nil,
	}).First(like).Error
}
