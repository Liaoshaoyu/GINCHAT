package router

import (
	"GINCHAT/docs"
	"GINCHAT/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	var r = gin.Default()

	// swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 静态资源
	r.Static("asset", "asset/")
	r.LoadHTMLGlob("templates/**/*")

	// 首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.GET("/chat", service.Chat)
	r.POST("/searchFriends", service.SearchFriends)

	// 用户
	r.GET("user/getUserList", service.GetUserList)
	r.POST("user/createrUser", service.CreaterUser)
	r.POST("user/deleteUser", service.DeleteUser)
	r.POST("user/updateUser", service.UpdateUser)
	r.POST("user/findUserByNameAndPwd", service.FindUserByNameAndPwd)
	r.POST("user/findByID", service.FindByID)

	// 发送消息
	r.GET("user/sendMsg", service.SendMsg)
	r.GET("user/sendUserMsg", service.SendUserMsg)

	// 添加
	r.POST("/contact/addfriend", service.AddFriend)
	r.POST("/contact/loadcommunity", service.LoadCommunity)

	// 404
	r.GET("/404", service.NotFound)
	return r
}
