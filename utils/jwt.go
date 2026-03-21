package utils

import (
	"2026-FM247-BackEnd/config"
	"2026-FM247-BackEnd/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Claims 结构体定义JWT的声明（payload）部分
// 继承jwt.StandardClaims，包含了JWT标准声明如过期时间、签发时间等
// UserID: 用户ID，用于标识用户身份
// Username: 用户名，用于标识用户身份
type Claims struct {
	UserID  uint   `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	Jti     string `json:"jti"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT令牌的函数
// 参数: user - 用户模型指针，包含用户ID和角色信息
// 返回值:
//
//	string - 生成的JWT令牌字符串
//	error - 错误信息，生成成功时为nil
//
// 功能: 根据用户信息生成一个有过期时间的JWT令牌
func GenerateToken(user *models.User) (string, error) {
	// 计算令牌的过期时间：当前时间 + 配置中指定的过期时长
	expirationTime := time.Now().Add(config.AppConfig.JWTExpire)

	// 创建Claims声明对象，包含自定义声明和标准声明
	claims := &Claims{
		UserID:  user.ID,          // 设置用户ID
		IsAdmin: user.IsAdmin,     // 设置是否为管理员
		Jti:     uuid.NewString(), // 设置唯一的JWT ID
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // 设置过期时间（Unix时间戳）
			IssuedAt:  time.Now().Unix(),     // 设置签发时间（Unix时间戳）
		},
	}

	// 使用HS256签名方法和自定义声明创建JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用配置中的密钥对令牌进行签名，返回签名后的令牌字符串
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// ValidateToken 验证JWT令牌的函数
// 参数: tokenString - 需要验证的JWT令牌字符串
// 返回值:
//
//	*Claims - 解析成功后的声明对象指针
//	error - 错误信息，验证成功时为nil
//
// 功能: 验证JWT令牌的有效性，包括签名验证和过期时间检查
func ValidateToken(tokenString string) (*Claims, error) {
	// 检查token的段数
	segments := strings.Split(tokenString, ".")
	if len(segments) != 3 {
		return nil, fmt.Errorf("token contains an invalid number of segments: %d", len(segments))
	}
	// 解析并验证令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// 验证回调函数，返回用于验证签名的密钥
		// 这里使用配置中的JWT密钥
		return []byte(config.AppConfig.JWTSecret), nil
	})

	// 如果解析过程中出现错误（如令牌格式错误、签名无效、已过期等），返回错误
	if err != nil {
		return nil, err
	}

	// 类型断言，将令牌的声明部分转换为自定义的Claims类型
	// 同时检查令牌是否有效（token.Valid）
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 令牌有效，返回解析出的声明信息
		return claims, nil
	}

	// 令牌无效（可能是类型断言失败或token.Valid为false）
	// 返回签名无效的错误
	return nil, jwt.ErrSignatureInvalid
}

// GetClaimsFromContext 从上下文中获取 Claims
func GetClaimsFromContext(c *gin.Context) (*Claims, error) {
	value, exists := c.Get("claims")
	if !exists {
		return nil, errors.New("用户信息不存在")
	}

	claims, ok := value.(*Claims)
	if !ok {
		return nil, fmt.Errorf("用户信息类型错误，期望 *Claims，实际: %T", value)
	}

	return claims, nil
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) (uint, error) {
	value, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("用户ID不存在")
	}

	userID, ok := value.(uint)
	if !ok {
		return 0, fmt.Errorf("用户ID类型错误，期望 uint，实际: %T", value)
	}

	return userID, nil
}

// GetJtiFromContext 从上下文中获取JTI
func GetJtiFromContext(c *gin.Context) (string, error) {
	value, exists := c.Get("jti")
	if !exists {
		return "", errors.New("令牌ID不存在")
	}

	jti, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("令牌ID类型错误，期望 string，实际: %T", value)
	}

	return jti, nil
}
