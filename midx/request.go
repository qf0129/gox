package midx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/constx"
	"github.com/rs/xid"
)

func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constx.KeyOfRequestId, xid.New().String())
	}
}
