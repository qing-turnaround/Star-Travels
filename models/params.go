package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// 定义请求的参数结构体

// ParamSignUp 注册
type ParamSignUp struct {
	Username   string `json:"username" binding:"required" example:"终生让步"` // 用户名字
	Password   string `json:"password" binding:"required" example:"12345"` // 用户密码
	RePassword string `json:"re_password" binding:"required,eqfield=Password" example:"12345"` // 再次确认密码
}

// ParamLogin 登录
type ParamLogin struct {
	Username string `json:"username" binding:"required" example:"终生让步"` // 用户名字
	Password string `json:"password" binding:"required" example:"12345"` // 用户密码
}

// ParamVote 投票请求携带数据
type ParamVote struct {
	// userId; 从请求中获取
	PostID string `json:"post_id" binding:"required"`
	//validator binding:oneof 表示该变量的值只能是其中一个（赞同为1，反对为-1）
	Direction string `json:"direction" binding:"oneof=0 1 -1"`
}

// ParamPostList 帖子请求url携带参数（form 可以 配合 c.ShouldBindQuery 来绑定 url中的参数, json是为了配合swagger文档的参数名字）
type ParamPostList struct {
	Page          int64  `form:"page" example:"1"`	// 查询第几页，默认第1页
	Size          int64  `form:"size" example:"10"` // 每一页帖子的数量，默认每一页10个帖子
	Order         string `form:"order" binding:"oneof=time score" example:"time"` // 查询的排序规则（根据时间获取投票数进行排序，填time 者 score，默认为time）
	CommunityName string `json:"community_name" form:"community_name" example:"成长的路口"` // 查询的帖子所在的社区名称
}

// ParamCreateCommunity 创建社区
type ParamCreateCommunity struct {
	CommunityID   int64  `json:"community_id"   binding:"required" example:"1"`
	CommunityName string `json:"community_name"   binding:"required" example:"Go"`
	Introduction  string `json:"introduction"  binding:"required" example:"Go语言"`
}

