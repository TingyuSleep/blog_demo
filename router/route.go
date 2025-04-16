package router

import (
	"blog_demo/controller"
	"blog_demo/logger"
	"blog_demo/middlewares"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	if err := controller.InitTrans("zh"); err != nil {
		log.Fatalf("初始化翻译器失败: %v", err)
	}
	r := gin.New()
	// 添加全局中间件，日志中间件，限流中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(time.Second*1, 1))

	// 测试路由
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	//使用路由组
	v1 := r.Group("/api/v1")

	//注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	//登录业务路由
	v1.POST("/login", controller.LoginHandler)
	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/postList", controller.GetPostListHandler)
		// 升级版：根据发布时间或者分数获取帖子列表
		v1.GET("/postList2", controller.GetPostListHandler2)
		v1.POST("/vote", controller.PostVoteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
