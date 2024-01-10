package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox"
)

func SetRequestUser(c *gin.Context, user any) {
	c.Set(gox.KeyOfRequestUser, user)
}

func GetRequestUser[T any](c *gin.Context) T {
	user, _ := c.Get(gox.KeyOfRequestUser)
	return user.(T)
}
