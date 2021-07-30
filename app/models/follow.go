package models

import (
	"context"
	"time"

	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/lib/db"
	"github.com/mises-id/sns/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Follow struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UID       uint64             `bson:"uid,omitempty"`
	FocusUID  uint64             `bson:"focus_uid,omitempty"`
	IsFriend  bool               `bson:"is_friend,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
	FromUser  *User
	ToUser    *User
}

func (a *Follow) BeforeCreate(ctx context.Context) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return nil
}

func ListFollow(ctx context.Context, uid uint64, relationType enum.RelationType, pageParams *pagination.TraditionalParams) ([]*Follow, pagination.Pagination, error) {
	follows := make([]*Follow, 0)
	chain := db.ODM(ctx)
	if relationType == enum.Fan {
		chain = chain.Where(bson.M{"focus_uid": uid})
	} else if relationType == enum.Following {
		chain = chain.Where(bson.M{"uid": uid})
	} else {
		chain = chain.Where(bson.M{"uid": uid, "is_friend": true})
	}
	paginator := pagination.NewTraditionalPaginator(pageParams.Page, pageParams.PerPage, chain)
	page, err := paginator.Paginate(&follows)
	if err != nil {
		return nil, nil, err
	}
	return follows, page, preloadFollowUser(ctx, follows)
}

func CreateFollow(ctx context.Context, uid, focusUID uint64, isFriend bool) (*Follow, error) {
	follow := &Follow{
		UID:      uid,
		FocusUID: focusUID,
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

func GetFollow(ctx context.Context, uid, focusUID uint64) (*Follow, error) {
	follow := &Follow{}
	result := db.DB().Collection("follows").FindOne(ctx, &bson.M{
		"uid":       uid,
		"focus_uid": focusUID,
	})
	err := result.Err()
	if err != nil {
		return nil, err
	}
	return follow, result.Decode(follow)
}

func DeleteFollow(ctx context.Context, uid, focusUID uint64) error {
	_, err := db.DB().Collection("follows").DeleteOne(ctx, bson.M{"uid": uid, "focus_uid": focusUID})
	return err
}

func preloadFollowUser(ctx context.Context, follows []*Follow) error {
	userIds := make([]uint64, 0)
	for _, follow := range follows {
		userIds = append(userIds, follow.UID, follow.FocusUID)
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
		follow.FromUser = userMap[follow.UID]
		follow.ToUser = userMap[follow.FocusUID]
	}
	return nil
}
