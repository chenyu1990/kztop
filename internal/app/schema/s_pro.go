package schema

import "time"

// Pro Pro对象
type Pro struct {
	MapName     string    `json:"map_name"`
	AuthID      string    `json:"auth_id"`
	Country     string    `json:"country"`
	Name        string    `json:"name"`
	Time        float64   `json:"time"`
	Weapon      string    `json:"weapon"`
	FinishCount int       `json:"finish_count"`
	Server      string    `json:"server"`
	Route       string    `json:"route"`
	Date        time.Time `json:"date"`
}

// ProQueryParam 查询条件
type ProQueryParam struct {
	MapName string `form:"mapname"`
	AuthID  string `form:"authid"`
}

// ProQueryOptions Pro对象查询可选参数项
type ProQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// ProQueryResult Pro对象查询结果
type ProQueryResult struct {
	Data       []*Pro
	PageResult *PaginationResult
}
