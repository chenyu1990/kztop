package ctl

import (
	"context"
	"kztop/internal/app/model"
	"kztop/internal/app/schema"
	"kztop/pkg/kreedz"
	"sort"
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

		if firstEachMap[record.Cate][record.RouteMapName()] == nil {
			firstEachMap[record.Cate][record.RouteMapName()] = record
		} else {
			if record.Cate != schema.WPN {
				if firstEachMap[record.Cate][record.RouteMapName()].Time > record.Time {
					firstEachMap[record.Cate][record.RouteMapName()] = record
				}
			} else {
				if speed := firstEachMap[record.Cate][record.RouteMapName()].Speed; (speed == record.Speed && firstEachMap[record.Cate][record.RouteMapName()].Time > record.Time) || speed > record.Speed {
					firstEachMap[record.Cate][record.RouteMapName()] = record
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