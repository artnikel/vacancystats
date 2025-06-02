package routes

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/artnikel/vacancystats/internal/config"
	"github.com/artnikel/vacancystats/internal/constants"
	"github.com/artnikel/vacancystats/internal/logging"
	"github.com/artnikel/vacancystats/internal/model"
	"github.com/artnikel/vacancystats/internal/utils"
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

func checkCorrectInput(selection string, list []string) int {
	number := 0
	firstInp := false
	for number < 1 || number > len(list) {
		if firstInp {
			utils.ClearConsole()
			fmt.Printf("\nincorrect input:\n%d", number)
		}
		fmt.Printf("\nselect %s:\n%s", selection, listToText(list))
		fmt.Fscan(os.Stdin, &number)
		firstInp = true
	}
	return number
}

func (r *Routes) Create(ctx context.Context) {
	newVacancy := &model.Vacancy{
		VacancyID:    uuid.New(),
		ResponseTime: time.Now(),
	}
	number := checkCorrectInput(constants.ResourceTypeInput, r.Config.Resource.ResourceList)
	newVacancy.Resource = r.Config.Resource.ResourceList[number-1]
	fmt.Println("\ncompany name:")
	fmt.Fscan(os.Stdin, &newVacancy.Company)
	number = checkCorrectInput(constants.StatusTypeInput, r.Config.Status.StatusList)
	newVacancy.Status = r.Config.Status.StatusList[number-1]
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
		fmt.Printf("%v", err)
	} else {
		fmt.Println("vacancy deleted")
	}
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
		fmt.Printf("%v", err)
	} else {
		fmt.Println("vacancy status updated")
	}
}
