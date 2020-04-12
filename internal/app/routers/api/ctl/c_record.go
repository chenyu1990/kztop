package ctl

import (
	"github.com/gin-gonic/gin"
	"kztop/internal/app/ginplus"
	"kztop/internal/app/model"
	"kztop/internal/app/schema"
	"kztop/pkg/kreedz"
	"net/http"
	"strings"
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

func (a *Record) Org(c *gin.Context) {
	organization := c.Param("org")
	params := schema.WorldRecordQueryParam{}
	if organization != "" {
		for i, name := range kreedz.SortName {
			if strings.ToLower(name) == strings.ToLower(organization) {
				params.Organization = kreedz.Organization(i)
			}
		}
	} else {
		params.Organization = kreedz.XtremeJumps
	}

	query, err := a.WorldRecordModel.Query(ginplus.NewContext(c), params)
	if err != nil {
		panic(err)
	}

	h := gin.H{
		"organizations": kreedz.Name[1:],
		"selectedOrg": params.Organization - 1,
		"orgSortName": kreedz.SortName[1:],
		"list": query.Data,
	}

	c.HTML(http.StatusOK, "record/index", h)
}

func (a *Record) Map(c *gin.Context) {
	mapname := c.Param("mapname")

	query, err := a.WorldRecordModel.Query(ginplus.NewContext(c), schema.WorldRecordQueryParam{
		MapName: mapname,
	})
	if err != nil {
		panic(err)
	}

	if len(query.Data) > 0 {
		ginplus.ResJSON(c, http.StatusOK, query.Data[0])
	} else {
		ginplus.ResJSON(c, http.StatusOK, schema.WorldRecord{})
	}
}
