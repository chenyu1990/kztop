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
	playersNub := make(map[string][]*schema.Nub)
	playersWpn := make(map[string][]*schema.Wpn)
	firstEachMapPro := make(map[string]*schema.Pro)
	firstEachMapNub := make(map[string]*schema.Nub)
	firstEachMapWpn := make(map[string]*schema.Wpn)
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
		if firstEachMapPro[record.MapName] == nil {
			firstEachMapPro[record.MapName] = record
		} else if firstEachMapPro[record.MapName].Time > record.Time {
			firstEachMapPro[record.MapName] = record
		}

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
		if firstEachMapNub[record.MapName] == nil {
			firstEachMapNub[record.MapName] = record
		} else if firstEachMapNub[record.MapName].Time > record.Time {
			firstEachMapNub[record.MapName] = record
		}

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
		if firstEachMapWpn[record.MapName] == nil {
			firstEachMapWpn[record.MapName] = record
		} else if speed := firstEachMapWpn[record.MapName].Speed; (speed == record.Speed && firstEachMapWpn[record.MapName].Time > record.Time) || speed > record.Speed {
			firstEachMapWpn[record.MapName] = record
		}

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

	countriesData, err := mCountry.Query(ctx, schema.CountryQueryParam{})
	if err != nil {
		panic(err)
	}
	countriesInfo := make(map[string]schema.Country)
	for _, countryInfo := range countriesData.Data {
		countriesInfo[strings.ToLower(countryInfo.SortName)] = *countryInfo
	}

	// 统计第一的数量
	for _, record := range firstEachMapPro {
		if kreedz.IsSteamID(record.AuthID) {
			playerStat[record.AuthID]["first"]++
		}
	}
	for _, record := range firstEachMapNub {
		if kreedz.IsSteamID(record.AuthID) {
			playerStat[record.AuthID]["first"]++
		}
	}
	for _, record := range firstEachMapWpn {
		if kreedz.IsSteamID(record.AuthID) {
			playerStat[record.AuthID]["first"]++
		}
	}

	a := &Top{
		ProModel:         mPro,
		NubModel:         mNub,
		WpnModel:         mWpn,
		CountryModel:     mCountry,
		WorldRecordModel: mWorldRecord,
		playersPro:       playersPro,
		playersNub:       playersNub,
		playersWpn:       playersWpn,
		playerStat:       playerStat,
		playerInfo:       playerInfo,
		firstEachMapPro:  firstEachMapPro,
		firstEachMapNub:  firstEachMapNub,
		firstEachMapWpn:  firstEachMapWpn,
		countries:        countries,
		countriesInfo:    countriesInfo,
	}
	a.UpdateStats("")
	return a
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
	NubModel              model.INub
	ProModel              model.IPro
	WpnModel              model.IWpn
	CountryModel          model.ICountry
	WorldRecordModel      model.IWorldRecord
	playersPro            map[string][]*schema.Pro
	playersNub            map[string][]*schema.Nub
	playersWpn            map[string][]*schema.Wpn
	firstEachMapPro       map[string]*schema.Pro
	firstEachMapNub       map[string]*schema.Nub
	firstEachMapWpn       map[string]*schema.Wpn
	playerStat            map[string]map[string]int64
	playerInfo            map[string]map[string]string
	playerStatSortByPro   []string
	playerStatSortByNub   []string
	playerStatSortByWpn   []string
	playerStatSortByFirst []string
	countries             map[string]int
	countriesInfo         map[string]schema.Country
	countriesSort         []string
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
		case "first":
			ginplus.ResSuccess(c, a.playerStatSortByFirst[bgn:end])
		}
		return
	}

	h := gin.H{
		"sort":                  sortField,
		"playerInfo":            a.playerInfo,
		"playerStat":            a.playerStat,
		"playerStatSortByPro":   a.playerStatSortByPro[:size],
		"playerStatSortByNub":   a.playerStatSortByNub[:size],
		"playerStatSortByWpn":   a.playerStatSortByWpn[:size],
		"playerStatSortByFirst": a.playerStatSortByFirst[:size],
		"countries":             a.countries,
		"countriesSort":         a.countriesSort,
		"countriesInfo":         a.countriesInfo,
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

	speed := 0
	var mapname, steamID string
	var time float64
	switch cate {
	case "pro":
		newRecord := schema.Pro{}
		if err := ginplus.ParseJSON(c, &newRecord); err != nil {
			ginplus.ResError(c, err)
			return
		}
		// 不支持盗版玩家进入排行，不需要验证nick了。
		if kreedz.IsSteamID(newRecord.AuthID) {
			ginplus.ResError(c, errors.New400Response("不支持盗版玩家进入排行"))
			return
		}

		if newRecord.Validation() == false {
			ginplus.ResError(c, errors.New400Response("dont hack me"))
			return
		}

		/*
			1. 更新数据库 done
			2. 更新国家/地区统计（当玩家数据不存在时） done
			3. 更新玩家信息 done
			4. 更新第一信息 done
			5. 更新排序 done
		*/
		ctx := ginplus.NewContext(c)
		query, err := a.ProModel.Query(ctx, schema.ProQueryParam{
			MapName: newRecord.MapName,
			AuthID:  newRecord.AuthID,
		})
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		if len(query.Data) == 0 {
			err = a.ProModel.Create(ctx, newRecord)
			if err != nil {
				ginplus.ResError(c, err)
				return
			}
			// 更新排行榜缓存
			a.playersPro[newRecord.AuthID] = append(a.playersPro[newRecord.AuthID], &newRecord)
			a.playerStat[newRecord.AuthID]["pro"]++
			a.playerStat[newRecord.AuthID]["total"]++
		} else {
			newRecord.FinishCount++
			err = a.ProModel.Update(ctx, &schema.ProQueryParam{
				MapName: newRecord.MapName,
				AuthID:  newRecord.AuthID,
			}, newRecord)
			if err != nil {
				ginplus.ResError(c, err)
				return
			}
			// 更新排行榜缓存
			for i, record := range a.playersPro[newRecord.AuthID] {
				if record.MapName == newRecord.MapName {
					a.playersPro[newRecord.AuthID][i] = &newRecord
					break
				}
			}
		}

		newCountry, newNick := a.UpdateInfo(newRecord.Country, newRecord.AuthID, newRecord.Name)
		if newCountry != "" || newNick != "" {
			go a.ProModel.UpdateInfo(ctx, schema.UpdateInfo{
				AuthID:  newRecord.AuthID,
				Country: newCountry,
				Name:    newNick,
			})
		}
		mapname = newRecord.MapName
		steamID = newRecord.AuthID
		time = newRecord.Time
	}

	// 更新统计一定要在最后
	a.SwapMapFirst(cate, mapname, steamID, time, speed)
	a.UpdateStats(cate)
}

func (a *Top) UpdateStats(cate string) {
	if cate == "" || cate == "pro" {
		proStat := make(map[string]int)
		for steamID, stat := range a.playerStat {
			proStat[steamID] = int(stat["pro"])
		}
		a.playerStatSortByPro = sortMapStringInt(proStat)
	}

	if cate == "" || cate == "nub" {
		nubStat := make(map[string]int)
		for steamID, stat := range a.playerStat {
			nubStat[steamID] = int(stat["nub"])
		}
		a.playerStatSortByNub = sortMapStringInt(nubStat)
	}

	if cate == "" || cate == "wpn" {
		wpnStat := make(map[string]int)
		for steamID, stat := range a.playerStat {
			wpnStat[steamID] = int(stat["wpn"])
		}
		a.playerStatSortByWpn = sortMapStringInt(wpnStat)
	}

	firstStat := make(map[string]int)
	for steamID, stat := range a.playerStat {
		firstStat[steamID] = int(stat["first"])
	}

	countriesStat := make(map[string]int)
	for country, count := range a.countries {
		countriesStat[country] = count
	}

	a.playerStatSortByFirst = sortMapStringInt(firstStat)
	a.countriesSort = sortMapStringInt(countriesStat)
}

func (a *Top) UpdateInfo(country, steamID, nick string) (updateCountry, updateNick string) {
	if country == "" {
		country = "n-"
		updateCountry = "n-"
	}

	if a.playerInfo[steamID] == nil {
		a.playerInfo[steamID] = make(map[string]string)

		updateNick = nick
		// 当玩家不存在时，新增国家/地区数量
		a.countries[country]++
	} else {
		if a.playerInfo[steamID]["country"] != country {
			updateCountry = country
			a.countries[a.playerInfo[steamID]["country"]]--
			a.countries[country]++
		}
		if a.playerInfo[steamID]["nick"] != nick {
			updateNick = nick
		}
	}
	a.playerInfo[steamID]["country"] = strings.ToLower(country)
	a.playerInfo[steamID]["steamID64"] = kreedz.SteamIDToSteamID64(steamID)
	a.playerInfo[steamID]["nick"] = nick
	return
}

func (a *Top) SwapMapFirst(cate, mapname, steamID string, time float64, speed int) {
	switch cate {
	case "pro":
		if _, ok := a.firstEachMapPro[mapname]; ok {
			if a.firstEachMapPro[mapname].AuthID != steamID && a.firstEachMapPro[mapname].Time > time {
				a.firstEachMapPro[mapname].AuthID = steamID
				a.firstEachMapPro[mapname].Time = time
				a.playerStat[a.firstEachMapPro[mapname].AuthID]["first"]--
				a.playerStat[steamID]["first"]++
			}
		} else {
			a.playerStat[steamID]["first"]++
		}
	case "nub":
		if _, ok := a.firstEachMapNub[mapname]; ok {
			if a.firstEachMapNub[mapname].AuthID != steamID && a.firstEachMapNub[mapname].Time > time {
				a.firstEachMapNub[mapname].AuthID = steamID
				a.firstEachMapNub[mapname].Time = time
				a.playerStat[a.firstEachMapNub[mapname].AuthID]["first"]--
				a.playerStat[steamID]["first"]++
			}
		} else {
			a.playerStat[steamID]["first"]++
		}
	case "wpn":
		if _, ok := a.firstEachMapWpn[mapname]; ok {
			if a.firstEachMapWpn[mapname].AuthID != steamID {
				if (speed == a.firstEachMapWpn[mapname].Speed && a.firstEachMapWpn[mapname].Time > time) || a.firstEachMapWpn[mapname].Speed > speed {
					a.firstEachMapWpn[mapname].AuthID = steamID
					a.firstEachMapWpn[mapname].Time = time
					a.playerStat[a.firstEachMapWpn[mapname].AuthID]["first"]--
					a.playerStat[steamID]["first"]++
				}
			}
		} else {
			a.playerStat[steamID]["first"]++
		}
	}
}
