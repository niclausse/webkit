package errorx

// 入参类错误:    4000~4999
// 服务端内部错误: 5000～5999
// 服务端数据错误: 6000~6999

var (
	InvalidParamErr = New(4000, "请求参数错误")

	SystemErr   = New(5000, "服务端错误")
	QueryDBErr  = New(5001, "数据查询失败")
	SaveDBErr   = New(5002, "数据保存失败")
	InsertDBErr = New(5003, "数据新增失败")
	UpdateDBErr = New(5004, "数据更新失败")
	DeleteDBErr = New(5005, "数据删除失败")

	DataNotExistErr = New(6000, "数据不存在")
)
