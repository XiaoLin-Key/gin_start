package controller

import (
	"gin_start/logger"
	"gin_start/logic"
	"gin_start/result"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
