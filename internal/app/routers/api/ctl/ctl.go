package ctl

import (
	"go.uber.org/dig"
)

// Inject 注入ctl
// 使用方式：
//   container := dig.New()
//   Inject(container)
//   container.Invoke(func(foo *ctl.Demo) {
//   })
func Inject(container *dig.Container) error {
	_ = container.Provide(NewServer)
	_ = container.Provide(NewTop)
	_ = container.Provide(NewRecord)
	return nil
}
