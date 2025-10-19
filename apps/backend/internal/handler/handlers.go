package handler

import (
	"github.com/sriniously/go-boilerplate/apps/backend/internal/server"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type MockAPIHandler struct {
	GetMockSchedules    func(echo.Context) error
	GetTodaySchedules    func(echo.Context) error
	GetScheduleById     func(echo.Context) error
	CreateSchedule      func(echo.Context) error
	UpdateScheduleStatus func(echo.Context) error
	SearchSchedules     func(echo.Context) error
	GetMockStats        func(echo.Context) error
	GetTaskStats         func(echo.Context) error
}

type Handlers struct {
	Health    *HealthHandler
	OpenAPI   *OpenAPIHandler
	EVV       *EVVHandler
	Swagger   *SwaggerHandler
	Mock      *MockAPIHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:    NewHealthHandler(s),
		OpenAPI:   NewOpenAPIHandler(s),
		EVV:       NewEVVHandler(services.ScheduleService, services.VisitService, services.TaskService),
		Swagger:   NewSwaggerHandler(),
		Mock: &MockAPIHandler{
			GetMockSchedules:    GetMockSchedules,
			GetTodaySchedules:    GetTodaySchedules,
			GetScheduleById:     GetScheduleById,
			CreateSchedule:      CreateSchedule,
			UpdateScheduleStatus: UpdateScheduleStatus,
			SearchSchedules:     SearchSchedules,
			GetMockStats:        GetMockStats,
			GetTaskStats:         GetTaskStats,
		},
	}
}
