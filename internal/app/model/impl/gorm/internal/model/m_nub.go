package model

import (
	"context"

	"kztop/internal/app/errors"
	"kztop/internal/app/model/impl/gorm/internal/entity"
	"kztop/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// NewNub 创建Nub存储实例
func NewNub(db *gorm.DB) *Nub {
	return &Nub{db}
}

// Nub Nub存储
type Nub struct {
	db *gorm.DB
}

func (a *Nub) getQueryOption(opts ...schema.NubQueryOptions) schema.NubQueryOptions {
	var opt schema.NubQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *Nub) Query(ctx context.Context, params schema.NubQueryParam, opts ...schema.NubQueryOptions) (*schema.NubQueryResult, error) {
	db := entity.GetNubDB(ctx, a.db)
	if v := params.MapName; v != "" {
		db = db.Where("mapname=?", v)
	}
	if v := params.AuthID; v != "" {
		db = db.Where("authid=?", v)
	}
	db = db.Order("`time` ASC")

	opt := a.getQueryOption(opts...)
	var list entity.Nubs
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.NubQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaNubs(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Nub) Get(ctx context.Context, recordID string, opts ...schema.NubQueryOptions) (*schema.Nub, error) {
	db := entity.GetNubDB(ctx, a.db).Where("record_id=?", recordID)
	var item entity.Nub
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaNub(), nil
}

// Create 创建数据
func (a *Nub) Create(ctx context.Context, item schema.Nub) error {
	Nub := entity.SchemaNub(item).ToNub()
	result := entity.GetNubDB(ctx, a.db).Create(Nub)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *Nub) Update(ctx context.Context, recordID string, item schema.Nub) error {
	Nub := entity.SchemaNub(item).ToNub()
	result := entity.GetNubDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(Nub)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateInfo 更新信息
func (a *Nub) UpdateInfo(ctx context.Context, info schema.UpdateInfo) error {
	result := entity.GetProDB(ctx, a.db).Where("authid=?", info.AuthID).Omit("authid").Updates(info)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *Nub) Delete(ctx context.Context, recordID string) error {
	result := entity.GetNubDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.Nub{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Nub) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetNubDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
