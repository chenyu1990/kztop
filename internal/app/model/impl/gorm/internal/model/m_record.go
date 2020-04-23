package model

import (
	"context"

	"github.com/jinzhu/gorm"
	"kztop/internal/app/errors"
	"kztop/internal/app/model/impl/gorm/internal/entity"
	"kztop/internal/app/schema"
)

// NewRecord 创建Record存储实例
func NewRecord(db *gorm.DB) *Record {
	return &Record{db}
}

// Record Record存储
type Record struct {
	db *gorm.DB
}

func (a *Record) getQueryOption(opts ...schema.RecordQueryOptions) schema.RecordQueryOptions {
	var opt schema.RecordQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *Record) where(db *gorm.DB, params *schema.RecordQueryParam) *gorm.DB {
	if v := params.Cate; v > 0 {
		db = db.Where("cate=?", v)
	}
	if v := params.MapName; v != "" {
		db = db.Where("mapname=?", v)
	}
	if v := params.SteamID; v != "" {
		db = db.Where("steam_id=?", v)
	}
	if v := params.Route; v != nil {
		db = db.Where("route=?", v)
	}

	return db
}

// Query 查询数据
func (a *Record) Query(ctx context.Context, params *schema.RecordQueryParam, opts ...schema.RecordQueryOptions) (*schema.RecordQueryResult, error) {
	db := entity.GetRecordDB(ctx, a.db)
	if params != nil {
		db = a.where(db, params)
	}

	opt := a.getQueryOption(opts...)
	if opt.OrderParam != nil {
		db = db.Order(Order(opt.OrderParam.Orders))
	}

	var list entity.Records
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RecordQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRecords(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Record) Get(ctx context.Context, recordID string, opts ...schema.RecordQueryOptions) (*schema.Record, error) {
	db := entity.GetRecordDB(ctx, a.db).Where("record_id=?", recordID)
	var item entity.Record
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaRecord(), nil
}

// Create 创建数据
func (a *Record) Create(ctx context.Context, item *schema.Record) error {
	Record := entity.SchemaRecord(*item).ToRecord()
	result := entity.GetRecordDB(ctx, a.db).Create(Record)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *Record) Update(ctx context.Context, params *schema.RecordQueryParam, item schema.Record) error {
	Record := entity.SchemaRecord(item).ToRecord()
	db := entity.GetRecordDB(ctx, a.db)
	db = a.where(db, params)
	result := db.Omit("cate", "mapname", "steam_id", "route").Updates(Record)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateInfo 更新信息
func (a *Record) UpdateInfo(ctx context.Context, info schema.UpdateInfo) error {
	result := entity.GetRecordDB(ctx, a.db).Where("steam_id=?", info.SteamID).Omit("steam_id").Updates(info)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *Record) Delete(ctx context.Context, recordID string) error {
	result := entity.GetRecordDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.Record{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Record) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetRecordDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
