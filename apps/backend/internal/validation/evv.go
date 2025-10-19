package validation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Schedule validation structures
type CreateScheduleRequest struct {
	ClientName string `json:"clientName" validate:"required,min=2,max=255"`
	ShiftTime  string `json:"shiftTime" validate:"required,min=5,max=20"` // Format: HH:MM-HH:MM
	Location   string `json:"location" validate:"required,min=2,max=255"`
}

type UpdateScheduleRequest struct {
	ClientName *string `json:"clientName,omitempty" validate:"omitempty,min=2,max=255"`
	ShiftTime  *string `json:"shiftTime,omitempty" validate:"omitempty,min=5,max=20"`
	Location   *string `json:"location,omitempty" validate:"omitempty,min=2,max=255"`
}

type UpdateScheduleStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=missed upcoming in_progress completed"`
}

// Visit validation structures
type StartVisitRequest struct {
	StartTime time.Time `json:"startTime" validate:"required"`
	StartLat  float64   `json:"startLat" validate:"required,min=-90,max=90"`
	StartLong float64   `json:"startLong" validate:"required,min=-180,max=180"`
}

type EndVisitRequest struct {
	EndTime time.Time `json:"endTime" validate:"required"`
	EndLat  float64   `json:"endLat" validate:"required,min=-90,max=90"`
	EndLong float64   `json:"endLong" validate:"required,min=-180,max=180"`
}

// Task validation structures
type CreateTaskRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
}

type UpdateTaskStatusRequest struct {
	Status string  `json:"status" validate:"required,oneof=pending completed not_completed"`
	Reason *string `json:"reason,omitempty"`
}

// Pagination validation
type PaginationQuery struct {
	Page  int `query:"page" validate:"min=1"`
	Limit int `query:"limit" validate:"min=1,max=100"`
}

type SearchQuery struct {
	Query string `query:"q" validate:"required,min=1"`
	Page  int    `query:"page" validate:"min=1"`
	Limit int    `query:"limit" validate:"min=1,max=100"`
}

// Validation methods for request bodies
func (r *CreateScheduleRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	// Additional validation for shift time format
	if !isValidShiftTime(r.ShiftTime) {
		return &CustomValidationErrors{
			{Field: "shiftTime", Message: "Shift time must be in HH:MM-HH:MM format"},
		}
	}

	return nil
}

func (r *UpdateScheduleRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *UpdateScheduleStatusRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *StartVisitRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *EndVisitRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *CreateTaskRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *UpdateTaskStatusRequest) Validate() error {
	validate := validator.New()

	// Additional validation for reason when status is not_completed
	if r.Status == "not_completed" && r.Reason == nil {
		return &CustomValidationErrors{
			{Field: "reason", Message: "Reason is required for not_completed tasks"},
		}
	}

	return validate.Struct(r)
}

func (r *PaginationQuery) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *SearchQuery) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Helper functions
func isValidShiftTime(shiftTime string) bool {
	// Simple validation for HH:MM-HH:MM format
	return len(shiftTime) >= 5 && len(shiftTime) <= 20 &&
		   shiftTime[2] == ':' && shiftTime[5] == '-'
}

// isValidCoordinates validates latitude and longitude
func isValidCoordinates(lat, long float64) bool {
	return lat >= -90 && lat <= 90 && long >= -180 && long <= 180
}

// isValidVisitStatus validates visit status
func isValidVisitStatus(status string) bool {
	validStatuses := map[string]bool{
		"not_started":  true,
		"in_progress":  true,
		"completed":    true,
	}
	return validStatuses[status]
}

// ValidateUUID validates a UUID string
func ValidateUUID(uuid string) bool {
	return IsValidUUID(uuid)
}

// ValidateCoordinates validates latitude and longitude
func ValidateCoordinates(lat, long float64) bool {
	return lat >= -90 && lat <= 90 && long >= -180 && long <= 180
}

// ValidateTime validates that time is not too far in the past or future
func ValidateTime(t time.Time) error {
	now := time.Now()

	// Not more than 1 hour in the future
	if t.After(now.Add(1 * time.Hour)) {
		return &CustomValidationErrors{
			{Field: "time", Message: "Time cannot be more than 1 hour in the future"},
		}
	}

	// Not more than 5 minutes in the past (allowing for clock drift)
	if t.Before(now.Add(-5 * time.Minute)) {
		return &CustomValidationErrors{
			{Field: "time", Message: "Time cannot be in the past"},
		}
	}

	return nil
}

// Custom validation error helper
func NewValidationError(field, message string) *CustomValidationErrors {
	return &CustomValidationErrors{
		{Field: field, Message: message},
	}
}