package errx

import (
	"fmt"
)

// common
var (
	Success            = New(0, "ok")
	RequestFailed      = New(100100, "请求失败")
	InvalidParams      = New(100101, "无效的参数")
	InvalidJsonParams  = New(100102, "无效的JSON参数")
	InvalidQueryParams = New(100103, "无效的查询参数")
	InvalidHeader      = New(100104, "无效的请求头")
	InvalidBody        = New(100105, "无效的请求体")

	CreateDataFailed    = New(100201, "创建数据失败")
	QueryDataFailed     = New(100202, "查询数据失败")
	UpdateDataFailed    = New(100203, "更新数据失败")
	DeleteDataFailed    = New(100204, "删除数据失败")
	TargetNotExists     = New(100205, "目标不存在")
	TargetAlreadyExists = New(100206, "目标已存在")
	FileDirNotExist     = New(100207, "目录不存在")
	FileNotExist        = New(100208, "文件不存在")
	FileReadFailed      = New(100209, "文件读取失败")
	FileWriteFailed     = New(100210, "文件写入失败")

	AuthFailed        = New(100401, "认证失败")
	InvalidToken      = New(100402, "无效的令牌")
	UserAlreadyExists = New(100403, "用户已存在")
	UserNotFound      = New(100404, "用户不存在")
	IncorrectPassword = New(100405, "密码不正确")
	HashPassword      = New(100406, "哈希密码失败")
	CreateToken       = New(100407, "创建令牌失败")

	PraseJsonError = New(100501, "解析json错误")
)

type Err interface {
	Error() string
	Msg() string
	Code() int
	String() string
	AddMsg(msg string) Err
	AddMsgf(format string, args ...any) Err
	AddErr(er error) Err
	Format(args ...any) Err
}

type err struct {
	code int
	msg  string
}

func New(code int, msg string) Err {
	return &err{code, msg}
}

func (e *err) Error() string {
	return e.msg
}

func (e *err) Msg() string {
	return e.msg
}

func (e *err) Code() int {
	return e.code
}

func (e *err) String() string {
	return fmt.Sprintf("Code: %v, Msg: %v", e.code, e.msg)
}

func (e *err) AddMsg(msg string) Err {
	return New(e.code, fmt.Sprintf("%s[%s]", e.msg, msg))
}

func (e *err) AddMsgf(format string, args ...any) Err {
	return New(e.code, fmt.Sprintf("%s[%s]", e.msg, fmt.Sprintf(format, args...)))
}

func (e *err) AddErr(er error) Err {
	return New(e.code, fmt.Sprintf("%s[%s]", e.msg, er.Error()))
}

func (e *err) Format(args ...any) Err {
	return New(e.code, fmt.Sprintf(e.msg, args...))
}

// type Err struct {
// 	Code int
// 	Msg  string
// }

// func (err *Err) Error() string {
// 	err2 := errors.New(err.Msg)
// 	return err.Msg
// }

// func (err *Err) String() string {
// 	return fmt.Sprintf("Code: %v, Msg: %v", err.Code, err.Msg)
// }

// func (err *Err) AddMsg(msg string) Err {
// 	return Err{Code: err.Code, Msg: fmt.Sprintf("%s[%s]", err.Msg, msg)}
// }

// func (err *Err) AddErr(er error) Err {
// 	return Err{Code: err.Code, Msg: fmt.Sprintf("%s[%s]", err.Msg, er.Error())}
// }

// func (err *Err) Format(args ...any) Err {
// 	return Err{Code: err.Code, Msg: fmt.Sprintf(err.Msg, args...)}
// }
