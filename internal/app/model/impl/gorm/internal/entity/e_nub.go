package entity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"kztop/internal/app/schema"
)

// GetNubDB 获取Nub存储
func GetNubDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, Nub{})
}

// SchemaNub Nub对象
type SchemaNub schema.Nub

// ToNub 转换为Nub实体
func (a SchemaNub) ToNub() *Nub {
	item := &Nub{
		MapName:     &a.MapName,
		AuthID:      &a.AuthID,
		Country:     &a.Country,
		Name:        &a.Name,
		Time:        &a.Time,
		Weapon:      &a.Weapon,
		FinishCount: &a.FinishCount,
		Server:      &a.Server,
		CheckPoints: &a.CheckPoints,
		GoChecks:    &a.GoChecks,
		Route:       &a.Route,
		Date:        &a.Date,
	}
	return item
}

// Nub Nub实体
type Nub struct {
	MapName     *string    `gorm:"column:mapname;size:64"`
	AuthID      *string    `gorm:"column:authid;size:64"`
	Country     *string    `gorm:"column:country;size:6"`
	Name        *string    `gorm:"column:name;size:64"`
	Time        *float64   `gorm:"column:time"`
	Weapon      *string    `gorm:"column:weapon;size:64"`
	FinishCount *int       `gorm:"column:fincnt"`
	Server      *string    `gorm:"column:server;size:64"`
	CheckPoints *int       `gorm:"column:checkpoints"`
	GoChecks    *int       `gorm:"column:gochecks"`
	Route       *string    `gorm:"column:route;size:16"`
	Date        *time.Time `gorm:"column:date"`
}

func (a Nub) String() string {
	return toString(a)
}

// TableName 表名
func (a Nub) TableName() string {
	return "kz_nub15"
}

// ToSchemaNub 转换为Nub对象
func (a Nub) ToSchemaNub() *schema.Nub {
	item := &schema.Nub{
		MapName:     *a.MapName,
		AuthID:      *a.AuthID,
		Country:     *a.Country,
		Name:        *a.Name,
		Time:        *a.Time,
		Weapon:      *a.Weapon,
		FinishCount: *a.FinishCount,
		Server:      *a.Server,
		CheckPoints: *a.CheckPoints,
		GoChecks:    *a.GoChecks,
		Route:       *a.Route,
		Date:        *a.Date,
	}
	return item
}

// Nubs Nub列表
type Nubs []*Nub

// ToSchemaNubs 转换为Nub对象列表
func (a Nubs) ToSchemaNubs() []*schema.Nub {
	list := make([]*schema.Nub, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaNub()
	}
	return list
}
