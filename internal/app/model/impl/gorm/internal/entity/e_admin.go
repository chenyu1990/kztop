package entity

import (
	"context"
	"github.com/jinzhu/gorm"
	"kztop/internal/app/schema"
)

// GetAdminDB 获取Admin存储
func GetAdminDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, Admin{})
}

// SchemaAdmin Admin对象
type SchemaAdmin schema.Admin

//// ToAdmin 转换为Admin实体
func (a SchemaAdmin) ToAdmin() *Admin {
	item := &Admin{
		Server:  &a.Server,
		SteamID: &a.SteamID,
		Access:  &a.Access,
		Valid:   &a.Valid,
	}
	return item
}

// Admin Admin实体
type Admin struct {
	Server  *string `gorm:"column:server;size:64"`
	SteamID *string `gorm:"column:steamid;size:32"`
	Access  *string `gorm:"column:access;size:26"`
	Valid   *bool   `gorm:"column:valid"`
}

func (a Admin) String() string {
	return toString(a)
}

// TableName 表名
func (a Admin) TableName() string {
	return "kz_admin"
}

// ToSchemaAdmin 转换为Admin对象
func (a Admin) ToSchemaAdmin() *schema.Admin {
	item := &schema.Admin{
		Server:  *a.Server,
		SteamID: *a.SteamID,
		Access:  *a.Access,
		Valid:   *a.Valid,
	}
	return item
}

// Admins Admin列表
type Admins []*Admin

// ToSchemaAdmins 转换为Admin对象列表
func (a Admins) ToSchemaAdmins() []*schema.Admin {
	list := make([]*schema.Admin, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaAdmin()
	}
	return list
}
