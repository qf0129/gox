package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/constx"
)

func SetRequestUser(c *gin.Context, user any, userId string) {
	c.Set(constx.KeyOfRequestUser, user)
	c.Set(constx.KeyOfRequestUserId, userId)
}

func GetRequestUserId(c *gin.Context) string {
	return c.GetString(constx.KeyOfRequestUserId)
}

func GetRequestUser[T any](c *gin.Context) T {
	user, _ := c.Get(constx.KeyOfRequestUser)
	return user.(T)
}
