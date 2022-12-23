package models

import (
	"GINCHAT/utils"
	"fmt"
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int // 关系类型	1-好友；2-群聊；3
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}
func SearchFriend(userId uint) []UserBasic {
	var contacts = make([]Contact, 0)
	var objIds = make([]uint64, 0)
	utils.Db.Where("owner_id = ? and type = 1", userId).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(v)
		objIds = append(objIds, uint64(v.TargetId))
	}
	var users = make([]UserBasic, 0)
	utils.Db.Where("id in ?", objIds).Find(&users)
	return users
}

func FindByID(userId uint) UserBasic {
	var user = UserBasic{}
	utils.Db.Where("id = ?", userId).First(&user)

	return user
}

//添加好友
func AddFriend(userId uint, targetId uint) (int, string) {
	user := UserBasic{}
	if targetId != 0 {
		user = FindByID(targetId)
		fmt.Println(targetId, " ", userId)
		if user.Salt != "" {
			if userId == user.ID {
				return -1, "不能加自己"
			}
			contact0 := Contact{}
			utils.Db.Where("owner_id =? and target_id =? and type=1",
				userId, targetId).Find(&contact0)
			if contact0.ID != 0 {
				return -1, "不能重复添加"
			}
			tx := utils.Db.Begin()
			//事务一旦开始，不论什么异常最终都会Rollback
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()
			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = targetId
			contact.Type = 1
			if err := utils.Db.Create(&contact).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			contact1 := Contact{}
			contact1.OwnerId = targetId
			//群管理
			//新建群
			contact1.TargetId = userId
			contact1.Type = 1
			if err := utils.Db.Create(&contact1).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			tx.Commit()
			return 0, "添加好友成功"
		}
		return -1, "没有找到此用户"
	}
	return -1, "好友ID不能为空"
}
