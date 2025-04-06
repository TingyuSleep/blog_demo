package controller

import (
	"blog_demo/logic"
	"blog_demo/models"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func PostVoteHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		// 类型断言
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 断言失败，说明是ShouldBindJSON反序列化时出错，如要求的是int类型，却传入string类型
			// ShouldBindJSON只能简单校验类型是否相符，不能验证required,oneof，这个是由validator校验的
			ResponseError(c, CodeInvalidParam)
			return
		}
		// 否则是validator校验时，参数不符合要求而出错，比如 required,oneof
		//翻译并去除掉错误提示中的结构体标识
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	//具体投票的业务逻辑
	if err = logic.VoteForPost(c, userID, p); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}
