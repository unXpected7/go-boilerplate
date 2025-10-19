package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/server"
)

type Repositories struct {
	Schedule *ScheduleRepository
	Visit     *VisitRepository
	Task      *TaskRepository
}

func NewRepositories(s *server.Server) *Repositories {
	var dbPool *pgxpool.Pool = s.DB.Pool
	_ = dbPool // Force usage of pgxpool import
	return &Repositories{
		Schedule: NewScheduleRepository(dbPool),
		Visit:     NewVisitRepository(dbPool),
		Task:      NewTaskRepository(dbPool),
	}
}
