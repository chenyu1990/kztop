package model

import (
	"context"
	"kztop/pkg/kreedz"

	"kztop/internal/app/errors"
	"kztop/internal/app/model/impl/gorm/internal/entity"
	"kztop/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// NewAdmin 创建Admin存储实例
func NewAdmin(db *gorm.DB) *Admin {
	return &Admin{db}
}

// Admin Admin存储
type Admin struct {
	db *gorm.DB
}

func (a *Admin) getQueryOption(opts ...schema.AdminQueryOptions) schema.AdminQueryOptions {
	var opt schema.AdminQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *Admin) where(db *gorm.DB, params *schema.AdminQueryParam) *gorm.DB {
	if v := params.Server; v != "" {
		db = db.Where("server=?", v)
	}
	if v := params.SteamID; v != "" {
		db = db.Where("steamid=?", v)
	}
	if v := params.Valid; v != nil {
		db = db.Where("valid=?", v)
	}

	return db
}

// Query 查询数据
func (a *Admin) Query(ctx context.Context, params schema.AdminQueryParam, opts ...schema.AdminQueryOptions) (*schema.AdminQueryResult, error) {
	db := entity.GetAdminDB(ctx, a.db)
	db = a.where(db, &params)

	opt := a.getQueryOption(opts...)
	var list entity.Admins
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.AdminQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaAdmins(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Admin) Get(ctx context.Context, params schema.AdminQueryParam) (*schema.Admin, error) {
	db := entity.GetAdminDB(ctx, a.db)
	db = a.where(db, &params)

	var item entity.Admin
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaAdmin(), nil
}

// Create 创建数据
func (a *Admin) Create(ctx context.Context, item schema.Admin) error {
	return ExecTrans(ctx, a.db, func(i context.Context) error {
		Admin := entity.SchemaAdmin(item).ToAdmin()
		result := entity.GetAdminDB(i, a.db).Create(Admin)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Update 更新数据
func (a *Admin) Update(ctx context.Context, organization kreedz.Organization, mapname string, item schema.Admin) error {
	return ExecTrans(ctx, a.db, func(i context.Context) error {
		Admin := entity.SchemaAdmin(item).ToAdmin()
		result := entity.GetAdminDB(i, a.db).
			Where("organization=?", organization).
			Where("mapname=?", mapname).
			Omit("organization", "mapname").Updates(Admin)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

// Delete 删除数据
func (a *Admin) Delete(ctx context.Context, typ uint) error {
	result := entity.GetAdminDB(ctx, a.db).
		Where("type=?", typ).
		Delete(entity.Admin{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Admin) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetAdminDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
