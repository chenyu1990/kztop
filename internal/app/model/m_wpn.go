package model

import (
	"context"

	"kztop/internal/app/schema"
)

// IWpn Wpn存储接口
type IWpn interface {
	// 查询数据
	Query(ctx context.Context, params schema.WpnQueryParam, opts ...schema.WpnQueryOptions) (*schema.WpnQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.WpnQueryOptions) (*schema.Wpn, error)
	// 创建数据
	Create(ctx context.Context, item schema.Wpn) error
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.Wpn) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
