package routes

import (
	"net/http"
	"time"
	"web_app/controller"
	_ "web_app/docs"
	"web_app/logger"
	"web_app/middlewares"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func SetUp(Mode string) *gin.Engine {
	if Mode == gin.ReleaseMode {
		// gin设置成发布模式：gin不在终端输出日志
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true),
		middlewares.RateLimitMiddleware(time.Microsecond*10, 10)) // 全网站限流

	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// swagger 文档（http://localhost:8080/swagger/index.html）
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1") // 此后用到v1的，访问路径前面都会有r.Group relativePath

	// 登陆 注册接口
	{
		v1.POST("/signup", controller.SignUpHandler) //注册
		v1.POST("/login", controller.LoginHandler)   //登录
	}

	// 获取帖子接口
	{
		v1.GET("/posts2", controller.GetPostListDetailHandler2)
		v1.GET("/post", controller.GetPostDetailHandler)
		v1.GET("/posts/", controller.GetPostListDetailHandler)

	}

	// 获取社区接口
	{
		v1.GET("/communities", controller.CommunityHandler)
		v1.GET("/community", controller.CommunityDetailHandler)
	}

	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		// 创建帖子
		v1.POST("/post", controller.CreatePostHandler)

		// 帖子投票
		v1.POST("/vote", controller.PostVoteHandler)

		// 创建社区（只有管理员“诸葛青”才能访问）
		v1.POST("/community/create", controller.CreateCommunityHandler)

	}
	// http://localhost:9999/debug/pprof/
	pprof.Register(r) // 注册pprof相关路由

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "抱歉，没有找到对应页面。再细心点吧！",
		})
	})
	return r
}
