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

// SignUpHandler 注册接口
// @Summary      用户注册接口
// @Description  处理用户注册逻辑，包含用户名、密码校验及重复密码一致性检查
// @Tags         用户相关
// @Accept       json
// @Produce      json
// @Param        data  body      models.ParamSignUp  true  "注册信息"
// @Success      200   {object}  result.ResponseData       "注册成功"
// @Failure      400   {object}  result.ResponseData       "参数错误"
// @Failure      500   {object}  result.ResponseData       "服务器繁忙"
// @Router       /signup [post]
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

// LoginHandler 登录接口
// @Summary      用户登录接口
// @Description  用户登录成功后会返回 userID、username 以及用于后续鉴权的 JWT Token
// @Tags         用户相关
// @Accept       json
// @Produce      json
// @Param        data  body      models.ParamLogin  true  "登录信息"
// @Success      200   {object}  result.ResponseData{data=models.ParamLoginResponse} "登录成功，返回用户信息及Token"
// @Failure      400   {object}  result.ResponseData                               "参数错误"
// @Failure      500   {object}  result.ResponseData                               "服务器繁忙"
// @Router       /login [post]
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
