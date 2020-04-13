package model

import (
	"context"
	"kztop/pkg/kreedz"

	"kztop/internal/app/schema"
)

// IRegion Region存储接口
type IRegion interface {
	// 查询数据
	Query(ctx context.Context, params schema.RegionQueryParam, opts ...schema.RegionQueryOptions) (*schema.RegionQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, params schema.RegionQueryParam) (*schema.Region, error)
	// 创建数据
	Create(ctx context.Context, item schema.Region) error
	// 更新数据
	Update(ctx context.Context, organization kreedz.Organization, mapname string, item schema.Region) error
	// 删除数据
	Delete(ctx context.Context, typ uint) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
