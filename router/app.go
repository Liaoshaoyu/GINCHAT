package router

import (
	"GINCHAT/docs"
	"GINCHAT/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("index", service.GetIndex)
	r.GET("user/getUserList", service.GetUserList)
	r.GET("user/createrUser", service.CreaterUser)
	r.GET("user/deleteUser", service.DeleteUser)
	r.POST("user/updateUser", service.UpdateUser)
	r.POST("user/findUserByNameAndPassword", service.FindUserByNameAndPassword)

	// 发送消息
	r.GET("user/sendMsg", service.SendMsg)
	return r
}
