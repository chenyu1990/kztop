package ctl

import (
	"github.com/gin-gonic/gin"
	"kztop/internal/app/ginplus"
	"kztop/internal/app/schema"
	"kztop/pkg/kreedz"
	"net/http"
)

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