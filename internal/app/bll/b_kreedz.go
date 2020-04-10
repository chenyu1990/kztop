package bll

import (
	"context"
	"kztop/pkg/kreedz"
	"time"
)

// IKreedz Kreedz业务逻辑接口
type IKreedz interface {
	RecordUpdate(ctx context.Context, organization kreedz.Organization, records []*kreedz.RecordInfo) error
	CreateRecord(ctx context.Context, organization kreedz.Organization, records []*kreedz.RecordInfo) error
	CreateNews(ctx context.Context, date time.Time, organization kreedz.Organization, news map[string]map[string][]*kreedz.RecordInfo) error
}
