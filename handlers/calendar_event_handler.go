package handler

import (
	"2026-FM247-BackEnd/service"
	"2026-FM247-BackEnd/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type CalendarEventServiceInterface interface {
	Create(userid uint, title string, date time.Time, gificon string) error
	Update(userid uint, id uint, title string, date time.Time, gificon string) error
	Delete(userid uint, id uint) error
	GetMonthEventsByUserID(userid uint, year int, month int) ([]string, error)
	GetEventByDate(dateStr string, userid uint) ([]service.CalendarEventInfo, error)
}

type CalendarEventHandler struct {
	calendarEventService CalendarEventServiceInterface
}

func NewCalendarEventHandler(calendarEventService CalendarEventServiceInterface) *CalendarEventHandler {
	return &CalendarEventHandler{calendarEventService: calendarEventService}
}

// CreateCalendarEvent 创建日历事件
// @Router /calendar/event [post]
func (h *CalendarEventHandler) CreateCalendarEvent(c *gin.Context) {
	var req CreateCalendarEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		FailWithMessage(c, "请求参数错误: "+err.Error())
		return
	}
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}
	eventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		FailWithMessage(c, "日期格式错误: "+err.Error())
		return
	}
	err = h.calendarEventService.Create(claims.UserID, req.Title, eventDate, req.Gificon)
	if err != nil {
		FailWithMessage(c, "创建事件失败: "+err.Error())
		return
	}
	OkWithMessage(c, "操作成功")
}

// UpdateCalendarEvent 更新日历事件
// @Router /calendar/event/{id} [put]
func (h *CalendarEventHandler) UpdateCalendarEvent(c *gin.Context) {
	var req UpdateCalendarEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		FailWithMessage(c, "请求参数错误: "+err.Error())
		return
	}
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}
	id := c.Param("id")
	eventID, err := utils.StringToUint(id)
	if err != nil {
		FailWithMessage(c, "无效的事件ID: "+err.Error())
		return
	}
	eventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		FailWithMessage(c, "日期格式错误: "+err.Error())
		return
	}
	err = h.calendarEventService.Update(claims.UserID, eventID, req.Title, eventDate, req.Gificon)
	if err != nil {
		FailWithMessage(c, "更新事件失败: "+err.Error())
		return
	}
	OkWithMessage(c, "操作成功")

}

// DeleteCalendarEvent 删除日历事件
// @Router /calendar/event/{id} [delete]
func (h *CalendarEventHandler) DeleteCalendarEvent(c *gin.Context) {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}
	id := c.Param("id")
	eventID, err := utils.StringToUint(id)
	if err != nil {
		FailWithMessage(c, "无效的事件ID: "+err.Error())
		return
	}
	err = h.calendarEventService.Delete(claims.UserID, eventID)
	if err != nil {
		FailWithMessage(c, "删除事件失败: "+err.Error())
		return
	}
	OkWithMessage(c, "删除成功")
}

// GetMonthEventsByUserID 获取用户的月度日历事件标红点
// @Router /calendar/event/{year}/{month} [get]
func (h *CalendarEventHandler) GetMonthEventsByUserID(c *gin.Context) {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}
	year := c.Param("year")
	month := c.Param("month")
	yearInt, err := utils.StringToInt(year)
	if err != nil {
		FailWithMessage(c, "无效的年份: "+err.Error())
		return
	}
	monthInt, err := utils.StringToInt(month)
	if err != nil {
		FailWithMessage(c, "无效的月份: "+err.Error())
		return
	}
	events, err := h.calendarEventService.GetMonthEventsByUserID(claims.UserID, yearInt, monthInt)
	if err != nil {
		FailWithMessage(c, "获取事件失败: "+err.Error())
		return
	}
	OkWithData(c, events)
}

// GetCalendarEventByDate 获取指定日期的日历事件
// @Router /calendar/event/{date} [get]
func (h *CalendarEventHandler) GetCalendarEventByDate(c *gin.Context) {
	dateStr := c.Param("date")
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}
	events, err := h.calendarEventService.GetEventByDate(dateStr, claims.UserID)
	if err != nil {
		FailWithMessage(c, "获取事件失败: "+err.Error())
		return
	}
	OkWithData(c, events)
}
