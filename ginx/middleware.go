package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/confx"
	"github.com/qf0129/gox/constx"
	"github.com/qf0129/gox/dbx"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/respx"
	"github.com/qf0129/gox/securex"
	"github.com/rs/xid"
)

func ReqIdMiddileWare(c *gin.Context) {
	c.Set(constx.KeyOfRequestId, xid.New().String())
}

// 验证token
func CheckTokenMiddileWare[T any](conf *confx.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		tk, err := c.Cookie(constx.KeyOfCookieToken)
		if err != nil {
			respx.ErrAuth(c, errx.InvalidToken)
			return
		}

		uid1, err := c.Cookie(constx.KeyOfCookieUserId)
		if err != nil {
			respx.ErrAuth(c, errx.InvalidToken)
			return
		}

		uid2, err := securex.ParseToken(tk, conf.EncryptSecret, int64(conf.CookieExpiredSeconds))
		if err != nil {
			respx.ErrAuth(c, errx.InvalidToken.AddErr(err))
			return
		}

		if uid1 != uid2 {
			respx.ErrAuth(c, errx.InvalidToken.AddErr(err))
			return
		}

		existsUser, err := dbx.QueryOneByPk[T](uid2)
		if err != nil {
			respx.ErrAuth(c, errx.UserNotFound.AddErr(err))
			return
		}
		SetRequestUser(c, existsUser, uid2)
		c.Next()
	}
}

// 跨域请求
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
