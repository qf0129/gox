package serverx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/pkg/dbx"
	"github.com/qf0129/gox/pkg/errx"
	"github.com/qf0129/gox/pkg/securex"
)

const (
	KeyOfRequestUser         = "ctx_req_user"
	KeyOfRequestUserPk       = "ctx_req_user_pk"
	KeyOfCookieToken         = "tk"
	KeyOfCookieUserPk        = "pk"
	KeyOfHeaderAuthorization = "Authorization"
)

func ClearCookie(c *gin.Context) {
	c.SetCookie(KeyOfCookieToken, "", -1, "/", "", false, true)
	c.SetCookie(KeyOfCookieUserPk, "", -1, "/", "", false, false)
}

func SetCookie(c *gin.Context, token, userPk, domain string, expiredSeconds int) {
	c.SetCookie(KeyOfCookieToken, token, expiredSeconds, "/", domain, false, true)
	c.SetCookie(KeyOfCookieUserPk, userPk, expiredSeconds, "/", domain, false, false)
}

func SetRequestUser(c *gin.Context, user any, userPk string) {
	c.Set(KeyOfRequestUser, user)
	c.Set(KeyOfRequestUserPk, userPk)
}

func GetRequestUserPk(c *gin.Context) string {
	return c.GetString(KeyOfRequestUserPk)
}

func GetRequestUser[T any](c *gin.Context) T {
	user, _ := c.Get(KeyOfRequestUser)
	return user.(T)
}

func MiddileWareCheckToken[T any](encryptSecret string, cookieExpiredSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tk, err := c.Cookie(KeyOfCookieToken)
		if err != nil {
			ResponseErr(c, errx.InvalidToken)
			return
		}

		cookiePk, err := c.Cookie(KeyOfCookieUserPk)
		if err != nil {
			ResponseErr(c, errx.InvalidToken)
			return
		}

		tokenPk, err := securex.ParseToken(tk, encryptSecret, int64(cookieExpiredSeconds))
		if err != nil {
			ResponseErr(c, errx.InvalidToken.AddErr(err))
			return
		}

		if cookiePk != tokenPk {
			ResponseErr(c, errx.InvalidToken.AddErr(err))
			return
		}

		existsUser, err := dbx.QueryOneByPk[T](tokenPk)
		if err != nil {
			ResponseErr(c, errx.UserNotFound.AddErr(err))
			return
		}
		SetRequestUser(c, existsUser, tokenPk)
		c.Next()
	}
}
