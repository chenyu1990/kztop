package schema

import (
	"fmt"
	"kztop/pkg/kreedz"
	"kztop/pkg/util"
	"time"
)

type Cate int

const (
	_ Cate = iota
	PRO
	NUB
	WPN
)

// Record Record对象
type Record struct {
	Cate        Cate      `json:"cate"`
	MapName     string    `json:"mapname"`
	SteamID     string    `json:"steam_id"`
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
	Hash        string    `json:"hash"`
}

func (a *Record) Validation() bool {
	// TODO Try toString()
	queryStr := fmt.Sprintf("%d%s%s%s%s%.2f%s%d%s%d%d%d%s",
		a.Cate,
		a.MapName,
		a.SteamID,
		a.Country,
		a.Name,
		a.Time,
		a.Weapon,
		a.FinishCount,
		a.Server,
		a.CheckPoints,
		a.GoChecks,
		a.Speed,
		a.Route,
	)
	hash := util.MD5HashString(util.MD5HashString(queryStr) + kreedz.SECRET_KEY)
	return hash == a.Hash
}

// RecordQueryParam 查询条件
type RecordQueryParam struct {
	Cate    Cate   `form:"cate"`
	MapName string `form:"mapname"`
	AuthID  string `form:"authid"`
}

// RecordQueryOptions Record对象查询可选参数项
type RecordQueryOptions struct {
	PageParam  *PaginationParam // 分页参数
	OrderParam *OrderParam      // 排序参数
}

// RecordQueryResult Record对象查询结果
type RecordQueryResult struct {
	Data       []*Record
	PageResult *PaginationResult
}
