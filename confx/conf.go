package confx

import (
	"log/slog"
	"os"

	"github.com/qf0129/gox/constx"
	"github.com/qf0129/gox/jsonx"
)

type Mysql struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type Sqlite struct {
	FilePath string
}

type Server struct {
	ListenAddr           string
	GinMode              string // debug,release,test
	LogLevel             string // panic,fatal,error,warn,info,debug,trace
	DBLogLevel           string // error,warn,info
	ReadTimeout          int64
	WriteTimeout         int64
	EncryptSecret        string
	CookieDomain         string
	CookieExpiredSeconds int
}

var (
	DefaultMysql  = Mysql{}
	DefaultSqlite = Sqlite{
		FilePath: constx.DefaultSqliteFile,
	}
	DefaultServer = Server{
		ListenAddr:           constx.DefaultListenAddr,
		GinMode:              constx.DefaultGinMode,
		LogLevel:             constx.DefaultLogLevel,
		DBLogLevel:           constx.DefaultDBLogLevel,
		ReadTimeout:          constx.DefaultReadTimeout,
		WriteTimeout:         constx.DefaultWriteTimeout,
		EncryptSecret:        constx.DefaultEncryptSecret,
		CookieExpiredSeconds: constx.DefaultCookieExpiredSeconds,
	}
)

type BaseConfig struct {
	Mysql  *Mysql
	Sqlite *Sqlite
	Server *Server
}

func ReadJsonConfig(target any) {
	data, err := os.ReadFile("conf.json")
	if err != nil {
		slog.Warn("Read config file failed, running with default config.")
		return
	}
	err = jsonx.Unmarshal(data, target)
	if err != nil {
		panic("UnmarshalConfigFailed: " + err.Error())
	}
}
