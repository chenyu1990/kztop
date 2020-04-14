package ctl

import (
	"github.com/gin-gonic/gin"
	"kztop/internal/app/errors"
	"kztop/internal/app/ginplus"
	"kztop/internal/app/schema"
	"kztop/pkg/Steam"
	"kztop/pkg/kreedz"
	"kztop/pkg/util"
	"net/http"
	"strconv"
)

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

	h := gin.H{
		"cate":  cate,
		"pCate": pCate,
	}
	if _, exist := a.players[cate][player]; !exist {

	} else {
		// 以完成次数排序
		finishCountStat := make(map[string]int)
		var pageSize uint64
		for i, record := range a.players[cate][player] {
			pageSize = uint64(i + 1)
			finishCountStat[record.MapName] = record.FinishCount
		}
		finishCountSort := sortMapStringInt(finishCountStat)

		list := make(map[string]*schema.Record)
		for _, record := range a.players[cate][player] {
			list[record.MapName] = record
		}

		h["list"] = list
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
					util.HandleHttpError(err)
					return
				}
				a.playerInfo[player]["avatarFull"] = profile.AvatarFull
				a.playerInfo[player]["onlineState"] = profile.OnlineState
				a.playerInfo[player]["visibilityState"] = Steam.VisibilityState[profile.VisibilityState]
			}(steamID64, player)
		}

		if pageSize > schema.PageSize {
			pageSize = schema.PageSize
		}
		h["sort"] = finishCountSort[:pageSize]
	}
	h["player"] = player
	h["stat"] = a.playerStat[player]
	h["info"] = a.playerInfo[player]
	h["region"] = a.regionsInfo
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