package storage

import (
	"context"
	"fmt"

	"github.com/artnikel/vacancystats/internal/model"
	"github.com/google/uuid"
	"github.com/recoilme/pudge"
)

type Pudge struct {
	pool *pudge.Db
}

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

func (storage *Pudge) Create(ctx context.Context, vacancy *model.Vacancy) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("error in pudge Create: %v", err)
	}
	return storage.pool.Set(vacancy.VacancyID, vacancy)
}

func (storage *Pudge) GetAll(ctx context.Context) ([]model.Vacancy, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("error in pudge GetAll: %v", err)
	}
	keys, err := storage.pool.Keys(nil, 0, 0, true)
	if err != nil {
		return nil, fmt.Errorf("error in pudge Keys: %v", err)
	}
	var vacancies []model.Vacancy
	for _, key := range keys {
		var vacancy model.Vacancy
		err := storage.pool.Get(key, &vacancy)
		if err != nil {
			continue
		}
		vacancies = append(vacancies, vacancy)
	}
	return vacancies, nil
}

func (storage *Pudge) Delete(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("error in pudge Delete: %v", err)
	}
	return storage.pool.Delete(id)
}

func (storage *Pudge) UpdateStatus(ctx context.Context, vacancy *model.Vacancy) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("error in pudge UpdateStatus: %v", err)
	}
	var updVacancy model.Vacancy
	err := storage.pool.Get(vacancy.VacancyID, &updVacancy)
	if err != nil {
		return  fmt.Errorf("error in pudge Get: %v", err)
	}
	updVacancy.Status = vacancy.Status
	return storage.pool.Set(vacancy.VacancyID, updVacancy)
}
