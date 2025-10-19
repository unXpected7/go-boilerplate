package service

import (
	"github.com/sriniously/go-boilerplate/apps/backend/internal/lib/job"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/repository"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/server"
)

type Services struct {
	Auth           *AuthService
	Job            *job.JobService
	ScheduleService *ScheduleService
	VisitService    *VisitService
	TaskService     *TaskService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	scheduleService := NewScheduleService(repos.Schedule, repos.Visit, repos.Task)
	visitService := NewVisitService(repos.Visit, repos.Schedule)
	taskService := NewTaskService(repos.Task, repos.Schedule)

	return &Services{
		Auth:           authService,
		Job:            s.Job,
		ScheduleService: scheduleService,
		VisitService:    visitService,
		TaskService:     taskService,
	}, nil
}
