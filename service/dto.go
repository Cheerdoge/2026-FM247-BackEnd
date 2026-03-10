package service

import "time"

// 用户信息dto
type UserInfo struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Gender     string    `json:"gender"`
	Telenum    string    `json:"telenum"`
	Avatar     string    `json:"avatarpath"`
	Experience int       `json:"experience"`
	Level      int       `json:"level"`
	CreatedAt  time.Time `json:"createdat"`
}

// 待办事项dto
type TodoInfo struct {
	ID    uint   `json:"id"`
	Event string `json:"event"`
}

// 每日学习数据dto
type DailyStudyDataInfo struct {
	Date      time.Time `json:"date"`
	StudyTime int       `json:"studytime"`
	Tomatoes  int       `json:"tomatoes"`
}

// 每月学习数据dto
type MonthlyStudyDataInfo struct {
	Date      time.Time `json:"date"`
	StudyTime int       `json:"studytime"`
	Tomatoes  int       `json:"tomatoes"`
}

// 总学习数据dto
type TotalStudyDataInfo struct {
	StudyTime int `json:"studytime"`
	Tomatoes  int `json:"tomatoes"`
}

// 音乐信息dto
type MusicInfo struct {
	Author   string `json:"author"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	FileURL  string `json:"url"`
}
