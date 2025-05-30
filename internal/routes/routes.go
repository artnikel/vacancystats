package routes

import (
	"context"
	"fmt"
	"log"

	"github.com/artnikel/vacancystats/internal/model"
)

type Storage interface {
	Create(ctx context.Context, vacancy *model.Vacancy) error
}

type Routes struct {
	storage Storage
}

func NewRoutes(storage Storage) *Routes {
	return &Routes{storage: storage}
}

func (r *Routes) Create(ctx context.Context) {
	newVacancy := &model.Vacancy{}
	fmt.Scanf("%s\n", newVacancy.Resource)
	fmt.Scanf("%s\n", newVacancy.Company)
	fmt.Scanf("%s\n", newVacancy.Result)
	err := r.storage.Create(ctx, newVacancy)
	if err != nil {
		log.Printf("%v", err)
	}
}
