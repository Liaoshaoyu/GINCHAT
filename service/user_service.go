package service

import (
	"GINCHAT/models"
	"github.com/gin-gonic/gin"
)

// GetUserList
// @Tags 获取用户列表
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	var data = make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})
}
