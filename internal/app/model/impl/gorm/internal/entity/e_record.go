package entity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"kztop/internal/app/schema"
)

// GetRecordDB 获取Record存储
func GetRecordDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, Record{})
}

// SchemaRecord Record对象
type SchemaRecord schema.Record

// ToRecord 转换为Record实体
func (a SchemaRecord) ToRecord() *Record {
	item := &Record{
		Cate:        &a.Cate,
		MapName:     &a.MapName,
		SteamID:     &a.SteamID,
		Region:      &a.Region,
		Nick:        &a.Nick,
		Time:        &a.Time,
		Weapon:      &a.Weapon,
		FinishCount: &a.FinishCount,
		Server:      &a.Server,
		CheckPoints: &a.CheckPoints,
		GoChecks:    &a.GoChecks,
		Route:       &a.Route,
		Speed:       &a.Speed,
		Date:        &a.Date,
	}
	return item
}

// Record Record实体
type Record struct {
	Cate        *schema.Cate `gorm:"column:cate;index:top"`
	MapName     *string      `gorm:"column:mapname;size:64;index:top"`
	SteamID     *string      `gorm:"column:steam_id;size:64"`
	Region      *string      `gorm:"column:region;size:6"`
	Nick        *string      `gorm:"column:nick;size:64"`
	Time        *float64     `gorm:"column:time"`
	Weapon      *string      `gorm:"column:weapon;size:64"`
	FinishCount *int         `gorm:"column:finish_count"`
	Server      *string      `gorm:"column:server;size:64"`
	CheckPoints *int         `gorm:"column:check_points"`
	GoChecks    *int         `gorm:"column:go_checks"`
	Route       *string      `gorm:"column:route;size:16"`
	Speed       *int         `gorm:"column:speed"`
	Date        *time.Time   `gorm:"column:date"`
}

func (a Record) String() string {
	return toString(a)
}

// TableName 表名
func (a Record) TableName() string {
	return "kz_record"
}

// ToSchemaRecord 转换为Record对象
func (a Record) ToSchemaRecord() *schema.Record {
	item := &schema.Record{
		Cate:        *a.Cate,
		MapName:     *a.MapName,
		SteamID:     *a.SteamID,
		Region:      *a.Region,
		Nick:        *a.Nick,
		Time:        *a.Time,
		Weapon:      *a.Weapon,
		FinishCount: *a.FinishCount,
		Server:      *a.Server,
		CheckPoints: *a.CheckPoints,
		GoChecks:    *a.GoChecks,
		Route:       *a.Route,
		Speed:       *a.Speed,
		Date:        *a.Date,
	}
	return item
}

// Records Record列表
type Records []*Record

// ToSchemaRecords 转换为Record对象列表
func (a Records) ToSchemaRecords() []*schema.Record {
	list := make([]*schema.Record, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRecord()
	}
	return list
}
