package controller

import (
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// PostVoteHandler 帖子投票接口
// @Summary      帖子投票接口
// @Description  可以为帖子进行投票
// @Tags         帖子相关接口
// @Accept       application/json
// @Produce      application/json
// @Param        Authorization  header  string            true  "Bearer 用户令牌"
// @Param        object         body   models.ParamVote  true  "投票参数"
// @Security     ApiKeyAuth
// @Success      200  {object}  swaggerResponse
// @Router       /vote [post]
func PostVoteHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(models.ParamVote)
	if err := c.ShouldBindJSON(p); err != nil {
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			zap.L().Error("postVoteHandler.shouldBindJSON failed", zap.Error(err))
			// ResponseError(c, CodeInvalidParam)
			ResponseWithData(c, CodeInvalidParam, p)
			return
		}
		ResponseWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans))) // 翻译错误，并没有去掉结构体
		return
	}

	// 获取当前请求用户的id
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 2. 具体投票的业务逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}
