package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"web_app/models"

	"go.uber.org/zap"
)

// dao层只涉及数据库的操作逻辑

const secret = "zhugeqing.top" // 用于加密解密的钥匙

/**** 读操作 ****/

// CheckUserExist 检查该用户名是的用户是否存在
func CheckUserExist(userName string) error {
	count := randomGetSlave().Where("user_name = ?", userName).Find(&models.User{}).RowsAffected
	if count > 0 {
		return ErrorUserExist
	}

	return nil
}

// Login 用户登录与数据库中的数据比较
func Login(user *models.User) (err error) {
	var originPassword = user.Password // 记录没有进行加密的密码
	err = randomGetSlave().Where("user_name = ?", user.UserName).Find(user).Error
	if user.Password != encryptPassword(originPassword) { // 比较数据库里的密码 和 加密用户输入的密码是否一致
		return ErrorInvalidPassword
	}
	return
}

// CheckToken 检查Token是否有效
func CheckToken(token string, userID int64) error {
	user := &models.User{}
	err := randomGetSlave().Where("user_id = ?", userID).Find(user).Error
	if user.Token != token {
		return errors.New("无效的token")
	}
	return err
}

// GetUserByID 通过用户ID查询用户信息
func GetUserByID(userID int64) (user *models.User, err error) {
	user = new(models.User)
	err = randomGetSlave().Where("user_id = ?", userID).Find(user).Error
	return
}

/* 写操作 */

// InsertUser 插入用户到数据库中
func InsertUser(user *models.User) (err error) {
	// 对密码加密
	user.Password = encryptPassword(user.Password)
	return master.Create(user).Error
}

// encryptPassword 将密码字符串加密
func encryptPassword(password string) string {
	h := md5.New() //返回一个hash.Hash 用于计算校验和
	h.Write([]byte(secret))
	// h.Sum(password)——计算password的md5值
	// 将得到的16进制的md5值进行编码成字符串
	return hex.EncodeToString(h.Sum([]byte(password)))
}


// UpdateToken 更新用户的token值（限制同一时间只有一台设备登录）
func UpdateToken(username, token string) error {
	zap.L().Info("token length is", zap.Int("len", len(token)))
	err := master.Model(models.User{}).Where("user_name = ?", username).
		UpdateColumn("token", token).Error
	return err
}



