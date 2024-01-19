package respx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/constx"
	"github.com/qf0129/gox/errx"
)

type RespBody struct {
	ReqId string `json:"req_id"`
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Data  any    `json:"data"`
}

func returnResp(c *gin.Context, httpCode int, code int, msg string, data any) {
	c.JSON(httpCode, &RespBody{
		ReqId: c.GetString(constx.KeyOfRequestId),
		Code:  code,
		Msg:   msg,
		Data:  data,
	})
	c.Abort()
}

func OK(c *gin.Context, data any) {
	returnResp(c, http.StatusOK, 0, "ok", data)
}

func Err(c *gin.Context, err *errx.Err) {
	returnResp(c, http.StatusOK, err.Code, err.Msg, nil)
}

func ErrRequest(c *gin.Context, msg string) {
	Err(c, errx.RequestFailed.AddMsg(msg))
}

func ErrAuth(c *gin.Context, err *errx.Err) {
	returnResp(c, http.StatusUnauthorized, err.Code, err.Msg, nil)
}
