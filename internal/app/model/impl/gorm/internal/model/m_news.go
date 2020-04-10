package model

import (
	"context"

	"kztop/internal/app/errors"
	"kztop/internal/app/model/impl/gorm/internal/entity"
	"kztop/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// NewNews 创建News存储实例
func NewNews(db *gorm.DB) *News {
	return &News{db}
}

// News News存储
type News struct {
	db *gorm.DB
}

func (a *News) getQueryOption(opts ...schema.NewsQueryOptions) schema.NewsQueryOptions {
	var opt schema.NewsQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *News) where(db *gorm.DB, params *schema.NewsQueryParam) *gorm.DB {
	if v := params.Organization; v > 0 {
		db = db.Where("organization=?", v)
	}
	if v := params.Period; v > 0 {
		db = db.Where("period=?", v)
	}
	return db
}

// Query 查询数据
func (a *News) Query(ctx context.Context, params schema.NewsQueryParam, opts ...schema.NewsQueryOptions) (*schema.NewsQueryResult, error) {
	db := entity.GetNewsDB(ctx, a.db)
	db = a.where(db, &params)

	opt := a.getQueryOption(opts...)
	var list entity.Newss
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.NewsQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaNewss(),
	}

	return qr, nil
}

func (a *News) Select(ctx context.Context, selectString string, params *schema.NewsQueryParam, sumResult interface{}) (bool, error) {
	db := entity.GetNewsDB(ctx, a.db)
	if params != nil {
		db = a.where(db, params)
	}
	result := db.Select(selectString).Scan(sumResult)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Get 查询指定数据
func (a *News) Get(ctx context.Context, recordID string, opts ...schema.NewsQueryOptions) (*schema.News, error) {
	db := entity.GetNewsDB(ctx, a.db).Where("record_id=?", recordID)
	var item entity.News
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaNews(), nil
}

// Create 创建数据
func (a *News) Create(ctx context.Context, item schema.News) error {
	return ExecTrans(ctx, a.db, func(i context.Context) error {
		News := entity.SchemaNews(item).ToNews()
		result := entity.GetNewsDB(i, a.db).Create(News)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Update 更新数据
func (a *News) Update(ctx context.Context, recordID string, item schema.News) error {
	return ExecTrans(ctx, a.db, func(i context.Context) error {
		News := entity.SchemaNews(item).ToNews()
		result := entity.GetNewsDB(i, a.db).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(News)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Delete 删除数据
func (a *News) Delete(ctx context.Context, recordID string) error {
	result := entity.GetNewsDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.News{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *News) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetNewsDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
