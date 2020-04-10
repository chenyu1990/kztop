package internal

import (
	"context"

	"kztop/internal/app/errors"
	"kztop/internal/app/model"
	"kztop/internal/app/schema"
)

// NewWpn 创建Wpn
func NewWpn(mWpn model.IWpn) *Wpn {
	return &Wpn{
		WpnModel: mWpn,
	}
}

// Wpn 示例程序
type Wpn struct {
	WpnModel model.IWpn
}

// Query 查询数据
func (a *Wpn) Query(ctx context.Context, params schema.WpnQueryParam, opts ...schema.WpnQueryOptions) (*schema.WpnQueryResult, error) {
	return a.WpnModel.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *Wpn) Get(ctx context.Context, recordID string, opts ...schema.WpnQueryOptions) (*schema.Wpn, error) {
	item, err := a.WpnModel.Get(ctx, recordID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}
