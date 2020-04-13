package entity

import (
	"context"
	"github.com/jinzhu/gorm"
	"kztop/internal/app/schema"
)

// GetRegionDB 获取Region存储
func GetRegionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, Region{})
}

// SchemaRegion Region对象
type SchemaRegion schema.Region

// ToRegion 转换为Region实体
func (a SchemaRegion) ToRegion() *Region {
	item := &Region{
		Chinese:  &a.Chinese,
		English:  &a.English,
		SortName: &a.SortName,
	}
	return item
}

// Region Region实体
type Region struct {
	Chinese  *string `gorm:"column:cn;size:64"`
	English  *string `gorm:"column:en;size:64"`
	SortName *string `gorm:"column:sort_name;size:3"`
}

func (a Region) String() string {
	return toString(a)
}

// TableName 表名
func (a Region) TableName() string {
	return "kz_region"
}

// ToSchemaRegion 转换为Region对象
func (a Region) ToSchemaRegion() *schema.Region {
	item := &schema.Region{
		Chinese:  *a.Chinese,
		English:  *a.English,
		SortName: *a.SortName,
	}
	return item
}

// Regions Region列表
type Regions []*Region

// ToSchemaRegions 转换为Region对象列表
func (a Regions) ToSchemaRegions() []*schema.Region {
	list := make([]*schema.Region, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRegion()
	}
	return list
}
