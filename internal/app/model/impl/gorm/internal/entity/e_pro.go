package entity

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"kztop/internal/app/schema"
)

// GetProDB 获取Pro存储
func GetProDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, Pro{})
}

// SchemaPro Pro对象
type SchemaPro schema.Pro

// ToPro 转换为Pro实体
func (a SchemaPro) ToPro() *Pro {
	item := &Pro{
		MapName:     &a.MapName,
		AuthID:      &a.AuthID,
		Country:     &a.Country,
		Name:        &a.Name,
		Time:        &a.Time,
		Weapon:      &a.Weapon,
		FinishCount: &a.FinishCount,
		Server:      &a.Server,
		Route:       &a.Route,
		Date:        &a.Date,
	}
	return item
}

// Pro Pro实体
type Pro struct {
	MapName     *string    `gorm:"column:mapname;size:64"`
	AuthID      *string    `gorm:"column:authid;size:64"`
	Country     *string    `gorm:"column:country;size:6"`
	Name        *string    `gorm:"column:name;size:64"`
	Time        *float64   `gorm:"column:time"`
	Weapon      *string    `gorm:"column:weapon;size:64"`
	FinishCount *int       `gorm:"column:fincnt"`
	Server      *string    `gorm:"column:server;size:64"`
	Route       *string    `gorm:"column:route;size:16"`
	Date        *time.Time `gorm:"column:date"`
}

func (a Pro) String() string {
	return toString(a)
}

// TableName 表名
func (a Pro) TableName() string {
	return "kz_pro15"
}

// ToSchemaPro 转换为Pro对象
func (a Pro) ToSchemaPro() *schema.Pro {
	item := &schema.Pro{
		MapName:     *a.MapName,
		AuthID:      *a.AuthID,
		Country:     *a.Country,
		Name:        *a.Name,
		Time:        *a.Time,
		Weapon:      *a.Weapon,
		FinishCount: *a.FinishCount,
		Server:      *a.Server,
		Route:       *a.Route,
		Date:        *a.Date,
	}
	return item
}

// Pros Pro列表
type Pros []*Pro

// ToSchemaPros 转换为Pro对象列表
func (a Pros) ToSchemaPros() []*schema.Pro {
	list := make([]*schema.Pro, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaPro()
	}
	return list
}
