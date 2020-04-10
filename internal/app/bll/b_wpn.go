package bll

import (
	"context"

	"kztop/internal/app/schema"
)

// IWpn Wpn业务逻辑接口
type IWpn interface {
	// 查询数据
	Query(ctx context.Context, params schema.WpnQueryParam, opts ...schema.WpnQueryOptions) (*schema.WpnQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.WpnQueryOptions) (*schema.Wpn, error)
}
