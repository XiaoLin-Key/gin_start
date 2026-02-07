package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"gin_start/models"
)

const secret = "XiaoLin"

// 检查指定用户名用户是否存在
func CheckUserExist(username string) (bool, error) {
	//查询用户
	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	err := db.Raw(sqlStr, username).Scan(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, ErrUserExist
	}
	return false, nil
}

// 加密密码
func EncryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}

// 插入一个新用户
func InsertUser(user *models.User) error {
	//对密码进行加密
	user.Password = EncryptPassword(user.Password)
	//执行sql入库
	sqlStr := "insert into user(user_id,username,password) values(?,?,?)"
	err := db.Exec(sqlStr, user.UserID, user.Username, user.Password).Error
	if err != nil {
		return err
	}
	return nil
}

// 登录校验
func Login(user *models.User) error {
	//查询用户
	oPassword := user.Password
	sqlStr := "select user_id,username,password from user where username = ?"
	err := db.Raw(sqlStr, user.Username).Scan(user).Error
	if err != nil {
		return err
	}
	if user.UserID == 0 {
		return ErrUserNotExist
	}
	//校验密码
	if user.Password != EncryptPassword(oPassword) {
		return ErrInvalidPassword
	}
	return nil
}

func GetUserByID(userID int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := "select user_id,username from user where user_id = ?"
	err = db.Raw(sqlStr, userID).Scan(user).Error
	return
}
