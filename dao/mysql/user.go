package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"web_app/models"
)

// dao层只涉及数据库的操作逻辑

const secret = "zhugeqing.top" // 用于加密解密的钥匙

// CheckUserExist 检查该用户名是的用户是否存在
func CheckUserExist(userName string) error {

	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, userName); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}

	return nil
}

// InsertUser 插入用户到数据库中
func InsertUser(user *models.User) (err error) {
	// 对密码加密
	user.Password = encryptPassword(user.Password)
	// 执行sql语句入库
	sqlStr := `insert into user(user_id, username, password)values(?, ?, ?)`
	if _, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password); err != nil {
		return
	}

	return nil
}

// encryptPassword 将密码字符串加密
func encryptPassword(password string) string {
	h := md5.New() //返回一个hash.Hash 用于计算校验和
	h.Write([]byte(secret))
	// h.Sum(password)——计算password的md5值
	// 将得到的16进制的md5值进行编码成字符串
	return hex.EncodeToString(h.Sum([]byte(password)))
}

// Login 用户登录与数据库中的数据比较
func Login(user *models.User) error {
	var originPassword = user.Password // 记录没有进行加密的密码
	sqlStr := `select user_id, username, password from user where username=?`
	if err := db.Get(user, sqlStr, user.Username); err != nil {
		return ErrorUserNotExist
	}
	if user.Password != encryptPassword(originPassword) { // 比较数据库里的秘密 和 加密用户输入的密码是否一致
		return ErrorInvalidPassword
	}

	return nil
}

// UpdateToken 更新用户的token值（限制同一时间只有一台设备登录）
func UpdateToken(username, token string) error {
	zap.L().Error("token length is", zap.Int("len", len(token)))
	sqlStr := `update user set token = ? where username = ?`
	if _, err := db.Exec(sqlStr, token, username); err != nil {
		return err
	}
	return nil
}

// CheckToken 检查Token是否有效
func CheckToken(token string, userID int64) error {
	sqlStr := `select token from user where user_id=?`
	var queryToken string
	if err := db.Get(&queryToken, sqlStr, userID); err != nil {
		return err
	}

	if queryToken != token {
		return errors.New("无效的token")
	}

	return nil
}

// GetUserByID 通过用户ID查询用户信息
func GetUserByID(userID int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id=?`
	err = db.Get(user, sqlStr, userID)
	return
}
