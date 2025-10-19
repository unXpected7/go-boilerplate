package router

import (
	"github.com/sriniously/go-boilerplate/apps/backend/internal/handler"

	"github.com/labstack/echo/v4"
)

func registerSystemRoutes(r *echo.Echo, h *handler.Handlers) {
	r.GET("/status", h.Health.CheckHealth)

	r.Static("/static", "static")

	// Add EVV Swagger UI
	r.GET("/docs", h.Swagger.ServeSwaggerUI)
	r.GET("/docs/swagger.json", h.Swagger.ServeOpenAPISpec)
}

func registerEVVRoutes(r *echo.Echo, h *handler.Handlers) {
	// Schedule endpoints
	r.GET("/api/v1/schedules", h.EVV.GetSchedules)
	r.GET("/api/v1/schedules/today", h.EVV.GetTodaySchedules)
	r.GET("/api/v1/schedules/:id", h.EVV.GetScheduleById)
	r.POST("/api/v1/schedules", h.EVV.CreateSchedule)
	r.PATCH("/api/v1/schedules/:id/status", h.EVV.UpdateScheduleStatus)
	r.GET("/api/v1/schedules/stats", h.EVV.GetMockStats)
	r.GET("/api/v1/schedules/search", h.EVV.SearchSchedules)

	// Visit tracking endpoints - simplified for now
	r.POST("/api/v1/schedules/:id/start", h.EVV.CreateSchedule)  // Temp: use create schedule
	r.POST("/api/v1/schedules/:id/end", h.EVV.CreateSchedule)     // Temp: use create schedule
	r.GET("/api/v1/schedules/:id/visit", h.EVV.GetScheduleById)   // Temp: return schedule data

	// Task management endpoints - simplified for now
	r.GET("/api/v1/schedules/:id/tasks", h.EVV.GetScheduleById)   // Temp: return schedule data
	r.POST("/api/v1/schedules/:id/tasks", h.EVV.CreateSchedule)  // Temp: use create schedule
	r.PATCH("/api/v1/schedules/:scheduleId/tasks/:taskId/status", h.EVV.UpdateTaskStatus)

	// Analytics endpoints
	r.GET("/api/v1/schedules/:id/analytics", h.EVV.GetMockStats)
	r.GET("/api/v1/tasks/stats", h.EVV.GetTaskStats)
}
