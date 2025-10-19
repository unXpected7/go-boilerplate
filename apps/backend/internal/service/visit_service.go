package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/model"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/repository"
)

type VisitService struct {
	visitRepo  *repository.VisitRepository
	scheduleRepo *repository.ScheduleRepository
}

func NewVisitService(visitRepo *repository.VisitRepository, scheduleRepo *repository.ScheduleRepository) *VisitService {
	return &VisitService{
		visitRepo:    visitRepo,
		scheduleRepo: scheduleRepo,
	}
}

// Start a visit for a schedule
func (v *VisitService) StartVisit(ctx context.Context, scheduleID uuid.UUID, startTime time.Time, startLat, startLong float64) (*model.Visit, error) {
	// Validate coordinates
	if !isValidCoordinates(startLat, startLong) {
		return nil, fmt.Errorf("invalid geolocation coordinates")
	}

	// Check if schedule exists
	_, err := v.scheduleRepo.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	// Check if visit already exists
	exists, err := v.visitRepo.VisitExistsForSchedule(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to check visit existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("visit already started for this schedule")
	}

	// Validate time is not in the past
	if startTime.Before(time.Now().Add(-5 * time.Minute)) {
		return nil, fmt.Errorf("start time cannot be in the past")
	}

	// Create visit
	visit, err := v.visitRepo.StartVisit(ctx, scheduleID, startTime, startLat, startLong)
	if err != nil {
		return nil, fmt.Errorf("failed to start visit: %w", err)
	}

	// Update schedule status to in_progress
	if err := v.scheduleRepo.UpdateScheduleStatus(ctx, scheduleID, "in_progress"); err != nil {
		return nil, fmt.Errorf("failed to update schedule status: %w", err)
	}

	return visit, nil
}

// End a visit
func (v *VisitService) EndVisit(ctx context.Context, scheduleID uuid.UUID, endTime time.Time, endLat, endLong float64) (*model.Visit, error) {
	// Validate coordinates
	if !isValidCoordinates(endLat, endLong) {
		return nil, fmt.Errorf("invalid geolocation coordinates")
	}

	// Get existing visit
	visit, err := v.visitRepo.GetVisitByScheduleID(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get visit: %w", err)
	}

	// Check if visit is already completed
	if visit.Status == "completed" {
		return nil, fmt.Errorf("visit is already completed")
	}

	// Validate end time is after start time
	if endTime.Before(visit.StartTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	// Validate end time is not too far in the future
	if endTime.After(time.Now().Add(1 * time.Hour)) {
		return nil, fmt.Errorf("end time cannot be more than 1 hour in the future")
	}

	// End visit
	updatedVisit, err := v.visitRepo.EndVisit(ctx, visit.ID, endTime, endLat, endLong)
	if err != nil {
		return nil, fmt.Errorf("failed to end visit: %w", err)
	}

	// Update schedule status to completed
	if err := v.scheduleRepo.UpdateScheduleStatus(ctx, scheduleID, "completed"); err != nil {
		return nil, fmt.Errorf("failed to update schedule status: %w", err)
	}

	return updatedVisit, nil
}

// Get visit by ID
func (v *VisitService) GetVisitByID(ctx context.Context, visitID uuid.UUID) (*model.Visit, error) {
	return v.visitRepo.GetVisitByID(ctx, visitID)
}

// Get visit by schedule ID
func (v *VisitService) GetVisitByScheduleID(ctx context.Context, scheduleID uuid.UUID) (*model.Visit, error) {
	return v.visitRepo.GetVisitByScheduleID(ctx, scheduleID)
}

// Update visit status
func (v *VisitService) UpdateVisitStatus(ctx context.Context, visitID uuid.UUID, status string) error {
	// Validate status
	if !isValidVisitStatus(status) {
		return fmt.Errorf("invalid visit status: %s", status)
	}

	return v.visitRepo.UpdateVisitStatus(ctx, visitID, status)
}

// Get visit statistics
func (v *VisitService) GetVisitStats(ctx context.Context) (*model.TaskStats, error) {
	return v.visitRepo.GetVisitStats(ctx)
}

// Get visit duration statistics
func (v *VisitService) GetVisitDurationStats(ctx context.Context) (map[string]interface{}, error) {
	return v.visitRepo.GetVisitDurationStats(ctx)
}

// Get visits by status
func (v *VisitService) GetVisitsByStatus(ctx context.Context, status string) ([]model.Visit, error) {
	return v.visitRepo.GetVisitsByStatus(ctx, status)
}

// Get active visits (in_progress)
func (v *VisitService) GetActiveVisits(ctx context.Context) ([]model.Visit, error) {
	return v.GetVisitsByStatus(ctx, "in_progress")
}

// Validate that a visit can be started for a schedule
func (v *VisitService) ValidateStartVisit(ctx context.Context, scheduleID uuid.UUID) error {
	// Check if schedule exists
	schedule, err := v.scheduleRepo.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return fmt.Errorf("failed to get schedule: %w", err)
	}

	// Check if visit already exists
	exists, err := v.visitRepo.VisitExistsForSchedule(ctx, scheduleID)
	if err != nil {
		return fmt.Errorf("failed to check visit existence: %w", err)
	}

	if exists {
		return fmt.Errorf("visit already started for this schedule")
	}

	// Check if schedule is in a valid status to start visit
	if schedule.Status == "completed" {
		return fmt.Errorf("schedule is already completed")
	}

	return nil
}

// Calculate visit duration and return formatted string
func (v *VisitService) CalculateVisitDuration(visit *model.Visit) string {
	if visit.EndTime == nil {
		return "In progress"
	}

	duration := visit.EndTime.Sub(visit.StartTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// Get average visit duration
func (v *VisitService) GetAverageVisitDuration(ctx context.Context) (time.Duration, error) {
	stats, err := v.GetVisitDurationStats(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get duration stats: %w", err)
	}

	if avgDuration, ok := stats["avg_duration"].(float64); ok {
		return time.Duration(avgDuration) * time.Minute, nil
	}

	return 0, fmt.Errorf("average duration not available")
}

// Validate visit data
func (v *VisitService) ValidateVisitData(scheduleID uuid.UUID, startLat, startLong, endLat, endLong float64, startTime, endTime *time.Time) error {
	// Validate coordinates
	if !isValidCoordinates(startLat, startLong) {
		return fmt.Errorf("invalid start coordinates")
	}

	if endTime != nil && !isValidCoordinates(endLat, endLong) {
		return fmt.Errorf("invalid end coordinates")
	}

	// Validate time sequence
	if startTime != nil && endTime != nil && endTime.Before(*startTime) {
		return fmt.Errorf("end time must be after start time")
	}

	return nil
}

// Helper function to check if coordinates are valid
func isValidCoordinates(lat, long float64) bool {
	return lat >= -90 && lat <= 90 && long >= -180 && long <= 180
}

// Helper function to check if visit status is valid
func isValidVisitStatus(status string) bool {
	validStatuses := map[string]bool{
		"not_started":  true,
		"in_progress":  true,
		"completed":    true,
	}
	return validStatuses[status]
}

// Get visit summary with calculated duration
func (v *VisitService) GetVisitSummary(ctx context.Context, scheduleID uuid.UUID) (*model.Visit, error) {
	visit, err := v.visitRepo.GetVisitByScheduleID(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get visit: %w", err)
	}

	// Calculate duration if not already calculated
	visit.CalculateDuration()

	return visit, nil
}