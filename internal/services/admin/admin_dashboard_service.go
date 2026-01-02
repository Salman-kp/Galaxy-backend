package admin

import (
	"time"

	"event-management-backend/internal/config"
	"event-management-backend/internal/domain/models"
)

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

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

//
// ---------------- SUMMARY ----------------
//
func (s *DashboardService) GetSummary() (*DashboardSummary, error) {
	var summary DashboardSummary
	db := config.DB

	if err := db.Model(&models.Event{}).
		Where("deleted_at IS NULL").
		Count(&summary.TotalEvents).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.Event{}).
		Where("status = ? AND deleted_at IS NULL", models.EventStatusCompleted).
		Count(&summary.CompletedEvents).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.Event{}).
		Where("status = ? AND deleted_at IS NULL", models.EventStatusOngoing).
		Count(&summary.OngoingEvents).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.Event{}).
		Where("status = ? AND deleted_at IS NULL", models.EventStatusUpcoming).
		Count(&summary.UpcomingEvents).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&models.User{}).
		Where("deleted_at IS NULL").
		Count(&summary.TotalUsers).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}

//
// ---------------- MONTHLY CHART ----------------
//
func (s *DashboardService) GetMonthlyEventChart(year int) ([]MonthlyEventCount, error) {
	var result []MonthlyEventCount

	err := config.DB.
		Table("events").
		Select(`
			TO_CHAR(date, 'YYYY-MM') AS month,
			COUNT(*) AS count
		`).
		Where(`
			EXTRACT(YEAR FROM date) = ?
			AND status != ?
			AND deleted_at IS NULL
		`, year, models.EventStatusCancelled).
		Group("month").
		Order("month ASC").
		Scan(&result).Error

	return result, err
}

//
// ---------------- DAILY CHART ----------------
//
func (s *DashboardService) GetDailyEventChart(year int, month int) ([]DailyEventCount, error) {
	var result []DailyEventCount

	err := config.DB.
		Table("events").
		Select(`
			TO_CHAR(date, 'YYYY-MM-DD') AS date,
			COUNT(*) AS count
		`).
		Where(`
			EXTRACT(YEAR FROM date) = ?
			AND EXTRACT(MONTH FROM date) = ?
			AND status != ?
			AND deleted_at IS NULL
		`, year, month, models.EventStatusCancelled).
		Group("date").
		Order("date ASC").
		Scan(&result).Error

	return result, err
}

//
// ---------------- TODAY EVENTS ----------------
//
func (s *DashboardService) GetTodayEvents() ([]models.Event, error) {
	var events []models.Event

	start := time.Now().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	err := config.DB.
		Model(&models.Event{}).
		Where(`
			date >= ? AND date < ?
			AND status != ?
			AND deleted_at IS NULL
		`, start, end, models.EventStatusCancelled).
		Order("reporting_time ASC").
		Find(&events).Error

	return events, err
}