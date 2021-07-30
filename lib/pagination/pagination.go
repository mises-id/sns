package pagination

type Pagination interface {
	BuildJSONResult() interface{}
	GetPerPage() int
}

type PageParams interface {
	isPagePrams()
}

type Paginator interface {
	Paginate(dataSource interface{}) (Pagination, error)
}
