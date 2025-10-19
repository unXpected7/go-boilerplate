package model

import (
	"github.com/google/uuid"
)

type Schedule struct {
	Base
	ClientName string    `json:"clientName" db:"client_name"`
	ShiftTime  string    `json:"shiftTime" db:"shift_time"`
	Location   string    `json:"location" db:"location"`
	Status     string    `json:"status" db:"status"`
	VisitID    *uuid.UUID `json:"visitId" db:"visit_id"`
}

type ScheduleWithVisit struct {
	Schedule
	Visit *Visit `json:"visit,omitempty" db:"visit"`
}

type ScheduleWithTasks struct {
	Schedule
	Visit  *Visit  `json:"visit,omitempty" db:"visit"`
	Tasks  []Task  `json:"tasks" db:"tasks"`
}

func (s *Schedule) TableName() string {
	return "schedules"
}