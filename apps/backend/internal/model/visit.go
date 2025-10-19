package model

import (
	"time"

	"github.com/google/uuid"
)

type Visit struct {
	Base
	ScheduleID      uuid.UUID `json:"scheduleId" db:"schedule_id"`
	StartTime       time.Time `json:"startTime" db:"start_time"`
	EndTime         *time.Time `json:"endTime" db:"end_time"`
	StartLatitude   float64   `json:"startLatitude" db:"start_latitude"`
	StartLongitude  float64   `json:"startLongitude" db:"start_longitude"`
	EndLatitude     *float64  `json:"endLatitude" db:"end_latitude"`
	EndLongitude    *float64  `json:"endLongitude" db:"end_longitude"`
	Status          string    `json:"status" db:"status"`
	DurationMinutes *int      `json:"durationMinutes" db:"duration_minutes"`
}

type VisitCreate struct {
	ScheduleID    uuid.UUID `json:"scheduleId" db:"schedule_id"`
	StartTime     time.Time `json:"startTime" db:"start_time"`
	StartLat      float64   `json:"startLat" db:"start_latitude"`
	StartLong     float64   `json:"startLong" db:"start_longitude"`
	Status        string    `json:"status" db:"status"`
}

type VisitUpdate struct {
	EndTime     *time.Time `json:"endTime" db:"end_time"`
	EndLat      *float64   `json:"endLat" db:"end_latitude"`
	EndLong     *float64   `json:"endLong" db:"end_longitude"`
	Status      string    `json:"status" db:"status"`
}

func (v *Visit) TableName() string {
	return "visits"
}

func (v *Visit) CalculateDuration() {
	if v.EndTime != nil && !v.StartTime.IsZero() {
		duration := v.EndTime.Sub(v.StartTime)
		minutes := int(duration.Minutes())
		v.DurationMinutes = &minutes
	}
}

func (v *Visit) GetDurationMinutes() int {
	if v.DurationMinutes != nil {
		return *v.DurationMinutes
	}
	return 0
}