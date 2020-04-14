package ctl

import (
	"kztop/internal/app/schema"
	"kztop/pkg/kreedz"
)

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
	a.playerInfo[steamID][schema.REGION] = region
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
				a.firstEachMap[newRecord.Cate][newRecord.MapName] = newRecord
			}
		}
	} else {
		a.playerStat[newRecord.SteamID][schema.FIRST]++
		a.firstEachMap[newRecord.Cate][newRecord.MapName] = newRecord
	}
}
