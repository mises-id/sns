package models

import (
	"context"
	"time"

	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/lib/db"
	"github.com/mises-id/sns/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Follow struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FromUID   uint64             `bson:"from_uid,omitempty"`
	ToUID     uint64             `bson:"to_uid,omitempty"`
	IsFriend  bool               `bson:"is_friend,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
	FromUser  *User              `bson:"-"`
	ToUser    *User              `bson:"-"`
}

func (a *Follow) BeforeCreate(ctx context.Context) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return nil
}

func ListFollow(ctx context.Context, uid uint64, relationType enum.RelationType, pageParams *pagination.QuickPagination) ([]*Follow, pagination.Pagination, error) {
	follows := make([]*Follow, 0)
	chain := db.ODM(ctx)
	if relationType == enum.Fan {
		chain = chain.Where(bson.M{"to_uid": uid})
	} else if relationType == enum.Following {
		chain = chain.Where(bson.M{"from_uid": uid})
	} else {
		chain = chain.Where(bson.M{"from_uid": uid, "is_friend": true})
	}
	paginator := pagination.NewQuickPaginator(pageParams.Limit, pageParams.NextID, chain)
	page, err := paginator.Paginate(&follows)
	if err != nil {
		return nil, nil, err
	}
	return follows, page, preloadFollowUser(ctx, follows)
}

func CreateFollow(ctx context.Context, fromUID, toUID uint64, isFriend bool) (*Follow, error) {
	follow := &Follow{
		FromUID:  fromUID,
		ToUID:    toUID,
		IsFriend: isFriend,
	}
	if err := follow.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	result, err := db.DB().Collection("follows").InsertOne(ctx, follow)
	if err != nil {
		return nil, err
	}
	follow.ID = result.InsertedID.(primitive.ObjectID)
	return follow, nil
}

func (f *Follow) SetFriend(ctx context.Context, isFriend bool) error {
	f.IsFriend = isFriend
	_, err := db.DB().Collection("follows").UpdateByID(ctx, f.ID, bson.M{"$set": bson.M{"is_friend": isFriend}})
	return err
}

func GetFollow(ctx context.Context, fromUID, toUID uint64) (*Follow, error) {
	follow := &Follow{}
	result := db.DB().Collection("follows").FindOne(ctx, &bson.M{
		"from_uid": fromUID,
		"to_uid":   toUID,
	})
	err := result.Err()
	if err != nil {
		return nil, err
	}
	return follow, result.Decode(follow)
}

func DeleteFollow(ctx context.Context, fromUID, toUID uint64) error {
	_, err := db.DB().Collection("follows").DeleteOne(ctx, bson.M{"from_uid": fromUID, "to_uid": toUID})
	return err
}

func ListFollowingUserIDs(ctx context.Context, uid uint64) ([]uint64, error) {
	cursor, err := db.DB().Collection("follows").Find(ctx, &bson.M{
		"from_uid": uid,
	}, &options.FindOptions{
		Projection: bson.M{"to_uid": 1},
	})
	if err != nil {
		return nil, err
	}
	follows := make([]*Follow, 0)
	if err = cursor.All(ctx, &follows); err != nil {
		return nil, err
	}
	ids := make([]uint64, len(follows))
	for i, follow := range follows {
		ids[i] = follow.ToUID
	}
	return ids, nil
}

func preloadFollowUser(ctx context.Context, follows []*Follow) error {
	userIds := make([]uint64, 0)
	for _, follow := range follows {
		userIds = append(userIds, follow.FromUID, follow.ToUID)
	}
	users := make([]*User, 0)
	err := db.ODM(ctx).Where(bson.M{"_id": bson.M{"$in": userIds}}).Find(&users).Error
	if err != nil {
		return err
	}
	userMap := make(map[uint64]*User)
	for _, user := range users {
		userMap[user.UID] = user
	}
	for _, follow := range follows {
		follow.FromUser = userMap[follow.FromUID]
		follow.ToUser = userMap[follow.ToUID]
	}
	return nil
}
