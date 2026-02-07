package controller

import (
	"errors"
	"gin_start/middlewares"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

func GetCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(middlewares.CtxUserID)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func GetPageInfo(c *gin.Context) (page, size int64) {
	//获取分页参数
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return
}
