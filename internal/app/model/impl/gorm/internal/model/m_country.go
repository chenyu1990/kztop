package model

import (
	"context"
	"kztop/pkg/kreedz"

	"kztop/internal/app/errors"
	"kztop/internal/app/model/impl/gorm/internal/entity"
	"kztop/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// NewCountry 创建Country存储实例
func NewCountry(db *gorm.DB) *Country {
	return &Country{db}
}

// Country Country存储
type Country struct {
	db *gorm.DB
}

func (a *Country) getQueryOption(opts ...schema.CountryQueryOptions) schema.CountryQueryOptions {
	var opt schema.CountryQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *Country) where(db *gorm.DB, params *schema.CountryQueryParam) *gorm.DB {
	if v := params.SortName; v != "" {
		db = db.Where("sort_name=?", v)
	}

	return db
}

// Query 查询数据
func (a *Country) Query(ctx context.Context, params schema.CountryQueryParam, opts ...schema.CountryQueryOptions) (*schema.CountryQueryResult, error) {
	db := entity.GetCountryDB(ctx, a.db)
	db = a.where(db, &params)

	opt := a.getQueryOption(opts...)
	var list entity.Countrys
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.CountryQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaCountrys(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Country) Get(ctx context.Context, params schema.CountryQueryParam) (*schema.Country, error) {
	db := entity.GetCountryDB(ctx, a.db)
	db = a.where(db, &params)

	var item entity.Country
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaCountry(), nil
}

// Create 创建数据
func (a *Country) Create(ctx context.Context, item schema.Country) error {
	return ExecTrans(ctx, a.db, func(i context.Context) error {
		Country := entity.SchemaCountry(item).ToCountry()
		result := entity.GetCountryDB(i, a.db).Create(Country)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Update 更新数据
func (a *Country) Update(ctx context.Context, organization kreedz.Organization, mapname string, item schema.Country) error {
	return ExecTrans(ctx, a.db, func(i context.Context) error {
		Country := entity.SchemaCountry(item).ToCountry()
		result := entity.GetCountryDB(i, a.db).
			Where("organization=?", organization).
			Where("mapname=?", mapname).
			Omit("organization", "mapname").Updates(Country)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Delete 删除数据
func (a *Country) Delete(ctx context.Context, typ uint) error {
	result := entity.GetCountryDB(ctx, a.db).
		Where("type=?", typ).
		Delete(entity.Country{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Country) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetCountryDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
