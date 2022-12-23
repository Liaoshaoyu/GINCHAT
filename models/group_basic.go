package models

import (
	"GINCHAT/utils"
	"fmt"
	"gorm.io/gorm"
)

type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    uint
	Type    int
	Desc    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    uint
	Type    int
	Desc    string
}

func (table *Community) TableName() string {
	return "community"
}

func LoadCommunity(ownerId uint) ([]*Community, string) {
	data := make([]*Community, 10)
	utils.Db.Where("owner_id = ? ", ownerId).Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	//utils.DB.Where()
	return data, "查询成功"
}
