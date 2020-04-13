package internal

import (
	"context"
	"encoding/json"
	"kztop/internal/app/schema"
	"kztop/pkg/kreedz"
	"time"

	"kztop/internal/app/model"
)

// NewKreedz 创建Kreedz
func NewKreedz(
	trans model.ITrans,
	mNews model.INews,
	mWorldRecord model.IWorldRecord,
) *Kreedz {
	return &Kreedz{
		TransModel:       trans,
		NewsModel:        mNews,
		WorldRecordModel: mWorldRecord,
	}
}

// Kreedz 示例程序
type Kreedz struct {
	TransModel       model.ITrans
	NewsModel        model.INews
	WorldRecordModel model.IWorldRecord
}

func (a *Kreedz) RecordUpdate(ctx context.Context, organization kreedz.Organization, records []*kreedz.RecordInfo) error {
	return ExecTrans(ctx, a.TransModel, func(i context.Context) error {
		for _, record := range records {
			exist, err := a.WorldRecordModel.Get(ctx, schema.WorldRecordQueryParam{
				MapName:      record.MapName,
				Organization: organization,
			})
			if err != nil {
				return nil
			}
			if exist != nil {
				a.WorldRecordModel.Update(ctx, organization, record.MapName, schema.WorldRecord{
					Holder: record.Holder,
					Region: record.Region,
					Time:   record.Time,
				})
			} else {
				a.WorldRecordModel.Create(ctx, schema.WorldRecord{
					Organization: organization,
					MapName:      record.MapName,
					Holder:       record.Holder,
					Region:       record.Region,
					Time:         record.Time,
				})
			}
		}

		return nil
	})
}

func (a *Kreedz) CreateRecord(ctx context.Context, organization kreedz.Organization, records []*kreedz.RecordInfo) error {
	for _, record := range records {
		if record == nil {
			continue
		}
		err := a.WorldRecordModel.Create(ctx, schema.WorldRecord{
			Organization: organization,
			MapName:      record.MapName,
			Holder:       record.Holder,
			Region:       record.Region,
			Time:         record.Time,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type mysqlResult struct {
	Period uint `json:"period"`
}

func (a *Kreedz) CreateNews(ctx context.Context, date time.Time, organization kreedz.Organization, news map[string]map[string][]*kreedz.RecordInfo) error {
	json, err := json.Marshal(news)
	if err != nil {
		panic(err)
	}

	result := mysqlResult{}
	_, err = a.NewsModel.Select(ctx, "max(period) as period", nil, &result)
	if err != nil {
		panic(err)
	}
	return a.NewsModel.Create(ctx, schema.News{
		Organization: organization,
		Period:       result.Period + 1, // TODO
		Data:         string(json),
		Date:         date,
	})
}
