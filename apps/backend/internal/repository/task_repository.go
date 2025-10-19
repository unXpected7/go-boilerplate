package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/model"
)

type TaskRepository struct {
	DB *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{DB: db}
}

// Get task by ID
func (r *TaskRepository) GetTaskByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	query := `SELECT * FROM tasks WHERE id = $1`

	var task model.Task
	err := r.DB.QueryRow(ctx, query, id).Scan(&task.ID, &task.ScheduleID, &task.Name, &task.Description, &task.Status, &task.Reason, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// Get tasks by schedule ID
func (r *TaskRepository) GetTasksByScheduleID(ctx context.Context, scheduleID uuid.UUID) ([]model.Task, error) {
	query := `SELECT * FROM tasks WHERE schedule_id = $1 ORDER BY created_at ASC`

	var tasks []model.Task
	rows, err := r.DB.Query(ctx, query, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.ScheduleID, &task.Name, &task.Description, &task.Status, &task.Reason, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tasks: %w", err)
	}

	return tasks, nil
}

// Create a new task
func (r *TaskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO tasks (id, schedule_id, name, description, status)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.DB.Exec(ctx, query, task.ID, task.ScheduleID, task.Name, task.Description, task.Status)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// Update task
func (r *TaskRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	query := `
		UPDATE tasks
		SET name = $1, description = $2, status = $3, reason = $4, completed_at = $5
		WHERE id = $6
	`

	_, err := r.DB.Exec(ctx, query, task.Name, task.Description, task.Status, task.Reason, task.CompletedAt, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// Update task status
func (r *TaskRepository) UpdateTaskStatus(ctx context.Context, taskID uuid.UUID, status string, reason *string) (*model.Task, error) {
	query := `
		UPDATE tasks
		SET status = $1, reason = $2, completed_at = CASE
			WHEN $1 = 'completed' THEN NOW()
			ELSE NULL
		END
		WHERE id = $3
		RETURNING *
	`

	var task model.Task
	err := r.DB.QueryRow(ctx, query, status, reason, taskID).Scan(&task.ID, &task.ScheduleID, &task.Name, &task.Description, &task.Status, &task.Reason, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	return &task, nil
}

// Get task statistics for a schedule
func (r *TaskRepository) GetTaskStatsBySchedule(ctx context.Context, scheduleID uuid.UUID) (*model.TaskStats, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'not_completed' THEN 1 END) as not_completed
		FROM tasks
		WHERE schedule_id = $1
	`

	var total, completed, pending, notCompleted int
	err := r.DB.QueryRow(ctx, query, scheduleID).Scan(&total, &completed, &pending, &notCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get task stats: %w", err)
	}

	stats := model.TaskStats{
		TotalTasks:      total,
		CompletedTasks:  completed,
		PendingTasks:    pending,
		NotCompletedTasks: notCompleted,
	}

	return &stats, nil
}

// Get task completion rate for a schedule
func (r *TaskRepository) GetTaskCompletionRate(ctx context.Context, scheduleID uuid.UUID) (float64, error) {
	query := `
		SELECT
			COUNT(CASE WHEN status = 'completed' THEN 1 END)::float / COUNT(*) * 100 as completion_rate
		FROM tasks
		WHERE schedule_id = $1
	`

	var completionRate float64
	err := r.DB.QueryRow(ctx, query, scheduleID).Scan(&completionRate)
	if err != nil {
		return 0, fmt.Errorf("failed to get task completion rate: %w", err)
	}

	return completionRate, nil
}

// Create multiple tasks for a schedule
func (r *TaskRepository) CreateBatchTasks(ctx context.Context, tasks []model.Task) error {
	if len(tasks) == 0 {
		return nil
	}

	query := `
		INSERT INTO tasks (id, schedule_id, name, description, status)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, task := range tasks {
		_, err := r.DB.Exec(ctx, query, task.ID, task.ScheduleID, task.Name, task.Description, task.Status)
		if err != nil {
			return fmt.Errorf("failed to create batch tasks: %w", err)
		}
	}

	return nil
}

// Delete task
func (r *TaskRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`

	_, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// Check if task exists
func (r *TaskRepository) TaskExists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)`

	var exists bool
	err := r.DB.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check task existence: %w", err)
	}

	return exists, nil
}

// Get tasks by status
func (r *TaskRepository) GetTasksByStatus(ctx context.Context, status string) ([]model.Task, error) {
	query := `SELECT * FROM tasks WHERE status = $1 ORDER BY created_at DESC`

	var tasks []model.Task
	rows, err := r.DB.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by status: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.ScheduleID, &task.Name, &task.Description, &task.Status, &task.Reason, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tasks: %w", err)
	}

	return tasks, nil
}

// Get incomplete tasks with reasons
func (r *TaskRepository) GetIncompleteTasks(ctx context.Context) ([]model.Task, error) {
	query := `SELECT * FROM tasks WHERE status = 'not_completed' AND reason IS NOT NULL ORDER BY created_at DESC`

	var tasks []model.Task
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get incomplete tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.ScheduleID, &task.Name, &task.Description, &task.Status, &task.Reason, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tasks: %w", err)
	}

	return tasks, nil
}

// Update task reason
func (r *TaskRepository) UpdateTaskReason(ctx context.Context, taskID uuid.UUID, reason string) error {
	query := `UPDATE tasks SET reason = $1 WHERE id = $2`

	_, err := r.DB.Exec(ctx, query, reason, taskID)
	if err != nil {
		return fmt.Errorf("failed to update task reason: %w", err)
	}

	return nil
}