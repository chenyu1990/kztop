package schema

// Admin Admin对象
type Admin struct {
	Server  string `json:"server"`
	SteamID string `json:"steamid"`
	Access  string `json:"access"`
	Valid   bool   `json:"valid"`
}

// AdminQueryParam 查询条件
type AdminQueryParam struct {
	Server  string `form:"server"`
	SteamID string `form:"steam_id"`
	Valid   *bool  `json:"valid"`
}

// AdminQueryOptions Admin对象查询可选参数项
type AdminQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// AdminQueryResult Admin对象查询结果
type AdminQueryResult struct {
	Data       []*Admin
	PageResult *PaginationResult
}

// Admin 用户对象列表
type Admins []*Admin