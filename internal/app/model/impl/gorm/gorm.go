package gorm

import (
	"time"

	"kztop/internal/app/config"
	"kztop/internal/app/model"
	"kztop/internal/app/model/impl/gorm/internal/entity"
	imodel "kztop/internal/app/model/impl/gorm/internal/model"
	"github.com/jinzhu/gorm"
	"go.uber.org/dig"

	// gorm存储注入
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Config 配置参数
type Config struct {
	Debug        bool
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

// NewDB 创建DB实例
func NewDB(c *Config) (*gorm.DB, error) {
	db, err := gorm.Open(c.DBType, c.DSN)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	return db, nil
}

// SetTablePrefix 设定表名前缀
func SetTablePrefix(prefix string) {
	entity.SetTablePrefix(prefix)
}

// AutoMigrate 自动映射数据表
func AutoMigrate(db *gorm.DB) error {
	if dbType := config.Global().Gorm.DBType; dbType == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate(
		new(entity.News),
		new(entity.Record),
		new(entity.WorldRecord),
		new(entity.Region),
		new(entity.Admin),
	).Error
}

// Inject 注入gorm实现
// 使用方式：
//   container := dig.New()
//   Inject(container)
//   container.Invoke(func(foo IDemo) {
//   })
func Inject(container *dig.Container) error {
	_ = container.Provide(imodel.NewTrans)
	_ = container.Provide(func(m *imodel.Trans) model.ITrans { return m })
	_ = container.Provide(imodel.NewDemo)
	_ = container.Provide(func(m *imodel.Demo) model.IDemo { return m })
	_ = container.Provide(imodel.NewMenu)
	_ = container.Provide(func(m *imodel.Menu) model.IMenu { return m })
	_ = container.Provide(imodel.NewRole)
	_ = container.Provide(func(m *imodel.Role) model.IRole { return m })
	_ = container.Provide(imodel.NewRecord)
	_ = container.Provide(func(m *imodel.Record) model.IRecord { return m })
	_ = container.Provide(imodel.NewNews)
	_ = container.Provide(func(m *imodel.News) model.INews { return m })
	_ = container.Provide(imodel.NewRegion)
	_ = container.Provide(func(m *imodel.Region) model.IRegion { return m })
	_ = container.Provide(imodel.NewAdmin)
	_ = container.Provide(func(m *imodel.Admin) model.IAdmin { return m })
	_ = container.Provide(imodel.NewWorldRecord)
	_ = container.Provide(func(m *imodel.WorldRecord) model.IWorldRecord { return m })
	return nil
}
