package service

import (
	"GINCHAT/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} Welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "templates/chat/head.html")
	if err != nil {
		panic(err)
	}
	err2 := ind.Execute(c.Writer, "index")
	if err2 != nil {
		panic(err2)
	}
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("templates/user/register.html")
	if err != nil {
		panic(err)
	}
	err2 := ind.Execute(c.Writer, "register")
	if err2 != nil {
		panic(err2)
	}
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles(
		"templates/chat/index.html",
		"templates/chat/head.html",
		"templates/chat/tabmenu.html",
		"templates/chat/concat.html",
		"templates/chat/group.html",
		"templates/chat/profile.html",
		"templates/chat/main.html",
		"templates/chat/foot.html",
	)
	if err != nil {
		panic(err)
	}
	var userId, _ = strconv.Atoi(c.Query("userId"))
	var token = c.Query("token")
	var user = models.UserBasic{}
	user.ID = uint(userId)
	user.Identify = token
	err2 := ind.Execute(c.Writer, user)
	if err2 != nil {
		panic(err2)
	}
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func NotFound(c *gin.Context) {
	ind, err := template.ParseFiles("templates/chat/404.html")
	if err != nil {
		panic(err)
	}
	err2 := ind.Execute(c.Writer, "404")
	if err2 != nil {
		panic(err2)
	}
}
