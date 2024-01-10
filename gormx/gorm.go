package gormx

import (
	"fmt"
	"log/slog"

	"github.com/qf0129/gox/confx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
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

var DefaultGormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
	Logger: logger.Default.LogMode(logger.Warn),
}

func Connect(opt *Option) *gorm.DB {
	if opt.ModelPrimaryKey == "" {
		opt.ModelPrimaryKey = "id"
	}
	if opt.QueryPageSize < 1 {
		opt.QueryPageSize = 10
	}
	if opt.GormConfig == nil {
		opt.GormConfig = DefaultGormConfig
	}
	Opt = opt

	var dbConn gorm.Dialector
	if opt.MysqlConfig.Host != "" {
		dbConn = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			opt.MysqlConfig.Username,
			opt.MysqlConfig.Password,
			opt.MysqlConfig.Host,
			opt.MysqlConfig.Port,
			opt.MysqlConfig.Database,
		))
		slog.Info(fmt.Sprintf("Connected DB with mysql:%s@%s", opt.MysqlConfig.Username, opt.MysqlConfig.Host))
	} else {
		if opt.SqliteConfig.FilePath == "" {
			opt.SqliteConfig.FilePath = "db.sqlite"
		}
		dbConn = sqlite.Open(opt.SqliteConfig.FilePath)
		slog.Info(fmt.Sprintf("Connected DB with sqlite:%s", opt.SqliteConfig.FilePath))
	}

	var err error
	DB, err = gorm.Open(dbConn, opt.GormConfig)
	if err != nil {
		panic("ConnectDatabaseFailed: " + err.Error())
	}

	if err := DB.AutoMigrate(opt.Models...); err != nil {
		panic("MigrateModelErr:" + err.Error())
	}
	return DB
}
