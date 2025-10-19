package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Base
	ScheduleID    uuid.UUID `json:"scheduleId" db:"schedule_id"`
	Name          string    `json:"name" db:"name"`
	Description   *string   `json:"description" db:"description"`
	Status        string    `json:"status" db:"status"`
	Reason        *string   `json:"reason" db:"reason"`
	CompletedAt   *time.Time `json:"completedAt" db:"completed_at"`
}

type TaskCreate struct {
	ScheduleID  uuid.UUID `json:"scheduleId" db:"schedule_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
}

type TaskUpdate struct {
	Status      string    `json:"status" db:"status"`
	Reason      *string   `json:"reason" db:"reason"`
	CompletedAt *time.Time `json:"completedAt" db:"completed_at"`
}

func (t *Task) TableName() string {
	return "tasks"
}

type TaskStats struct {
	TotalTasks      int `json:"totalTasks"`
	CompletedTasks  int `json:"completedTasks"`
	PendingTasks    int `json:"pendingTasks"`
	NotCompletedTasks int `json:"notCompletedTasks"`
}

func (t *Task) IsCompleted() bool {
	return t.Status == "completed"
}

func (t *Task) IsNotCompleted() bool {
	return t.Status == "not_completed"
}

func (t *Task) IsPending() bool {
	return t.Status == "pending"
}