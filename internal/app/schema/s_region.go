package schema

// Region Region对象
type Region struct {
	Chinese  string `json:"chinese"`
	English  string `json:"english"`
	SortName string `json:"sort_name"`
}

// RegionQueryParam 查询条件
type RegionQueryParam struct {
	SortName string `form:"sort_name"`
}

// RegionQueryOptions Region对象查询可选参数项
type RegionQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// RegionQueryResult Region对象查询结果
type RegionQueryResult struct {
	Data       []*Region
	PageResult *PaginationResult
}
