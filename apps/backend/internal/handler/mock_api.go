package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Helper functions for creating pointers
func intPtr(i int) *int {
	return &i
}

// strPtr is defined in evv_handler.go to avoid duplication

// Mock API handlers for fallback scenarios
func GetMockSchedules(c echo.Context) error {
	// Generate mock data
	mockSchedules := struct {
		Schedules []struct {
			ID          string    `json:"id"`
			ClientName  string    `json:"clientName"`
			ShiftTime   string    `json:"shiftTime"`
			Location    string    `json:"location"`
			Status      string    `json:"status"`
			Visit       *struct {
				ID         string `json:"id"`
				ScheduleID string `json:"scheduleId"`
				Status     string `json:"status"`
			} `json:"visit,omitempty"`
			Tasks       []struct {
				ID          string `json:"id"`
				ScheduleID  string `json:"scheduleId"`
				Name        string `json:"name"`
				Description string `json:"description,omitempty"`
				Status      string `json:"status"`
				Reason      string `json:"reason,omitempty"`
				CompletedAt string `json:"completedAt,omitempty"`
				CreatedAt   string `json:"createdAt"`
				UpdatedAt   string `json:"updatedAt"`
			} `json:"tasks"`
			CreatedAt string `json:"createdAt"`
			UpdatedAt string `json:"updatedAt"`
		} `json:"schedules"`
	}{
		Schedules: []struct {
			ID          string    `json:"id"`
			ClientName  string    `json:"clientName"`
			ShiftTime   string    `json:"shiftTime"`
			Location    string    `json:"location"`
			Status      string    `json:"status"`
			Visit       *struct {
				ID         string `json:"id"`
				ScheduleID string `json:"scheduleId"`
				Status     string `json:"status"`
			} `json:"visit,omitempty"`
			Tasks       []struct {
				ID          string `json:"id"`
				ScheduleID  string `json:"scheduleId"`
				Name        string `json:"name"`
				Description string `json:"description,omitempty"`
				Status      string `json:"status"`
				Reason      string `json:"reason,omitempty"`
				CompletedAt string `json:"completedAt,omitempty"`
				CreatedAt   string `json:"createdAt"`
				UpdatedAt   string `json:"updatedAt"`
			} `json:"tasks"`
			CreatedAt string `json:"createdAt"`
			UpdatedAt string `json:"updatedAt"`
		}{
			{
				ID:         "1",
				ClientName: "John Smith",
				ShiftTime:  "09:00 - 12:00",
				Location:   "123 Main St, Anytown",
				Status:     "upcoming",
				Visit: &struct {
					ID         string `json:"id"`
					ScheduleID string `json:"scheduleId"`
					Status     string `json:"status"`
				}{
					ID:         "visit1",
					ScheduleID: "1",
					Status:     "not_started",
				},
				Tasks: []struct {
					ID          string `json:"id"`
					ScheduleID  string `json:"scheduleId"`
					Name        string `json:"name"`
					Description string `json:"description,omitempty"`
					Status      string `json:"status"`
					Reason      string `json:"reason,omitempty"`
					CompletedAt string `json:"completedAt,omitempty"`
					CreatedAt   string `json:"createdAt"`
					UpdatedAt   string `json:"updatedAt"`
				}{
					{
						ID:          "task1",
						ScheduleID:  "1",
						Name:        "Morning Medication",
						Description: "Administer morning medication",
						Status:      "pending",
						CreatedAt:   "2025-10-18T08:00:00Z",
						UpdatedAt:   "2025-10-18T08:00:00Z",
					},
				},
				CreatedAt: "2025-10-18T08:00:00Z",
				UpdatedAt: "2025-10-18T08:00:00Z",
			},
		},
	}

	return c.JSON(http.StatusOK, mockSchedules)
}

func GetMockStats(c echo.Context) error {
	mockStats := map[string]interface{}{
		"total":           8,
		"upcoming":        3,
		"in_progress":     1,
		"completed":       4,
		"missed":          1,
		"completion_rate": 62.5,
		"average_duration": 165,
		"today_schedules": 5,
	}

	return c.JSON(http.StatusOK, mockStats)
}

func GetTodaySchedules(c echo.Context) error {
	// Return mock data for today's schedules
	todaySchedules := struct {
		Schedules []struct {
			ID          string    `json:"id"`
			ClientName  string    `json:"clientName"`
			ShiftTime   string    `json:"shiftTime"`
			Location    string    `json:"location"`
			Status      string    `json:"status"`
			Visit       *struct {
				ID         string `json:"id"`
				ScheduleID string `json:"scheduleId"`
				Status     string `json:"status"`
			} `json:"visit,omitempty"`
			Tasks       []struct {
				ID          string `json:"id"`
				ScheduleID  string `json:"scheduleId"`
				Name        string `json:"name"`
				Description string `json:"description,omitempty"`
				Status      string `json:"status"`
				Reason      string `json:"reason,omitempty"`
				CompletedAt string `json:"completedAt,omitempty"`
				CreatedAt   string `json:"createdAt"`
				UpdatedAt   string `json:"updatedAt"`
			} `json:"tasks"`
			CreatedAt string `json:"createdAt"`
			UpdatedAt string `json:"updatedAt"`
		} `json:"schedules"`
	}{
		Schedules: []struct {
			ID          string    `json:"id"`
			ClientName  string    `json:"clientName"`
			ShiftTime   string    `json:"shiftTime"`
			Location    string    `json:"location"`
			Status      string    `json:"status"`
			Visit       *struct {
				ID         string `json:"id"`
				ScheduleID string `json:"scheduleId"`
				Status     string `json:"status"`
			} `json:"visit,omitempty"`
			Tasks       []struct {
				ID          string `json:"id"`
				ScheduleID  string `json:"scheduleId"`
				Name        string `json:"name"`
				Description string `json:"description,omitempty"`
				Status      string `json:"status"`
				Reason      string `json:"reason,omitempty"`
				CompletedAt string `json:"completedAt,omitempty"`
				CreatedAt   string `json:"createdAt"`
				UpdatedAt   string `json:"updatedAt"`
			} `json:"tasks"`
			CreatedAt string `json:"createdAt"`
			UpdatedAt string `json:"updatedAt"`
		}{
			{
				ID:         "2",
				ClientName: "Jane Doe",
				ShiftTime:  "10:00 - 14:00",
				Location:   "456 Oak Ave, Somewhere",
				Status:     "in_progress",
				Visit: &struct {
					ID         string `json:"id"`
					ScheduleID string `json:"scheduleId"`
					Status     string `json:"status"`
				}{
					ID:         "visit2",
					ScheduleID: "2",
					Status:     "in_progress",
				},
				Tasks: []struct {
					ID          string `json:"id"`
					ScheduleID  string `json:"scheduleId"`
					Name        string `json:"name"`
					Description string `json:"description,omitempty"`
					Status      string `json:"status"`
					Reason      string `json:"reason,omitempty"`
					CompletedAt string `json:"completedAt,omitempty"`
					CreatedAt   string `json:"createdAt"`
					UpdatedAt   string `json:"updatedAt"`
				}{
					{
						ID:          "task2",
						ScheduleID:  "2",
						Name:        "Vital Signs Check",
						Status:      "in_progress",
						CreatedAt:   "2025-10-18T09:00:00Z",
						UpdatedAt:   "2025-10-18T09:00:00Z",
					},
				},
				CreatedAt: "2025-10-18T09:00:00Z",
				UpdatedAt: "2025-10-18T09:00:00Z",
			},
		},
	}

	return c.JSON(http.StatusOK, todaySchedules)
}

func GetScheduleById(c echo.Context) error {
	id := c.Param("id")

	var schedule struct {
		ID          string    `json:"id"`
		ClientName  string    `json:"clientName"`
		ShiftTime   string    `json:"shiftTime"`
		Location    string    `json:"location"`
		Status      string    `json:"status"`
		Visit       *struct {
			ID           string  `json:"id"`
			ScheduleID   string  `json:"scheduleId"`
			StartTime    string  `json:"startTime,omitempty"`
			EndTime      string  `json:"endTime,omitempty"`
			StartLocation *struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
				Accuracy  *int    `json:"accuracy,omitempty"`
				Timestamp *string `json:"timestamp,omitempty"`
			} `json:"startLocation,omitempty"`
			EndLocation  *struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
				Accuracy  *int    `json:"accuracy,omitempty"`
				Timestamp *string `json:"timestamp,omitempty"`
			} `json:"endLocation,omitempty"`
			Status         string  `json:"status"`
			DurationMinutes *int    `json:"durationMinutes,omitempty"`
			CreatedAt      string  `json:"createdAt"`
			UpdatedAt      string  `json:"updatedAt"`
		} `json:"visit,omitempty"`
		Tasks       []struct {
			ID          string `json:"id"`
			ScheduleID  string `json:"scheduleId"`
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
			Status      string `json:"status"`
			Reason      string `json:"reason,omitempty"`
			CompletedAt string `json:"completedAt,omitempty"`
			CreatedAt   string `json:"createdAt"`
			UpdatedAt   string `json:"updatedAt"`
		} `json:"tasks"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}

	switch id {
	case "1":
		schedule = struct {
			ID          string    `json:"id"`
			ClientName  string    `json:"clientName"`
			ShiftTime   string    `json:"shiftTime"`
			Location    string    `json:"location"`
			Status      string    `json:"status"`
			Visit       *struct {
				ID           string  `json:"id"`
				ScheduleID   string  `json:"scheduleId"`
				StartTime    string  `json:"startTime,omitempty"`
				EndTime      string  `json:"endTime,omitempty"`
				StartLocation *struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
					Accuracy  *int    `json:"accuracy,omitempty"`
					Timestamp *string `json:"timestamp,omitempty"`
				} `json:"startLocation,omitempty"`
				EndLocation  *struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
					Accuracy  *int    `json:"accuracy,omitempty"`
					Timestamp *string `json:"timestamp,omitempty"`
				} `json:"endLocation,omitempty"`
				Status         string  `json:"status"`
				DurationMinutes *int    `json:"durationMinutes,omitempty"`
				CreatedAt      string  `json:"createdAt"`
				UpdatedAt      string  `json:"updatedAt"`
			} `json:"visit,omitempty"`
			Tasks       []struct {
				ID          string `json:"id"`
				ScheduleID  string `json:"scheduleId"`
				Name        string `json:"name"`
				Description string `json:"description,omitempty"`
				Status      string `json:"status"`
				Reason      string `json:"reason,omitempty"`
				CompletedAt string `json:"completedAt,omitempty"`
				CreatedAt   string `json:"createdAt"`
				UpdatedAt   string `json:"updatedAt"`
			} `json:"tasks"`
			CreatedAt string `json:"createdAt"`
			UpdatedAt string `json:"updatedAt"`
		}{
			ID:         "1",
			ClientName: "John Smith",
			ShiftTime:  "09:00 - 12:00",
			Location:   "123 Main St, Anytown",
			Status:     "upcoming",
			Visit: nil, // TODO: Fix mock data structure
			Tasks: nil, // TODO: Fix mock data structure
			CreatedAt: "2025-10-18T08:00:00Z",
			UpdatedAt: "2025-10-18T12:00:00Z",
		}
	default:
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Schedule not found"})
	}

	return c.JSON(http.StatusOK, schedule)
}

func CreateSchedule(c echo.Context) error {
	var requestBody struct {
		ClientName string `json:"clientName"`
		ShiftTime  string `json:"shiftTime"`
		Location   string `json:"location"`
	}

	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Return mock created schedule
	createdSchedule := struct {
		ID          string    `json:"id"`
		ClientName  string    `json:"clientName"`
		ShiftTime   string    `json:"shiftTime"`
		Location    string    `json:"location"`
		Status      string    `json:"status"`
		Visit       *struct {
			ID         string `json:"id"`
			ScheduleID string `json:"scheduleId"`
			Status     string `json:"status"`
		} `json:"visit,omitempty"`
		Tasks       []struct {
			ID          string `json:"id"`
			ScheduleID  string `json:"scheduleId"`
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
			Status      string `json:"status"`
			Reason      string `json:"reason,omitempty"`
			CompletedAt string `json:"completedAt,omitempty"`
			CreatedAt   string `json:"createdAt"`
			UpdatedAt   string `json:"updatedAt"`
		} `json:"tasks"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		ID:        "3",
		ClientName: requestBody.ClientName,
		ShiftTime:  requestBody.ShiftTime,
		Location:   requestBody.Location,
		Status:     "upcoming",
		Visit: &struct {
			ID         string `json:"id"`
			ScheduleID string `json:"scheduleId"`
			Status     string `json:"status"`
		}{
			ID:         "visit3",
			ScheduleID: "3",
			Status:     "not_started",
		},
		Tasks:     []struct {
			ID          string `json:"id"`
			ScheduleID  string `json:"scheduleId"`
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
			Status      string `json:"status"`
			Reason      string `json:"reason,omitempty"`
			CompletedAt string `json:"completedAt,omitempty"`
			CreatedAt   string `json:"createdAt"`
			UpdatedAt   string `json:"updatedAt"`
		}{},
		CreatedAt: "2025-10-18T10:00:00Z",
		UpdatedAt: "2025-10-18T10:00:00Z",
	}

	return c.JSON(http.StatusCreated, createdSchedule)
}

func UpdateScheduleStatus(c echo.Context) error {
	id := c.Param("id")

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.Bind(&statusUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Return mock updated schedule
	updatedSchedule := struct {
		ID          string    `json:"id"`
		ClientName  string    `json:"clientName"`
		ShiftTime   string    `json:"shiftTime"`
		Location    string    `json:"location"`
		Status      string    `json:"status"`
		Visit       *struct {
			ID         string `json:"id"`
			ScheduleID string `json:"scheduleId"`
			Status     string `json:"status"`
		} `json:"visit,omitempty"`
		Tasks       []struct {
			ID          string `json:"id"`
			ScheduleID  string `json:"scheduleId"`
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
			Status      string `json:"status"`
			Reason      string `json:"reason,omitempty"`
			CompletedAt string `json:"completedAt,omitempty"`
			CreatedAt   string `json:"createdAt"`
			UpdatedAt   string `json:"updatedAt"`
		} `json:"tasks"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		ID:        id,
		ClientName: "John Smith",
		ShiftTime:  "09:00 - 12:00",
		Location:   "123 Main St, Anytown",
		Status:     statusUpdate.Status,
		Visit: &struct {
			ID         string `json:"id"`
			ScheduleID string `json:"scheduleId"`
			Status     string `json:"status"`
		}{
			ID:         "visit1",
			ScheduleID: "1",
			Status:     statusUpdate.Status,
		},
		Tasks:     []struct {
			ID          string `json:"id"`
			ScheduleID  string `json:"scheduleId"`
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
			Status      string `json:"status"`
			Reason      string `json:"reason,omitempty"`
			CompletedAt string `json:"completedAt,omitempty"`
			CreatedAt   string `json:"createdAt"`
			UpdatedAt   string `json:"updatedAt"`
		}{},
		CreatedAt: "2025-10-18T08:00:00Z",
		UpdatedAt: "2025-10-18T10:00:00Z",
	}

	return c.JSON(http.StatusOK, updatedSchedule)
}

func SearchSchedules(c echo.Context) error {
	// Mock search results
	searchResults := struct {
		Schedules []struct {
			ID          string    `json:"id"`
			ClientName  string    `json:"clientName"`
			ShiftTime   string    `json:"shiftTime"`
			Location    string    `json:"location"`
			Status      string    `json:"status"`
			Tasks       []struct {
				ID          string `json:"id"`
				ScheduleID  string `json:"scheduleId"`
				Name        string `json:"name"`
				Description string `json:"description,omitempty"`
				Status      string `json:"status"`
				Reason      string `json:"reason,omitempty"`
				CompletedAt string `json:"completedAt,omitempty"`
				CreatedAt   string `json:"createdAt"`
				UpdatedAt   string `json:"updatedAt"`
			} `json:"tasks"`
			CreatedAt string `json:"createdAt"`
			UpdatedAt string `json:"updatedAt"`
		} `json:"schedules"`
	}{
		Schedules: []struct {
			ID          string    `json:"id"`
			ClientName  string    `json:"clientName"`
			ShiftTime   string    `json:"shiftTime"`
			Location    string    `json:"location"`
			Status      string    `json:"status"`
			Tasks       []struct {
				ID          string `json:"id"`
				ScheduleID  string `json:"scheduleId"`
				Name        string `json:"name"`
				Description string `json:"description,omitempty"`
				Status      string `json:"status"`
				Reason      string `json:"reason,omitempty"`
				CompletedAt string `json:"completedAt,omitempty"`
				CreatedAt   string `json:"createdAt"`
				UpdatedAt   string `json:"updatedAt"`
			} `json:"tasks"`
			CreatedAt string `json:"createdAt"`
			UpdatedAt string `json:"updatedAt"`
		}{
			{
				ID:         "1",
				ClientName: "John Smith",
				ShiftTime:  "09:00 - 12:00",
				Location:   "123 Main St, Anytown",
				Status:     "upcoming",
				Tasks:      []struct {
					ID          string `json:"id"`
					ScheduleID  string `json:"scheduleId"`
					Name        string `json:"name"`
					Description string `json:"description,omitempty"`
					Status      string `json:"status"`
					Reason      string `json:"reason,omitempty"`
					CompletedAt string `json:"completedAt,omitempty"`
					CreatedAt   string `json:"createdAt"`
					UpdatedAt   string `json:"updatedAt"`
				}{},
				CreatedAt: "2025-10-18T08:00:00Z",
				UpdatedAt: "2025-10-18T08:00:00Z",
			},
		},
	}

	return c.JSON(http.StatusOK, searchResults)
}

func GetTaskStats(c echo.Context) error {
	mockTaskStats := map[string]interface{}{
		"total_tasks":                     15,
		"completed_tasks":                 11,
		"pending_tasks":                   3,
		"not_completed_tasks":              1,
		"completion_rate":                 73.3,
		"tasks_by_category": map[string]int{
			"medication":    4,
			"vital_signs":  3,
			"therapy":       2,
			"personal_care": 4,
			"other":         2,
		},
		"tasks_by_status": map[string]int{
			"completed":     11,
			"pending":      3,
			"not_completed": 1,
		},
		"average_completion_time_minutes": 45,
	}

	return c.JSON(http.StatusOK, mockTaskStats)
}