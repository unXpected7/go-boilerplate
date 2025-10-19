package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/service"
)

type EVVHandler struct {
	scheduleService *service.ScheduleService
	visitService    *service.VisitService
	taskService     *service.TaskService
}

// Mock data generator methods - simplified
func (h *EVVHandler) getMockSchedulesResponse(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"schedules": []map[string]interface{}{
			{
				"id":          "1",
				"clientName":  "John Smith",
				"shiftTime":   "09:00 - 12:00",
				"location":    "123 Main St, Anytown",
				"status":      "upcoming",
				"visitId":     nil,
				"visit":       nil,
				"tasks":       []map[string]interface{}{},
			},
		},
	})
}

func NewEVVHandler(scheduleService *service.ScheduleService, visitService *service.VisitService, taskService *service.TaskService) *EVVHandler {
	return &EVVHandler{
		scheduleService: scheduleService,
		visitService:    visitService,
		taskService:     taskService,
	}
}

// Get all schedules with pagination and filtering
func (h *EVVHandler) GetSchedules(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	_ = c.QueryParam("status") // Keep parameter for future use

	// Return mock data for now
	return h.getMockSchedulesResponse(c)
}

// Get today's schedules
func (h *EVVHandler) GetTodaySchedules(c echo.Context) error {
	// Return mock data for now
	todaySchedules := map[string]interface{}{
		"data":  []map[string]interface{}{},
		"total": 0,
	}

	return c.JSON(http.StatusOK, todaySchedules)
}

// Get schedule by ID
func (h *EVVHandler) GetScheduleById(c echo.Context) error {
	id := c.Param("id")
	
	// Return mock data for now
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":         id,
		"clientName": "John Smith",
		"shiftTime":  "09:00 - 12:00",
		"location":   "123 Main St, Anytown",
		"status":     "upcoming",
		"visit":      nil,
		"tasks":      []map[string]interface{}{},
	})
}

// Create schedule
func (h *EVVHandler) CreateSchedule(c echo.Context) error {
	// Return success response for now
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Schedule created successfully",
		"id":      "new-schedule-id",
	})
}

// Update schedule status
func (h *EVVHandler) UpdateScheduleStatus(c echo.Context) error {
	id := c.Param("id")
	
	// Return success response for now
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Schedule status updated successfully",
		"id":      id,
	})
}

// Search schedules
func (h *EVVHandler) SearchSchedules(c echo.Context) error {
	// Return mock data for now
	return h.getMockSchedulesResponse(c)
}

// Get mock stats
func (h *EVVHandler) GetMockStats(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"totalSchedules": 10,
		"completedTasks": 5,
		"pendingTasks":   3,
		"upcomingVisits": 2,
	})
}

// Get task stats
func (h *EVVHandler) GetTaskStats(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"totalTasks":      25,
		"completedTasks":  15,
		"pendingTasks":    8,
		"notCompletedTasks": 2,
	})
}

// Update task status
func (h *EVVHandler) UpdateTaskStatus(c echo.Context) error {
	scheduleId := c.Param("scheduleId")
	taskId := c.Param("taskId")

	// Return success response for now
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Task status updated successfully",
		"scheduleId": scheduleId,
		"taskId": taskId,
	})
}
