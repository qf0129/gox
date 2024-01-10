package confx

import (
	"encoding/json"
	"log/slog"
	"os"
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
	ListenAddr         string
	Domain             string
	GinMode            string // debug,release,test
	LogLevel           string // panic,fatal,error,warn,info,debug,trace
	DBLogLevel         string // error,warn,info
	ReadTimeout        int64
	WriteTimeout       int64
	TokenExpiredSecond int
	EncryptSecret      string
}

var (
	DefaultMysql  = Mysql{}
	DefaultSqlite = Sqlite{
		FilePath: "db.sqlite",
	}
	DefaultServer = Server{
		ListenAddr:         ":8888",
		Domain:             "",
		GinMode:            "debug",
		LogLevel:           "debug",
		DBLogLevel:         "warn",
		ReadTimeout:        60,
		WriteTimeout:       60,
		TokenExpiredSecond: 3600,
		EncryptSecret:      "DEFAULT_SECRET",
	}
)

type BaseConfig struct {
	Mysql  *Mysql
	Sqlite *Sqlite
	Server *Server
}

func ReadJsonConfig(target any) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		slog.Warn("Read config file failed, running with default config.")
		return
	}
	err = json.Unmarshal(data, target)
	if err != nil {
		panic("UnmarshalConfigFailed: " + err.Error())
	}
}
