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
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identify      string
	ClientIp      string
	ClientPort    string
	Salt          string
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

func FindUserByNameAndPassword(name, password string) UserBasic {
	var user = UserBasic{}
	utils.Db.Where("name = ? and pass_word = ?", name, password).First(&user)

	// token加密
	var str = fmt.Sprintf("%d", time.Now().Unix())
	utils.Db.Model(&user).Where("id = ?", user.ID).Update("identify", utils.MD5Encode(str))
	return user
}

func FindUserByName(name string) UserBasic {
	var user = UserBasic{}
	utils.Db.Where("name = ?", name).First(&user)
	return user
}
func FindUserByPhone(phone string) UserBasic {
	var user = UserBasic{}
	utils.Db.Where("phone = ?", phone).First(&user)
	return user
}
func FindUserByEmail(email string) UserBasic {
	var user = UserBasic{}
	utils.Db.Where("email = ?", email).First(&user)
	return user
}

func CreaterUser(user UserBasic) *gorm.DB {
	return utils.Db.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.Db.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.Db.Model(user).Updates(
		UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}
