package controller

import (
	"gin_start/logger"
	"gin_start/logic"
	"gin_start/models"
	myValidator "gin_start/pkg/validator"
	"gin_start/result"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func CreatePostHandler(c *gin.Context) {
	//获取参数并校验
	p := new(models.Post)
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
	//创建帖子
	userID, err := GetCurrentUser(c)
	if err != nil {
		result.ResponseError(c, result.CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	if err := logic.CreatePost(p); err != nil {
		logger.Lg.Error("logic.CreatePost() failed", zap.Error(err))
		result.ResponseError(c, result.CodeServerBusy)
		return
	}

	//返回响应
	result.ResponseSuccess(c, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	//获取帖子id
	idStr := c.Param("id")
	if idStr == "" {
		logger.Lg.Error("idStr is empty")
		result.ResponseError(c, result.CodeInvalidParam)
		return
	}
	//将字符串转换为int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Lg.Error("strconv.ParseInt() failed, err: %v", zap.Error(err))
		result.ResponseError(c, result.CodeInvalidParam)
		return
	}
	//查询帖子详情
	data, err := logic.GetPostById(id)
	if err != nil {
		logger.Lg.Error("logic.GetPostById() failed, err: %v", zap.Error(err))
		result.ResponseError(c, result.CodeServerBusy)
		return
	}
	result.ResponseSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	page, size := GetPageInfo(c)
	data, err := logic.GetPostList(page, size)
	if err != nil {
		logger.Lg.Error("logic.GetPostList() failed, err: %v", zap.Error(err))
		result.ResponseError(c, result.CodeServerBusy)
		return
	}
	result.ResponseSuccess(c, data)
}
