package routes

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/artnikel/vacancystats/internal/config"
	"github.com/artnikel/vacancystats/internal/logging"
	"github.com/artnikel/vacancystats/internal/model"
	"github.com/google/uuid"
)

type Storage interface {
	Create(ctx context.Context, vacancy *model.Vacancy) error
	GetAll(ctx context.Context) ([]model.Vacancy, error)
}

type Routes struct {
	storage Storage
	Logger  *logging.Logger
	Config  *config.Config
}

func NewRoutes(storage Storage, logger *logging.Logger, cfg *config.Config) *Routes {
	return &Routes{storage: storage, Logger: logger, Config: cfg}
}

func (r *Routes) Create(ctx context.Context) {
	newVacancy := &model.Vacancy{
		VacancyID:    uuid.New(),
		ResponseTime: time.Now(),
	}
	var resourceNumber int
	fmt.Println("\nselect resource:", r.Config.Resource.ResourceList)
	fmt.Fscan(os.Stdin, &resourceNumber)
	newVacancy.Resource = r.Config.Resource.ResourceList[resourceNumber]
	fmt.Println("\ncompany name:")
	fmt.Fscan(os.Stdin, &newVacancy.Company)
	var resultNumber int
	fmt.Println("\nselect result:", r.Config.Result.ResultList)
	fmt.Fscan(os.Stdin, &resultNumber)
	newVacancy.Result = r.Config.Result.ResultList[resultNumber]
	err := r.storage.Create(ctx, newVacancy)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
	}
}

func (r *Routes) GetStats(ctx context.Context) {
	vacancies, err := r.storage.GetAll(ctx)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
	}
	for _, vacancy := range vacancies {
		fmt.Println(vacancy)
	}

}
