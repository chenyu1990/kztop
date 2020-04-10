package job

import (
	"context"
	"kztop/internal/app/config"
	"sync"

	//"context"
	"kztop/internal/app/bll"
	"kztop/pkg/kreedz"
)

func NewKreedzJob(
	bKreedz bll.IKreedz,
) *KreedzJob {
	return &KreedzJob{
		KreedzBll: bKreedz,
	}
}

type KreedzJob struct {
	KreedzBll bll.IKreedz
	updateWorldRecordMutex sync.Mutex
}

func (a *KreedzJob) UpdateWorldRecord() {
	a.updateWorldRecordMutex.Lock()
	defer a.updateWorldRecordMutex.Unlock()

	ctx := context.Background()

	wr := kreedz.WorldRecord{}
	wr.Context = ctx

	cfg := config.Global()
	var organizations []kreedz.Organization
	if cfg.RunMode == "debug" {
		organizations = []kreedz.Organization{kreedz.DebugWorldRecord}
	} else {
		organizations = []kreedz.Organization{kreedz.XtremeJumps, kreedz.CosyClimbing}
	}

	for _, organization := range organizations {
		wr.Organization = organization
		first, records := wr.FirstSync()
		if first == true {
			err := a.KreedzBll.CreateRecord(ctx, wr.Organization, records)
			if err != nil {
				panic(err)
			}
		} else {
			if wr.CheckUpdate(&organization) == true {
				err := a.KreedzBll.RecordUpdate(ctx, wr.Organization, wr.NewRecords)
				if err != nil {
					panic(err)
				}
				err = a.KreedzBll.CreateNews(ctx, wr.NewsDate, wr.Organization, wr.News)
				if err != nil {
					panic(err)
				}
				wr.CopyFile()
			}
		}
	}
}
