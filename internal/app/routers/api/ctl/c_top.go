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
	mRegion model.IRegion,
	mRecord model.IRecord,
	mWorldRecord model.IWorldRecord,
) *Top {
	// map[steamID]recordInfo
	players := make(map[schema.Cate]map[string][]*schema.Record)
	players[schema.PRO] = make(map[string][]*schema.Record)
	players[schema.NUB] = make(map[string][]*schema.Record)
	players[schema.WPN] = make(map[string][]*schema.Record)
	firstEachMap := make(map[schema.Cate]map[string]*schema.Record)
	firstEachMap[schema.PRO] = make(map[string]*schema.Record)
	firstEachMap[schema.NUB] = make(map[string]*schema.Record)
	firstEachMap[schema.WPN] = make(map[string]*schema.Record)
	// map[steamID][firstCount/nub/pro/wpn/total/visitor/gotGreat]int64
	playerStat := make(map[string]map[schema.Cate]int64)
	playerInfo := make(map[string]map[string]string)
	regions := make(map[string]int)

	ctx := context.Background()

	records, err := mRecord.Query(ctx, nil)
	if err != nil {
		panic(err)
	}
	for _, record := range records.Data {
		if kreedz.IsSteamID(record.SteamID) == false {
			continue
		}

		if firstEachMap[record.Cate][record.MapName] == nil {
			firstEachMap[record.Cate][record.MapName] = record
		} else {
			if record.Cate != schema.WPN {
				if firstEachMap[record.Cate][record.MapName].Time > record.Time {
					firstEachMap[record.Cate][record.MapName] = record
				}
			} else {
				if speed := firstEachMap[record.Cate][record.MapName].Speed; (speed == record.Speed && firstEachMap[record.Cate][record.MapName].Time > record.Time) || speed > record.Speed {
					firstEachMap[record.Cate][record.MapName] = record
				}
			}
		}

		if playerStat[record.SteamID] == nil {
			playerStat[record.SteamID] = make(map[schema.Cate]int64)
			playerInfo[record.SteamID] = make(map[string]string)
			// Regions 在地区统计完成后统计，部分纪录存在玩家存在两个地区
			// regions[record.Region]++
		}

		players[record.Cate][record.SteamID] = append(players[record.Cate][record.SteamID], record)
		playerStat[record.SteamID][record.Cate]++
		playerStat[record.SteamID][schema.TOTAL]++
		playerStat[record.SteamID][schema.FinishCount] += int64(record.FinishCount)
		playerInfo[record.SteamID][schema.STEAMID64] = kreedz.SteamIDToSteamID64(record.SteamID)
		playerInfo[record.SteamID][schema.NICK] = record.Nick
		// 矫正历史纪录的错误
		if playerInfo[record.SteamID][schema.REGION] == "" {
			if record.Region == "" {
				playerInfo[record.SteamID][schema.REGION] = "n-"
			} else {
				playerInfo[record.SteamID][schema.REGION] = strings.ToLower(record.Region)
			}
		} else {
			if record.Region != "" && record.Region != "n-" {
				playerInfo[record.SteamID][schema.REGION] = strings.ToLower(record.Region)
			}
		}
	}

	regionsData, err := mRegion.Query(ctx, schema.RegionQueryParam{})
	if err != nil {
		panic(err)
	}
	regionsInfo := make(map[string]schema.Region)
	for _, regionInfo := range regionsData.Data {
		regionsInfo[strings.ToLower(regionInfo.SortName)] = *regionInfo
	}

	// 统计第一的数量
	for _, records := range firstEachMap {
		for _, record := range records {
			playerStat[record.SteamID][schema.FIRST]++
		}
	}

	// 统计地区
	for _, player := range playerInfo {
		regions[player["region"]]++
	}

	a := &Top{
		RegionModel:      mRegion,
		RecordModel:      mRecord,
		WorldRecordModel: mWorldRecord,
		players:          players,
		playerStat:       playerStat,
		playerInfo:       playerInfo,
		firstEachMap:     firstEachMap,
		regions:          regions,
		regionsInfo:      regionsInfo,
		playerStatSort:   make(map[schema.Cate][]string),
	}
	a.UpdateStats(schema.NULL)
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
	RegionModel      model.IRegion
	RecordModel      model.IRecord
	WorldRecordModel model.IWorldRecord
	players          map[schema.Cate]map[string][]*schema.Record
	firstEachMap     map[schema.Cate]map[string]*schema.Record
	playerStat       map[string]map[schema.Cate]int64
	playerInfo       map[string]map[string]string
	playerStatSort   map[schema.Cate][]string
	regions          map[string]int
	regionsInfo      map[string]schema.Region
	regionsSort      []string
}

func (a *Top) Player(c *gin.Context) {
	pCate := c.Param("cate")
	cate := schema.GetCate(pCate)
	if cate < schema.PRO || cate > schema.WPN {
		cate = schema.NUB
	}
	steamID64 := c.Param("player")
	player := kreedz.SteamID64ToSteamID(steamID64)
	if kreedz.IsSteamID(player) == false {
		ginplus.ResError(c, errors.New400Response("not like a valid steam_id?"))
		return
	}

	if _, exist := a.players[cate][player]; !exist {
		ginplus.ResError(c, errors.New400Response("user not exist"))
		return
	}

	// 以完成次数排序
	finishCountStat := make(map[string]int)
	for _, record := range a.players[cate][player] {
		finishCountStat[record.MapName] = record.FinishCount
	}
	finishCountSort := sortMapStringInt(finishCountStat)

	list := make(map[string]*schema.Record)
	for _, record := range a.players[cate][player] {
		list[record.MapName] = record
	}
	h := gin.H{
		"cate":  cate,
		"pCate": pCate,
		"list":  list,
	}

	pageStr := c.Query("page")
	if pageStr != "" {
		page, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			ginplus.ResError(c, errors.New400Response("页面错误"))
			return
		}

		bgn := (page - 1) * schema.PageSize
		end := bgn + schema.PageSize
		h["sort"] = finishCountSort[bgn:end]
		c.HTML(http.StatusOK, "top/player_more", h)
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

	h["player"] = player
	h["stat"] = a.playerStat[player]
	h["info"] = a.playerInfo[player]
	h["region"] = a.regionsInfo
	h["sort"] = finishCountSort[:schema.PageSize]
	h["steamVisibilityState"] = Steam.VisibilityState
	c.HTML(http.StatusOK, "top/player", h)
}

func (a *Top) Players(c *gin.Context) {
	sortField := c.Query("sort")
	cate := schema.GetCate(sortField)
	if cate == schema.NULL {
		cate = schema.PRO
		sortField = "pro"
	}

	h := gin.H{
		"cate":        cate,
		"sortField":   sortField,
		"playerInfo":  a.playerInfo,
		"playerStat":  a.playerStat,
		"regions":     a.regions,
		"regionsSort": a.regionsSort,
		"regionsInfo": a.regionsInfo,
	}

	pageStr := c.Query("page")
	if pageStr != "" {
		page, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			ginplus.ResError(c, errors.New400Response("页面错误"))
			return
		}

		bgn := (page - 1) * schema.PageSize
		end := bgn + schema.PageSize
		h["playerStatSort"] = a.playerStatSort[cate][bgn:end]
		c.HTML(http.StatusOK, "top/players_more", h)
		return
	}

	h["playerStatSort"] = a.playerStatSort[cate][:schema.PageSize]
	c.HTML(http.StatusOK, "top/players", h)
}

func (a *Top) Top(c *gin.Context) {
	mapname := c.Param("mapname")
	pCate := c.Param("cate")
	cate := schema.GetCate(pCate)

	ctx := ginplus.NewContext(c)

	record, err := a.WorldRecordModel.Get(ctx, schema.WorldRecordQueryParam{
		MapName:       mapname,
		Organizations: []kreedz.Organization{kreedz.XtremeJumps, kreedz.CosyClimbing, kreedz.WorldSurf},
	})
	if err != nil {
		panic(err)
	}

	h := gin.H{
		"mapname": mapname,
		"cate":    pCate,
	}
	if record != nil {
		h["record"] = record
		h["wr"] = kreedz.Name[record.Organization]
	} else {
		h["record"] = schema.WorldRecord{
			Holder: "n/a",
			Region: "n-a",
			Time:   0,
		}
		h["wr"] = ""
	}

		orders := schema.Orders{}
		if cate == schema.WPN {
			orders["speed"] = schema.OrderASC
		}
		orders["time"] = schema.OrderASC
		list, err := a.RecordModel.Query(ctx, &schema.RecordQueryParam{
			Cate:    cate,
			MapName: mapname,
		}, schema.RecordQueryOptions{
			PageParam: &schema.PaginationParam{
				PageSize: 100,
			},
			OrderParam: &schema.OrderParam{
				Orders: orders,
			},
		})
		if err != nil {
		ginplus.ResError(c, err)
		return
	}
	h["list"] = list.Data
	h["total"] = list.PageResult.Total

	c.HTML(http.StatusOK, "top/index", h)
}

func (a *Top) UpdateRecord(c *gin.Context) {
	newRecord := schema.Record{}
	if err := ginplus.ParseJSON(c, &newRecord); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if newRecord.Speed == 250 {
		if newRecord.CheckPoints == 0 && newRecord.GoChecks == 0 {
			newRecord.Cate = schema.PRO
		} else {
			newRecord.Cate = schema.NUB
		}
	} else {
		newRecord.Cate = schema.WPN
	}
	cate := newRecord.Cate
	// 不支持盗版玩家进入排行，不需要验证nick了。
	if kreedz.IsSteamID(newRecord.SteamID) {
		ginplus.ResError(c, errors.New400Response("不支持盗版玩家进入排行"))
		return
	}

	if newRecord.Validation() == false {
		ginplus.ResError(c, errors.New400Response("dont hack me"))
		return
	}

	/*
		1. 更新数据库 done
		2. 更新地区统计（当玩家数据不存在时） done
		3. 更新玩家信息 done
		4. 更新第一信息 done
		5. 更新排序 done
	*/
	ctx := ginplus.NewContext(c)
	query, err := a.RecordModel.Query(ctx, &schema.RecordQueryParam{
		Cate:    cate,
		MapName: newRecord.MapName,
		SteamID: newRecord.SteamID,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if len(query.Data) == 0 {
		err = a.RecordModel.Create(ctx, &newRecord)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}

		// 更新排行榜缓存
		a.players[cate][newRecord.SteamID] = append(a.players[cate][newRecord.SteamID], &newRecord)
		// 只有没有时，才更新
		a.playerStat[newRecord.SteamID][cate]++
		a.playerStat[newRecord.SteamID][schema.TOTAL]++
	} else {
		newRecord.FinishCount++
		err = a.RecordModel.Update(ctx, &schema.RecordQueryParam{
			Cate:    cate,
			MapName: newRecord.MapName,
			SteamID: newRecord.SteamID,
		}, newRecord)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}

		// 更新排行榜缓存
		for i, record := range a.players[cate][newRecord.SteamID] {
			if record.MapName == newRecord.MapName {
				a.players[cate][newRecord.SteamID][i] = &newRecord
				break
			}
		}
	}
	// 有完成就自增
	a.playerStat[newRecord.SteamID][schema.FinishCount]++

	// 更新用户信息
	newRegion, newNick := a.UpdateInfo(newRecord.Region, newRecord.SteamID, newRecord.Nick)
	if newRegion != "" || newNick != "" {
		go a.RecordModel.UpdateInfo(ctx, schema.UpdateInfo{
			SteamID: newRecord.SteamID,
			Region:  newRegion,
			Nick:    newNick,
		})
	}

	// 更新统计一定要在最后
	a.SwapMapFirst(&newRecord)
	a.UpdateStats(cate)
}

func (a *Top) UpdateStats(cate schema.Cate) {
	m := make(map[schema.Cate]map[string]int)
	m[schema.PRO] = make(map[string]int)
	m[schema.NUB] = make(map[string]int)
	m[schema.WPN] = make(map[string]int)
	for steamID, stat := range a.playerStat {
		if cate == schema.NULL {
			m[schema.PRO][steamID] = int(stat[schema.PRO])
			m[schema.NUB][steamID] = int(stat[schema.NUB])
			m[schema.WPN][steamID] = int(stat[schema.WPN])
		} else {
			m[cate][steamID] = int(stat[cate])
		}
	}
	if cate != schema.NULL {
		a.playerStatSort[cate] = sortMapStringInt(m[cate])
	} else {
		a.playerStatSort[schema.PRO] = sortMapStringInt(m[schema.PRO])
		a.playerStatSort[schema.NUB] = sortMapStringInt(m[schema.NUB])
		a.playerStatSort[schema.WPN] = sortMapStringInt(m[schema.WPN])
	}

	firstStat := make(map[string]int)
	for steamID, stat := range a.playerStat {
		firstStat[steamID] = int(stat[schema.FIRST])
	}

	regionsStat := make(map[string]int)
	for region, count := range a.regions {
		regionsStat[region] = count
	}

	a.playerStatSort[schema.FIRST] = sortMapStringInt(firstStat)
	a.regionsSort = sortMapStringInt(regionsStat)
}

func (a *Top) UpdateInfo(region, steamID, nick string) (updateRegion, updateNick string) {
	if region == "" {
		region = "n-"
		updateRegion = "n-"
	}

	if a.playerInfo[steamID] == nil {
		a.playerInfo[steamID] = make(map[string]string)

		updateNick = nick
		// 当玩家不存在时，新增地区数量
		a.regions[region]++
	} else {
		if a.playerInfo[steamID][schema.REGION] != region {
			updateRegion = region
			a.regions[a.playerInfo[steamID][schema.REGION]]--
			a.regions[region]++
		}
		if a.playerInfo[steamID][schema.NICK] != nick {
			updateNick = nick
		}
	}
	a.playerInfo[steamID][schema.REGION] = strings.ToLower(region)
	a.playerInfo[steamID][schema.STEAMID64] = kreedz.SteamIDToSteamID64(steamID)
	a.playerInfo[steamID][schema.NICK] = nick
	return
}

func (a *Top) SwapMapFirst(newRecord *schema.Record) {
	if _, ok := a.firstEachMap[newRecord.Cate][newRecord.MapName]; ok {
		if a.firstEachMap[newRecord.Cate][newRecord.MapName].SteamID != newRecord.SteamID {
			swap := false
			if newRecord.Cate == schema.WPN {
				swap = (newRecord.Speed == a.firstEachMap[newRecord.Cate][newRecord.MapName].Speed && a.firstEachMap[newRecord.Cate][newRecord.MapName].Time > newRecord.Time) || a.firstEachMap[newRecord.Cate][newRecord.MapName].Speed > newRecord.Speed
			} else {
				swap = a.firstEachMap[newRecord.Cate][newRecord.MapName].Time > newRecord.Time
			}
			if swap {
				a.playerStat[a.firstEachMap[newRecord.Cate][newRecord.MapName].SteamID][schema.FIRST]--
				a.playerStat[newRecord.SteamID][schema.FIRST]++
			}
		}
	} else {
		a.playerStat[newRecord.SteamID][schema.FIRST]++
	}

	a.firstEachMap[newRecord.Cate][newRecord.MapName].SteamID = newRecord.SteamID
	a.firstEachMap[newRecord.Cate][newRecord.MapName].Time = newRecord.Time
}
