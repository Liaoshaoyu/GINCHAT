package service

import (
	"GINCHAT/models"
	"GINCHAT/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 获取用户列表
// @Tags 用户模块
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	var data = make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0, // 0: "成功" | -1: "失败"
		"message": "成功",
		"data":    data,
	})
}

// FindUserByNameAndPwd
// @Summary 获取用户列表
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/findUserByNameAndPassword [post]
func FindUserByNameAndPwd(c *gin.Context) {
	var data = models.UserBasic{}
	//var name = c.Query("name")
	//var password = c.Query("password")
	var name = c.Request.FormValue("name")
	var password = c.Request.FormValue("password")
	var user = models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(-1, gin.H{
			"code":    -1, // 0: "成功" | -1: "失败"
			"message": "用户不存在！",
		})
		return
	}
	var flag = utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(-1, gin.H{
			"code":    -1, // 0: "成功" | -1: "失败"
			"message": "密码不正确！",
		})
		return
	}
	data = models.FindUserByNameAndPassword(name, user.PassWord)
	c.JSON(200, gin.H{
		"code":    0, // 0: "成功" | -1: "失败"
		"message": "成功",
		"data":    data,
	})
}

// CreaterUser
// @Summary 注册用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/createrUser [post]
func CreaterUser(c *gin.Context) {
	var user = models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	var password = c.Request.FormValue("password")
	var repassword = c.Request.FormValue("repassword")

	var salt = fmt.Sprintf("%06d", rand.Int31())
	//var res = models.FindUserByName(user.Name)
	var res = models.FindUserByNameAndPassword(user.Name, password)

	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(-1, gin.H{
			"code":    -1, // 0: "成功" | -1: "失败"
			"message": "用户名或密码不能为空！",
		})
		return
	}

	if res.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1, // 0: "成功" | -1: "失败"
			"message": "用户名已注册！",
		})
		return
	}
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1, // 0: "成功" | -1: "失败"
			"message": "两次密码不一致！",
		})
		return
	}

	//user.PassWord = password
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	user.LoginTime = time.Now()
	user.LogoutTime = time.Now()
	user.HeartBeatTime = time.Now()

	var str = fmt.Sprintf("%d", time.Now().Unix())
	user.Identify = utils.MD5Encode(str)

	models.CreaterUser(user)

	c.JSON(200, gin.H{
		"code":    0, // 0: "成功" | -1: "失败"
		"message": "用户注册成功！",
	})

}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "用户id"
// @Success 200 {string} json{"code", "message"}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context) {
	var user = models.UserBasic{}
	var id, _ = strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)

	c.JSON(200, gin.H{
		"code":    0, // 0: "成功" | -1: "失败"
		"message": "用户删除成功！",
	})
}

// UpdateUser
// @Summary 更新用户
// @Tags 用户模块
// @param id formData string false "用户id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code", "message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	var user = models.UserBasic{}
	var id, _ = strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	var _, err = govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(-1, gin.H{
			"code":    -1, // 0: "成功" | -1: "失败"
			"message": "参数格式错误！",
		})
		return
	}
	models.UpdateUser(user)
	c.JSON(200, gin.H{
		"code":    0, // 0: "成功" | -1: "失败"
		"message": "用户更新成功！",
	})
}

// 防止跨域站点伪造请求
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", c)
	var ws, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		fmt.Print("")
		//var msg, err = utils.Subscribe(c, utils.PublishKey)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//var tm = time.Now().Format("2006-01-02 15:04:05")
		//var m = fmt.Sprintf("[ws][%s]:[%s]", tm, msg)
		//err = ws.WriteMessage(1, []byte(m))
		//if err != nil {
		//	fmt.Println(err)
		//}
	}
}

func SendUserMsg(c *gin.Context) {
	fmt.Println("service/SendUserMsg context: ", c)
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	var userId, _ = strconv.Atoi(c.Request.FormValue("userId"))
	var users = models.SearchFriend(uint(userId))

	utils.RespOKList(c.Writer, users, len(users))
}

func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	// name := c.Request.FormValue("name")
	data := models.FindByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
} //dao层之前已写好

func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetId, _ := strconv.Atoi(c.Request.FormValue("targetId"))
	code, msg := models.AddFriend(uint(userId), uint(targetId))
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	// name := c.Request.FormValue("name")
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}
