package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"html"
	"html/template"
	"kztop/internal/app/middleware"
	"kztop/internal/app/routers/api/ctl"
	"kztop/pkg/auth"
	"kztop/pkg/kreedz"
)

// RegisterRouter 注册/api路由
func RegisterRouter(app *gin.Engine, container *dig.Container) error {
	err := ctl.Inject(container)
	if err != nil {
		return err
	}

	return container.Invoke(func(
		a auth.Auther,
		e *casbin.SyncedEnforcer,
		cServer *ctl.Server,
		cTop *ctl.Top,
		cRecord *ctl.Record,
	) error {
		// 请求频率限制中间件
		app.Use(middleware.RateLimiterMiddleware())
		app.SetFuncMap(template.FuncMap{
			"SteamIDToSteamID64": kreedz.SteamIDToSteamID64,
			"IsSteamID":          kreedz.IsSteamID,
			"Inc":                kreedz.Inc,
			"Mod":                kreedz.Mod,
			"IsSlowly":           kreedz.IsSlowly,
			"SubFloatRtnString":  kreedz.SubFloatRtnString,
			"FormatDate":         kreedz.FormatDate,
			"ToLower":            kreedz.ToLower,
			"SecondsToMinutes":   kreedz.SecondsToMinutes,
			"UnescapeString":     html.UnescapeString,
		})
		app.LoadHTMLGlob("templates/**/*")

		top := app.Group("/top")
		{
			top.GET("/:cate/:mapname", cTop.Top)
			top.PUT("/:cate", cTop.UpdateRecord)
		}

		player := app.Group("/player")
		{
			player.GET("/:player", cTop.Player)
			player.GET("/:player/:cate", cTop.Player)
		}

		players := app.Group("/players")
		{
			players.GET("", cTop.Players)
		}

		record := app.Group("/record")
		{
			record.GET("/org", cRecord.Org)
			record.GET("/org/:org", cRecord.Org)
			record.GET("/map/:mapname", cRecord.Map)
		}

		ips := app.Group("/ips")
		{
			ips.GET("", cServer.List)
		}

		return nil
	})
}
