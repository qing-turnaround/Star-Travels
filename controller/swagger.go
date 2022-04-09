package controller

type swaggerResponse struct {
	Code    ResCode     `json:"code"`    // 业务响应状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 返回数据
}

type swaggerPostRequest struct {
	CommunityID int64  `json:"community_id" binding:"required" example:"4"`    // 社区ID
	Title       string `json:"title" binding:"required" example:"人生没有什么放不下！"`  // 标题
	Content     string `json:"content" binding:"required" example:"知足，知止，便是福"` // 内容
}

