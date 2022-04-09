package controller

import (
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	adminName = "诸葛青"
)

// --- 跟社区相关的 ---

// CreateCommunityHandler 创建社区接口
// @Summary      创建社区接口
// @Description  可以创建社区（只有管理员身份有权限）
// @Tags         社区相关接口
// @Accept       application/json
// @Produce      application/json
// @Param        Authorization  header  string                       true  "Bearer 用户令牌"
// @Param        ParamCreateCommunity         body   models.ParamCreateCommunity  true  "创建社区参数"
// @Security     ApiKeyAuth
// @Success      200  {object}  swaggerResponse "成功"
// @Router       /community/create [post]
func CreateCommunityHandler(c *gin.Context) {
	// 1.校验身份（只有管理员诸葛青才能进行操作）
	v, _ := c.Get(CtxUserNameKey)
	userName := v.(string)
	if userName != adminName {
		zap.L().Error("CreateCommunityHandler failed", zap.String("username", userName))
		ResponseError(c, CodeNotAdmin)
		return
	}

	// 2. 获取参数
	p := new(models.ParamCreateCommunity)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("CreateCommunityHandler ShouldBind failed",
			zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 3. 创建社区
	if err := logic.CreateCommunity(p); err != nil {
		zap.L().Error("logic.CreateCommunity failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// CommunityHandler 查看社区接口
// @Summary      查看社区接口
// @Description  通过该接口获取所有社区
// @Tags         社区相关接口
// @Produce      application/json
// @Success      200  {object}  swaggerResponse
// @Router       /communities [get]
func CommunityHandler(c *gin.Context) {
	// 1. 查询到所有的社区（community_id, community_name） 以列表形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务器报错暴露给外面
		return
	}
	ResponseSuccess(c, data)

}

// CommunityDetailHandler 根据社区名字来获取社区详情接口
// @Summary      根据社区名字来获取社区详情接口
// @Description  通过该社区可获取对应名字的社区详细信息
// @Tags         社区相关接口
// @Produce      application/json
// @Param	name query string  true "社区名字"
// @Success      200  {object}  swaggerResponse
// @Router       /community [get]
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区名字
	communityName := c.Query("name") //解析路径参数

	// 2. 获取社区详情
	data, err := logic.GetCommunityDetailByName(communityName)
	if err != nil {
		zap.L().Error("logic.CommunityDetailHandler failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
