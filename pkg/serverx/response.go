package serverx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/errx"
)

type RespBody struct {
	ReqId string `json:",omitempty"`
	Code  int
	Msg   string
	Data  interface{}
}

func Response(c *gin.Context, httpCode int, code int, msg string, data interface{}, reqId string) {
	c.JSON(httpCode, &RespBody{
		ReqId: reqId,
		Code:  code,
		Msg:   msg,
		Data:  data,
	})
	c.Abort()
}

func ResponseOK(c *gin.Context, data interface{}) {
	Response(c, http.StatusOK, 0, "", data, "")
}

func ResponseErr(c *gin.Context, err errx.Err) {
	Response(c, http.StatusOK, err.Code(), err.Msg(), nil, "")
}
