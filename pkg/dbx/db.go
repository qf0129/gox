package dbx

import (
	"fmt"
	"log/slog"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB     *gorm.DB
	Option *DBOption
)

type MysqlConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type SqliteConfig struct {
	DBFile string
}

type DBOption struct {
	Mysql           *MysqlConfig
	Sqlite          *SqliteConfig
	Gorm            *gorm.Config
	MigrateModels   []any
	ModelPrimaryKey string
	DefaultPageSize int
}

const (
	DefaultMysqlHost       = "localhost"
	DefaultMysqlPort       = 3306
	DefaultMysqlUsername   = "root"
	DefaultMysqlPassword   = "root"
	DefaultSqliteFile      = "sqlite.db"
	DefaultModelPrimaryKey = "id"
	DefaultPageSize        = 10
)

func Connect(opts ...*DBOption) {
	Option = loadDefaultDbOption(opts)
	var dbConn gorm.Dialector
	if Option.Mysql != nil {
		dbConn = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s",
			Option.Mysql.Username, Option.Mysql.Password,
			Option.Mysql.Host, Option.Mysql.Port,
			Option.Mysql.Database,
		))
		slog.Info("### Connected MySQL", slog.String("host", Option.Mysql.Host), slog.Int("port", Option.Mysql.Port), slog.String("user", Option.Mysql.Username), slog.String("db", Option.Mysql.Database))
	} else {
		dbConn = sqlite.Open(Option.Sqlite.DBFile)
		slog.Info("### Connected SQLite", slog.String("db", Option.Sqlite.DBFile))
	}
	var err error
	DB, err = gorm.Open(dbConn, Option.Gorm)
	if err != nil {
		panic("ConnectDatabaseFailed: " + err.Error())
	}
	if Option.MigrateModels != nil {
		if err := DB.AutoMigrate(Option.MigrateModels...); err != nil {
			panic("MigrateModelsErr:" + err.Error())
		}
	}
}

func loadDefaultDbOption(opts []*DBOption) *DBOption {
	var opt *DBOption
	if len(opts) == 0 {
		opt = &DBOption{}
	} else {
		opt = opts[0]
	}
	if opt.Mysql != nil {
		if opt.Mysql.Host == "" {
			opt.Mysql.Host = DefaultMysqlHost
		}
		if opt.Mysql.Port == 0 {
			opt.Mysql.Port = DefaultMysqlPort
		}
		if opt.Mysql.Username == "" {
			opt.Mysql.Username = DefaultMysqlUsername
		}
		if opt.Mysql.Password == "" {
			opt.Mysql.Password = DefaultMysqlPassword
		}
	} else {
		if opt.Sqlite == nil {
			opt.Sqlite = &SqliteConfig{}
		}
		if opt.Sqlite.DBFile == "" {
			opt.Sqlite.DBFile = DefaultSqliteFile
		}
	}
	if opt.Gorm == nil {
		opt.Gorm = &gorm.Config{}
	}
	if opt.Gorm.NamingStrategy == nil {
		opt.Gorm.NamingStrategy = schema.NamingStrategy{
			SingularTable: true,
		}
	}
	if opt.ModelPrimaryKey == "" {
		opt.ModelPrimaryKey = DefaultModelPrimaryKey
	}
	if opt.DefaultPageSize == 0 {
		opt.DefaultPageSize = DefaultPageSize
	}
	return opt
}
