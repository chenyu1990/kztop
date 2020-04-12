package ctl

import (
	"github.com/gin-gonic/gin"
	"kztop/internal/app/ginplus"
	"kztop/internal/app/model"
	"kztop/internal/app/schema"
	"net/http"
)

// NewRecord 创建record控制器
func NewRecord(
	mWorldRecord model.IWorldRecord,
) *Record {
	return &Record{
		WorldRecordModel: mWorldRecord,
	}
}

type Record struct {
	WorldRecordModel model.IWorldRecord
}

func (a *Record) Query(c *gin.Context) {
	mapname := c.Param("mapname")
	params := schema.WorldRecordQueryParam{}
	if mapname != "" {
		params.MapName = mapname
	}

	query, err := a.WorldRecordModel.Query(ginplus.NewContext(c), params)
	if err != nil {
		panic(err)
	}

	h := gin.H{
		"list": query.Data,
	}

	if mapname != "" {
		ginplus.ResJSON(c, http.StatusOK, query.Data[0])
	} else {
		c.HTML(http.StatusOK, "record/index", h)
	}
}
