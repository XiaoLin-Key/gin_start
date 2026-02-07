package controller

import (
	"gin_start/logic"
	"gin_start/models"
	myValidator "gin_start/pkg/validator"
	"gin_start/result"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

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
