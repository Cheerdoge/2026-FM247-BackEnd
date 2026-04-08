package repository

import (
	"2026-FM247-BackEnd/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CalendarEventRepository struct {
	db *gorm.DB
}

func NewCalendarEventRepository(db *gorm.DB) *CalendarEventRepository {
	return &CalendarEventRepository{db: db}
}

func (r *CalendarEventRepository) Create(event *models.CalendarEvent) error {
	return r.db.Create(event).Error
}

func (r *CalendarEventRepository) Update(id uint, event map[string]interface{}) error {
	err := r.db.Model(&models.CalendarEvent{}).Where("id = ?", id).Updates(event).Error
	return err
}

func (r *CalendarEventRepository) Delete(id uint) error {
	return r.db.Delete(&models.CalendarEvent{}, id).Error
}

func (r *CalendarEventRepository) GetEventByDate(dateStr string, userid uint) ([]models.CalendarEvent, error) {
	var events []models.CalendarEvent
	err := r.db.Where("date = ? AND user_id = ?", dateStr, userid).Find(&events).Error
	if err != nil {
		return []models.CalendarEvent{}, err
	}
	if len(events) == 0 {
		return []models.CalendarEvent{}, errors.New("今日无事件")
	}
	return events, nil
}

func (r *CalendarEventRepository) GetMonthEventsByUserID(userid uint, startdate time.Time, enddate time.Time) ([]string, error) {
	var eventDates []string
	err := r.db.Model(&models.CalendarEvent{}).Where("user_id = ? AND date >= ? AND date <= ?", userid, startdate, enddate).
		Distinct("date").Pluck("DATE_FORMAT(date, '%Y-%m-%d')", &eventDates).Error
	return eventDates, err
}

func (r *CalendarEventRepository) GetUserID(id uint) (uint, error) {
	var event models.CalendarEvent
	err := r.db.First(&event, id).Error
	if err != nil {
		return 0, err
	}
	return event.UserID, nil
}
