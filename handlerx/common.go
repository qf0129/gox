package handlerx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/respx"
)

func HealthHandler(c *gin.Context) {
	respx.OK(c, "ok")
}
