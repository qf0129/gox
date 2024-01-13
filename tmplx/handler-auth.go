package tmplx

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/gox/confx"
	"github.com/qf0129/gox/constx"
	"github.com/qf0129/gox/dbx"
	"github.com/qf0129/gox/errx"
	"github.com/qf0129/gox/ginx"
	"github.com/qf0129/gox/respx"
	"github.com/qf0129/gox/securex"
	"github.com/qf0129/gox/validx"
)

type AuthRequestBody struct {
	Username string `validate:"gte=2,lte=50" json:"username"`
	Password string `validate:"gte=2,lte=50" json:"password"`
}

// 用户登录接口
func HandleSignIn(conf *confx.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequestBody
		if err := c.ShouldBindJSON(&req); err != nil {
			respx.Err(c, errx.InvalidParams)
			return
		}

		existUser, er := dbx.QueryOneByMap[User](map[string]any{"username": req.Username})
		if er != nil {
			respx.Err(c, errx.UserNotFound)
			return
		}

		if !securex.VerifyPassword(req.Password, existUser.Password) {
			respx.Err(c, errx.IncorrectPassword)
			return
		}

		token, er := securex.CreateToken(existUser.Id, conf.EncryptSecret)
		if er != nil {
			respx.Err(c, errx.CreateToken)
			return
		}

		setAuthCookie(c, token, existUser.Id, conf)
		respx.OK(c, gin.H{"token": token})
	}
}

// 用户注册接口
func HandleSignUp(conf *confx.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequestBody
		if err := c.ShouldBindJSON(&req); err != nil {
			respx.Err(c, errx.InvalidParams)
			return
		}

		if err := validx.Validate(&req); err != nil {
			respx.Err(c, err)
			return
		}

		existUser, _ := dbx.QueryOneByMap[User](map[string]any{"username": req.Username})
		if existUser.Id != "" {
			respx.Err(c, errx.UserAlreadyExists)
			return
		}

		psdHash, er := securex.HashPassword(req.Password)
		if er != nil {
			respx.Err(c, errx.HashPassword)
			return
		}

		u := &User{
			Username: req.Username,
			Password: psdHash,
		}

		if er = dbx.Create[User](u); er != nil {
			respx.Err(c, errx.CreateUser.AddErr(er))
			return
		}
		token, er := securex.CreateToken(u.Id, conf.EncryptSecret)
		if er != nil {
			respx.Err(c, errx.CreateToken.AddErr(er))
			return
		}
		setAuthCookie(c, token, u.Id, conf)
		respx.OK(c, gin.H{dbx.Opt.ModelPrimaryKey: u.Id})
	}
}

// 用户退出登陆接口
func HandleSignOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		delAuthCookie(c)
		respx.OK(c, nil)
	}
}

func delAuthCookie(c *gin.Context) {
	c.SetCookie(constx.KeyOfCookieToken, "", -1, "/", "", false, true)
	c.SetCookie(constx.KeyOfCookieUserId, "", -1, "/", "", false, false)
}

func setAuthCookie(c *gin.Context, token, userId string, conf *confx.Server) {
	c.SetCookie(constx.KeyOfCookieToken, token, conf.CookieExpiredSeconds, "/", conf.CookieDomain, false, true)
	c.SetCookie(constx.KeyOfCookieUserId, userId, conf.CookieExpiredSeconds, "/", conf.CookieDomain, false, false)
}

type UpdatePasswordBody struct {
	OldPsd string `json:"old_psd" validate:"gte=2,lte=50"`
	NewPsd string `json:"new_psd" validate:"gte=2,lte=50"`
}

func UpdatePassword(c *gin.Context) {
	user := ginx.GetRequestUser[User](c)

	var req *UpdatePasswordBody
	if err := c.ShouldBindJSON(&req); err != nil {
		respx.Err(c, errx.InvalidJsonParams)
		return
	}

	if err := validx.Validate(req); err != nil {
		respx.Err(c, err)
		return
	}

	if !user.CheckPassword(req.OldPsd) {
		respx.Err(c, errx.IncorrectPassword)
		return
	}

	user.SetPassword(req.NewPsd)
	dbx.DB.Save(user)
	respx.OK(c, true)
}
