package pagination

import (
	"encoding/base64"

	"github.com/mises-id/sns/lib/db/odm"
	"go.mongodb.org/mongo-driver/bson"
)

type PageQuickParams struct {
	Limit  int64  `json:"limit" query:"limit"`
	NextID string `json:"next_id" query:"next_id"`
}

func (*PageQuickParams) isPagePrams() {}

func (p *PageQuickParams) GetLimit() int64 {
	if p.Limit <= 0 || p.Limit > 200 {
		return 50
	}
	return p.Limit
}

type QuickPagination struct {
	Limit  int64  `json:"limit" query:"limit"`
	NextID string `json:"next_id" query:"next_id"`
}

type QuickPaginator struct {
	Limit  int64   `json:"-"`
	NextID string  `json:"-"`
	DB     *odm.DB `json:"-"`
}

func NewQuickPaginator(limit int64, nextID string, db *odm.DB) Paginator {
	if limit == 0 {
		limit = 50
	}

	return &QuickPaginator{
		Limit:  limit,
		NextID: nextID,
		DB:     db,
	}
}

type nextItem struct {
	ID string `bson:"_id,omitempty"`
}

func (p *QuickPaginator) Paginate(dataSource interface{}) (Pagination, error) {
	db := p.DB
	var err error
	if p.NextID != "" {
		db = db.Where(bson.M{"_id": bson.M{"$lte": p.NextID}})
	}
	err = db.Sort(bson.M{"_id": -1}).Limit(p.Limit).Find(dataSource).Error
	if err != nil {
		return nil, err
	}

	items := make([]*nextItem, 0)
	if err = db.Skip(p.Limit).Limit(1).Find(&items).Error; err != nil {
		return nil, err
	}
	pageToken := ""
	if len(items) > 0 {
		pageToken = base64.StdEncoding.EncodeToString([]byte(items[0].ID))
	}
	return &QuickPagination{
		Limit:  p.Limit,
		NextID: pageToken,
	}, nil
}

func (p *QuickPagination) BuildJSONResult() interface{} {
	return p
}

func (p *QuickPagination) GetPageSize() int {
	return int(p.Limit)
}
