package middlewares

import (
	"github.com/gin-gonic/gin"
	"strings"
	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/pkg/jwt"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization: Bearer xxxxxxxxx.xxxxxxx.xxxxx
		// 这里的具体实现方式要依据实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization") // 获取请求头里面的Authorization字段
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort() // 验证失败，调用 Abort 以确保这个请求的其他函数不会被调用。
			return
		}
		var token string
		// 为了方便 swagger UI测试接口
		if authHeader[:6] == "Bearer" {
			// 按空格分割（第一部分为Bearer, 第二部分为token字符串）
			parts := strings.SplitN(authHeader, " ", 2)
			token = parts[1]
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				controller.ResponseError(c, controller.CodeInvalidToken)
				c.Abort()
				return
			}
		} else {
			token = authHeader
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(token)
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// 检验用户的token是否与数据库的一致（限制一台设备登录）
		if err = mysql.CheckToken(token, mc.UserID); err != nil {
			// 表示useID的token 与 当前token不匹配
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前请求的user信息保存到请求的上下文c上
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Set(controller.CtxUserNameKey, mc.Username)
		// c.Next() 调用后续处理函数
		c.Next() // 后续的处理函数可以用过c.Get(ctxUserIDKey)来获取当前请求的用户信息
	}
}
