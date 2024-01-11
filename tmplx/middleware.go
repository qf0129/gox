package tmplx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/confx"
	"github.com/qf0129/gox/constx"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/ginx"
	"github.com/qf0129/gox/gormx/daox"
	"github.com/qf0129/gox/respx"
	"github.com/qf0129/gox/securex"
)

// 验证token
func CheckTokenMiddleware(conf *confx.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		tk, err := c.Cookie(constx.KeyOfCookieToken)
		if tk == "" || err != nil {
			respx.Err(c, errx.InvalidToken)
			return
		}

		uid, err := securex.ParseToken(tk, conf.EncryptSecret, int64(conf.CookieExpiredSeconds))
		if err != nil {
			respx.Err(c, errx.InvalidToken.AddErr(err))
			return
		}

		existsUser, err := daox.QueryOneByPk[User](uid)
		if err != nil {
			respx.Err(c, errx.UserNotFound.AddErr(err))
			return
		}
		ginx.SetRequestUser(c, existsUser)
		c.Next()
	}
}
