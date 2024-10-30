package image

import "gocode/backend/backend/text-to-picture/models/user"

//用户历史查询
type UserQuery struct {
	ID       int            `json:"id" gorm:"primarykey"`
	UserName string            `json:"user_name" gorm:"not null"`
	Params   string         `json:"params"`
	Result   string         `json:"result"`
	Time     string         `json:"time"`
	User     user.UserLogin `gorm:"foreignKey:UserName;references:UserName"`
}

//用户收藏查询
type FavoritedImage struct {
	ID       int            `json:"id" gorm:"primarykey"`
	UserName string         `json:"user_name" gorm:"not null"`
	Result   string         `json:"result"`
	User     user.UserLogin `gorm:"foreignKey:UserName;references:UserName"`
}
