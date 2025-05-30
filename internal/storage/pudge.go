package repository

import (
	"context"
	"fmt"

	"github.com/artnikel/vacancystats/internal/model"
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

