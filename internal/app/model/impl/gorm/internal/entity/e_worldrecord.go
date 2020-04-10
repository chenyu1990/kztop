package entity

import (
	"context"
	"github.com/jinzhu/gorm"
	"kztop/internal/app/schema"
	"kztop/pkg/kreedz"
)

// GetWorldRecordDB 获取WorldRecord存储
func GetWorldRecordDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return getDBWithModel(ctx, defDB, WorldRecord{})
}

// SchemaWorldRecord WorldRecord对象
type SchemaWorldRecord schema.WorldRecord

// ToWorldRecord 转换为WorldRecord实体
func (a SchemaWorldRecord) ToWorldRecord() *WorldRecord {
	item := &WorldRecord{
		MapName:      &a.MapName,
		Holder:       &a.Holder,
		Country:      &a.Country,
		Time:         &a.Time,
		Organization: &a.Organization,
	}
	return item
}

// WorldRecord WorldRecord实体
type WorldRecord struct {
	MapName      *string              `gorm:"column:mapname;size:64;unique_index:record"`
	Holder       *string              `gorm:"column:holder;size:64"`
	Country      *string              `gorm:"column:country;size:3"`
	Time         *float64             `gorm:"column:time"`
	Organization *kreedz.Organization `gorm:"column:organization;size:1;unique_index:record"`
}

func (a WorldRecord) String() string {
	return toString(a)
}

// TableName 表名
func (a WorldRecord) TableName() string {
	return "kz_record"
}

// ToSchemaWorldRecord 转换为WorldRecord对象
func (a WorldRecord) ToSchemaWorldRecord() *schema.WorldRecord {
	item := &schema.WorldRecord{
		MapName:      *a.MapName,
		Holder:       *a.Holder,
		Country:      *a.Country,
		Time:         *a.Time,
		Organization: *a.Organization,
	}
	return item
}

// WorldRecords WorldRecord列表
type WorldRecords []*WorldRecord

// ToSchemaWorldRecords 转换为WorldRecord对象列表
func (a WorldRecords) ToSchemaWorldRecords() []*schema.WorldRecord {
	list := make([]*schema.WorldRecord, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaWorldRecord()
	}
	return list
}
