package internal

import (
	"context"

	"kztop/internal/app/errors"
	"kztop/internal/app/model"
	"kztop/internal/app/schema"
)

// NewNub 创建Nub
func NewNub(mNub model.INub) *Nub {
	return &Nub{
		NubModel: mNub,
	}
}

// Nub 示例程序
type Nub struct {
	NubModel model.INub
}

// Query 查询数据
func (a *Nub) Query(ctx context.Context, params schema.NubQueryParam, opts ...schema.NubQueryOptions) (*schema.NubQueryResult, error) {
	return a.NubModel.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *Nub) Get(ctx context.Context, recordID string, opts ...schema.NubQueryOptions) (*schema.Nub, error) {
	item, err := a.NubModel.Get(ctx, recordID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}