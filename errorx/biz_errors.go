package errorx

const (
	NoParamInvalid = 1 //参数错误
	NoSystemError  = 2 //服务内部错误

	NoDBInsertErr = 10
	NoDBUpdateErr = 11
	NoDBFindErr   = 12
)

var ErrMSG = map[int]string{
	NoParamInvalid: "请求参数不合理",
	NoSystemError:  "服务异常， 请稍后重试",

	NoDBInsertErr: "db insert error",
	NoDBUpdateErr: "db update error",
	NoDBFindErr:   "db find error",
}

var (
	ParamInvalid = New(NoParamInvalid, ErrMSG[NoParamInvalid])
	SystemError  = New(NoSystemError, ErrMSG[NoSystemError])

	DBInsertErr = New(NoDBInsertErr, ErrMSG[NoDBInsertErr])
	DBUpdateErr = New(NoDBUpdateErr, ErrMSG[NoDBUpdateErr])
	DBFindErr   = New(NoDBFindErr, ErrMSG[NoDBFindErr])
)
