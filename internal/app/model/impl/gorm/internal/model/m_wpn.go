package model

import (
	"context"

	"kztop/internal/app/errors"
	"kztop/internal/app/model/impl/gorm/internal/entity"
	"kztop/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// NewWpn 创建Wpn存储实例
func NewWpn(db *gorm.DB) *Wpn {
	return &Wpn{db}
}

// Wpn Wpn存储
type Wpn struct {
	db *gorm.DB
}

func (a *Wpn) getQueryOption(opts ...schema.WpnQueryOptions) schema.WpnQueryOptions {
	var opt schema.WpnQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *Wpn) Query(ctx context.Context, params schema.WpnQueryParam, opts ...schema.WpnQueryOptions) (*schema.WpnQueryResult, error) {
	db := entity.GetWpnDB(ctx, a.db)
	if v := params.MapName; v != "" {
		db = db.Where("mapname=?", v)
	}
	if v := params.AuthID; v != "" {
		db = db.Where("authid=?", v)
	}
	db = db.Order("speed ASC, `time` ASC")

	opt := a.getQueryOption(opts...)
	var list entity.Wpns
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.WpnQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaWpns(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Wpn) Get(ctx context.Context, recordID string, opts ...schema.WpnQueryOptions) (*schema.Wpn, error) {
	db := entity.GetWpnDB(ctx, a.db).Where("record_id=?", recordID)
	var item entity.Wpn
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaWpn(), nil
}

// Create 创建数据
func (a *Wpn) Create(ctx context.Context, item schema.Wpn) error {
	Wpn := entity.SchemaWpn(item).ToWpn()
	result := entity.GetWpnDB(ctx, a.db).Create(Wpn)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *Wpn) Update(ctx context.Context, recordID string, item schema.Wpn) error {
	Wpn := entity.SchemaWpn(item).ToWpn()
	result := entity.GetWpnDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(Wpn)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateInfo 更新信息
func (a *Wpn) UpdateInfo(ctx context.Context, info schema.UpdateInfo) error {
	result := entity.GetProDB(ctx, a.db).Where("authid=?", info.AuthID).Omit("authid").Updates(info)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *Wpn) Delete(ctx context.Context, recordID string) error {
	result := entity.GetWpnDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.Wpn{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Wpn) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetWpnDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
