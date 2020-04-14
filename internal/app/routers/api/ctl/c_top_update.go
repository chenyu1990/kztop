package ctl

import (
	"github.com/gin-gonic/gin"
	"kztop/internal/app/errors"
	"kztop/internal/app/ginplus"
	"kztop/internal/app/schema"
	"kztop/pkg/kreedz"
	"strconv"
	"strings"
	"time"
)

func (a *Top) UpdateRecord(c *gin.Context) {
	newRecord := schema.Record{}
	if err := ginplus.ParseJSON(c, &newRecord); err != nil {
		ginplus.ResError(c, err)
		return
	}

	// AMXX 的 json 不支持 float
	t, err := strconv.ParseFloat(newRecord.TimeString, 10)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	newRecord.Time = t

	if newRecord.Validation() == false {
		ginplus.ResError(c, errors.New400Response("dont hack me"))
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
	if kreedz.IsSteamID(newRecord.SteamID) == false {
		ginplus.ResError(c, errors.New400Response("不支持盗版玩家进入排行"))
		return
	}

	if a.playerStat[newRecord.SteamID] == nil {
		a.playerStat[newRecord.SteamID] = make(map[schema.Cate]int64)
	}
	/*
		1. 更新数据库 done
		2. 更新地区统计（当玩家数据不存在时） done
		3. 更新玩家信息 done
		4. 更新第一信息 done
		5. 更新排序 done
		6. 返回最新排行
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
	newRecord.Date = time.Now().UTC()
	var oldRecord *schema.Record
	if len(query.Data) == 0 {
		newRecord.FinishCount = 1
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
		oldRecord = query.Data[0]
		newRecord.FinishCount = oldRecord.FinishCount + 1
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
	newRegion, newNick := a.UpdateInfo(strings.ToLower(newRecord.Region), newRecord.SteamID, newRecord.Nick)
	if newRegion != "" || newNick != "" {
		go a.RecordModel.UpdateInfo(ctx, schema.UpdateInfo{
			SteamID: newRecord.SteamID,
			Region:  newRegion,
			Nick:    newNick,
		})
	}

	// 历史纪录不存在，或旧纪录过慢，才有可能成为第一
	if oldRecord == nil || oldRecord.Time > newRecord.Time {
		go a.SwapMapFirst(&newRecord)
	}
	// 更新统计一定要在 SwapMapFirst 之后
	go a.UpdateStats(cate)

	orders := schema.Orders{}
	if cate == schema.WPN {
		orders["speed"] = schema.OrderASC
	}
	orders["time"] = schema.OrderASC
	query, err = a.RecordModel.Query(ctx, &schema.RecordQueryParam{
		Cate:    cate,
		MapName: newRecord.MapName,
	}, schema.RecordQueryOptions{
		PageParam:  &schema.PaginationParam{
			PageSize:  100,
		},
		OrderParam: &schema.OrderParam{
			Orders: orders,
		},
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 因为历史纪录存在非正版玩家的数据，所以此处没有做全局缓存
	// 而且每张地图100个排名，缓存数据量比较大，15个排名的可以考虑
	var rank int
	for i, record := range query.Data {
		if record.SteamID == newRecord.SteamID {
			rank = i + 1
		}
	}

	ginplus.ResSuccess(c, rank)
}