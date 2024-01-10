package crudx

import "github.com/gin-gonic/gin"

func CreateCrudRoute[T any](group *gin.RouterGroup, modelName string) {
	group.GET("/"+modelName, QueryManyHandler[T]())
	group.POST("/"+modelName, CreateModelHandler[T]())
	group.PUT("/"+modelName+"/:id", UpdateHandler[T]())
	group.DELETE("/"+modelName+"/:id", DeleteHandler[T]())
}
