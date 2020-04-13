package model

import (
	"context"

	"kztop/internal/app/schema"
)

// IRecord Record存储接口
type IRecord interface {
	// 查询数据
	Query(ctx context.Context, params *schema.RecordQueryParam, opts ...schema.RecordQueryOptions) (*schema.RecordQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, recordID string, opts ...schema.RecordQueryOptions) (*schema.Record, error)
	// 创建数据
	Create(ctx context.Context, item *schema.Record) error
	// 更新数据
	Update(ctx context.Context, params *schema.RecordQueryParam, item schema.Record) error
	UpdateInfo(ctx context.Context, info schema.UpdateInfo) error
	// 删除数据
	Delete(ctx context.Context, recordID string) error
	// 更新状态
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
