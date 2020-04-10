package bll

import (
	"context"

	"kztop/internal/app/schema"
)

// INub Nub业务逻辑接口
type INub interface {
	// 查询数据
	Query(ctx context.Context, params schema.NubQueryParam, opts ...schema.NubQueryOptions) (*schema.NubQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.NubQueryOptions) (*schema.Nub, error)
}
