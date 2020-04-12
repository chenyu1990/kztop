package data

import (
	"context"
	"go.uber.org/dig"
	"kztop/internal/app/bll"
	"kztop/internal/app/model"
	"kztop/internal/app/schema"
	"kztop/pkg/util"
)

// initAdmin 初始化合约数据
func InitAdminData(ctx context.Context, container *dig.Container) error {
	return container.Invoke(func(trans bll.ITrans, admin model.IAdmin) error {
		// 检查是否存在数据，如果不存在则初始化
		result, err := admin.Query(ctx, schema.AdminQueryParam{}, schema.AdminQueryOptions{
			PageParam: &schema.PaginationParam{PageIndex: 1, PageSize: 1},
		})
		if err != nil {
			return err
		} else if len(result.Data) > 0 {
			return nil
		}

		var data schema.Admins
		err = util.JSONUnmarshal([]byte(adminData), &data)
		if err != nil {
			return err
		}

		return createAdmin(ctx, trans, admin, data)
	})
}

func createAdmin(ctx context.Context, trans bll.ITrans, admin model.IAdmin, list schema.Admins) error {
	return trans.Exec(ctx, func(ctx context.Context) error {
		for _, item := range list {
			sItem := schema.Admin{
				Server:  item.Server,
				SteamID: item.SteamID,
				Access:  item.Access,
				Valid:   item.Valid,
			}
			err := admin.Create(ctx, sItem)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

const adminData = `
[
	{
		"server": "ALL",
		"steamid": "STEAM_0:0:33403241",
		"access": "abcdefghijklmnopqrstu",
		"valid": true
	}
]
`