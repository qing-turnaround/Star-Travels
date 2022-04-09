package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	// 每 fillInterval 时间添加一个令牌。初始桶内令牌数为cap
	bucket := ratelimit.NewBucket(fillInterval, cap)
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
