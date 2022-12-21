package main

import (
	"GINCHAT/router"
	"GINCHAT/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()

	var r = router.Router()
	err := r.Run(":8081")
	if err != nil {
		panic("无法在指定端口启动服务。")
	}
}
