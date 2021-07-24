package models

import (
	"context"

	"github.com/mises-id/sns/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Counter struct {
	ID  string `bson:"_id,omitempty"`
	Seq uint64 `bson:"seq,omitempty"`
}

func getNextSeq(ctx context.Context, id string) (uint64, error) {
	result := db.DB().Collection("counters").FindOneAndUpdate(ctx, bson.M{"_id": id},
		bson.M{"$inc": bson.M{"seq": 1}})
	err := result.Err()
	if err != nil {
		if mongo.ErrNoDocuments == err {
			return firstSeq(ctx, id)
		}
		return 0, result.Err()
	}
	counter := &Counter{}
	return counter.Seq, result.Decode(counter)
}

func firstSeq(ctx context.Context, id string) (uint64, error) {
	counter := &Counter{
		ID:  id,
		Seq: 1,
	}
	_, err := db.DB().Collection("counters").InsertOne(ctx, counter)
	return counter.Seq, err
}
