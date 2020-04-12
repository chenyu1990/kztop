package model

import (
	"context"

	"kztop/internal/app/errors"
	"kztop/internal/app/model/impl/gorm/internal/entity"
	"kztop/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// NewPro 创建Pro存储实例
func NewPro(db *gorm.DB) *Pro {
	return &Pro{db}
}

// Pro Pro存储
type Pro struct {
	db *gorm.DB
}

func (a *Pro) getQueryOption(opts ...schema.ProQueryOptions) schema.ProQueryOptions {
	var opt schema.ProQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *Pro) where(db *gorm.DB, params *schema.ProQueryParam) *gorm.DB {
	if v := params.MapName; v != "" {
		db = db.Where("mapname=?", v)
	}
	if v := params.AuthID; v != "" {
		db = db.Where("authid=?", v)
	}

	return db
}

// Query 查询数据
func (a *Pro) Query(ctx context.Context, params schema.ProQueryParam, opts ...schema.ProQueryOptions) (*schema.ProQueryResult, error) {
	db := entity.GetProDB(ctx, a.db)
	db = a.where(db, &params)
	db = db.Order("`time` ASC")

	opt := a.getQueryOption(opts...)
	var list entity.Pros
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.ProQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaPros(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Pro) Get(ctx context.Context, recordID string, opts ...schema.ProQueryOptions) (*schema.Pro, error) {
	db := entity.GetProDB(ctx, a.db).Where("record_id=?", recordID)
	var item entity.Pro
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaPro(), nil
}

// Create 创建数据
func (a *Pro) Create(ctx context.Context, item schema.Pro) error {
	Pro := entity.SchemaPro(item).ToPro()
	result := entity.GetProDB(ctx, a.db).Create(Pro)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *Pro) Update(ctx context.Context, params *schema.ProQueryParam, item schema.Pro) error {
	Pro := entity.SchemaPro(item).ToPro()
	db := entity.GetProDB(ctx, a.db)
	db = a.where(db, params)
	result := db.Omit("mapname", "authid").Updates(Pro)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *Pro) Delete(ctx context.Context, recordID string) error {
	result := entity.GetProDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.Pro{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Pro) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetProDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
