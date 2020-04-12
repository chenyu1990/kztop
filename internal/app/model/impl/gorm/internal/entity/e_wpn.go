package entity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"kztop/internal/app/schema"
)

// GetWpnDB 获取Wpn存储
func GetWpnDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, Wpn{})
}

// SchemaWpn Wpn对象
type SchemaWpn schema.Wpn

// ToWpn 转换为Wpn实体
func (a SchemaWpn) ToWpn() *Wpn {
	item := &Wpn{
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
		Speed:       &a.Speed,
		Date:        &a.Date,
	}
	return item
}

// Wpn Wpn实体
type Wpn struct {
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
	Speed       *int       `gorm:"column:speed"`
	Date        *time.Time `gorm:"column:date"`
}

func (a Wpn) String() string {
	return toString(a)
}

// TableName 表名
func (a Wpn) TableName() string {
	return "kz_wpn15"
}

// ToSchemaWpn 转换为Wpn对象
func (a Wpn) ToSchemaWpn() *schema.Wpn {
	item := &schema.Wpn{
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
		Speed:       *a.Speed,
		Date:        *a.Date,
	}
	return item
}

// Wpns Wpn列表
type Wpns []*Wpn

// ToSchemaWpns 转换为Wpn对象列表
func (a Wpns) ToSchemaWpns() []*schema.Wpn {
	list := make([]*schema.Wpn, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaWpn()
	}
	return list
}
