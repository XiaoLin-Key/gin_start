package controller

import (
	"gin_start/logger"
	"gin_start/logic"
	"gin_start/result"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 社区列表接口
// @Security ApiKeyAuth
// @Summary      获取所有社区列表
// @Description  查询并返回系统中所有可用的社区信息
// @Tags         社区相关
// @Accept       json
// @Produce      json
// @Success      200  {object}  result.ResponseData{data=[]models.Community}  "成功返回社区列表"
// @Failure      500  {object}  result.ResponseData                          "服务器繁忙"
// @Router       /community [get]
func CommunityHandler(c *gin.Context) {
	//查询到所有社区
	data, err := logic.GetCommunityList()
	if err != nil {
		logger.Lg.Error("logic.GetCommunityList() failed, err: %v", zap.Error(err))
		result.ResponseError(c, result.CodeServerBusy)
		return
	}
	result.ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区详情接口
// @Security ApiKeyAuth
// @Summary      根据 ID 获取社区详情
// @Description  通过路径参数传入社区 ID，获取该社区的详细信息
// @Tags         社区相关
// @Accept       json
// @Produce      json
// @Param        id   path      int64  true  "社区ID"
// @Success      200  {object}  result.ResponseData{data=models.CommunityDetail} "成功返回社区详情"
// @Failure      400  {object}  result.ResponseData                             "参数错误"
// @Failure      500  {object}  result.ResponseData                             "服务器繁忙"
// @Router       /community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	//获取社区id
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
	//查询社区详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		logger.Lg.Error("logic.GetCommunityDetail() failed, err: %v", zap.Error(err))
		result.ResponseError(c, result.CodeServerBusy)
		return
	}
	result.ResponseSuccess(c, data)
}
