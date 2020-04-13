package schema

import "kztop/pkg/kreedz"

// WorldRecord WorldRecord对象
type WorldRecord struct {
	MapName      string              `json:"mapname"`
	Holder       string              `json:"holder"`
	Region       string              `json:"region"`
	Time         float64             `json:"time"`
	Organization kreedz.Organization `json:"organization"`
}

// WorldRecordQueryParam 查询条件
type WorldRecordQueryParam struct {
	MapName       string                `form:"mapname"`
	Holder        string                `form:"holder"`
	Region        string                `form:"region"`
	Organization  kreedz.Organization   `form:"organization"`
	Organizations []kreedz.Organization `-`
}

// WorldRecordQueryOptions WorldRecord对象查询可选参数项
type WorldRecordQueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// WorldRecordQueryResult WorldRecord对象查询结果
type WorldRecordQueryResult struct {
	Data       []*WorldRecord
	PageResult *PaginationResult
}
