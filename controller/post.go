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

// CreatePostHandler 创建帖子接口
// @Security ApiKeyAuth
// @Summary      创建帖子
// @Description  用户登录后可以发布新帖子，包含标题、内容及所属社区
// @Tags         帖子相关
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        post  body      models.Post  true  "帖子信息"
// @Success      200   {object}  result.ResponseData "成功"
// @Failure      400   {object}  result.ResponseData "请求参数错误"
// @Failure      401   {object}  result.ResponseData "需要登录"
// @Router       /post [post]
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

// GetPostDetailHandler 获取帖子详情接口
// @Security ApiKeyAuth
// @Summary      获取帖子详情
// @Description  根据路径中的帖子 ID 查询详细信息
// @Tags         帖子相关
// @Accept       json
// @Produce      json
// @Param        id   path      int64  true  "帖子ID"
// @Success      200  {object}  result.ResponseData{data=models.PostDetailVO} "成功"
// @Failure      400  {object}  result.ResponseData "参数错误"
// @Router       /post/{id} [get]
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

// GetPostListHandler 获取帖子列表接口
// @Security ApiKeyAuth
// @Summary      分页获取帖子列表
// @Description  简单的分页获取帖子列表，按创建时间排序
// @Tags         帖子相关
// @Accept       json
// @Produce      json
// @Param        page  query     int    false  "页码"
// @Param        size  query     int    false  "每页数量"
// @Success      200   {object}  result.ResponseData{data=[]models.PostDetailVO} "成功"
// @Router       /posts [get]
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

// GetPostListHandler2 帖子列表升级版接口
// @Security ApiKeyAuth
// @Summary      分页获取帖子列表(可排序)
// @Description  支持按照时间或分数（热度）排序分页获取帖子
// @Tags         帖子相关
// @Accept       json
// @Produce      json
// @Param        object  query     models.ParamPostList  false  "查询参数"
// @Success      200     {object}  result.ResponseData{data=[]models.PostDetailVO} "成功"
// @Router       /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	p := models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	//获取分页参数
	if err := c.ShouldBindQuery(&p); err != nil {
		logger.Lg.Error("c.ShouldBindQuery() failed, err: %v", zap.Error(err))
		result.ResponseError(c, result.CodeInvalidParam)
		return
	}
	data, err := logic.GetPostListNew(&p)
	if err != nil {
		logger.Lg.Error("logic.GetPostListNew() failed, err: %v", zap.Error(err))
		result.ResponseError(c, result.CodeServerBusy)
		return
	}
	result.ResponseSuccess(c, data)
}
