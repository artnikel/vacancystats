package model

import (
	"time"

	"github.com/google/uuid"
)

type Vacancy struct {
	VacancyID    uuid.UUID
	Resource     string
	Company      string
	Result       string
	ResponseTime time.Time
}
