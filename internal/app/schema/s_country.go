package schema

// Country Country对象
type Country struct {
	Chinese  string `json:"chinese"`
	English  string `json:"english"`
	SortName string `json:"sort_name"`
}

// CountryQueryParam 查询条件
type CountryQueryParam struct {
	SortName string `form:"sort_name"`
}

// CountryQueryOptions Country对象查询可选参数项
type CountryQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// CountryQueryResult Country对象查询结果
type CountryQueryResult struct {
	Data       []*Country
	PageResult *PaginationResult
}
