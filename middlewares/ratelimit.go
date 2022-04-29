package middlewares

import (
	"log"
	"net/http"
	"time"
	"web_app/settings"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() func(c *gin.Context) {
	fillInterval := time.Microsecond * time.Duration(settings.Conf.FileInterval) // 添加速度
	log.Println(fillInterval, "fdsfds")
	Cap := settings.Conf.Cap // 容量

	// 每 fillInterval 时间添加一个令牌。初始桶内令牌数为 Cap
	bucket := ratelimit.NewBucket(fillInterval, Cap)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit...
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusOK, "流量限制，请稍后重试...")
			c.Abort() // 停止该请求
			return
		}
		c.Next()
	}
}
