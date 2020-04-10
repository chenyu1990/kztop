package entity

import (
	"context"
	"github.com/jinzhu/gorm"
	"kztop/internal/app/schema"
	"kztop/pkg/kreedz"
	"time"
)

// GetNewsDB 获取News存储
func GetNewsDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, News{})
}

// SchemaNews News对象
type SchemaNews schema.News

// ToNews 转换为News实体
func (a SchemaNews) ToNews() *News {
	item := &News{
		Organization: &a.Organization,
		Period:       &a.Period,
		Data:         &a.Data,
		Date:         &a.Date,
	}
	return item
}

// News News实体
type News struct {
	Organization *kreedz.Organization `gorm:"column:organization;size:1"`
	Period       *uint                `gorm:"column:period"`
	Data         *string              `gorm:"column:data;type:blob"`
	Date         *time.Time           `gorm:"column:date"`
}

func (a News) String() string {
	return toString(a)
}

// TableName 表名
func (a News) TableName() string {
	return "kz_news"
}

// ToSchemaNews 转换为News对象
func (a News) ToSchemaNews() *schema.News {
	item := &schema.News{
		Organization: *a.Organization,
		Period:       *a.Period,
		Data:         *a.Data,
		Date:         *a.Date,
	}
	return item
}

// Newss News列表
type Newss []*News

// ToSchemaNewss 转换为News对象列表
func (a Newss) ToSchemaNewss() []*schema.News {
	list := make([]*schema.News, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaNews()
	}
	return list
}
