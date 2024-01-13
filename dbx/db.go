package dbx

import (
	"fmt"
	"log/slog"

	"github.com/glebarez/sqlite"
	"github.com/qf0129/gox/confx"
	"github.com/qf0129/gox/constx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB  *gorm.DB
	Opt *Option
)

type Option struct {
	MysqlConfig  *confx.Mysql
	SqliteConfig *confx.Sqlite
	Models       []any

	GormConfig      *gorm.Config
	ModelPrimaryKey string
	QueryPageSize   int
}

func initOption(opt *Option) {
	if opt.ModelPrimaryKey == "" {
		opt.ModelPrimaryKey = constx.DefaultModelPrimaryKey
	}
	if opt.QueryPageSize < 1 {
		opt.QueryPageSize = constx.DefaultQueryPageSize
	}
	if opt.GormConfig == nil {
		opt.GormConfig = &gorm.Config{}
	}
	if opt.GormConfig.NamingStrategy == nil {
		opt.GormConfig.NamingStrategy = schema.NamingStrategy{SingularTable: true}
	}
	if opt.GormConfig.Logger == nil {
		opt.GormConfig.Logger = logger.Default.LogMode(logger.Error)
	}
}

func Connect(opt *Option) *gorm.DB {
	initOption(opt)

	var dbConn gorm.Dialector
	if opt.MysqlConfig != nil && opt.MysqlConfig.Host != "" {
		dbConn = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			opt.MysqlConfig.Username,
			opt.MysqlConfig.Password,
			opt.MysqlConfig.Host,
			opt.MysqlConfig.Port,
			opt.MysqlConfig.Database,
		))
		slog.Info("### Connected MySQL", slog.String("host", opt.MysqlConfig.Host), slog.String("user", opt.MysqlConfig.Username))
	} else {
		if opt.SqliteConfig == nil {
			opt.SqliteConfig = &confx.Sqlite{}
		}
		if opt.SqliteConfig.FilePath == "" {
			opt.SqliteConfig.FilePath = constx.DefaultSqliteFile
		}
		dbConn = sqlite.Open(opt.SqliteConfig.FilePath)
		slog.Info("### Connected SQLite", slog.String("path", opt.SqliteConfig.FilePath))
	}

	var err error
	DB, err = gorm.Open(dbConn, opt.GormConfig)
	if err != nil {
		panic("ConnectDatabaseFailed: " + err.Error())
	}

	if err := DB.AutoMigrate(opt.Models...); err != nil {
		panic("MigrateModelErr:" + err.Error())
	}

	Opt = opt
	return DB
}
