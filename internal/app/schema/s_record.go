package schema

import (
	"fmt"
	"kztop/pkg/kreedz"
	"kztop/pkg/util"
	"time"
)

type Cate int

const (
	_ Cate = iota // 不要影响 PRO / NUB / WPN 的值，影响某些逻辑
	PRO
	NUB
	WPN
	NULL        = 0
	FIRST       = 100
	TOTAL       = 101
	FinishCount = 102
)

type PlayerInfo string

const (
	_         PlayerInfo = ""
	REGION               = "region"
	NICK                 = "nick"
	STEAMID64            = "steamID64"
)

const PageSize = uint64(30)

func GetCate(cate string) Cate {
	switch cate {
	case "pro":
		return PRO
	case "nub":
		return NUB
	case "wpn":
		return WPN
	case "first":
		return FIRST
	}
	return NULL
}

func (a *Cate) ToString() string {
	switch *a {
	case PRO:
		return "pro"
	case NUB:
		return "nub"
	case WPN:
		return "wpn"
	case FIRST:
		return "FIRST"
	}
	return ""
}

// Record Record对象
type Record struct {
	Cate        Cate      `json:"cate"`
	MapName     string    `json:"mapname"`
	SteamID     string    `json:"steam_id"`
	Region      string    `json:"region"`
	Nick        string    `json:"nick"`
	Time        float64   `json:"-"`
	TimeString  string    `json:"time"`
	Weapon      string    `json:"weapon"`
	FinishCount int       `json:"-"`
	Server      string    `json:"server"`
	CheckPoints int       `json:"check_points"`
	GoChecks    int       `json:"go_checks"`
	Speed       int       `json:"speed"`
	Route       string    `json:"route"`
	Date        time.Time `json:"date"`
	Hash        string    `json:"hash"`
}

type UpdateInfo struct {
	SteamID string `json:"steam_id"`
	Region  string `json:"region"`
	Nick    string `json:"nick"`
}

func (a *Record) Validation() bool {
	queryStr := fmt.Sprintf("%d%s%s%s%s%s%s%s%d%d%d%s",
		a.Cate,
		a.MapName,
		a.SteamID,
		a.Region,
		a.Nick,
		a.TimeString,
		a.Weapon,
		a.Server,
		a.CheckPoints,
		a.GoChecks,
		a.Speed,
		a.Route,
	)
	hash := util.MD5HashString(util.MD5HashString(queryStr) + kreedz.SECRET_KEY)
	return hash == a.Hash
}

func (a *Record) RouteMapName() string {
	if a.Route != "" {
		return fmt.Sprintf("%s[%s]", a.MapName, a.Route)
	}
	return a.MapName
}

// RecordQueryParam 查询条件
type RecordQueryParam struct {
	Cate    Cate   `form:"cate"`
	MapName string `form:"mapname"`
	SteamID string `form:"steam_id"`
	Route   string `form:"route"`
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
