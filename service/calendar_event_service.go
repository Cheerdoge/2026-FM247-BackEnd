package service

import (
	"2026-FM247-BackEnd/models"
	"errors"
	"time"
)

type CalendarEventRepository interface {
	Create(event *models.CalendarEvent) error
	Update(id uint, event map[string]interface{}) error
	Delete(id uint) error
	GetEventByDate(dateStr string, userid uint) ([]models.CalendarEvent, error)
	GetMonthEventsByUserID(userid uint, startdate time.Time, enddate time.Time) ([]string, error)
	GetUserID(id uint) (uint, error)
}

type CalendarEventService struct {
	repo CalendarEventRepository
}

func NewCalendarEventService(calendarEventRepo CalendarEventRepository) *CalendarEventService {
	return &CalendarEventService{repo: calendarEventRepo}
}

func (s *CalendarEventService) Create(userid uint, title string, date time.Time, gificon string) error {
	event := &models.CalendarEvent{
		UserID:  userid,
		Title:   title,
		Date:    date,
		Gificon: gificon,
	}
	return s.repo.Create(event)
}

func (s *CalendarEventService) Update(userid uint, id uint, title string, date time.Time, gificon string) error {
	event := make(map[string]interface{})
	if title != "" {
		event["title"] = title
	}
	if !date.IsZero() {
		event["date"] = date
	}
	if gificon != "" {
		event["gificon"] = gificon
	}
	eventUserID, err := s.repo.GetUserID(id)
	if err != nil {
		return err
	}
	if eventUserID != userid {
		return errors.New("无权限修改该事件")
	}
	return s.repo.Update(id, event)
}

func (s *CalendarEventService) Delete(userid uint, id uint) error {
	eventUserID, err := s.repo.GetUserID(id)
	if err != nil {
		return err
	}
	if eventUserID != userid {
		return errors.New("无权限删除该事件")
	}
	return s.repo.Delete(id)
}

func (s *CalendarEventService) GetMonthEventsByUserID(userid uint, year int, month int) ([]string, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)
	eventDates, err := s.repo.GetMonthEventsByUserID(userid, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return eventDates, nil
}

func (s *CalendarEventService) GetEventByDate(dateStr string, userid uint) ([]CalendarEventInfo, error) {
	events, err := s.repo.GetEventByDate(dateStr, userid)
	if err != nil {
		return []CalendarEventInfo{}, err
	}
	var eventInfos []CalendarEventInfo
	for _, event := range events {
		eventInfos = append(eventInfos, CalendarEventInfo{
			ID:      event.ID,
			Title:   event.Title,
			Date:    event.Date,
			Gificon: event.Gificon,
		})
	}
	return eventInfos, nil
}
