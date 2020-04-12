package schema

import (
	"fmt"
	"kztop/pkg/kreedz"
	"kztop/pkg/util"
	"time"
)

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
	Hash        string    `json:"hash"`
}

func (a *Nub) Validation() bool {
	queryStr := fmt.Sprintf("%s%s%s%s%.2f%s%d%s%d%d%s",
		a.MapName,
		a.AuthID,
		a.Country,
		a.Name,
		a.Time,
		a.Weapon,
		a.FinishCount,
		a.Server,
		a.CheckPoints,
		a.GoChecks,
		a.Route,
	)
	hash := util.MD5HashString(util.MD5HashString(queryStr) + kreedz.SECRET_KEY)
	return hash == a.Hash
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
