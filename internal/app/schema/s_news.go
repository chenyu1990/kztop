package schema

import (
	"kztop/pkg/kreedz"
	"time"
)

// News News对象
type News struct {
	Organization kreedz.Organization `json:"organization"`
	Period       uint                `json:"period"`
	Data         string              `json:"data"`
	Date         time.Time           `json:"date"`
}

// NewsQueryParam 查询条件
type NewsQueryParam struct {
	Organization kreedz.Organization `form:"organization"`
	Period       uint                `form:"period"`
}

// NewsQueryOptions News对象查询可选参数项
type NewsQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// NewsQueryResult News对象查询结果
type NewsQueryResult struct {
	Data       []*News
	PageResult *PaginationResult
}
