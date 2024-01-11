package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/constx"
)

func SetRequestUser(c *gin.Context, user any) {
	c.Set(constx.KeyOfRequestUser, user)
}

func GetRequestUser[T any](c *gin.Context) T {
	user, _ := c.Get(constx.KeyOfRequestUser)
	return user.(T)
}
