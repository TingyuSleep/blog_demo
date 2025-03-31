package controller

import (
	"blog_demo/middlewares"
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

// GetCurrentUser 获取当前登录用户的ID
func GetCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, exist := c.Get(middlewares.CtxUserIDKey)
	if !exist {
		err = ErrorUserNotLogin
		return
	}
	//对接口(c.Get会返回一个接口类型值)进行类型断言
	userID, ok := uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}
