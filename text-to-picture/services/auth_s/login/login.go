package login

import (
	"errors"
	middlewire "gocode/backend/backend/text-to-picture/middlewire/jwt"
	models "gocode/backend/backend/text-to-picture/models/init"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID         int
	Name       string
	Email      string
	Password   string
	IsVerified bool
}

// 注册
// 注册
func Register(c *gin.Context) {
	// 定义用于接收 JSON 数据的结构体
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 解析 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		return
	}

	// 数据验证
	if len(input.Name) == 0 || len(input.Email) == 0 || len(input.Password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "参数不完整或密码长度不足"})
		return
	}

	// 检查邮箱是否已存在
	var user User
	result := models.DB.Where("email = ?", input.Email).First(&user)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "用户邮箱已存在"})
		return
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库查询错误"})
		return
	}

	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "密码加密失败"})
		return
	}

	// 创建用户
	newUser := User{
		Name:       input.Name,
		Email:      input.Email,
		Password:   string(hashedPassword),
		IsVerified: true,
	}
	if err := models.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "用户创建失败"})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
	})
}

// 登录
func Login(c *gin.Context) {
	// 定义用于接收 JSON 数据的结构体
	var input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// 解析 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		return
	}

	// 查找用户
	var user User
	result := models.DB.Where("name = ?", input.Name).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库查询错误"})
		return
	}

	// 验证密码将用户的明文密码与数据库中的哈希密码对比
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "密码错误"})
		return
	}

	// 生成 JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middlewire.Claims{
		Username: input.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middlewire.JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成 token 错误"})
		return
	}

	// 登录成功
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   tokenString,
	})
}
