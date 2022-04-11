package controller

type swaggerResponse struct {
	Code    ResCode     `json:"code"`    // 业务响应状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 返回数据
}
