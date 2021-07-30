package pagination

import (
	"github.com/mises-id/sns/lib/db/odm"
	"go.mongodb.org/mongo-driver/bson"
)

type TraditionalParams struct {
	Page    int64 `json:"page" query:"page"`
	PerPage int64 `json:"per_page" query:"per_page"`
}

func (*TraditionalParams) isPagePrams() {}

type TraditionalPagination struct {
	TotalRecords int64 `json:"total_records" query:"total_records"`
	TotalPages   int64 `json:"total_pages" query:"total_pages"`
	PerPage      int64 `json:"per_page" query:"per_page"`
	CurrentPage  int64 `json:"current_page" query:"current_page"`
}

type TraditionalPaginator struct {
	Page    int64   `json:"-"`
	PerPage int64   `json:"-"`
	Offset  int64   `json:"-"`
	DB      *odm.DB `json:"-"`
}

func DefaultTraditionalParams() *TraditionalParams {
	return &TraditionalParams{
		Page:    1,
		PerPage: 50,
	}
}

func NewTraditionalParams(page, perPage int64) *TraditionalParams {
	if page < 1 {
		page = 1
	}
	if perPage < 2 {
		perPage = 50
	}
	return &TraditionalParams{
		Page:    page,
		PerPage: perPage,
	}
}

func NewTraditionalPaginator(page, perPage int64, db *odm.DB) Paginator {
	if page == 0 {
		page = 1
	}
	if perPage == 0 {
		perPage = 50
	}
	offset := (page - 1) * perPage

	return &TraditionalPaginator{
		Page:    page,
		PerPage: perPage,
		Offset:  offset,
		DB:      db,
	}
}

func (p *TraditionalPaginator) Paginate(dataSource interface{}) (Pagination, error) {
	db := p.DB

	var count int64
	err := db.Model(dataSource).Count(&count).Error
	if err != nil {
		return nil, err
	}
	err = db.Sort(bson.M{"_id": -1}).Limit(p.PerPage).Skip(p.Offset).Find(dataSource).Error
	if err != nil {
		return nil, err
	}
	totalPages := count / p.PerPage
	if count%int64(p.PerPage) > 0 {
		totalPages++
	}

	return &TraditionalPagination{
		TotalRecords: count,
		TotalPages:   totalPages,
		PerPage:      p.PerPage,
		CurrentPage:  p.Page,
	}, nil
}

func (p *TraditionalPagination) BuildJSONResult() interface{} {
	return p
}

func (p *TraditionalPagination) GetPerPage() int {
	return int(p.PerPage)
}

func (p *TraditionalPagination) SetPageToken(lastRecordID uint64) {
}
