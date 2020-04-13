package model

import (
	"context"
	"kztop/pkg/kreedz"

	"kztop/internal/app/errors"
	"kztop/internal/app/model/impl/gorm/internal/entity"
	"kztop/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// NewRegion 创建Region存储实例
func NewRegion(db *gorm.DB) *Region {
	return &Region{db}
}

// Region Region存储
type Region struct {
	db *gorm.DB
}

func (a *Region) getQueryOption(opts ...schema.RegionQueryOptions) schema.RegionQueryOptions {
	var opt schema.RegionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *Region) where(db *gorm.DB, params *schema.RegionQueryParam) *gorm.DB {
	if v := params.SortName; v != "" {
		db = db.Where("sort_name=?", v)
	}

	return db
}

// Query 查询数据
func (a *Region) Query(ctx context.Context, params schema.RegionQueryParam, opts ...schema.RegionQueryOptions) (*schema.RegionQueryResult, error) {
	db := entity.GetRegionDB(ctx, a.db)
	db = a.where(db, &params)

	opt := a.getQueryOption(opts...)
	var list entity.Regions
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RegionQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRegions(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Region) Get(ctx context.Context, params schema.RegionQueryParam) (*schema.Region, error) {
	db := entity.GetRegionDB(ctx, a.db)
	db = a.where(db, &params)

	var item entity.Region
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaRegion(), nil
}

// Create 创建数据
func (a *Region) Create(ctx context.Context, item schema.Region) error {
	return ExecTrans(ctx, a.db, func(i context.Context) error {
		Region := entity.SchemaRegion(item).ToRegion()
		result := entity.GetRegionDB(i, a.db).Create(Region)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Update 更新数据
func (a *Region) Update(ctx context.Context, organization kreedz.Organization, mapname string, item schema.Region) error {
	return ExecTrans(ctx, a.db, func(i context.Context) error {
		Region := entity.SchemaRegion(item).ToRegion()
		result := entity.GetRegionDB(i, a.db).
			Where("organization=?", organization).
			Where("mapname=?", mapname).
			Omit("organization", "mapname").Updates(Region)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Delete 删除数据
func (a *Region) Delete(ctx context.Context, typ uint) error {
	result := entity.GetRegionDB(ctx, a.db).
		Where("type=?", typ).
		Delete(entity.Region{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Region) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetRegionDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
