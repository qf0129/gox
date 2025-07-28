package errx

import (
	"fmt"
)

// common
var (
	Success            = New(0, "ok")
	RequestFailed      = New(400100, "请求失败")
	InvalidParams      = New(400101, "无效的参数")
	InvalidJsonParams  = New(400102, "无效的JSON参数")
	InvalidQueryParams = New(400103, "无效的查询参数")
	InvalidHeader      = New(400104, "无效的请求头")
	InvalidBody        = New(400105, "无效的请求体")

	CreateDataFailed    = New(400201, "创建数据失败")
	QueryDataFailed     = New(400202, "查询数据失败")
	UpdateDataFailed    = New(400203, "更新数据失败")
	DeleteDataFailed    = New(400204, "删除数据失败")
	TargetNotExists     = New(400205, "目标不存在")
	TargetAlreadyExists = New(400206, "目标已存在")
	FileDirNotExist     = New(400207, "目录不存在")
	FileNotExist        = New(400208, "文件不存在")
	FileReadFailed      = New(400209, "文件读取失败")
	FileWriteFailed     = New(400210, "文件写入失败")

	AuthFailed        = New(401001, "认证失败")
	InvalidToken      = New(401002, "无效的令牌")
	TokenIsExpired    = New(401003, "令牌已过期")
	UserAlreadyExists = New(402001, "用户已存在")
	UserNotFound      = New(402002, "用户不存在")
	IncorrectPassword = New(402003, "密码不正确")
	HashPassword      = New(402004, "哈希密码失败")
	CreateToken       = New(402005, "创建令牌失败")

	PraseJsonError = New(400501, "解析json错误")
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
