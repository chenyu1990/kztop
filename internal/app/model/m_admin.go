package model

import (
	"context"
	"kztop/pkg/kreedz"

	"kztop/internal/app/schema"
)

// IAdmin Admin存储接口
type IAdmin interface {
	// 查询数据
	Query(ctx context.Context, params schema.AdminQueryParam, opts ...schema.AdminQueryOptions) (*schema.AdminQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, params schema.AdminQueryParam) (*schema.Admin, error)
	// 创建数据
	Create(ctx context.Context, item schema.Admin) error
	// 更新数据
	Update(ctx context.Context, organization kreedz.Organization, mapname string, item schema.Admin) error
	// 删除数据
	Delete(ctx context.Context, typ uint) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
