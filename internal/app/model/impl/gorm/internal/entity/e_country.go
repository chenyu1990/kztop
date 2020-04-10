package entity

import (
	"context"
	"github.com/jinzhu/gorm"
	"kztop/internal/app/schema"
)

// GetCountryDB 获取Country存储
func GetCountryDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, Country{})
}

// SchemaCountry Country对象
type SchemaCountry schema.Country

// ToCountry 转换为Country实体
func (a SchemaCountry) ToCountry() *Country {
	item := &Country{
		Chinese:  &a.Chinese,
		English:  &a.English,
		SortName: &a.SortName,
	}
	return item
}

// Country Country实体
type Country struct {
	Chinese  *string `gorm:"column:cn;size:64"`
	English  *string `gorm:"column:en;size:64"`
	SortName *string `gorm:"column:sort_name;size:3"`
}

func (a Country) String() string {
	return toString(a)
}

// TableName 表名
func (a Country) TableName() string {
	return "kz_country"
}

// ToSchemaCountry 转换为Country对象
func (a Country) ToSchemaCountry() *schema.Country {
	item := &schema.Country{
		Chinese:  *a.Chinese,
		English:  *a.English,
		SortName: *a.SortName,
	}
	return item
}

// Countrys Country列表
type Countrys []*Country

// ToSchemaCountrys 转换为Country对象列表
func (a Countrys) ToSchemaCountrys() []*schema.Country {
	list := make([]*schema.Country, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaCountry()
	}
	return list
}
