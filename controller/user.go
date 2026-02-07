package controller

import (
	"errors"
	"gin_start/dao/mysql"
	"gin_start/logic"
	"gin_start/models"
	"gin_start/result"

	myValidator "gin_start/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SignUpHandler(c *gin.Context) {
	//参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 获取 validator.ValidationErrors 类型的错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			result.ResponseError(c, result.CodeInvalidParam)
			return
		}
		// 使用我们定义的全局 Trans 进行翻译
		result.ResponseErrorWithMsg(c, result.CodeInvalidParam, errs.Translate(myValidator.Trans))
		return
	}
	//参数校验
	if p.Password != p.RePassword {
		result.ResponseErrorWithMsg(c, result.CodeInvalidParam, "两次密码不一致")
		return
	}
	//业务逻辑处理
	if err := logic.SignUp(p); err != nil {
		if errors.Is(err, mysql.ErrUserExist) {
			result.ResponseError(c, result.CodeUserExist)
			return
		}
		result.ResponseError(c, result.CodeServerBusy)
		return
	}
	//返回响应
	result.ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	//参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 获取 validator.ValidationErrors 类型的错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			result.ResponseError(c, result.CodeInvalidParam)
			return
		}
		// 使用我们定义的全局 Trans 进行翻译
		result.ResponseErrorWithMsg(c, result.CodeInvalidParam, errs.Translate(myValidator.Trans))
		return
	}
	//业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		if errors.Is(err, mysql.ErrUserNotExist) {
			result.ResponseError(c, result.CodeUserNotExist)
			return
		}
		if errors.Is(err, mysql.ErrInvalidPassword) {
			result.ResponseError(c, result.CodeInvalidPassword)
			return
		}
		result.ResponseError(c, result.CodeServerBusy)
		return
	}
	//返回响应
	result.ResponseSuccess(c, gin.H{
		"user_id":  user.UserID,
		"username": user.Username,
		"token":    user.Token,
	})
}
