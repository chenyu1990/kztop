package bll

import (
	"context"

	"kztop/internal/app/schema"
)

// IPro Pro业务逻辑接口
type IPro interface {
	// 查询数据
	Query(ctx context.Context, params schema.ProQueryParam, opts ...schema.ProQueryOptions) (*schema.ProQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.ProQueryOptions) (*schema.Pro, error)
}
