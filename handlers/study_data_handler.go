package handler

import (
	"2026-FM247-BackEnd/service"
	"2026-FM247-BackEnd/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type StudyDataService interface {
	AddStudyData(userID uint, date time.Time, studyTime int, tomatoes int) (bool, int, int, error)
	GetDailyStudyData(userID uint, date time.Time) (service.DailyStudyDataInfo, string)
	GetMonthlyStudyData(userID uint, date time.Time) (service.MonthlyStudyDataInfo, string)
	GetTotalStudyData(userID uint) (service.TotalStudyDataInfo, string)
	GetWeekStudyData(userID uint, date time.Time) ([]service.DailyStudyDataInfo, string)
	GetMonthStudyData(userID uint, date time.Time) ([]service.DailyStudyDataInfo, string)
	GetYearStudyData(userID uint, date time.Time) ([]service.MonthlyStudyDataInfo, string)
}

type StudyDataHandler struct {
	service StudyDataService
}

func NewStudyDataHandler(service StudyDataService) *StudyDataHandler {
	return &StudyDataHandler{service: service}
}

// AddStudyData 增加学习数据
// @Router /api/studydata [post]
func (h *StudyDataHandler) AddStudyData(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	// 2. 绑定请求参数
	var req AddStudyDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		FailWithMessage(c, "请求参数有误: "+err.Error())
		return
	}
	// 3. 增加学习数据
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(loc)
	isLevelUp, nowLevel, expToNextLevel, err := h.service.AddStudyData(claims.UserID, t, req.StudyTime, req.Tomatoes)
	if err != nil {
		FailWithMessage(c, err.Error())
		return
	}
	if isLevelUp {
		Ok(c,
			"学习数据记录成功，恭喜你升级到"+utils.IntToString(nowLevel)+"！距离下一级还差"+utils.IntToString(expToNextLevel),
			gin.H{"level_up": true, "new_level": nowLevel, "exp_to_next_level": expToNextLevel},
		)
		return
	}
	// 4. 返回结果
	Ok(c, "学习数据记录成功", gin.H{"level_up": false, "current_level": nowLevel, "exp_to_next_level": expToNextLevel})
}

// GetDailyStudyData 获取今日学习数据
// @Router /api/studydata/daily/:date [get]
func (h *StudyDataHandler) GetDailyStudyData(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	// 2. 获取每日学习数据
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(loc)
	data, msg := h.service.GetDailyStudyData(claims.UserID, t)
	if msg != "" {
		FailWithMessage(c, "获取每日学习数据失败: "+msg)
		return
	}
	// 3. 返回结果
	OkWithData(c, data)
}

// GetTotalStudyData 获取总学习数据
// @Router /api/studydata/total [get]
func (h *StudyDataHandler) GetTotalStudyData(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	// 2. 获取总学习数据
	data, msg := h.service.GetTotalStudyData(claims.UserID)
	if msg != "" {
		FailWithMessage(c, "获取总学习数据失败: "+msg)
		return
	}
	// 3. 返回结果
	OkWithData(c, data)
}

// GetWeekStudyData 获取本周每日学习数据
// @Router /api/studydata/weekly/:date [get]
func (h *StudyDataHandler) GetWeekStudyData(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	// 2. 获取本周学习数据
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(loc)
	data, msg := h.service.GetWeekStudyData(claims.UserID, t)
	if msg != "" {
		FailWithMessage(c, "获取本周学习数据失败: "+msg)
		return
	}
	// 3. 返回结果
	OkWithData(c, data)
}

// GetMonthlyStudyData 获取本月每日学习数据
// @Router /api/studydata/monthly/:date [get]
func (h *StudyDataHandler) GetMonthlyStudyData(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	// 2. 获取每月学习数据
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(loc)
	data, msg := h.service.GetMonthStudyData(claims.UserID, t)
	if msg != "" {
		FailWithMessage(c, "获取每月学习数据失败: "+msg)
		return
	}
	// 3. 返回结果
	OkWithData(c, data)
}

// GetMonthStudyData 获取本年每月学习数据
// @Router /api/studydata/yearly/:date [get]
func (h *StudyDataHandler) GetYearStudyData(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}
	// 2. 获取本年学习数据
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(loc)
	data, msg := h.service.GetYearStudyData(claims.UserID, t)
	if msg != "" {
		FailWithMessage(c, "获取本年学习数据失败: "+msg)
		return
	}
	// 3. 返回结果
	OkWithData(c, data)
}
