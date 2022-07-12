package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID int32 `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"colunm:add_time"`
	UpdatedAt time.Time `gorm:"colunm:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

//md5 信息摘要算法
type User struct {
	BaseModel
	Mobile string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string `gorm:"type:varchar(100);not null"`
	NickName string `gorm:"type:varchar(20);not null"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender string `gorm:"column:gender;default:male;type:varchar(6) comment 'female 表示女 male表示男'"`
	Role int `gorm:"column:role;default:1;type:int(1) comment '1 表示普通用户 2表示管理员'"`
}

