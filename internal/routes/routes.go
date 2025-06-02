package routes

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/artnikel/vacancystats/internal/config"
	"github.com/artnikel/vacancystats/internal/logging"
	"github.com/artnikel/vacancystats/internal/model"
	"github.com/google/uuid"
)

type Storage interface {
	Create(ctx context.Context, vacancy *model.Vacancy) error
	GetAll(ctx context.Context) ([]model.Vacancy, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, vacancy *model.Vacancy) error
}

type Routes struct {
	storage Storage
	Logger  *logging.Logger
	Config  *config.Config
}

func NewRoutes(storage Storage, logger *logging.Logger, cfg *config.Config) *Routes {
	return &Routes{storage: storage, Logger: logger, Config: cfg}
}

func listToText(list []string) (text string) {
	for index, value := range list {
		text += strconv.Itoa(index+1) + ") " + value + "\n"
	}
	return text
}

func (r *Routes) Create(ctx context.Context) {
	newVacancy := &model.Vacancy{
		VacancyID:    uuid.New(),
		ResponseTime: time.Now(),
	}
	var resourceNumber int
	fmt.Printf("\nselect resource:\n%s", listToText(r.Config.Resource.ResourceList))
	fmt.Fscan(os.Stdin, &resourceNumber)
	newVacancy.Resource = r.Config.Resource.ResourceList[resourceNumber-1]
	fmt.Println("\ncompany name:")
	fmt.Fscan(os.Stdin, &newVacancy.Company)
	var statusNumber int
	fmt.Printf("\nselect status:\n%s", listToText(r.Config.Status.StatusList))
	fmt.Fscan(os.Stdin, &statusNumber)
	newVacancy.Status = r.Config.Status.StatusList[statusNumber-1]
	err := r.storage.Create(ctx, newVacancy)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
	}
	fmt.Println("vacancy info added")
}

func (r *Routes) GetStats(ctx context.Context) {
	vacancies, err := r.storage.GetAll(ctx)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
	}
	fmt.Println("id  /  resource  /  company  /  status  /  time")
	fmt.Println("-----------------------------------------------")
	for _, vacancy := range vacancies {
		fmt.Printf("%s / %s / %s / %s / %s\n",
			vacancy.VacancyID.String(),
			vacancy.Resource,
			vacancy.Company,
			vacancy.Status,
			vacancy.ResponseTime.Format(time.RFC1123))
	}
}

func (r *Routes) Delete(ctx context.Context) {
	var idText string
	fmt.Println("\ninput id:")
	fmt.Fscan(os.Stdin, &idText)
	id, err := uuid.Parse(idText)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
	}
	err = r.storage.Delete(ctx, id)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
	}
	fmt.Println("vacancy deleted")
}

func (r *Routes) UpdateStatus(ctx context.Context) {
	var updVacancy model.Vacancy
	var idText string
	fmt.Println("\ninput id:")
	fmt.Fscan(os.Stdin, &idText)
	id, err := uuid.Parse(idText)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
	}
	updVacancy.VacancyID = id
	var statusNumber int
	fmt.Printf("\nselect status:\n%s", listToText(r.Config.Status.StatusList))
	fmt.Fscan(os.Stdin, &statusNumber)
	updVacancy.Status = r.Config.Status.StatusList[statusNumber-1]
	err = r.storage.UpdateStatus(ctx, &updVacancy)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
	}
	fmt.Println("vacancy status updated")
}
