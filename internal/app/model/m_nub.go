package model

import (
	"context"

	"kztop/internal/app/schema"
)

// INub Nub存储接口
type INub interface {
	// 查询数据
	Query(ctx context.Context, params schema.NubQueryParam, opts ...schema.NubQueryOptions) (*schema.NubQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.NubQueryOptions) (*schema.Nub, error)
	// 创建数据
	Create(ctx context.Context, item schema.Nub) error
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.Nub) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
