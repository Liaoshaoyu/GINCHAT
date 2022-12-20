package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("config app initialized.")
	//fmt.Println("config mysql", viper.Get("mysql"))
}

func InitMysql() {
	Db, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{})
	fmt.Println("config mysql initialized.")
	//user := models.UserBasic{}
	//db.Find(&user)
	//fmt.Println(user)
}
