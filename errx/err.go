package errx

import "fmt"

// common
var (
	RequestFailed       = &Err{Code: 100100, Msg: "请求失败"}
	CreateDataFailed    = &Err{Code: 100101, Msg: "创建数据失败"}
	QueryDataFailed     = &Err{Code: 100102, Msg: "查询数据失败"}
	UpdateDataFailed    = &Err{Code: 100103, Msg: "更新数据失败"}
	DeleteDataFailed    = &Err{Code: 100104, Msg: "删除数据失败"}
	TargetNotExists     = &Err{Code: 100105, Msg: "目标数据不存在"}
	TargetAlreadyExists = &Err{Code: 100106, Msg: "目标数据已存在"}

	InvalidParams     = &Err{Code: 100201, Msg: "无效的参数"}
	InvalidJsonParams = &Err{Code: 100202, Msg: "无效的JSON参数"}
	InvalidPathParams = &Err{Code: 100203, Msg: "无效的路径参数"}
	InvalidHeader     = &Err{Code: 100204, Msg: "无效的请求头"}

	ValidateParamFailed = &Err{Code: 100211, Msg: "校验参数失败"}
	ParseParamFailed    = &Err{Code: 100212, Msg: "解析参数失败"}
	NoPermission        = &Err{Code: 100213, Msg: "没有权限"}

	AuthFailed        = &Err{Code: 100401, Msg: "认证失败"}
	InvalidToken      = &Err{Code: 100402, Msg: "无效的令牌"}
	UserAlreadyExists = &Err{Code: 100403, Msg: "用户已存在"}
	UserNotFound      = &Err{Code: 100404, Msg: "用户不存在"}
	IncorrectPassword = &Err{Code: 100405, Msg: "密码不正确"}
	HashPassword      = &Err{Code: 100406, Msg: "哈希密码失败"}
	CreateToken       = &Err{Code: 100407, Msg: "创建令牌失败"}
)

type Err struct {
	subErr error `json:"-"`
	Code   int
	Msg    string
}

func (err *Err) String() string {
	return fmt.Sprintf("Code: %v, Message: %v, Error:%v", err.Code, err.Msg, err.subErr)
}

func (err *Err) AddMsg(msg string) *Err {
	return &Err{Code: err.Code, Msg: err.Msg + ":" + msg}
}

func (err *Err) AddErr(er error) *Err {
	return &Err{Code: err.Code, Msg: err.Msg + ":" + er.Error()}
}

func (err *Err) Args(args ...any) *Err {
	return &Err{Code: err.Code, Msg: fmt.Sprintf(err.Msg, args...)}
}
