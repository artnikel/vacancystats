// Package model provides the data models used in the application
package model

import (
	"time"

	"github.com/google/uuid"
)

// Vacancy struct consists of the fields for responding to a vacancy
type Vacancy struct {
	VacancyID    uuid.UUID
	Resource     string
	Company      string
	Status       string
	ResponseTime time.Time
}
