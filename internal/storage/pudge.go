// Package storage provides a Pudge-based implementation for storing vacancies
package storage

import (
	"context"
	"fmt"

	"github.com/artnikel/vacancystats/internal/model"
	"github.com/google/uuid"
	"github.com/recoilme/pudge"
)

// Pudge represents a wrapper around the Pudge key-value store
type Pudge struct {
	pool *pudge.Db
}

// NewPudge initializes and returns a new Pudge storage instance
func NewPudge(dbFolder string) (*Pudge, error) {
	cfg := &pudge.Config{SyncInterval: 0}
	dbVacancies, err := pudge.Open(dbFolder, cfg)
	if err != nil {
		return nil, err
	}
	return &Pudge{
		pool: dbVacancies,
	}, nil
}

// Create adds a new vacancy to the Pudge store
func (storage *Pudge) Create(ctx context.Context, vacancy *model.Vacancy) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%v", err)
	}
	return storage.pool.Set(vacancy.VacancyID, vacancy)
}

// GetAll retrieves all vacancies from the Pudge store
func (storage *Pudge) GetAll(ctx context.Context) ([]*model.Vacancy, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	keys, err := storage.pool.Keys(nil, 0, 0, true)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	vacancies := make([]*model.Vacancy, 0)
	for _, key := range keys {
		var vacancy model.Vacancy
		err := storage.pool.Get(key, &vacancy)
		if err != nil {
			continue
		}
		vacancies = append(vacancies, &vacancy)
	}
	return vacancies, nil
}

// Delete removes a vacancy from the Pudge store by ID
func (storage *Pudge) Delete(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%v", err)
	}
	return storage.pool.Delete(id)
}

// UpdateStatus updates the status field of an existing vacancy
func (storage *Pudge) UpdateStatus(ctx context.Context, vacancy *model.Vacancy) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%v", err)
	}
	var updVacancy model.Vacancy
	err := storage.pool.Get(vacancy.VacancyID, &updVacancy)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	updVacancy.Status = vacancy.Status
	return storage.pool.Set(vacancy.VacancyID, updVacancy)
}
