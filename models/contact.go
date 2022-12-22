package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int // 关系类型
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}
