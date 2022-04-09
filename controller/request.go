package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strconv"
)

var (
	CtxUserIDKey      = "userID"
	CtxUserNameKey    = "userName"
	ErrorUserNotLogin = errors.New("用户未登录")
)

// GetCurrentUser 获取当前用户的userID
func GetCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}

	return
}

// GetPageInfo 获取分页的参数
func GetPageInfo(c *gin.Context) (int64, int64) {
	pageNumStr := c.Query("page")
	pageSizeStr := c.Query("size")

	page, err := strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		page = viper.GetInt64("paramPostList.page")
	}
	size, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		size = viper.GetInt64("paramPostList.size")

	}
	return page, size
}
