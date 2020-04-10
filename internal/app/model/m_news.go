package model

import (
	"context"

	"kztop/internal/app/schema"
)

// INews News存储接口
type INews interface {
	// 查询数据
	Query(ctx context.Context, params schema.NewsQueryParam, opts ...schema.NewsQueryOptions) (*schema.NewsQueryResult, error)
	Select(ctx context.Context, selectString string, params *schema.NewsQueryParam, sumResult interface{}) (bool, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.NewsQueryOptions) (*schema.News, error)
	// 创建数据
	Create(ctx context.Context, item schema.News) error
	// 更新数据
	Update(ctx context.Context, recordID string, item schema.News) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
