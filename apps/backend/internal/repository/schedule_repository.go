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

type ScheduleRepository struct {
	DB *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{DB: db}
}

// Get all schedules with pagination and filtering
func (r *ScheduleRepository) GetSchedules(ctx context.Context, page, limit int, status string) (*model.PaginatedResponse[model.Schedule], error) {
	query := `
		SELECT * FROM schedules
		WHERE ($1 = '' OR status = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	var schedules []model.Schedule
	offset := (page - 1) * limit
	rows, err := r.DB.Query(ctx, query, status, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedules: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schedule model.Schedule
		if err := rows.Scan(&schedule.ID, &schedule.ClientName, &schedule.ShiftTime, &schedule.Location, &schedule.Status, &schedule.VisitID, &schedule.CreatedAt, &schedule.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate schedules: %w", err)
	}

	// Get total count for pagination
	countQuery := `
		SELECT COUNT(*) FROM schedules
		WHERE ($1 = '' OR status = $2)
	`

	var total int
	err = r.DB.QueryRow(ctx, countQuery, status, status).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule count: %w", err)
	}

	totalPages := (total + limit - 1) / limit

	return &model.PaginatedResponse[model.Schedule]{
		Data:       schedules,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// Get today's schedules
func (r *ScheduleRepository) GetTodaySchedules(ctx context.Context) ([]model.Schedule, error) {
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT * FROM schedules
		WHERE created_at::date = $1::date
		ORDER BY shift_time ASC
	`

	var schedules []model.Schedule
	rows, err := r.DB.Query(ctx, query, today)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's schedules: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schedule model.Schedule
		if err := rows.Scan(&schedule.ID, &schedule.ClientName, &schedule.ShiftTime, &schedule.Location, &schedule.Status, &schedule.VisitID, &schedule.CreatedAt, &schedule.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate schedules: %w", err)
	}

	return schedules, nil
}

// Get schedule by ID
func (r *ScheduleRepository) GetScheduleByID(ctx context.Context, id uuid.UUID) (*model.Schedule, error) {
	query := `SELECT * FROM schedules WHERE id = $1`

	var schedule model.Schedule
	err := r.DB.QueryRow(ctx, query, id).Scan(&schedule.ID, &schedule.ClientName, &schedule.ShiftTime, &schedule.Location, &schedule.Status, &schedule.VisitID, &schedule.CreatedAt, &schedule.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("schedule not found")
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	return &schedule, nil
}

// Get schedule with visit and tasks
func (r *ScheduleRepository) GetScheduleWithDetails(ctx context.Context, id uuid.UUID) (*model.ScheduleWithTasks, error) {
	query := `
		SELECT s.*, v.*, t.*
		FROM schedules s
		LEFT JOIN visits v ON s.visit_id = v.id
		LEFT JOIN tasks t ON s.id = t.schedule_id
		WHERE s.id = $1
		ORDER BY t.created_at ASC
	`

	type result struct {
		ScheduleID    uuid.UUID
		ClientName    string
		ShiftTime     time.Time
		Location      string
		Status        string
		VisitID       *uuid.UUID
		CreatedAt     time.Time
		UpdatedAt     time.Time
		VisitID_      uuid.UUID      `db:"visit_id"`
		VisitStart    time.Time      `db:"visit_start"`
		VisitEnd      *time.Time     `db:"visit_end"`
		VisitStatus   string         `db:"visit_status"`
		TaskID        uuid.UUID      `db:"task_id"`
		TaskName      string         `db:"task_name"`
		TaskDesc      string         `db:"task_description"`
		TaskStatus    string         `db:"task_status"`
		TaskReason    *string        `db:"task_reason"`
		TaskCompleted *time.Time     `db:"task_completed_at"`
	}

	var results []result
	rows, err := r.DB.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule with details: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var res result
		if err := rows.Scan(&res.ScheduleID, &res.ClientName, &res.ShiftTime, &res.Location, &res.Status, &res.VisitID, &res.CreatedAt, &res.UpdatedAt,
			&res.VisitID_, &res.VisitStart, &res.VisitEnd, &res.VisitStatus,
			&res.TaskID, &res.TaskName, &res.TaskDesc, &res.TaskStatus, &res.TaskReason, &res.TaskCompleted); err != nil {
			return nil, fmt.Errorf("failed to scan result: %w", err)
		}
		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate results: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("schedule not found")
	}

	// Build the response
	schedule := model.Schedule{}
	schedule.ID = results[0].ScheduleID
	schedule.ClientName = results[0].ClientName
	schedule.ShiftTime = results[0].ShiftTime.Format("2006-01-02 15:04:05") // Convert time.Time to string
	schedule.Location = results[0].Location
	schedule.Status = results[0].Status
	schedule.VisitID = results[0].VisitID

	tasks := make([]model.Task, 0)
	var visit *model.Visit

	// Process results to build tasks and visit
	for _, res := range results {
		if res.TaskID != uuid.Nil {
			task := model.Task{}
			task.ID = res.TaskID
			task.ScheduleID = res.ScheduleID
			task.Name = res.TaskName
			if res.TaskDesc != "" {
		desc := res.TaskDesc
		task.Description = &desc
	}
			task.Status = res.TaskStatus
			task.Reason = res.TaskReason
			task.CompletedAt = res.TaskCompleted
			tasks = append(tasks, task)
		}

		if res.VisitID_ != uuid.Nil && visit == nil {
			visit = &model.Visit{}
			visit.ID = res.VisitID_
			visit.ScheduleID = res.ScheduleID
			visit.StartTime = res.VisitStart
			visit.EndTime = res.VisitEnd
			visit.Status = res.VisitStatus
		}
	}

	return &model.ScheduleWithTasks{
		Schedule: schedule,
		Tasks:    tasks,
	}, nil
}

// Create a new schedule
func (r *ScheduleRepository) CreateSchedule(ctx context.Context, schedule *model.Schedule) error {
	query := `
		INSERT INTO schedules (id, client_name, shift_time, location, status)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.DB.Exec(ctx, query, schedule.ID, schedule.ClientName, schedule.ShiftTime, schedule.Location, schedule.Status)
	if err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	return nil
}

// Update schedule
func (r *ScheduleRepository) UpdateSchedule(ctx context.Context, schedule *model.Schedule) error {
	query := `
		UPDATE schedules
		SET client_name = $1, shift_time = $2, location = $3, status = $4, visit_id = $5
		WHERE id = $6
	`

	_, err := r.DB.Exec(ctx, query, schedule.ClientName, schedule.ShiftTime, schedule.Location, schedule.Status, schedule.VisitID, schedule.ID)
	if err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}

	return nil
}

// Update schedule status
func (r *ScheduleRepository) UpdateScheduleStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE schedules SET status = $1 WHERE id = $2`

	_, err := r.DB.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update schedule status: %w", err)
	}

	return nil
}

// Get schedule statistics
func (r *ScheduleRepository) GetScheduleStats(ctx context.Context) (*model.TaskStats, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'not_completed' THEN 1 END) as not_completed
		FROM schedules
	`

	var total, completed, pending, notCompleted int
	err := r.DB.QueryRow(ctx, query).Scan(&total, &completed, &pending, &notCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule stats: %w", err)
	}

	stats := model.TaskStats{
		TotalTasks:      total,
		CompletedTasks:  completed,
		PendingTasks:    pending,
		NotCompletedTasks: notCompleted,
	}

	return &stats, nil
}

// Search schedules by client name or location
func (r *ScheduleRepository) SearchSchedules(ctx context.Context, queryStr string, page, limit int) (*model.PaginatedResponse[model.Schedule], error) {
	query := `
		SELECT * FROM schedules
		WHERE LOWER(client_name) LIKE LOWER($1) OR LOWER(location) LIKE LOWER($1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	searchPattern := "%" + queryStr + "%"

	var schedules []model.Schedule
	offset := (page - 1) * limit
	rows, err := r.DB.Query(ctx, query, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search schedules: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schedule model.Schedule
		if err := rows.Scan(&schedule.ID, &schedule.ClientName, &schedule.ShiftTime, &schedule.Location, &schedule.Status, &schedule.VisitID, &schedule.CreatedAt, &schedule.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate schedules: %w", err)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*) FROM schedules
		WHERE LOWER(client_name) LIKE LOWER($1) OR LOWER(location) LIKE LOWER($1)
	`

	var total int
	err = r.DB.QueryRow(ctx, countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get search count: %w", err)
	}

	totalPages := (total + limit - 1) / limit

	return &model.PaginatedResponse[model.Schedule]{
		Data:       schedules,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}