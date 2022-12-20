package models

import (
	"GINCHAT/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string
	Email         string
	Identify      string
	ClientIp      string
	ClientPort    string
	LoginTime     time.Time
	LogoutTime    time.Time
	HeartBeatTime time.Time
	IsLogout      bool
	DeviceInfo    string
}

func (self *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	var data = make([]*UserBasic, 10)
	utils.Db.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}
