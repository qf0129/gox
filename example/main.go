package main

import (
	"net/http"

	"github.com/qf0129/gox/pkg/crudx"
	"github.com/qf0129/gox/pkg/dbx"
	"github.com/qf0129/gox/pkg/serverx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	Id   int64  `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(255);not null"`
	Age  int    `gorm:"type:int;not null"`
}

func main() {
	dbx.Connect(&dbx.DBOption{
		Sqlite: &dbx.SqliteConfig{
			DBFile: "db.sqlite",
		},
		MigrateModels: []any{&User{}},
		Gorm: &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	})
	app := serverx.NewApp()
	crudModels := map[string]crudx.CrudModel{
		"user": {Model: &User{}, Methods: "crud"},
	}
	app.AddApi(
		&serverx.ApiInfo{Method: http.MethodPost, Path: "crud", Handler: crudx.CrudHandler(crudModels)},
	)
	app.Run()
}
