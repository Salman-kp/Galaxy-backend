package admin

import (
	"time"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/models"
)

// DASHBOARD RESPONSE STRUCTS

type DashboardSummary struct {
	TotalEvents     int64 `json:"total_events"`
	CompletedEvents int64 `json:"completed_events"`
	OngoingEvents   int64 `json:"ongoing_events"`
	UpcomingEvents  int64 `json:"upcoming_events"`
	TotalUsers      int64 `json:"total_users"`
}

type MonthlyEventCount struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type DailyEventCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type TodayEvent struct {
	ID            uint   `json:"id"`
	EventName     string `json:"event_name"`
	Date          string `json:"date"`
	TimeSlot      string `json:"time_slot"`
	ReportingTime string `json:"reporting_time"`
	Status        string `json:"status"`
}


type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}


// 3.1 TOP 5 BOXES

func (s *DashboardService) GetSummary() (*DashboardSummary, error) {
	var summary DashboardSummary

	db := config.DB

	db.Model(&models.Event{}).Count(&summary.TotalEvents)

	db.Model(&models.Event{}).
		Where("status = ?", models.EventStatusCompleted).
		Count(&summary.CompletedEvents)

	db.Model(&models.Event{}).
		Where("status = ?", models.EventStatusOngoing).
		Count(&summary.OngoingEvents)

	db.Model(&models.Event{}).
		Where("status = ?", models.EventStatusUpcoming).
		Count(&summary.UpcomingEvents)

	db.Model(&models.User{}).Count(&summary.TotalUsers)

	return &summary, nil
}


// 3.2 MONTH-BASED CHART

func (s *DashboardService) GetMonthlyEventChart(year int) ([]MonthlyEventCount, error) {
	var result []MonthlyEventCount

	err := config.DB.
		Table("events").
		Select(`
			TO_CHAR(date, 'YYYY-MM') as month,
			COUNT(*) as count
		`).
		Where("EXTRACT(YEAR FROM date) = ?", year).
		Group("month").
		Order("month ASC").
		Scan(&result).Error

	return result, err
}


// (Admin clicks month → dates update)

func (s *DashboardService) GetDailyEventChart(year int, month int) ([]DailyEventCount, error) {
	var result []DailyEventCount

	err := config.DB.
		Table("events").
		Select(`
			TO_CHAR(date, 'YYYY-MM-DD') as date,
			COUNT(*) as count
		`).
		Where(`
			EXTRACT(YEAR FROM date) = ?
			AND EXTRACT(MONTH FROM date) = ?
		`, year, month).
		Group("date").
		Order("date ASC").
		Scan(&result).Error

	return result, err
}


// 3.3 TODAY'S WORK LIST

func (s *DashboardService) GetTodayEvents() ([]models.Event, error) {
	var events []models.Event

	today := time.Now().Truncate(24 * time.Hour)

	err := config.DB.
		Model(&models.Event{}).
		Where("date = ?", today).
		Order("date ASC, reporting_time ASC").
		Find(&events).Error

	return events, err
}