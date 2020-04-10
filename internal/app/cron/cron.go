package cron

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/dig"
	"kztop/internal/app/config"
	"kztop/internal/app/cron/job"
)

func Init(container *dig.Container) {
	err := Inject(container)
	if err != nil {
		return
	}

	_ = container.Invoke(func(
		jKreedz *job.KreedzJob,
	) error {
		cfg := config.Global()
		var c *cron.Cron
		if cfg.RunMode == "debug" {
			c = cron.New(cron.WithSeconds(), cron.WithChain(
				cron.Recover(cron.DefaultLogger),
			))
			// Seconds | Minutes | Hours | Day of month | Month | Day of week
			//_, err = c.AddFunc("*/1 * * * * *", jKreedz.UpdateWorldRecord)
		} else {
			c = cron.New(cron.WithChain(
				cron.Recover(cron.DefaultLogger),
			))
			// Minutes | Hours | Day of month | Month | Day of week
			_, err = c.AddFunc("* */1 * * *", jKreedz.UpdateWorldRecord)
		}

		if c == nil {
			panic("计划任务出错")
		}

		c.Start()

		return nil
	})
}

func Inject(container *dig.Container) error {
	_ = container.Provide(job.NewKreedzJob)
	return nil
}