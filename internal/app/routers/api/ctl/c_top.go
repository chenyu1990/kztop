package ctl

import (
	"context"
	"github.com/gin-gonic/gin"
	"kztop/internal/app/errors"
	"kztop/internal/app/ginplus"
	"kztop/internal/app/model"
	"kztop/internal/app/schema"
	"kztop/pkg/Steam"
	"kztop/pkg/kreedz"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// NewTop 创建top控制器
func NewTop(
	mPro model.IPro,
	mNub model.INub,
	mWpn model.IWpn,
	mCountry model.ICountry,
	mWorldRecord model.IWorldRecord,
) *Top {
	// map[steamID]recordInfo
	playersPro := make(map[string][]*schema.Pro)
	playersWpn := make(map[string][]*schema.Wpn)
	playersNub := make(map[string][]*schema.Nub)
	// map[steamID][firstCount/nub/pro/wpn/total/visitor/gotGreat]int64
	playerStat := make(map[string]map[string]int64)
	playerInfo := make(map[string]map[string]string)
	countries := make(map[string]int)

	ctx := context.Background()

	pro, err := mPro.Query(ctx, schema.ProQueryParam{})
	if err != nil {
		panic(err)
	}
	for _, record := range pro.Data {
		if kreedz.IsSteamID(record.AuthID) == false {
			continue
		}
		if playerStat[record.AuthID] == nil {
			playerStat[record.AuthID] = make(map[string]int64)
			playerInfo[record.AuthID] = make(map[string]string)
			if record.Country == "" {
				record.Country = "n-"
			}
			countries[record.Country]++
		}

		playersPro[record.AuthID] = append(playersPro[record.AuthID], record)
		playerStat[record.AuthID]["pro"]++
		playerStat[record.AuthID]["total"]++
		playerInfo[record.AuthID]["country"] = strings.ToLower(record.Country)
		playerInfo[record.AuthID]["steamID64"] = kreedz.SteamIDToSteamID64(record.AuthID)
		playerInfo[record.AuthID]["nick"] = record.Name

		// TODO first / visitor / gotGreat
	}

	nub, err := mNub.Query(ctx, schema.NubQueryParam{})
	if err != nil {
		panic(err)
	}
	for _, record := range nub.Data {
		if kreedz.IsSteamID(record.AuthID) == false {
			continue
		}
		if playerStat[record.AuthID] == nil {
			playerStat[record.AuthID] = make(map[string]int64)
			playerInfo[record.AuthID] = make(map[string]string)
			if record.Country == "" {
				record.Country = "n-"
			}
			countries[record.Country]++
		}

		playersNub[record.AuthID] = append(playersNub[record.AuthID], record)
		playerStat[record.AuthID]["nub"]++
		playerStat[record.AuthID]["total"]++
		playerInfo[record.AuthID]["country"] = strings.ToLower(record.Country)
		playerInfo[record.AuthID]["steamID64"] = kreedz.SteamIDToSteamID64(record.AuthID)
		playerInfo[record.AuthID]["nick"] = record.Name
		// TODO first / visitor / gotGreat
	}

	wpn, err := mWpn.Query(ctx, schema.WpnQueryParam{})
	if err != nil {
		panic(err)
	}
	for _, record := range wpn.Data {
		if kreedz.IsSteamID(record.AuthID) == false {
			continue
		}
		if playerStat[record.AuthID] == nil {
			playerStat[record.AuthID] = make(map[string]int64)
			playerInfo[record.AuthID] = make(map[string]string)
			if record.Country == "" {
				record.Country = "n-"
			}
			countries[record.Country]++
		}

		playersWpn[record.AuthID] = append(playersWpn[record.AuthID], record)
		playerStat[record.AuthID]["wpn"]++
		playerStat[record.AuthID]["total"]++
		playerInfo[record.AuthID]["country"] = strings.ToLower(record.Country)
		playerInfo[record.AuthID]["steamID64"] = kreedz.SteamIDToSteamID64(record.AuthID)
		playerInfo[record.AuthID]["nick"] = record.Name
		// TODO first / visitor / gotGreat
	}

	proStat := make(map[string]int)
	for steamID, stat := range playerStat {
		proStat[steamID] = int(stat["pro"])
	}

	nubStat := make(map[string]int)
	for steamID, stat := range playerStat {
		nubStat[steamID] = int(stat["nub"])
	}

	wpnStat := make(map[string]int)
	for steamID, stat := range playerStat {
		wpnStat[steamID] = int(stat["wpn"])
	}

	countriesStat := make(map[string]int)
	for country, count := range countries {
		countriesStat[country] = count
	}

	countriesData, err := mCountry.Query(ctx, schema.CountryQueryParam{})
	if err != nil {
		panic(err)
	}
	countriesInfo := make(map[string]schema.Country)
	for _, countryInfo := range countriesData.Data {
		countriesInfo[strings.ToLower(countryInfo.SortName)] = *countryInfo
	}

	return &Top{
		ProModel:            mPro,
		NubModel:            mNub,
		WpnModel:            mWpn,
		CountryModel:        mCountry,
		WorldRecordModel:    mWorldRecord,
		playersPro:          playersPro,
		playersNub:          playersNub,
		playersWpn:          playersWpn,
		playerStat:          playerStat,
		playerInfo:          playerInfo,
		playerStatSortByPro: sortMapStringInt(proStat),
		playerStatSortByNub: sortMapStringInt(nubStat),
		playerStatSortByWpn: sortMapStringInt(wpnStat),
		countries:           countries,
		countriesInfo:       countriesInfo,
		countriesSort:       sortMapStringInt(countriesStat),
	}
}

func sortMapStringInt(values map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range values {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	sorted := make([]string, len(values))
	for i, kv := range ss {
		sorted[i] = kv.Key
	}
	return sorted
}

type Top struct {
	NubModel            model.INub
	ProModel            model.IPro
	WpnModel            model.IWpn
	CountryModel        model.ICountry
	WorldRecordModel    model.IWorldRecord
	playersPro          map[string][]*schema.Pro
	playersNub          map[string][]*schema.Nub
	playersWpn          map[string][]*schema.Wpn
	playerStat          map[string]map[string]int64
	playerInfo          map[string]map[string]string
	playerStatSortByPro []string
	playerStatSortByNub []string
	playerStatSortByWpn []string
	countries           map[string]int
	countriesInfo       map[string]schema.Country
	countriesSort       []string
}

func (a *Top) Player(c *gin.Context) {
	cate := c.Param("type")
	if cate == "" {
		cate = "nub"
	}
	steamID64 := c.Param("player")
	player := kreedz.SteamID64ToSteamID(steamID64)
	if kreedz.IsSteamID(player) == false {
		ginplus.ResError(c, errors.New400Response("not like a valid steam_id?"))
		return
	}

	h := gin.H{
		"player": player,
		"cate":   cate,
		"stat":   a.playerStat[player],
		"info":   a.playerInfo[player],
	}

	var exist bool
	switch cate {
	case "pro":
		_, exist = a.playersPro[player]
		if exist {
			h["list"] = a.playersPro[player]
		}
	case "nub":
		_, exist = a.playersNub[player]
		if exist {
			h["list"] = a.playersNub[player]
		}
	case "wpn":
		_, exist = a.playersWpn[player]
		if exist {
			h["list"] = a.playersWpn[player]
		}
	}
	if !exist {
		ginplus.ResError(c, errors.New400Response("user not exist"))
		return
	}

	if _, ok := a.playerInfo[player]["avatar"]; !ok {
		go func(steamID64 string, player string) {
			profile, err := Steam.GetProfile(steamID64)
			if err != nil {
				return
			}
			a.playerInfo[player]["avatarFull"] = profile.AvatarFull
			a.playerInfo[player]["onlineState"] = profile.OnlineState
			a.playerInfo[player]["visibilityState"] = Steam.VisibilityState[profile.VisibilityState]
		}(steamID64, player)
	}

	c.HTML(http.StatusOK, "top/player", h)
}

func (a *Top) Players(c *gin.Context) {
	sortField := c.Query("sort")
	if sortField == "" {
		sortField = "pro"
	}

	size := uint64(30)
	pageStr := c.Query("page")
	if pageStr != "" {
		page, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			ginplus.ResError(c, errors.New400Response("页面错误"))
			return
		}

		bgn := (page - 1) * size
		end := bgn + size
		switch sortField {
		case "pro":
			ginplus.ResSuccess(c, a.playerStatSortByPro[bgn:end])
		case "nub":
			ginplus.ResSuccess(c, a.playerStatSortByNub[bgn:end])
		case "wpn":
			ginplus.ResSuccess(c, a.playerStatSortByWpn[bgn:end])
		}
		return
	}

	h := gin.H{
		"sort":                sortField,
		"playerInfo":          a.playerInfo,
		"playerStat":          a.playerStat,
		"playerStatSortByPro": a.playerStatSortByPro[:size],
		"playerStatSortByNub": a.playerStatSortByNub[:size],
		"playerStatSortByWpn": a.playerStatSortByWpn[:size],
		"countries":           a.countries,
		"countriesSort":       a.countriesSort,
		"countriesInfo":       a.countriesInfo,
	}
	c.HTML(http.StatusOK, "top/players", h)
}

func (a *Top) Top(c *gin.Context) {
	mapname := c.Param("mapname")
	cate := c.Param("cate")

	ctx := ginplus.NewContext(c)

	record, err := a.WorldRecordModel.Get(ctx, schema.WorldRecordQueryParam{
		MapName:       mapname,
		Organizations: []kreedz.Organization{kreedz.XtremeJumps, kreedz.CosyClimbing},
	})
	if err != nil {
		panic(err)
	}

	h := gin.H{
		"mapname": mapname,
		"cate":    cate,
	}
	if record != nil {
		h["record"] = record
		h["wr"] = kreedz.Name[record.Organization]
	} else {
		h["record"] = schema.WorldRecord{
			Holder:  "n/a",
			Country: "n-a",
			Time:    0,
		}
		h["wr"] = ""
	}

	switch cate {
	case "nub":
		list, _ := a.NubModel.Query(ctx, schema.NubQueryParam{
			MapName: mapname,
		}, schema.NubQueryOptions{PageParam: &schema.PaginationParam{
			PageSize: 100,
		}})
		h["list"] = list.Data
		h["total"] = list.PageResult.Total
	case "pro":
		list, _ := a.ProModel.Query(ctx, schema.ProQueryParam{
			MapName: mapname,
		}, schema.ProQueryOptions{PageParam: &schema.PaginationParam{
			PageSize: 100,
		}})
		h["list"] = list.Data
		h["total"] = list.PageResult.Total
	case "wpn":
		list, _ := a.WpnModel.Query(ctx, schema.WpnQueryParam{
			MapName: mapname,
		}, schema.WpnQueryOptions{PageParam: &schema.PaginationParam{
			PageSize: 100,
		}})
		h["list"] = list.Data
		h["total"] = list.PageResult.Total
	}

	c.HTML(http.StatusOK, "top/index", h)
}

func (a *Top) UpdateRecord(c *gin.Context) {
	cate := c.Param("cate")

	switch cate {
	case "pro":
		record := schema.Pro{}
		if err := ginplus.ParseJSON(c, &record); err != nil {
			ginplus.ResError(c, err)
			return
		}
		if record.Validation() == false {
			ginplus.ResError(c, errors.New400Response("error hash"))
			return
		}

		ctx := ginplus.NewContext(c)
		// 不支持盗版玩家进入排行，不需要验证nick了。
		if kreedz.IsSteamID(record.AuthID) {
			query, err := a.ProModel.Query(ctx, schema.ProQueryParam{
				MapName: record.MapName,
				AuthID:  record.AuthID,
			})
			if err != nil {
				ginplus.ResError(c, errors.New400Response("error hash"))
				return
			}
			if len(query.Data) == 0 {
				err = a.ProModel.Create(ctx, record)
				if err != nil {
					ginplus.ResError(c, errors.New400Response("error hash"))
					return
				}

			} else {
				record.FinishCount++
				err = a.ProModel.Update(ctx, &schema.ProQueryParam{
					MapName: record.MapName,
					AuthID:  record.AuthID,
				}, record)
				if err != nil {
					ginplus.ResError(c, errors.New400Response("error hash"))
					return
				}

			}
		}

	}
}
