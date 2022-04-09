package controller

import (
	"errors"
	"fmt"
	"net/http"
	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 用户注册接口
// @Summary      用户注册接口
// @Description  通过该接口进行注册账号
// @Tags         用户相关接口
// @Accept       application/json
// @Produce      application/json
// @Param        ParamSignUp  body     models.ParamSignUp  true  "用户注册参数"
// @Success      200     {object}  swaggerResponse “成功”
// @Router       /signup [post]
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		fmt.Println(err)
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseWithMsg(c, CodeInvalidParam, errs.Translate(trans)) //将错误翻译成中文
		return
	}

	// 2. 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
		}
		ResponseError(c, CodeServerBusy)
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册失败！",
		})
		return
	}
	// 3. 返回响应
	//ResponseWithMsg(c, CodeSuccess, "注册成功！")
	ResponseSuccess(c, "注册成功！") //
}

// LoginHandler 用户登录接口
// @Summary      用户登录接口
// @Description  通过该接口进行登录账号，获得 token
// @Tags         用户相关接口
// @Accept       application/json
// @Produce      application/json
// @Param        ParamLogin  body     models.ParamLogin  true  "用户登录参数"
// @Success      200     {object}  swaggerResponse “成功”
// @Router       /login [post]
func LoginHandler(c *gin.Context) {
	// 1.获取参数 和 参数校验
	p := new(models.ParamLogin)
	// 参数绑定
	if err := c.ShouldBindJSON(p); err != nil {
		fmt.Println(err)
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans))) //将错误翻译成中文
		return
	}

	// 2.业务处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username:", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
		} else if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
		} else {
			ResponseError(c, CodeServerBusy)
		}
		return
	}
	// 修改数据库中token的值，便于限制同一时间只能有一台设备登录
	if err := mysql.UpdateToken(p.Username, user.Token); err != nil {
		zap.L().Error("alter token failed", zap.Error(err))
		return
	}

	// 3.返回响应
	//ResponseWithMsg(c, CodeSuccess, token)
	ResponseSuccess(c, gin.H{
		"user_id":  fmt.Sprintf("%d", user.UserID),
		"username": user.Username,
		"token":    user.Token,
	})
}
