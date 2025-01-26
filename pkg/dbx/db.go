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

func ConnectDB(opt *DBOption) {
	loadDefaultDbOption(opt)
	Option = opt
	var dbConn gorm.Dialector
	if opt.Mysql != nil {
		dbConn = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			opt.Mysql.Username, opt.Mysql.Password,
			opt.Mysql.Host, opt.Mysql.Port,
			opt.Mysql.Database,
		))
		slog.Info("### Connected MySQL", slog.String("host", opt.Mysql.Host), slog.Int("port", opt.Mysql.Port), slog.String("user", opt.Mysql.Username), slog.String("db", opt.Mysql.Database))
	} else {
		dbConn = sqlite.Open(opt.Sqlite.DBFile)
		slog.Info("### Connected SQLite", slog.String("db", opt.Sqlite.DBFile))
	}
	var err error
	DB, err = gorm.Open(dbConn, opt.Gorm)
	if err != nil {
		panic("ConnectDatabaseFailed: " + err.Error())
	}
	if opt.MigrateModels != nil {
		if err := DB.AutoMigrate(opt.MigrateModels...); err != nil {
			panic("MigrateModelsErr:" + err.Error())
		}
	}
}

func loadDefaultDbOption(opt *DBOption) {
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
}
