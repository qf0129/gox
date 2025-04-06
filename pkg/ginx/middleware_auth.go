package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/dbx"
	"github.com/qf0129/gox/pkg/errx"
	"github.com/qf0129/gox/pkg/securex"
)

const (
	KeyOfRequestUser         = "ctx_reqUser"
	KeyOfRequestUserId       = "ctx_reqUserId"
	KeyOfCookieToken         = "t"
	KeyOfCookieUserId        = "u"
	KeyOfHeaderAuthorization = "Authorization"
)

func SetRequestUser(c *gin.Context, user any, userId string) {
	c.Set(KeyOfRequestUser, user)
	c.Set(KeyOfRequestUserId, userId)
}

func GetRequestUserId(c *gin.Context) string {
	return c.GetString(KeyOfRequestUserId)
}

func GetRequestUser[T any](c *gin.Context) T {
	user, _ := c.Get(KeyOfRequestUser)
	return user.(T)
}

func MiddileWareCheckToken[T any](encryptSecret string, cookieExpiredSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tk, err := c.Cookie(KeyOfCookieToken)
		if err != nil {
			ResponseFailed(c, errx.InvalidToken)
			return
		}

		uid1, err := c.Cookie(KeyOfCookieUserId)
		if err != nil {
			ResponseFailed(c, errx.InvalidToken)
			return
		}

		uid2, err := securex.ParseToken(tk, encryptSecret, int64(cookieExpiredSeconds))
		if err != nil {
			ResponseFailed(c, errx.InvalidToken.AddErr(err))
			return
		}

		if uid1 != uid2 {
			ResponseFailed(c, errx.InvalidToken.AddErr(err))
			return
		}

		existsUser, err := dbx.QueryOneByPk[T](uid2)
		if err != nil {
			ResponseFailed(c, errx.UserNotFound.AddErr(err))
			return
		}
		SetRequestUser(c, existsUser, uid2)
		c.Next()
	}
}
