package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/model"
)

// MockDatabase represents a mock database for testing
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	args = m.Called(ctx, query, args)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockDatabase) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	args = m.Called(ctx, query, args)
	return args.Get(0).(*sql.Rows), args.Error(1)
}

func (m *MockDatabase) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	args = m.Called(ctx, query, args)
	return args.Get(0).(*sql.Row)
}

func (m *MockDatabase) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*sql.Tx), args.Error(1)
}

// Test ScheduleRepository
func TestScheduleRepository_GetSchedules(t *testing.T) {
	db := &MockDatabase{}
	repo := NewScheduleRepository(db)

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
		{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ClientName: "Jane Smith",
			ShiftTime:  "10:00-18:00",
			Location:   "456 Oak Ave",
			Status:     "completed",
		},
	}

	// Mock database responses
	db.On("SelectContext", ctx, mock.Anything,
		"SELECT * FROM schedules WHERE ($1 = '' OR status = $2) ORDER BY created_at DESC LIMIT $3 OFFSET $4",
		"", "", 10, 0).Return(expectedSchedules, nil)

	db.On("QueryRowContext", ctx,
		"SELECT COUNT(*) FROM schedules WHERE ($1 = '' OR status = $2)",
		"", "").Return(mock.NewResult(int64(len(expectedSchedules)), 0))

	result, err := repo.GetSchedules(ctx, 1, 10, "")

	assert.NoError(t, err)
	assert.Equal(t, len(expectedSchedules), len(result.Data))
	assert.Equal(t, "John Doe", result.Data[0].ClientName)
	assert.Equal(t, "Jane Smith", result.Data[1].ClientName)
	db.AssertExpectations(t)
}

func TestScheduleRepository_GetSchedules_WithStatusFilter(t *testing.T) {
	db := &MockDatabase{}
	repo := NewScheduleRepository(db)

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
			Status:     "completed",
		},
	}

	db.On("SelectContext", ctx, mock.Anything,
		"SELECT * FROM schedules WHERE ($1 = '' OR status = $2) ORDER BY created_at DESC LIMIT $3 OFFSET $4",
		"completed", "completed", 10, 0).Return(expectedSchedules, nil)

	db.On("QueryRowContext", ctx,
		"SELECT COUNT(*) FROM schedules WHERE ($1 = '' OR status = $2)",
		"completed", "completed").Return(mock.NewResult(int64(len(expectedSchedules)), 0))

	result, err := repo.GetSchedules(ctx, 1, 10, "completed")

	assert.NoError(t, err)
	assert.Equal(t, len(expectedSchedules), len(result.Data))
	assert.Equal(t, "completed", result.Data[0].Status)
	db.AssertExpectations(t)
}

func TestScheduleRepository_GetTodaySchedules(t *testing.T) {
	db := &MockDatabase{}
	repo := NewScheduleRepository(db)

	ctx := context.Background()
	today := time.Now().Format("2006-01-02")
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

	db.On("SelectContext", ctx, mock.Anything,
		"SELECT * FROM schedules WHERE created_at::date = $1::date ORDER BY shift_time ASC",
		today).Return(expectedSchedules, nil)

	result, err := repo.GetTodaySchedules(ctx)

	assert.NoError(t, err)
	assert.Equal(t, len(expectedSchedules), len(result))
	assert.Equal(t, "John Doe", result[0].ClientName)
	db.AssertExpectations(t)
}

func TestScheduleRepository_GetScheduleByID(t *testing.T) {
	db := &MockDatabase{}
	repo := NewScheduleRepository(db)

	ctx := context.Background()
	scheduleID := uuid.New()
	expectedSchedule := &model.Schedule{
		Base: model.Base{
			ID:        scheduleID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ClientName: "John Doe",
		ShiftTime:  "09:00-17:00",
		Location:   "123 Main St",
		Status:     "upcoming",
	}

	db.On("GetContext", ctx, mock.Anything,
		"SELECT * FROM schedules WHERE id = $1",
		scheduleID).Return(expectedSchedule, nil)

	result, err := repo.GetScheduleByID(ctx, scheduleID)

	assert.NoError(t, err)
	assert.Equal(t, expectedSchedule.ID, result.ID)
	assert.Equal(t, "John Doe", result.ClientName)
	db.AssertExpectations(t)
}

func TestScheduleRepository_GetScheduleByID_NotFound(t *testing.T) {
	db := &MockDatabase{}
	repo := NewScheduleRepository(db)

	ctx := context.Background()
	scheduleID := uuid.New()

	db.On("GetContext", ctx, mock.Anything,
		"SELECT * FROM schedules WHERE id = $1",
		scheduleID).Return(nil, sql.ErrNoRows)

	_, err := repo.GetScheduleByID(ctx, scheduleID)

	assert.Error(t, err)
	assert.Equal(t, "schedule not found", err.Error())
	db.AssertExpectations(t)
}

func TestScheduleRepository_CreateSchedule(t *testing.T) {
	db := &MockDatabase{}
	repo := NewScheduleRepository(db)

	ctx := context.Background()
	schedule := &model.Schedule{
		ID:          uuid.New(),
		ClientName:  "John Doe",
		ShiftTime:  "09:00-17:00",
		Location:   "123 Main St",
		Status:     "upcoming",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	db.On("ExecContext", ctx,
		"INSERT INTO schedules (id, client_name, shift_time, location, status) VALUES ($1, $2, $3, $4, $5)",
		schedule.ID, schedule.ClientName, schedule.ShiftTime, schedule.Location, schedule.Status).Return(mock.NewResult(1, 1), nil)

	err := repo.CreateSchedule(ctx, schedule)

	assert.NoError(t, err)
	db.AssertExpectations(t)
}

// Test VisitRepository
func TestVisitRepository_GetVisitByScheduleID(t *testing.T) {
	db := &MockDatabase{}
	repo := NewVisitRepository(db)

	ctx := context.Background()
	scheduleID := uuid.New()
	expectedVisit := &model.Visit{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ScheduleID:      scheduleID,
		StartTime:       time.Now(),
		StartLatitude:   40.7128,
		StartLongitude:  -74.0060,
		Status:          "in_progress",
	}

	db.On("GetContext", ctx, mock.Anything,
		"SELECT * FROM visits WHERE schedule_id = $1 ORDER BY created_at DESC LIMIT 1",
		scheduleID).Return(expectedVisit, nil)

	result, err := repo.GetVisitByScheduleID(ctx, scheduleID)

	assert.NoError(t, err)
	assert.Equal(t, expectedVisit.ScheduleID, result.ScheduleID)
	assert.Equal(t, "in_progress", result.Status)
	db.AssertExpectations(t)
}

func TestVisitRepository_StartVisit(t *testing.T) {
	db := &MockDatabase{}
	repo := NewVisitRepository(db)

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

	db.On("ExecContext", ctx,
		"INSERT INTO visits (id, schedule_id, start_time, start_latitude, start_longitude, status) VALUES ($1, $2, $3, $4, $5, $6)",
		mock.Anything, scheduleID, startTime, startLat, startLong, "in_progress").Return(mock.NewResult(1, 1), nil)

	result, err := repo.StartVisit(ctx, scheduleID, startTime, startLat, startLong)

	assert.NoError(t, err)
	assert.Equal(t, scheduleID, result.ScheduleID)
	assert.Equal(t, startTime, result.StartTime)
	assert.Equal(t, startLat, result.StartLatitude)
	assert.Equal(t, startLong, result.StartLongitude)
	db.AssertExpectations(t)
}

func TestVisitRepository_EndVisit(t *testing.T) {
	db := &MockDatabase{}
	repo := NewVisitRepository(db)

	ctx := context.Background()
	visitID := uuid.New()
	endTime := time.Now()
	endLat := 40.7589
	endLong := -73.9851

	db.On("ExecContext", ctx,
		"UPDATE visits SET end_time = $1, end_latitude = $2, end_longitude = $3, status = $4 WHERE id = $5",
		endTime, endLat, endLong, "completed", visitID).Return(mock.NewResult(1, 1), nil)

	db.On("GetContext", ctx, mock.Anything,
		"SELECT * FROM visits WHERE id = $1",
		visitID).Return(&model.Visit{
		Base: model.Base{
			ID:        visitID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ScheduleID:      uuid.New(),
		StartTime:       time.Now().Add(-2 * time.Hour),
		EndTime:         &endTime,
		StartLatitude:   40.7128,
		StartLongitude:  -74.0060,
		EndLatitude:     &endLat,
		EndLongitude:    &endLong,
		Status:          "completed",
		DurationMinutes: new(int), // Duration will be calculated
	}, nil)

	result, err := repo.EndVisit(ctx, visitID, endTime, endLat, endLong)

	assert.NoError(t, err)
	assert.NotNil(t, result.EndTime)
	assert.Equal(t, "completed", result.Status)
	db.AssertExpectations(t)
}

func TestVisitRepository_VisitExistsForSchedule(t *testing.T) {
	db := &MockDatabase{}
	repo := NewVisitRepository(db)

	ctx := context.Background()
	scheduleID := uuid.New()

	db.On("QueryRowContext", ctx,
		"SELECT EXISTS(SELECT 1 FROM visits WHERE schedule_id = $1)",
		scheduleID).Return(mock.NewResult(int64(1), 0), nil).Run(func(args mock.Arguments) {
		// Mock the Scan method to return true
		result := args.Get(0).(*interface{})
		*result = true
	})

	exists, err := repo.VisitExistsForSchedule(ctx, scheduleID)

	assert.NoError(t, err)
	assert.True(t, exists)
	db.AssertExpectations(t)
}

// Test TaskRepository
func TestTaskRepository_GetTasksByScheduleID(t *testing.T) {
	db := &MockDatabase{}
	repo := NewTaskRepository(db)

	ctx := context.Background()
	scheduleID := uuid.New()
	expectedTasks := []model.Task{
		{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ScheduleID: scheduleID,
			Name:        "Medication Administration",
			Status:      "pending",
		},
		{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ScheduleID: scheduleID,
			Name:        "Personal Care",
			Status:      "completed",
		},
	}

	db.On("SelectContext", ctx, mock.Anything,
		"SELECT * FROM tasks WHERE schedule_id = $1 ORDER BY created_at ASC",
		scheduleID).Return(expectedTasks, nil)

	result, err := repo.GetTasksByScheduleID(ctx, scheduleID)

	assert.NoError(t, err)
	assert.Equal(t, len(expectedTasks), len(result))
	assert.Equal(t, "Medication Administration", result[0].Name)
	assert.Equal(t, "Personal Care", result[1].Name)
	db.AssertExpectations(t)
}

func TestTaskRepository_CreateTask(t *testing.T) {
	db := &MockDatabase{}
	repo := NewTaskRepository(db)

	ctx := context.Background()
	task := &model.Task{
		ID:        uuid.New(),
		ScheduleID: uuid.New(),
		Name:      "Medication Administration",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db.On("ExecContext", ctx,
		"INSERT INTO tasks (id, schedule_id, name, description, status) VALUES ($1, $2, $3, $4, $5)",
		task.ID, task.ScheduleID, task.Name, (*string)(nil), task.Status).Return(mock.NewResult(1, 1), nil)

	err := repo.CreateTask(ctx, task)

	assert.NoError(t, err)
	db.AssertExpectations(t)
}

func TestTaskRepository_UpdateTaskStatus(t *testing.T) {
	db := &MockDatabase{}
	repo := NewTaskRepository(db)

	ctx := context.Background()
	taskID := uuid.New()
	status := "completed"
	reason := "Patient refused medication"
	completedAt := time.Now()

	db.On("GetContext", ctx, mock.Anything,
		"SELECT * FROM tasks WHERE id = $1",
		taskID).Return(&model.Task{
		Base: model.Base{
			ID:        taskID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ScheduleID: uuid.New(),
		Name:       "Medication Administration",
		Status:     "pending",
	}, nil)

	db.On("ExecContext", ctx,
		"UPDATE tasks SET status = $1, reason = $2, completed_at = CASE WHEN $1 = 'completed' THEN NOW() ELSE NULL END WHERE id = $3 RETURNING *",
		status, reason, taskID).Return(mock.NewResult(1, 1), nil).Run(func(args mock.Arguments) {
		// Mock the returned result
		result := &model.Task{
			Base: model.Base{
				ID:        taskID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ScheduleID:  uuid.New(),
			Name:        "Medication Administration",
			Status:      status,
			Reason:      &reason,
			CompletedAt: &completedAt,
		}
		args.Get(0).(interface{}) = result
	})

	result, err := repo.UpdateTaskStatus(ctx, taskID, status, &reason)

	assert.NoError(t, err)
	assert.Equal(t, status, result.Status)
	assert.Equal(t, reason, *result.Reason)
	assert.NotNil(t, result.CompletedAt)
	db.AssertExpectations(t)
}

// Test error cases
func TestRepository_ErrorHandling(t *testing.T) {
	db := &MockDatabase{}
	repo := NewScheduleRepository(db)

	ctx := context.Background()
	scheduleID := uuid.New()

	// Test database error
	db.On("GetContext", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)

	_, err := repo.GetScheduleByID(ctx, scheduleID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get schedule")

	db.AssertExpectations(t)
}

// Test pagination edge cases
func TestScheduleRepository_Pagination(t *testing.T) {
	db := &MockDatabase{}
	repo := NewScheduleRepository(db)

	ctx := context.Background()
	expectedSchedules := []model.Schedule{
		{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ClientName: "Client 1",
			ShiftTime:  "09:00-17:00",
			Location:   "Location 1",
			Status:     "upcoming",
		},
	}

	// Test page 1, limit 1
	db.On("SelectContext", ctx, mock.Anything,
		"SELECT * FROM schedules WHERE ($1 = '' OR status = $2) ORDER BY created_at DESC LIMIT $3 OFFSET $4",
		"", "", 1, 0).Return(expectedSchedules, nil)

	db.On("QueryRowContext", ctx,
		"SELECT COUNT(*) FROM schedules WHERE ($1 = '' OR status = $2)",
		"", "").Return(mock.NewResult(int64(5), 0))

	result, err := repo.GetSchedules(ctx, 1, 1, "")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result.Data))
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 1, result.Limit)
	assert.Equal(t, 5, result.Total)
	assert.Equal(t, 5, result.TotalPages) // ceil(5/1) = 5

	db.AssertExpectations(t)
}

// Test search functionality
func TestScheduleRepository_SearchSchedules(t *testing.T) {
	db := &MockDatabase{}
	repo := NewScheduleRepository(db)

	ctx := context.Background()
	query := "john"
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

	db.On("SelectContext", ctx, mock.Anything,
		"SELECT * FROM schedules WHERE LOWER(client_name) LIKE LOWER($1) OR LOWER(location) LIKE LOWER($1) ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		"%john%", 10, 0).Return(expectedSchedules, nil)

	db.On("QueryRowContext", ctx,
		"SELECT COUNT(*) FROM schedules WHERE LOWER(client_name) LIKE LOWER($1) OR LOWER(location) LIKE LOWER($1)",
		"%john%").Return(mock.NewResult(int64(1), 0))

	result, err := repo.SearchSchedules(ctx, query, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, len(expectedSchedules), len(result.Data))
	assert.Equal(t, "John Doe", result.Data[0].ClientName)
	db.AssertExpectations(t)
}