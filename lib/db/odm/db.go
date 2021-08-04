package odm

import (
	"context"
	"log"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	db *mongo.Database
}

type collectioner interface {
	CollectionName() string
}

func NewClient(conn *mongo.Database) *Client {
	return &Client{db: conn}
}

func (c *Client) NewSession(ctx context.Context) *DB {
	return &DB{
		ctx:     ctx,
		db:      c.db,
		options: &options.FindOptions{},
	}
}

type DB struct {
	db             *mongo.Database
	ctx            context.Context
	Error          error
	collectionName string
	out            interface{}
	options        *options.FindOptions
	condition      bson.M
}

func (db *DB) Collection(collectionName string) *DB {
	db.collectionName = collectionName
	return db
}

func (db *DB) Model(model interface{}) *DB {
	db.collectionName = ""
	db.reflectCollectionName(model)
	return db
}

func (db *DB) Where(condition bson.M) *DB {
	if db.condition == nil {
		db.condition = condition
		return db
	}
	for key, value := range condition {
		db.condition[key] = value

	}
	return db
}

func (db *DB) Count(c *int64) *DB {
	countDocs, err := db.db.Collection(db.reflectCollectionName()).CountDocuments(db.ctx, db.condition)
	db.Error = err
	*c = countDocs
	return db
}

func (db *DB) Sort(sort interface{}) *DB {
	db.options.Sort = sort
	return db
}

func (db *DB) Limit(limit int64) *DB {
	db.options.Limit = &limit
	return db
}

func (db *DB) Skip(skip int64) *DB {
	db.options.Skip = &skip
	return db
}

func (db *DB) First(out interface{}, conditions ...bson.M) *DB {
	db.out = out
	for _, condition := range conditions {
		db = db.Where(condition)
	}
	result := db.db.Collection(db.reflectCollectionName()).FindOne(db.ctx, db.condition, &options.FindOneOptions{Sort: bson.M{"_id": 1}})
	db.Error = result.Err()
	if db.Error != nil {
		return db
	}
	db.Error = result.Decode(out)
	return db
}

func (db *DB) Last(out interface{}, conditions ...bson.M) *DB {
	db.out = out
	for _, condition := range conditions {
		db = db.Where(condition)
	}
	result := db.db.Collection(db.reflectCollectionName()).FindOne(db.ctx, db.condition, &options.FindOneOptions{Sort: bson.M{"_id": -1}})
	db.Error = result.Err()
	if db.Error != nil {
		return db
	}
	db.Error = result.Decode(out)
	return db
}

func (db *DB) Find(out interface{}, conditions ...bson.M) *DB {
	db.out = out
	for _, condition := range conditions {
		db = db.Where(condition)
	}
	var result *mongo.Cursor
	if db.options == nil {
		result, db.Error = db.db.Collection(db.reflectCollectionName()).Find(db.ctx, db.condition)
	} else {
		result, db.Error = db.db.Collection(db.reflectCollectionName()).Find(db.ctx, db.condition, db.options)
	}
	if db.Error != nil {
		return db
	}
	db.Error = result.All(db.ctx, out)
	return db
}

func (db *DB) reflectCollectionName(outs ...interface{}) string {
	if db.collectionName != "" {
		return db.collectionName
	}
	out := db.out
	if len(outs) > 0 {
		out = outs[0]
	}
	t := reflect.TypeOf(out)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	} else {
		log.Fatal("input must be ptr")
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	value := reflect.New(t)
	item, ok := value.Interface().(collectioner)
	if ok {
		db.collectionName = item.CollectionName()
	} else {
		// TODO convert city -> cities, bus -> buses, box-boxes, brush-brushes
		db.collectionName = strings.ToLower(t.Name()) + "s"
	}
	return db.collectionName
}
