package models

import (
	"time"

	"gorm.io/gorm"
)

type TotalStudyData struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id" gorm:"uniqueIndex"`
	StudyTime int            `json:"study_time"`
	Tomatoes  int            `json:"tomatoes"`
}

type DailyStudyData struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id" gorm:"uniqueIndex:idx_user_date"` //联合索引,保证每天一条记录
	Date      time.Time      `json:"date" gorm:"uniqueIndex:idx_user_date"`
	StudyTime int            `json:"study_time"`
	Tomatoes  int            `json:"tomatoes"`
}

type MonthlyStudyData struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id" gorm:"uniqueIndex:idx_user_month"` //联合索引,保证每月一条记录
	Month     time.Time      `json:"month" gorm:"uniqueIndex:idx_user_month"`
	StudyTime int            `json:"study_time"`
	Tomatoes  int            `json:"tomatoes"`
}

type Note struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id"`
	Title     string         `json:"title"`
	Date      time.Time      `json:"date"`
	Content   string         `json:"content"`
}

type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string    `gorm:"type:varchar(50);not null" json:"username"`
	Email      string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Telenum    string    `gorm:"type:varchar(20)" json:"telenum"`
	Password   string    `gorm:"type:varchar(255);not null" json:"-"` // 密码哈希，不返回给前端
	Gender     string    `gorm:"type:varchar(10);default:'草履虫'" json:"gender"`
	Experience int       `gorm:"default:0" json:"experience"` // 经验值
	Level      int       `gorm:"default:1" json:"level"`      // 等级
	Avatar     string    `gorm:"type:varchar(500);default:'default-avatar.png'" json:"avatar_path"`
	IsAdmin    bool      `gorm:"default:false" json:"is_admin"` // 是否为管理员
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 扩展字段（根据需求添加）
	LastLoginAt *time.Time `json:"last_login_at"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	// Settings    string     `gorm:"type:json" json:"settings"` // 用户设置，JSON格式存储
}

// Todo 待办事项表
type Todo struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`           // 关联用户ID
	Event     string    `gorm:"type:varchar(255);not null" json:"event"` // 事件
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`        // 创建时间
	User      User      `gorm:"foreignKey:UserID" json:"-"`
}

// Music 音乐
type Music struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Author     string    `gorm:"type:varchar(255);not null" json:"author"` // 音乐作者
	Title      string    `gorm:"type:varchar(255);not null" json:"title"`  // 音乐标题
	Duration   int       `gorm:"not null" json:"duration"`                 // 音乐时长，单位为秒
	FileURL    string    `gorm:"type:varchar(500);not null" json:"url"`    // 音乐URL
	UploaderID uint      `json:"uploader_id"`                              // 上传者 ID, 0 表示系统音乐
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`         // 创建时间
}

// TokenBlacklist 黑名单表，用于存储被禁用的 JWT 令牌
type TokenBlacklist struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Jti       string         `json:"jti" gorm:"type:varchar(255);uniqueIndex"`
	ExpiresAt time.Time      `json:"expires_at"`
}

// AmbientSound 环境音效表
type AmbientSound struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null"`
	URL       string         `json:"url" gorm:"type:varchar(500);not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// CalendarEvent 日历事件表
type CalendarEvent struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"index:idx_user_date"`
	Title     string         `json:"title" gorm:"type:varchar(255);not null"`
	Gificon   string         `json:"gificon" gorm:"type:varchar(255)"`
	Date      time.Time      `json:"date" gorm:"type:date;index:idx_user_date"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Gif Gif表，用于存储GIF图标
type Gif struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null;uniqueIndex"`
	URL       string         `json:"url" gorm:"type:varchar(500);not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
