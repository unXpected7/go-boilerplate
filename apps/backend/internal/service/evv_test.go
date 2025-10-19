package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/model"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/repository"
)

// Mock repositories for testing
type MockScheduleRepository struct {
	mock.Mock
}

func (m *MockScheduleRepository) GetSchedules(ctx context.Context, page, limit int, status string) (*model.PaginatedResponse[model.Schedule], error) {
	args := m.Called(ctx, page, limit, status)
	return args.Get(0).(*model.PaginatedResponse[model.Schedule]), args.Error(1)
}

func (m *MockScheduleRepository) GetTodaySchedules(ctx context.Context) ([]model.Schedule, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Schedule), args.Error(1)
}

func (m *MockScheduleRepository) GetScheduleByID(ctx context.Context, id uuid.UUID) (*model.Schedule, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Schedule), args.Error(1)
}

func (m *MockScheduleRepository) GetScheduleWithDetails(ctx context.Context, id uuid.UUID) (*model.ScheduleWithTasks, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.ScheduleWithTasks), args.Error(1)
}

func (m *MockScheduleRepository) CreateSchedule(ctx context.Context, schedule *model.Schedule) error {
	args := m.Called(ctx, schedule)
	return args.Error(0)
}

func (m *MockScheduleRepository) UpdateSchedule(ctx context.Context, schedule *model.Schedule) error {
	args := m.Called(ctx, schedule)
	return args.Error(0)
}

func (m *MockScheduleRepository) UpdateScheduleStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockScheduleRepository) GetScheduleStats(ctx context.Context) (*model.TaskStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.TaskStats), args.Error(1)
}

func (m *MockScheduleRepository) SearchSchedules(ctx context.Context, queryStr string, page, limit int) (*model.PaginatedResponse[model.Schedule], error) {
	args := m.Called(ctx, queryStr, page, limit)
	return args.Get(0).(*model.PaginatedResponse[model.Schedule]), args.Error(1)
}

type MockVisitRepository struct {
	mock.Mock
}

func (m *MockVisitRepository) GetVisitByID(ctx context.Context, id uuid.UUID) (*model.Visit, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Visit), args.Error(1)
}

func (m *MockVisitRepository) GetVisitByScheduleID(ctx context.Context, scheduleID uuid.UUID) (*model.Visit, error) {
	args := m.Called(ctx, scheduleID)
	return args.Get(0).(*model.Visit), args.Error(1)
}

func (m *MockVisitRepository) CreateVisit(ctx context.Context, visit *model.Visit) error {
	args := m.Called(ctx, visit)
	return args.Error(0)
}

func (m *MockVisitRepository) UpdateVisit(ctx context.Context, visit *model.Visit) error {
	args := m.Called(ctx, visit)
	return args.Error(0)
}

func (m *MockVisitRepository) UpdateVisitStatus(ctx context.Context, visitID uuid.UUID, status string) error {
	args := m.Called(ctx, visitID, status)
	return args.Error(0)
}

func (m *MockVisitRepository) StartVisit(ctx context.Context, scheduleID uuid.UUID, startTime time.Time, startLat, startLong float64) (*model.Visit, error) {
	args := m.Called(ctx, scheduleID, startTime, startLat, startLong)
	return args.Get(0).(*model.Visit), args.Error(1)
}

func (m *MockVisitRepository) EndVisit(ctx context.Context, visitID uuid.UUID, endTime time.Time, endLat, endLong float64) (*model.Visit, error) {
	args := m.Called(ctx, visitID, endTime, endLat, endLong)
	return args.Get(0).(*model.Visit), args.Error(1)
}

func (m *MockVisitRepository) GetVisitStats(ctx context.Context) (*model.TaskStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.TaskStats), args.Error(1)
}

func (m *MockVisitRepository) GetVisitsByStatus(ctx context.Context, status string) ([]model.Visit, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]model.Visit), args.Error(1)
}

func (m *MockVisitRepository) VisitExistsForSchedule(ctx context.Context, scheduleID uuid.UUID) (bool, error) {
	args := m.Called(ctx, scheduleID)
	return args.Bool(0), args.Error(1)
}

func (m *MockVisitRepository) GetVisitDurationStats(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetTaskByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTasksByScheduleID(ctx context.Context, scheduleID uuid.UUID) ([]model.Task, error) {
	args := m.Called(ctx, scheduleID)
	return args.Get(0).([]model.Task), args.Error(1)
}

func (m *MockTaskRepository) CreateTask(ctx context.Context, task *model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) UpdateTaskStatus(ctx context.Context, taskID uuid.UUID, status string, reason *string) (*model.Task, error) {
	args := m.Called(ctx, taskID, status, reason)
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) GetTaskStatsBySchedule(ctx context.Context, scheduleID uuid.UUID) (*model.TaskStats, error) {
	args := m.Called(ctx, scheduleID)
	return args.Get(0).(*model.TaskStats), args.Error(1)
}

func (m *MockTaskRepository) GetTaskCompletionRate(ctx context.Context, scheduleID uuid.UUID) (float64, error) {
	args := m.Called(ctx, scheduleID)
	return args.Float64(0), args.Error(1)
}

func (m *MockTaskRepository) CreateBatchTasks(ctx context.Context, tasks []model.Task) error {
	args := m.Called(ctx, tasks)
	return args.Error(0)
}

func (m *MockTaskRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTaskRepository) TaskExists(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockTaskRepository) GetTasksByStatus(ctx context.Context, status string) ([]model.Task, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]model.Task), args.Error(1)
}

func (m *MockTaskRepository) GetIncompleteTasks(ctx context.Context) ([]model.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Task), args.Error(1)
}

func (m *MockTaskRepository) UpdateTaskReason(ctx context.Context, taskID uuid.UUID, reason string) error {
	args := m.Called(ctx, taskID, reason)
	return args.Error(0)
}

func (m *MockTaskRepository) GetOverallTaskStats(ctx context.Context) (*model.TaskStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.TaskStats), args.Error(1)
}

// Test ScheduleService
func TestScheduleService_GetSchedules(t *testing.T) {
	mockScheduleRepo := new(MockScheduleRepository)
	mockVisitRepo := new(MockVisitRepository)
	mockTaskRepo := new(MockTaskRepository)

	scheduleService := NewScheduleService(mockScheduleRepo, mockVisitRepo, mockTaskRepo)

	ctx := context.Background()
	expectedSchedules := []model.Schedule{
		{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ClientName: "John Doe",
			ShiftTime:  "09:00-17:00",
			Location:   "123 Main St",
			Status:     "upcoming",
		},
	}

	expectedResponse := &model.PaginatedResponse[model.Schedule]{
		Data:       expectedSchedules,
		Page:       1,
		Limit:      10,
		Total:      1,
		TotalPages: 1,
	}

	mockScheduleRepo.On("GetSchedules", ctx, 1, 10, "").Return(expectedResponse, nil)

	result, err := scheduleService.GetSchedules(ctx, 1, 10, "")

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	mockScheduleRepo.AssertExpectations(t)
}

func TestScheduleService_GetScheduleByID(t *testing.T) {
	mockScheduleRepo := new(MockScheduleRepository)
	mockVisitRepo := new(MockVisitRepository)
	mockTaskRepo := new(MockTaskRepository)

	scheduleService := NewScheduleService(mockScheduleRepo, mockVisitRepo, mockTaskRepo)

	ctx := context.Background()
	scheduleID := uuid.New()
	expectedSchedule := &model.ScheduleWithTasks{
		Schedule: model.Schedule{
			Base: model.Base{
				ID:        scheduleID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ClientName: "John Doe",
			ShiftTime:  "09:00-17:00",
			Location:   "123 Main St",
			Status:     "upcoming",
		},
		Tasks: []model.Task{},
	}

	mockScheduleRepo.On("GetScheduleWithDetails", ctx, scheduleID).Return(expectedSchedule, nil)

	result, err := scheduleService.GetScheduleByID(ctx, scheduleID)

	assert.NoError(t, err)
	assert.Equal(t, expectedSchedule, result)
	mockScheduleRepo.AssertExpectations(t)
}

func TestScheduleService_CreateSchedule(t *testing.T) {
	mockScheduleRepo := new(MockScheduleRepository)
	mockVisitRepo := new(MockVisitRepository)
	mockTaskRepo := new(MockTaskRepository)

	scheduleService := NewScheduleService(mockScheduleRepo, mockVisitRepo, mockTaskRepo)

	ctx := context.Background()
	expectedSchedule := &model.Schedule{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ClientName: "John Doe",
		ShiftTime:  "09:00-17:00",
		Location:   "123 Main St",
		Status:     "upcoming",
	}

	mockScheduleRepo.On("CreateSchedule", ctx, expectedSchedule).Return(nil)

	result, err := scheduleService.CreateSchedule(ctx, "John Doe", "09:00-17:00", "123 Main St")

	assert.NoError(t, err)
	assert.Equal(t, expectedSchedule.ClientName, result.ClientName)
	assert.Equal(t, expectedSchedule.ShiftTime, result.ShiftTime)
	assert.Equal(t, expectedSchedule.Location, result.Location)
	mockScheduleRepo.AssertExpectations(t)
}

// Test VisitService
func TestVisitService_StartVisit(t *testing.T) {
	mockVisitRepo := new(MockVisitRepository)
	mockScheduleRepo := new(MockScheduleRepository)

	visitService := NewVisitService(mockVisitRepo, mockScheduleRepo)

	ctx := context.Background()
	scheduleID := uuid.New()
	startTime := time.Now()
	startLat := 40.7128
	startLong := -74.0060

	expectedVisit := &model.Visit{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ScheduleID:      scheduleID,
		StartTime:       startTime,
		StartLatitude:   startLat,
		StartLongitude:  startLong,
		Status:          "in_progress",
	}

	mockScheduleRepo.On("GetScheduleByID", ctx, scheduleID).Return(&model.Schedule{
		Base: model.Base{
			ID:        scheduleID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ClientName: "John Doe",
		ShiftTime:  "09:00-17:00",
		Location:   "123 Main St",
		Status:     "upcoming",
	}, nil)

	mockVisitRepo.On("VisitExistsForSchedule", ctx, scheduleID).Return(false, nil)
	mockVisitRepo.On("StartVisit", ctx, scheduleID, startTime, startLat, startLong).Return(expectedVisit, nil)
	mockScheduleRepo.On("UpdateScheduleStatus", ctx, scheduleID, "in_progress").Return(nil)

	result, err := visitService.StartVisit(ctx, scheduleID, startTime, startLat, startLong)

	assert.NoError(t, err)
	assert.Equal(t, expectedVisit.ScheduleID, result.ScheduleID)
	assert.Equal(t, expectedVisit.StartTime, result.StartTime)
	mockVisitRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
}

func TestVisitService_EndVisit(t *testing.T) {
	mockVisitRepo := new(MockVisitRepository)
	mockScheduleRepo := new(MockScheduleRepository)

	visitService := NewVisitService(mockVisitRepo, mockScheduleRepo)

	ctx := context.Background()
	scheduleID := uuid.New()
	visitID := uuid.New()
	startTime := time.Now().Add(-2 * time.Hour)
	endTime := time.Now()
	startLat := 40.7128
	startLong := -74.0060
	endLat := 40.7589
	endLong := -73.9851

	existingVisit := &model.Visit{
		Base: model.Base{
			ID:        visitID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ScheduleID:      scheduleID,
		StartTime:       startTime,
		StartLatitude:   startLat,
		StartLongitude:  startLong,
		Status:          "in_progress",
	}

	expectedVisit := &model.Visit{
		Base: model.Base{
			ID:        visitID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ScheduleID:      scheduleID,
		StartTime:       startTime,
		EndTime:         &endTime,
		StartLatitude:   startLat,
		StartLongitude:  startLong,
		EndLatitude:     &endLat,
		EndLongitude:    &endLong,
		Status:          "completed",
	}

	mockVisitRepo.On("GetVisitByScheduleID", ctx, scheduleID).Return(existingVisit, nil)
	mockVisitRepo.On("EndVisit", ctx, visitID, endTime, endLat, endLong).Return(expectedVisit, nil)
	mockScheduleRepo.On("UpdateScheduleStatus", ctx, scheduleID, "completed").Return(nil)

	result, err := visitService.EndVisit(ctx, scheduleID, endTime, endLat, endLong)

	assert.NoError(t, err)
	assert.Equal(t, expectedVisit.ScheduleID, result.ScheduleID)
	assert.NotNil(t, result.EndTime)
	mockVisitRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
}

// Test TaskService
func TestTaskService_CreateTask(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockScheduleRepo := new(MockScheduleRepository)

	taskService := NewTaskService(mockTaskRepo, mockScheduleRepo)

	ctx := context.Background()
	scheduleID := uuid.New()
	taskName := "Medication Administration"
	taskDescription := "Administer prescribed medication"

	expectedTask := &model.Task{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ScheduleID:  scheduleID,
		Name:        taskName,
		Description: &taskDescription,
		Status:      "pending",
	}

	mockTaskRepo.On("CreateTask", ctx, expectedTask).Return(nil)

	result, err := taskService.CreateTask(ctx, scheduleID, taskName, taskDescription)

	assert.NoError(t, err)
	assert.Equal(t, expectedTask.ScheduleID, result.ScheduleID)
	assert.Equal(t, expectedTask.Name, result.Name)
	assert.Equal(t, expectedTask.Status, result.Status)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_UpdateTaskStatus(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockScheduleRepo := new(MockScheduleRepository)

	taskService := NewTaskService(mockTaskRepo, mockScheduleRepo)

	ctx := context.Background()
	taskID := uuid.New()
	scheduleID := uuid.New()
	status := "completed"
	completedAt := time.Now()

	expectedTask := &model.Task{
		Base: model.Base{
			ID:        taskID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ScheduleID:  scheduleID,
		Name:        "Medication Administration",
		Status:      status,
		CompletedAt: &completedAt,
	}

	mockTaskRepo.On("GetTaskByID", ctx, taskID).Return(expectedTask, nil)
	mockTaskRepo.On("UpdateTaskStatus", ctx, taskID, status, (*string)(nil)).Return(expectedTask, nil)

	result, err := taskService.UpdateTaskStatus(ctx, taskID, status, nil)

	assert.NoError(t, err)
	assert.Equal(t, status, result.Status)
	assert.NotNil(t, result.CompletedAt)
	mockTaskRepo.AssertExpectations(t)
}

func TestTaskService_UpdateTaskStatusWithReason(t *testing.T) {
	mockTaskRepo := new(MockTaskRepository)
	mockScheduleRepo := new(MockScheduleRepository)

	taskService := NewTaskService(mockTaskRepo, mockScheduleRepo)

	ctx := context.Background()
	taskID := uuid.New()
	scheduleID := uuid.New()
	status := "not_completed"
	reason := "Patient refused medication"

	mockTaskRepo.On("GetTaskByID", ctx, taskID).Return(&model.Task{
		Base: model.Base{
			ID:        taskID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ScheduleID: scheduleID,
		Name:       "Medication Administration",
		Status:     "pending",
	}, nil)
	mockTaskRepo.On("UpdateTaskStatus", ctx, taskID, status, &reason).Return(&model.Task{
		Base: model.Base{
			ID:        taskID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ScheduleID: scheduleID,
		Name:       "Medication Administration",
		Status:     status,
		Reason:     &reason,
	}, nil)

	result, err := taskService.UpdateTaskStatus(ctx, taskID, status, &reason)

	assert.NoError(t, err)
	assert.Equal(t, status, result.Status)
	assert.Equal(t, reason, *result.Reason)
	mockTaskRepo.AssertExpectations(t)
}

// Test validation functions
func TestValidationFunctions(t *testing.T) {
	t.Run("Valid coordinates", func(t *testing.T) {
		assert.True(t, isValidCoordinates(40.7128, -74.0060)) // New York
		assert.True(t, isValidCoordinates(51.5074, -0.1278))  // London
		assert.True(t, isValidCoordinates(-33.8688, 151.2093)) // Sydney
	})

	t.Run("Invalid coordinates", func(t *testing.T) {
		assert.False(t, isValidCoordinates(91.0, 0.0))    // Invalid latitude
		assert.False(t, isValidCoordinates(-91.0, 0.0))   // Invalid latitude
		assert.False(t, isValidCoordinates(0.0, 181.0))   // Invalid longitude
		assert.False(t, isValidCoordinates(0.0, -181.0))  // Invalid longitude
	})

	t.Run("Valid status", func(t *testing.T) {
		assert.True(t, isValidVisitStatus("not_started"))
		assert.True(t, isValidVisitStatus("in_progress"))
		assert.True(t, isValidVisitStatus("completed"))
		assert.False(t, isValidVisitStatus("invalid"))
	})

	t.Run("Valid shift time", func(t *testing.T) {
		assert.True(t, isValidShiftTime("09:00-17:00"))
		assert.True(t, isValidShiftTime("08:30-16:30"))
		assert.False(t, isValidShiftTime("09:00-17"))
		assert.False(t, isValidShiftTime("9:00-17:00"))
		assert.False(t, isValidShiftTime("invalid-time"))
	})
}