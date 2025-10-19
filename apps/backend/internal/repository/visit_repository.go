package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/model"
)

type VisitRepository struct {
	DB *pgxpool.Pool
}

func NewVisitRepository(db *pgxpool.Pool) *VisitRepository {
	return &VisitRepository{DB: db}
}

// Get visit by ID
func (r *VisitRepository) GetVisitByID(ctx context.Context, id uuid.UUID) (*model.Visit, error) {
	query := `SELECT * FROM visits WHERE id = $1`

	var visit model.Visit
	err := r.DB.QueryRow(ctx, query, id).Scan(&visit.ID, &visit.ScheduleID, &visit.StartTime, &visit.EndTime, &visit.StartLatitude, &visit.StartLongitude, &visit.EndLatitude, &visit.EndLongitude, &visit.Status, &visit.DurationMinutes, &visit.CreatedAt, &visit.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("visit not found")
		}
		return nil, fmt.Errorf("failed to get visit: %w", err)
	}

	return &visit, nil
}

// Get visit by schedule ID
func (r *VisitRepository) GetVisitByScheduleID(ctx context.Context, scheduleID uuid.UUID) (*model.Visit, error) {
	query := `SELECT * FROM visits WHERE schedule_id = $1 ORDER BY created_at DESC LIMIT 1`

	var visit model.Visit
	err := r.DB.QueryRow(ctx, query, scheduleID).Scan(&visit.ID, &visit.ScheduleID, &visit.StartTime, &visit.EndTime, &visit.StartLatitude, &visit.StartLongitude, &visit.EndLatitude, &visit.EndLongitude, &visit.Status, &visit.DurationMinutes, &visit.CreatedAt, &visit.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no visit found for schedule")
		}
		return nil, fmt.Errorf("failed to get visit: %w", err)
	}

	return &visit, nil
}

// Create a new visit
func (r *VisitRepository) CreateVisit(ctx context.Context, visit *model.Visit) error {
	query := `
		INSERT INTO visits (id, schedule_id, start_time, start_latitude, start_longitude, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.DB.Exec(ctx, query, visit.ID, visit.ScheduleID, visit.StartTime, visit.StartLatitude, visit.StartLongitude, visit.Status)
	if err != nil {
		return fmt.Errorf("failed to create visit: %w", err)
	}

	return nil
}

// Update visit
func (r *VisitRepository) UpdateVisit(ctx context.Context, visit *model.Visit) error {
	query := `
		UPDATE visits
		SET end_time = $1, end_latitude = $2, end_longitude = $3, status = $4, duration_minutes = $5
		WHERE id = $6
	`

	_, err := r.DB.Exec(ctx, query, visit.EndTime, visit.EndLatitude, visit.EndLongitude, visit.Status, visit.DurationMinutes, visit.ID)
	if err != nil {
		return fmt.Errorf("failed to update visit: %w", err)
	}

	return nil
}

// Update visit status
func (r *VisitRepository) UpdateVisitStatus(ctx context.Context, visitID uuid.UUID, status string) error {
	query := `UPDATE visits SET status = $1 WHERE id = $2`

	_, err := r.DB.Exec(ctx, query, status, visitID)
	if err != nil {
		return fmt.Errorf("failed to update visit status: %w", err)
	}

	return nil
}

// Start visit (set start time and location)
func (r *VisitRepository) StartVisit(ctx context.Context, scheduleID uuid.UUID, startTime time.Time, startLat, startLong float64) (*model.Visit, error) {
	visit := &model.Visit{
		ScheduleID:    scheduleID,
		StartTime:     startTime,
		StartLatitude: startLat,
		StartLongitude: startLong,
		Status:        "in_progress",
	}

	query := `
		INSERT INTO visits (id, schedule_id, start_time, start_latitude, start_longitude, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	visit.ID = uuid.New()

	_, err := r.DB.Exec(ctx, query, visit.ID, visit.ScheduleID, visit.StartTime, visit.StartLatitude, visit.StartLongitude, visit.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to start visit: %w", err)
	}

	return visit, nil
}

// End visit (set end time and location)
func (r *VisitRepository) EndVisit(ctx context.Context, visitID uuid.UUID, endTime time.Time, endLat, endLong float64) (*model.Visit, error) {
	visit := &model.Visit{
		EndTime:     &endTime,
		EndLatitude: &endLat,
		EndLongitude: &endLong,
		Status:      "completed",
	}

	query := `
		UPDATE visits
		SET end_time = $1, end_latitude = $2, end_longitude = $3, status = $4
		WHERE id = $5
	`

	_, err := r.DB.Exec(ctx, query, visit.EndTime, visit.EndLatitude, visit.EndLongitude, visit.Status, visit.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to end visit: %w", err)
	}

	// Get updated visit with calculated duration
	updatedVisit, err := r.GetVisitByID(ctx, visitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated visit: %w", err)
	}

	return updatedVisit, nil
}

// Get visit statistics
func (r *VisitRepository) GetVisitStats(ctx context.Context) (*model.TaskStats, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'not_started' THEN 1 END) as not_completed
		FROM visits
	`

	var total, completed, pending, notCompleted int
	err := r.DB.QueryRow(ctx, query).Scan(&total, &completed, &pending, &notCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get visit stats: %w", err)
	}

	stats := &model.TaskStats{
		TotalTasks:      total,
		CompletedTasks:  completed,
		PendingTasks:    pending,
		NotCompletedTasks: notCompleted,
	}

	return stats, nil
}

// Get visit by status
func (r *VisitRepository) GetVisitsByStatus(ctx context.Context, status string) ([]model.Visit, error) {
	query := `SELECT * FROM visits WHERE status = $1 ORDER BY created_at DESC`

	var visits []model.Visit
	rows, err := r.DB.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get visits by status: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var visit model.Visit
		if err := rows.Scan(&visit.ID, &visit.ScheduleID, &visit.StartTime, &visit.EndTime, &visit.StartLatitude, &visit.StartLongitude, &visit.EndLatitude, &visit.EndLongitude, &visit.Status, &visit.DurationMinutes, &visit.CreatedAt, &visit.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan visit: %w", err)
		}
		visits = append(visits, visit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate visits: %w", err)
	}

	return visits, nil
}

// Check if visit exists for schedule
func (r *VisitRepository) VisitExistsForSchedule(ctx context.Context, scheduleID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM visits WHERE schedule_id = $1)`

	var exists bool
	err := r.DB.QueryRow(ctx, query, scheduleID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check visit existence: %w", err)
	}

	return exists, nil
}

// Get visit duration statistics
func (r *VisitRepository) GetVisitDurationStats(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT
			AVG(duration_minutes) as avg_duration,
			MIN(duration_minutes) as min_duration,
			MAX(duration_minutes) as max_duration,
			COUNT(*) as total_completed
		FROM visits
		WHERE duration_minutes IS NOT NULL
	`

	var avgDuration, minDuration, maxDuration float64
	var totalCompleted int64

	err := r.DB.QueryRow(ctx, query).Scan(&avgDuration, &minDuration, &maxDuration, &totalCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get visit duration stats: %w", err)
	}

	stats := map[string]interface{}{
		"avg_duration":    avgDuration,
		"min_duration":    minDuration,
		"max_duration":    maxDuration,
		"total_completed": totalCompleted,
	}

	return stats, nil
}