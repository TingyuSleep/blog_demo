package router

import (
	"blog_demo/controller"
	"blog_demo/logger"
	"blog_demo/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if err := controller.InitTrans("zh"); err != nil {

	}
	r := gin.New()
	//使用日志中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//注册业务路由
	r.POST("/signup", controller.SignUpHandler)
	//登录业务路由
	r.POST("/login", controller.LoginHandler)

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		//如果是登录的用户，判断请求头中是否有 有效的JWT
		isLogin := true
		c.Request.Header.Get("Authorization")
		if isLogin {
			c.String(http.StatusOK, "pong")
		} else {
			//否则返回 请先登录
			c.String(http.StatusOK, "请登录")
		}
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
