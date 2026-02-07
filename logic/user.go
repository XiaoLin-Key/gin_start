package logic

import (
	"errors"
	"gin_start/dao/mysql"
	"gin_start/models"
	"gin_start/pkg/jwt"
	"gin_start/pkg/snowflake"
)

// 注册用户
func SignUp(p *models.ParamSignUp) (err error) {
	//判断用户存不存在
	if exist, err := mysql.CheckUserExist(p.Username); err == nil && exist {
		//用户存在
		return errors.New("用户已存在")
	}
	//生成UID
	userID := snowflake.GenID()
	//构造一个User实例
	user := models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//保存数据库
	if err := mysql.InsertUser(&user); err != nil {
		return err
	}
	return nil
}

// 登录用户
func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//判断密码是否正确
	if err = mysql.Login(user); err != nil {
		//密码错误
		return nil, err
	}
	//密码正确
	//生成token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return user, nil
}
