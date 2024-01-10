package midx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/ginx"
	"github.com/qf0129/gox/gormx/daox"
	"github.com/qf0129/gox/respx"
	"github.com/qf0129/gox/securex"
)

// 验证token
func CheckToken[T any](secretKey string, tokenExpiredSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tk, err := c.Cookie(gox.KeyOfCookieToken)
		if err != nil {
			respx.Err(c, errx.InvalidToken)
			return
		}

		uid, err := securex.ParseToken(tk, secretKey, int64(tokenExpiredSeconds))
		if err != nil {
			respx.Err(c, errx.InvalidToken.AddErr(err))
			return
		}

		existsUser, err := daox.QueryOneByPk[T](uid)
		if err != nil {
			respx.Err(c, errx.UserNotFound.AddErr(err))
			return
		}
		ginx.SetRequestUser(c, existsUser)
		c.Next()
	}
}
