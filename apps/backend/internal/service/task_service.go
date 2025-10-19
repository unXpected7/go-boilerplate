package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/model"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/repository"
)

type TaskService struct {
	taskRepo     *repository.TaskRepository
	scheduleRepo *repository.ScheduleRepository
}

func NewTaskService(taskRepo *repository.TaskRepository, scheduleRepo *repository.ScheduleRepository) *TaskService {
	return &TaskService{
		taskRepo:     taskRepo,
		scheduleRepo: scheduleRepo,
	}
}

// Create a new task
func (t *TaskService) CreateTask(ctx context.Context, scheduleID uuid.UUID, name, description string) (*model.Task, error) {
	task := &model.Task{
		Base: model.Base{
			BaseWithId: model.BaseWithId{
				ID: uuid.New(),
			},
			BaseWithCreatedAt: model.BaseWithCreatedAt{
				CreatedAt: time.Now(),
			},
			BaseWithUpdatedAt: model.BaseWithUpdatedAt{
				UpdatedAt: time.Now(),
			},
		},
		ScheduleID:  scheduleID,
		Name:        name,
		Description: &description,
		Status:      "pending",
	}

	if err := t.taskRepo.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// Create multiple tasks for a schedule
func (t *TaskService) CreateBatchTasks(ctx context.Context, scheduleID uuid.UUID, tasks []model.TaskCreate) error {
	taskModels := make([]model.Task, len(tasks))
	for i, taskCreate := range tasks {
		taskModels[i] = model.Task{
			Base: model.Base{
				BaseWithId: model.BaseWithId{
					ID: uuid.New(),
				},
				BaseWithCreatedAt: model.BaseWithCreatedAt{
					CreatedAt: time.Now(),
				},
				BaseWithUpdatedAt: model.BaseWithUpdatedAt{
					UpdatedAt: time.Now(),
				},
			},
			ScheduleID:  scheduleID,
			Name:        taskCreate.Name,
			Description: taskCreate.Description,
			Status:      "pending",
		}
	}

	if err := t.taskRepo.CreateBatchTasks(ctx, taskModels); err != nil {
		return fmt.Errorf("failed to create batch tasks: %w", err)
	}

	return nil
}

// Get task by ID
func (t *TaskService) GetTaskByID(ctx context.Context, taskID uuid.UUID) (*model.Task, error) {
	return t.taskRepo.GetTaskByID(ctx, taskID)
}

// Get tasks by schedule ID
func (t *TaskService) GetTasksByScheduleID(ctx context.Context, scheduleID uuid.UUID) ([]model.Task, error) {
	return t.taskRepo.GetTasksByScheduleID(ctx, scheduleID)
}

// Update task status
func (t *TaskService) UpdateTaskStatus(ctx context.Context, taskID uuid.UUID, status string, reason *string) (*model.Task, error) {
	// Validate status
	if !isValidTaskStatus(status) {
		return nil, fmt.Errorf("invalid task status: %s", status)
	}

	// For not_completed tasks, reason is required
	if status == "not_completed" && (reason == nil || *reason == "") {
		return nil, fmt.Errorf("reason is required for not_completed tasks")
	}

	// Validate that task exists
	exists, err := t.taskRepo.TaskExists(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to check task existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	// Update task
	updatedTask, err := t.taskRepo.UpdateTaskStatus(ctx, taskID, status, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	return updatedTask, nil
}

// Update task details
func (t *TaskService) UpdateTask(ctx context.Context, taskID uuid.UUID, name, description string, status string, reason *string) (*model.Task, error) {
	// Validate status
	if status != "" && !isValidTaskStatus(status) {
		return nil, fmt.Errorf("invalid task status: %s", status)
	}

	// For not_completed tasks, reason is required
	if status == "not_completed" && (reason == nil || *reason == "") {
		return nil, fmt.Errorf("reason is required for not_completed tasks")
	}

	// Get existing task
	task, err := t.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Update fields
	if name != "" {
		task.Name = name
	}
	if description != "" {
		task.Description = &description
	}
	if status != "" {
		task.Status = status
	}
	if reason != nil {
		task.Reason = reason
	}
	task.UpdatedAt = time.Now()

	// Update task
	if err := t.taskRepo.UpdateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// Get task statistics for a schedule
func (t *TaskService) GetTaskStatsBySchedule(ctx context.Context, scheduleID uuid.UUID) (*model.TaskStats, error) {
	return t.taskRepo.GetTaskStatsBySchedule(ctx, scheduleID)
}

// Get task completion rate for a schedule
func (t *TaskService) GetTaskCompletionRate(ctx context.Context, scheduleID uuid.UUID) (float64, error) {
	return t.taskRepo.GetTaskCompletionRate(ctx, scheduleID)
}

// Get tasks by status
func (t *TaskService) GetTasksByStatus(ctx context.Context, status string) ([]model.Task, error) {
	return t.taskRepo.GetTasksByStatus(ctx, status)
}

// Get incomplete tasks with reasons
func (t *TaskService) GetIncompleteTasks(ctx context.Context) ([]model.Task, error) {
	return t.taskRepo.GetIncompleteTasks(ctx)
}

// Delete task
func (t *TaskService) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	// Check if task exists
	exists, err := t.taskRepo.TaskExists(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to check task existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("task not found")
	}

	if err := t.taskRepo.DeleteTask(ctx, taskID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// Validate task update permissions
func (t *TaskService) ValidateTaskUpdate(ctx context.Context, taskID uuid.UUID, scheduleID uuid.UUID) error {
	// Check if task exists and belongs to the schedule
	task, err := t.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	if task.ScheduleID != scheduleID {
		return fmt.Errorf("task does not belong to the specified schedule")
	}

	return nil
}

// Calculate overall task statistics for all schedules
func (t *TaskService) GetOverallTaskStats(ctx context.Context) (*model.TaskStats, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'not_completed' THEN 1 END) as not_completed
		FROM tasks
	`

	var total, completed, pending, notCompleted int
	err := t.taskRepo.DB.QueryRow(ctx, query).Scan(&total, &completed, &pending, &notCompleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get overall task stats: %w", err)
	}

	stats := &model.TaskStats{
		TotalTasks:      total,
		CompletedTasks:  completed,
		PendingTasks:    pending,
		NotCompletedTasks: notCompleted,
	}

	return stats, nil
}

// Get tasks that require attention (not completed with reasons)
func (t *TaskService) GetTasksRequiringAttention(ctx context.Context) ([]model.Task, error) {
	incompleteTasks, err := t.GetIncompleteTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get incomplete tasks: %w", err)
	}

	// Filter tasks that have reasons (require attention)
	var attentionTasks []model.Task
	for _, task := range incompleteTasks {
		if task.Reason != nil && *task.Reason != "" {
			attentionTasks = append(attentionTasks, task)
		}
	}

	return attentionTasks, nil
}

// Update task reason only
func (t *TaskService) UpdateTaskReason(ctx context.Context, taskID uuid.UUID, reason string) error {
	// Validate reason is not empty
	if reason == "" {
		return fmt.Errorf("reason cannot be empty")
	}

	return t.taskRepo.UpdateTaskReason(ctx, taskID, reason)
}

// Generate task report for a schedule
func (t *TaskService) GenerateTaskReport(ctx context.Context, scheduleID uuid.UUID) (map[string]interface{}, error) {
	// Get tasks for schedule
	tasks, err := t.GetTasksByScheduleID(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	// Get statistics
	stats, err := t.GetTaskStatsBySchedule(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task stats: %w", err)
	}

	// Calculate completion rate
	completionRate, err := t.GetTaskCompletionRate(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get completion rate: %w", err)
	}

	// Get schedule details
	schedule, err := t.scheduleRepo.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	report := map[string]interface{}{
		"scheduleId":    schedule.ID,
		"clientName":    schedule.ClientName,
		"shiftTime":     schedule.ShiftTime,
		"location":      schedule.Location,
		"totalTasks":    stats.TotalTasks,
		"completed":     stats.CompletedTasks,
		"pending":       stats.PendingTasks,
		"notCompleted":  stats.NotCompletedTasks,
		"completionRate": completionRate,
		"tasks":         tasks,
		"generatedAt":   time.Now().Format(time.RFC3339),
	}

	return report, nil
}

// Helper function to check if task status is valid
func isValidTaskStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":      true,
		"completed":    true,
		"not_completed": true,
	}
	return validStatuses[status]
}

// Mark all pending tasks as not completed (for bulk operations)
func (t *TaskService) MarkPendingTasksAsNotCompleted(ctx context.Context, scheduleID uuid.UUID, reason string) error {
	// Get pending tasks
	tasks, err := t.GetTasksByScheduleID(ctx, scheduleID)
	if err != nil {
		return fmt.Errorf("failed to get tasks: %w", err)
	}

	// Update each pending task
	for _, task := range tasks {
		if task.Status == "pending" {
			_, err := t.UpdateTaskStatus(ctx, task.ID, "not_completed", &reason)
			if err != nil {
				return fmt.Errorf("failed to update task %s: %w", task.ID, err)
			}
		}
	}

	return nil
}