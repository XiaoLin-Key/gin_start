package controller

import (
	"gin_start/logic"
	"gin_start/models"
	myValidator "gin_start/pkg/validator"
	"gin_start/result"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PostVoteHandler 帖子投票接口
// @Security ApiKeyAuth
// @Summary      帖子投票
// @Description  用户给帖子投票，支持赞成票(1)、反对票(-1)和取消投票(0)
// @Tags         投票相关
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        data  body      models.ParamVoteData  true  "投票数据"
// @Success      200   {object}  result.ResponseData        "投票成功"
// @Failure      400   {object}  result.ResponseData        "请求参数错误"
// @Failure      401   {object}  result.ResponseData        "需要登录"
// @Failure      500   {object}  result.ResponseData        "服务端繁忙"
// @Router       /vote [post]
func PostVoteHandler(c *gin.Context) {
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			result.ResponseError(c, result.CodeInvalidParam)
			return
		}
		//翻译错误信息
		result.ResponseErrorWithMsg(c, result.CodeInvalidParam, errs.Translate(myValidator.Trans))
		return
	}
	userID, err := GetCurrentUser(c)
	if err != nil {
		result.ResponseError(c, result.CodeInvalidParam)
		return
	}
	if err := logic.VoteForPost(userID, p); err != nil {
		result.ResponseError(c, result.CodeServerBusy)
		return
	}
	result.ResponseSuccess(c, nil)

}
