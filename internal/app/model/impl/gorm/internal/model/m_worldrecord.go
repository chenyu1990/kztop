package model

import (
	"context"
	"kztop/pkg/kreedz"

	"kztop/internal/app/errors"
	"kztop/internal/app/model/impl/gorm/internal/entity"
	"kztop/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// NewWorldRecord 创建WorldRecord存储实例
func NewWorldRecord(db *gorm.DB) *WorldRecord {
	return &WorldRecord{db}
}

// WorldRecord WorldRecord存储
type WorldRecord struct {
	db *gorm.DB
}

func (a *WorldRecord) getQueryOption(opts ...schema.WorldRecordQueryOptions) schema.WorldRecordQueryOptions {
	var opt schema.WorldRecordQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *WorldRecord) where(db *gorm.DB, params *schema.WorldRecordQueryParam) *gorm.DB {
	if v := params.Organization; v > 0 {
		db = db.Where("organization=?", v)
	}
	if v := params.MapName; v != "" {
		db = db.Where("mapname=?", v)
	}
	if v := params.Holder; v != "" {
		db = db.Where("holder=?", v)
	}
	if v := params.Country; v != "" {
		db = db.Where("country=?", v)
	}
	if v := params.Organizations; v != nil && len(v) > 1 {
		db = db.Where("organization=? or organization=?", v[0], v[1])
	}


	return db
}

// Query 查询数据
func (a *WorldRecord) Query(ctx context.Context, params schema.WorldRecordQueryParam, opts ...schema.WorldRecordQueryOptions) (*schema.WorldRecordQueryResult, error) {
	db := entity.GetWorldRecordDB(ctx, a.db)
	db = a.where(db, &params)
	db = db.Order("`time` ASC")

	opt := a.getQueryOption(opts...)
	var list entity.WorldRecords
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.WorldRecordQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaWorldRecords(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *WorldRecord) Get(ctx context.Context, params schema.WorldRecordQueryParam) (*schema.WorldRecord, error) {
	db := entity.GetWorldRecordDB(ctx, a.db)
	db = a.where(db, &params)

	var item entity.WorldRecord
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaWorldRecord(), nil
}

// Create 创建数据
func (a *WorldRecord) Create(ctx context.Context, item schema.WorldRecord) error {
	return ExecTrans(ctx, a.db, func(i context.Context) error {
		WorldRecord := entity.SchemaWorldRecord(item).ToWorldRecord()
		result := entity.GetWorldRecordDB(i, a.db).Create(WorldRecord)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Update 更新数据
func (a *WorldRecord) Update(ctx context.Context, organization kreedz.Organization, mapname string, item schema.WorldRecord) error {
	return ExecTrans(ctx, a.db, func(i context.Context) error {
		WorldRecord := entity.SchemaWorldRecord(item).ToWorldRecord()
		result := entity.GetWorldRecordDB(i, a.db).
			Where("organization=?", organization).
			Where("mapname=?", mapname).
			Omit("organization", "mapname").Updates(WorldRecord)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Delete 删除数据
func (a *WorldRecord) Delete(ctx context.Context, typ uint) error {
	result := entity.GetWorldRecordDB(ctx, a.db).
		Where("type=?", typ).
		Delete(entity.WorldRecord{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *WorldRecord) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetWorldRecordDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
