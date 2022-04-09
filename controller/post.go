package controller

import (
	"strconv"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子接口
// @Summary      创建帖子的接口
// @Description  通过该接口来创建帖子
// @Tags         帖子相关接口
// @Accept       application/json
// @Produce      application/json
// @Param        Authorization  header  string       true  "Bearer 用户令牌"
// @Param        Post         body  swaggerPostRequest  false  "创建帖子参数"
// @Security     ApiKeyAuth
// @Success      200  {object}   swaggerResponse "成功"
// @Router       /post [post]
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("c.shouldBindJSON failed:", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从c 中取得请求用户的ID
	useID, err := GetCurrentUser(c)
	if err != nil {
		zap.L().Error("controller.CreatePost.GetCurrentUser failed:", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = useID
	// 2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的接口
// @Summary      获取帖子详情的接口
// @Description  通过urlID参数来获取帖子详情的接口
// @Tags         帖子相关接口
// @Accept       application/json
// @Produce      application/json
// @Param        id         query   int64   true  "查询参数"
// @Security     ApiKeyAuth
// @Success      200  {object}  swaggerResponse "成功"
// @Router       /post [get]
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取参数（从URL中获取帖子的id）
	idStr := c.Query("id")

	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid id", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 根据id取出帖子数据（查数据库）
	data, err := logic.GetPostByID(postID)
	if err != nil {
		zap.L().Error("logic.GetPostDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListDetailHandler 获取帖子列表详细信息的接口
// @Summary      获取帖子列表详细信息的接口
// @Description  通过url参数获取帖子列表详细信息的接口
// @Tags         帖子相关接口
// @Accept       application/json
// @Produce      application/json
// @Param        page         query   int64   false  "第几页"
// @Param        size        query   int64   false  "每页多少个帖子"
// @Success      200  {object}  swaggerResponse "成功"
// @Router       /posts [get]
func GetPostListDetailHandler(c *gin.Context) {
	// 1. 获取参数（获取分页的参数）
	page, size := GetPageInfo(c)
	// 2. 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostListDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListDetailHandler2 升级版获取帖子接口
// @Summary      升级版获取帖子接口
// @Description  可根据社区名称（默认为空）和 帖子排序规则来获取帖子（也可以填page和size参数）
// @Tags         帖子相关接口
// @Accept       application/json
// @Produce      application/json
// @Param        ParamPostList         query   models.ParamPostList  false  "帖子查询参数"
// @Security     ApiKeyAuth
// @Success      200  {object}  swaggerResponse "成功"
// @Router       /posts2 [get]
func GetPostListDetailHandler2(c *gin.Context) {
	// 1. 获取参数（从url中获取参数）/api/v1/posts2?page=1&size=10&order=time
	p := &models.ParamPostList{
		Page:  viper.GetInt64("paramPostList.page"), // 默认参数第一页
		Size:  viper.GetInt64("paramPostList.size"), // 默认每一页10个帖子
		Order: models.OrderTime,                     // 默认按时间排序
		CommunityName: "",
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListDetailHandler2 c.ShouldBindQuery failed",
			zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 获取数据
	//data, err := logic.GetPostList2(p)
	data, err := logic.GetPostListNew(p)
	if err != nil {
		zap.L().Error("logic.GetPostListDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetCommunityPostListHandler 根据社区来请求帖子
//func GetCommunityPostListHandler(c *gin.Context) {
//	// 1. 获取参数（从url中获取参数）/api/v1/posts2?page=1&size=10&order=time
//	p := &models.ParamPostList{
//		Page:  viper.GetInt64("paramPostList.page"), // 默认参数第一页
//		Size:  viper.GetInt64("paramPostList.size"), // 默认每一页10个帖子
//		Order: models.OrderTime,                     // 默认按时间排序
//	}
//	// c.ShouldBind 根据请求的数据类型选择相应的方法来获取数据
//	if err := c.ShouldBind(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler with invalid params",
//			zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	// 2. 获取数据
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("GetCommunityPostListHandler logic.GetCommunityPostList failed",
//			zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//
//	ResponseSuccess(c, data)
//}
