package model

import (
	"context"
	"kztop/pkg/kreedz"

	"kztop/internal/app/schema"
)

// ICountry Country存储接口
type ICountry interface {
	// 查询数据
	Query(ctx context.Context, params schema.CountryQueryParam, opts ...schema.CountryQueryOptions) (*schema.CountryQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, params schema.CountryQueryParam) (*schema.Country, error)
	// 创建数据
	Create(ctx context.Context, item schema.Country) error
	// 更新数据
	Update(ctx context.Context, organization kreedz.Organization, mapname string, item schema.Country) error
	// 删除数据
	Delete(ctx context.Context, typ uint) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
