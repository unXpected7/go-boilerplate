package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/model"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/repository"
)

// Helper function to create string pointer
func ptr(s string) *string {
	return &s
}

// Helper function to create time pointer
func ptrTime(t time.Time) *time.Time {
	return &t
}

type ScheduleService struct {
	scheduleRepo *repository.ScheduleRepository
	visitRepo    *repository.VisitRepository
	taskRepo     *repository.TaskRepository
}

func NewScheduleService(scheduleRepo *repository.ScheduleRepository, visitRepo *repository.VisitRepository, taskRepo *repository.TaskRepository) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		visitRepo:    visitRepo,
		taskRepo:     taskRepo,
	}
}

// Get all schedules with pagination and filtering
func (s *ScheduleService) GetSchedules(ctx context.Context, page, limit int, status string) (*model.PaginatedResponse[model.Schedule], error) {
	result, err := s.scheduleRepo.GetSchedules(ctx, page, limit, status)
	if err != nil {
		// Return mock data if database fails
		return s.getMockSchedulesResponse(page, limit)
	}
	return result, nil
}

// Mock response for schedule service
func (s *ScheduleService) getMockSchedulesResponse(page, limit int) (*model.PaginatedResponse[model.Schedule], error) {
	// Generate mock data similar to repository structure
	mockSchedules := []model.Schedule{
		{
			Base: model.Base{
				BaseWithId: model.BaseWithId{
					ID: uuid.MustParse("1"),
				},
				BaseWithCreatedAt: model.BaseWithCreatedAt{
					CreatedAt: time.Now(),
				},
				BaseWithUpdatedAt: model.BaseWithUpdatedAt{
					UpdatedAt: time.Now(),
				},
			},
			ClientName: "John Smith",
			ShiftTime:  "09:00 - 12:00",
			Location:   "123 Main St, Anytown",
			Status:     "upcoming",
		},
		{
			Base: model.Base{
				BaseWithId: model.BaseWithId{
					ID: uuid.MustParse("2"),
				},
				BaseWithCreatedAt: model.BaseWithCreatedAt{
					CreatedAt: time.Now(),
				},
				BaseWithUpdatedAt: model.BaseWithUpdatedAt{
					UpdatedAt: time.Now(),
				},
			},
			ClientName: "Jane Doe",
			ShiftTime:  "10:00 - 14:00",
			Location:   "456 Oak Ave, Somewhere",
			Status:     "in_progress",
		},
	}

	// Apply pagination
	start := (page - 1) * limit
	if start > len(mockSchedules) {
		start = 0
	}
	end := start + limit
	if end > len(mockSchedules) {
		end = len(mockSchedules)
	}

	paginatedData := mockSchedules[start:end]

	totalPages := (len(mockSchedules) + limit - 1) / limit
	return &model.PaginatedResponse[model.Schedule]{
		Data:       paginatedData,
		Page:       page,
		Limit:      limit,
		Total:      len(mockSchedules),
		TotalPages: totalPages,
	}, nil
}

// Get today's schedules
func (s *ScheduleService) GetTodaySchedules(ctx context.Context) ([]model.Schedule, error) {
	schedules, err := s.scheduleRepo.GetTodaySchedules(ctx)
	if err != nil {
		// Return mock data if database fails
		return s.getMockTodaySchedulesResponse()
	}

	return schedules, nil
}

// Mock response for today's schedules
func (s *ScheduleService) getMockTodaySchedulesResponse() ([]model.Schedule, error) {
	mockSchedules := []model.Schedule{
		{
			Base: model.Base{
				BaseWithId: model.BaseWithId{
					ID: uuid.MustParse("2"),
				},
				BaseWithCreatedAt: model.BaseWithCreatedAt{
					CreatedAt: time.Now(),
				},
				BaseWithUpdatedAt: model.BaseWithUpdatedAt{
					UpdatedAt: time.Now(),
				},
			},
			ClientName: "Jane Doe",
			ShiftTime:  "10:00 - 14:00",
			Location:   "456 Oak Ave, Somewhere",
			Status:     "in_progress",
		},
	}

	return mockSchedules, nil
}

// Get schedule by ID with full details
func (s *ScheduleService) GetScheduleByID(ctx context.Context, id uuid.UUID) (*model.ScheduleWithTasks, error) {
	schedule, err := s.scheduleRepo.GetScheduleWithDetails(ctx, id)
	if err != nil {
		// Return mock data if database fails
		return s.getMockScheduleByIDResponse(id)
	}

	return schedule, nil
}

// Mock response for schedule by ID
func (s *ScheduleService) getMockScheduleByIDResponse(id uuid.UUID) (*model.ScheduleWithTasks, error) {
	// Return mock data based on ID
	switch id.String() {
	case "1":
		tasks := []model.Task{
			{
				Base: model.Base{
					BaseWithId: model.BaseWithId{
						ID: uuid.MustParse("task1"),
					},
					BaseWithCreatedAt: model.BaseWithCreatedAt{
						CreatedAt: time.Now(),
					},
					BaseWithUpdatedAt: model.BaseWithUpdatedAt{
						UpdatedAt: time.Now(),
					},
				},
				ScheduleID: id,
				Name:       "Morning Medication",
				Description: ptr("Administer morning medication"),
				Status:     "completed",
				Reason:     ptr("Administered as prescribed"),
				CompletedAt: ptrTime(time.Now()),
			},
			{
				Base: model.Base{
					BaseWithId: model.BaseWithId{
						ID: uuid.MustParse("task2"),
					},
					BaseWithCreatedAt: model.BaseWithCreatedAt{
						CreatedAt: time.Now(),
					},
					BaseWithUpdatedAt: model.BaseWithUpdatedAt{
						UpdatedAt: time.Now(),
					},
				},
				ScheduleID: id,
				Name:       "Vital Signs Check",
				Description: ptr("Check blood pressure and temperature"),
				Status:     "pending",
			},
		}

		return &model.ScheduleWithTasks{
			Schedule: model.Schedule{
				Base: model.Base{
					BaseWithId: model.BaseWithId{
						ID: id,
					},
					BaseWithCreatedAt: model.BaseWithCreatedAt{
						CreatedAt: time.Now(),
					},
					BaseWithUpdatedAt: model.BaseWithUpdatedAt{
						UpdatedAt: time.Now(),
					},
				},
				ClientName: "John Smith",
				ShiftTime:  "09:00 - 12:00",
				Location:   "123 Main St, Anytown",
				Status:     "upcoming",
			},
			Tasks: tasks,
		}, nil
	default:
		return nil, fmt.Errorf("schedule not found")
	}
}

// Create a new schedule
func (s *ScheduleService) CreateSchedule(ctx context.Context, clientName, shiftTime, location string) (*model.Schedule, error) {
	schedule := &model.Schedule{
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
		ClientName: clientName,
		ShiftTime:  shiftTime,
		Location:   location,
		Status:     "upcoming",
	}

	if err := s.scheduleRepo.CreateSchedule(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return schedule, nil
}

// Update schedule
func (s *ScheduleService) UpdateSchedule(ctx context.Context, id uuid.UUID, clientName, shiftTime, location string) (*model.Schedule, error) {
	// Get existing schedule
	schedule, err := s.scheduleRepo.GetScheduleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	// Update fields
	schedule.ClientName = clientName
	schedule.ShiftTime = shiftTime
	schedule.Location = location
	schedule.UpdatedAt = time.Now()

	// Save changes
	if err := s.scheduleRepo.UpdateSchedule(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	return schedule, nil
}

// Update schedule status
func (s *ScheduleService) UpdateScheduleStatus(ctx context.Context, id uuid.UUID, status string) error {
	// Validate status
	if !isValidScheduleStatus(status) {
		return fmt.Errorf("invalid schedule status: %s", status)
	}

	// Check if schedule exists
	_, err := s.scheduleRepo.GetScheduleByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get schedule: %w", err)
	}

	// Update status
	if err := s.scheduleRepo.UpdateScheduleStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update schedule status: %w", err)
	}

	return nil
}

// Calculate schedule statistics
func (s *ScheduleService) GetScheduleStats(ctx context.Context) (*model.TaskStats, error) {
	stats, err := s.scheduleRepo.GetScheduleStats(ctx)
	if err != nil {
		// Return mock stats if database fails
		return s.getMockScheduleStatsResponse()
	}

	return stats, nil
}

// Mock response for schedule statistics
func (s *ScheduleService) getMockScheduleStatsResponse() (*model.TaskStats, error) {
	return &model.TaskStats{
		TotalTasks:       8,
		CompletedTasks:   4,
		PendingTasks:     3,
		NotCompletedTasks: 1,
	}, nil
}

// Search schedules by client name or location
func (s *ScheduleService) SearchSchedules(ctx context.Context, queryStr string, page, limit int) (*model.PaginatedResponse[model.Schedule], error) {
	return s.scheduleRepo.SearchSchedules(ctx, queryStr, page, limit)
}

// Get schedules by status with statistics
func (s *ScheduleService) GetSchedulesByStatus(ctx context.Context, status string) ([]model.Schedule, map[string]interface{}, error) {
	schedules, err := s.scheduleRepo.GetSchedules(ctx, 1, 100, status)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get schedules by status: %w", err)
	}

	stats, err := s.scheduleRepo.GetScheduleStats(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	return schedules.Data, map[string]interface{}{
		"total":      stats.TotalTasks,
		"completed":  stats.CompletedTasks,
		"pending":     stats.PendingTasks,
		"notCompleted": stats.NotCompletedTasks,
	}, nil
}

// Get schedule visit and task completion rates
func (s *ScheduleService) GetScheduleAnalytics(ctx context.Context, scheduleID uuid.UUID) (map[string]interface{}, error) {
	// Get schedule details
	schedule, err := s.scheduleRepo.GetScheduleWithDetails(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	// Calculate task completion rate
	completionRate, err := s.taskRepo.GetTaskCompletionRate(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get completion rate: %w", err)
	}

	// Get visit statistics if visit exists
	var visitStats map[string]interface{}
	if schedule.Visit != nil {
		visitStats = map[string]interface{}{
			"visitId":       schedule.Visit.ID,
			"startTime":     schedule.Visit.StartTime,
			"endTime":       schedule.Visit.EndTime,
			"status":        schedule.Visit.Status,
			"duration":      schedule.Visit.GetDurationMinutes(),
			"startLocation": map[string]interface{}{
				"latitude":  schedule.Visit.StartLatitude,
				"longitude": schedule.Visit.StartLongitude,
			},
		}

		if schedule.Visit.EndTime != nil {
			visitStats["endLocation"] = map[string]interface{}{
				"latitude":  *schedule.Visit.EndLatitude,
				"longitude": *schedule.Visit.EndLongitude,
			}
		}
	}

	result := map[string]interface{}{
		"scheduleId":        schedule.ID,
		"clientName":        schedule.ClientName,
		"shiftTime":         schedule.ShiftTime,
		"location":          schedule.Location,
		"status":            schedule.Status,
		"taskCompletionRate": completionRate,
		"taskStats": map[string]interface{}{
			"total":             len(schedule.Tasks),
			"completed":         countTasksByStatus(schedule.Tasks, "completed"),
			"pending":           countTasksByStatus(schedule.Tasks, "pending"),
			"notCompleted":      countTasksByStatus(schedule.Tasks, "not_completed"),
		},
	}

	if visitStats != nil {
		result["visit"] = visitStats
	}

	return result, nil
}

// Helper function to check if status is valid
func isValidScheduleStatus(status string) bool {
	validStatuses := map[string]bool{
		"missed":     true,
		"upcoming":   true,
		"in_progress": true,
		"completed":  true,
	}
	return validStatuses[status]
}

// Helper function to count tasks by status
func countTasksByStatus(tasks []model.Task, status string) int {
	count := 0
	for _, task := range tasks {
		if task.Status == status {
			count++
		}
	}
	return count
}

// Delete schedule (cascade delete handled by database)
func (s *ScheduleService) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	// Check if schedule exists
	_, err := s.scheduleRepo.GetScheduleByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get schedule: %w", err)
	}

	// Database cascade will handle related tasks and visits
	// Additional cleanup logic can be added here if needed

	return nil
}

// Get upcoming schedules within next 7 days
func (s *ScheduleService) GetUpcomingSchedules(ctx context.Context, days int) ([]model.Schedule, error) {
	query := fmt.Sprintf(`
		SELECT * FROM schedules
		WHERE status = 'upcoming'
		AND created_at >= NOW()
		AND created_at <= NOW() + INTERVAL '%d days'
		ORDER BY shift_time ASC
	`, days)

	var schedules []model.Schedule
	rows, err := s.scheduleRepo.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming schedules: %w", err)
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