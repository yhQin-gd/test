package user_r

import (
	"errors"
	"fmt"
	"gocode/backend/backend/text-to-picture/models/image"
	 models "gocode/backend/backend/text-to-picture/models/user"
	"regexp"

	"gorm.io/gorm"
)

// 正则表达式验证邮箱格式
func isValidEmail(email string) bool {

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// 向用户登录表插入数据
func InsertUserLogin(db *gorm.DB, user *models.UserLogin) error {
	if user.UserName == "" {
		return fmt.Errorf("名字为空")
	}
	if user.Email == "" {
		return fmt.Errorf("邮箱为空")
	}
	if user.ID == 0 {
		return fmt.Errorf("id为空")
	}
	if len(user.Password) < 6 {
		return fmt.Errorf("密码少于6位")
	}
	if isValidEmail(user.Email) == false {
		return fmt.Errorf("邮箱格式不正确")
	}
	var existingUserLogin models.UserLogin

	result := db.Where("UserName = ?", user.UserName).First(&existingUserLogin)
	if result.Error == nil {
		return fmt.Errorf("用户名已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询用户名时发生错误: %v", result.Error)
	}

	result = db.Where("Email = ?", user.Email).First(&existingUserLogin)
	if result.Error == nil {
		return fmt.Errorf("邮箱已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询邮箱时发生错误: %v", result.Error)
	}

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("插入用户登录表失败: %v", err)
	}

	return nil
}

// 向用户查询表插入数据
func InsertUserQuery(db *gorm.DB, user *image.UserQuery) error {
	if err := InsertUserLogin(db, &user.User); err != nil {
		return err
	}
	if user.Params == "" {
		return fmt.Errorf("参数为空")
	}
	if user.Result == "" {
		return fmt.Errorf("结果为空")
	}
	if user.Time == "" {
		return fmt.Errorf("时间参数为空")
	}

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("插入用户查询表失败: %v", err)
	}

	return nil

}

// 向用户收藏表插入数据
func InsertFavoritedImage(db *gorm.DB, user *image.FavoritedImage) error {
	if err := InsertUserLogin(db, &user.User); err != nil {
		return err
	}
	if user.Result == "" {
		return fmt.Errorf("结果为空")
	}

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("插入用户收藏表失败: %v", err)
	}
	return nil
}
