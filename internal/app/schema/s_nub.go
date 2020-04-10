package schema

import "time"

// Nub Nub对象
type Nub struct {
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
	Route       string    `json:"route"`
	Date        time.Time `json:"date"`
}

// NubQueryParam 查询条件
type NubQueryParam struct {
	MapName string `form:"mapname"`
	AuthID  string `form:"authid"`
}

// NubQueryOptions Nub对象查询可选参数项
type NubQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// NubQueryResult Nub对象查询结果
type NubQueryResult struct {
	Data       []*Nub
	PageResult *PaginationResult
}
