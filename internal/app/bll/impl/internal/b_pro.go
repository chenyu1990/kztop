package internal

import (
	"context"

	"kztop/internal/app/errors"
	"kztop/internal/app/model"
	"kztop/internal/app/schema"
)

// NewPro 创建Pro
func NewPro(mPro model.IPro) *Pro {
	return &Pro{
		ProModel: mPro,
	}
}

// Pro 示例程序
type Pro struct {
	ProModel model.IPro
}

// Query 查询数据
func (a *Pro) Query(ctx context.Context, params schema.ProQueryParam, opts ...schema.ProQueryOptions) (*schema.ProQueryResult, error) {
	return a.ProModel.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *Pro) Get(ctx context.Context, recordID string, opts ...schema.ProQueryOptions) (*schema.Pro, error) {
	item, err := a.ProModel.Get(ctx, recordID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}
