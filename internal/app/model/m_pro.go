package model

import (
	"context"

	"kztop/internal/app/schema"
)

// IPro Pro存储接口
type IPro interface {
	// 查询数据
	Query(ctx context.Context, params schema.ProQueryParam, opts ...schema.ProQueryOptions) (*schema.ProQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.ProQueryOptions) (*schema.Pro, error)
	// 创建数据
	Create(ctx context.Context, item schema.Pro) error
	// 更新数据
	Update(ctx context.Context, params *schema.ProQueryParam, item schema.Pro) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
