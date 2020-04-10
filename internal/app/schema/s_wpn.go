package schema

import "time"

// Wpn Wpn对象
type Wpn struct {
	MapName     string    `json:"map_name"`
	AuthID      string    `json:"auth_id"`
	Country     string    `json:"country"`
	Name        string    `json:"name"`
	Time        float64   `json:"time"`
	Weapon      string    `json:"weapon"`
	FinishCount int       `json:"finish_count"`
	Server      string    `json:"server"`
	CheckPoints int       `json:"check_points"`
	GoChecks    int       `json:"go_checks"`
	Speed       int       `json:"speed"`
	Route       string    `json:"route"`
	Date        time.Time `json:"date"`
}

// WpnQueryParam 查询条件
type WpnQueryParam struct {
	MapName string `form:"mapname"`
	AuthID  string `form:"authid"`
}

// WpnQueryOptions Wpn对象查询可选参数项
type WpnQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// WpnQueryResult Wpn对象查询结果
type WpnQueryResult struct {
	Data       []*Wpn
	PageResult *PaginationResult
}
