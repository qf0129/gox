package respx

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

const RequestIdKey = "REQID"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(RequestIdKey, xid.New().String())
	}
}
