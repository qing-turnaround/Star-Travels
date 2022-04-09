package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

// SignUp 注册逻辑处理函数
func SignUp(p *models.ParamSignUp) (err error) {
	// 1.判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return
	}
	if err != nil {
		return
	}
	// 2.生成用户ID
	userID := snowflake.GetID()
	// 构造一个User实例
	u := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3.保持进入数据库
	return mysql.InsertUser(u)
}

// Login 登录逻辑处理函数
func Login(p *models.ParamLogin) (user *models.User, err error) {
	// 构造一个User实例
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	if err = mysql.Login(user); err != nil {
		return nil, err
	}

	// 生成JWT（json web token）
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return
}
