// Package routes handles CLI interaction and delegates storage operations
package routes

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/artnikel/vacancystats/internal/config"
	"github.com/artnikel/vacancystats/internal/constants"
	"github.com/artnikel/vacancystats/internal/logging"
	"github.com/artnikel/vacancystats/internal/model"
	"github.com/artnikel/vacancystats/internal/utils"
	"github.com/google/uuid"
)

// Storage defines the interface for storage layer interactions
type Storage interface {
	Create(ctx context.Context, vacancy *model.Vacancy) error
	GetAll(ctx context.Context) ([]*model.Vacancy, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, vacancy *model.Vacancy) error
}

// Routes holds dependencies and provides CLI-based operations
type Routes struct {
	storage Storage
	Logger  *logging.Logger
	Config  *config.Config
}

// NewRoutes creates and returns a new Routes instance
func NewRoutes(storage Storage, logger *logging.Logger, cfg *config.Config) *Routes {
	return &Routes{storage: storage, Logger: logger, Config: cfg}
}

// listToText converts a list of strings to a numbered string list
func listToText(list []string) (text string) {
	for index, value := range list {
		text += strconv.Itoa(index+1) + ") " + value + "\n"
	}
	return text
}

// checkCorrectInput validates numeric user input for selection from a list
func checkCorrectInput(selection string, list []string) int {
	number := 0
	firstInp := false
	for number < 1 || number > len(list) {
		if firstInp {
			utils.ClearConsole()
			fmt.Printf("\nincorrect input:\n%d", number)
		}
		fmt.Printf("\nselect %s:\n%s", selection, listToText(list))
		_, err := fmt.Fscan(os.Stdin, &number)
		if err != nil {
			fmt.Printf("\ninput error:\n%v", err)
		}
		firstInp = true
	}
	return number
}

// Create prompts user input and creates a new vacancy record
func (r *Routes) Create(ctx context.Context) {
	newVacancy := &model.Vacancy{
		VacancyID:    uuid.New(),
		ResponseTime: time.Now(),
	}
	number := checkCorrectInput(constants.ResourceTypeInput, r.Config.Resource.ResourceList)
	newVacancy.Resource = r.Config.Resource.ResourceList[number-1]
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\ncompany name:")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\ninput error:\n%v", err)
	} else {
		newVacancy.Company = strings.TrimSpace(input)
	}
	number = checkCorrectInput(constants.StatusTypeInput, r.Config.Status.StatusList)
	newVacancy.Status = r.Config.Status.StatusList[number-1]
	err = r.storage.Create(ctx, newVacancy)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
	}
	fmt.Println("vacancy info added")
}

// GetResponses retrieves and prints all vacancy records
func (r *Routes) GetResponses(ctx context.Context) {
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

// GetStats calculates response statistics in percentages
func (r *Routes) GetStats(ctx context.Context) {
	vacancies, err := r.storage.GetAll(ctx)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
	}
	countMap := make(map[string]int, len(r.Config.Status.StatusList))
	for _, vacancy := range vacancies {
		countMap[vacancy.Status]++
	}
	totalCount := len(vacancies)
	fmt.Printf("\ntotal count of responses: %d\n", totalCount)
	for status, count := range countMap {
		fmt.Printf("%s - %.2f%% \n", status, float64(count)/float64(totalCount)*100) //nolint:mnd // percentage calculation
	}
}

// Delete requests for ID and removes the corresponding vacancy
func (r *Routes) Delete(ctx context.Context) {
	var idText string
	fmt.Println("\ninput id:")
	_, err := fmt.Fscan(os.Stdin, &idText)
	if err != nil {
		fmt.Printf("\ninput error:\n%v", err)
	}
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

// UpdateStatus requests for ID and updates the status of a vacancy
func (r *Routes) UpdateStatus(ctx context.Context) {
	var updVacancy model.Vacancy
	var idText string
	fmt.Println("\ninput id:")
	_, err := fmt.Fscan(os.Stdin, &idText)
	if err != nil {
		fmt.Printf("\ninput error:\n%v", err)
	}
	id, err := uuid.Parse(idText)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
	}
	updVacancy.VacancyID = id
	number := checkCorrectInput(constants.StatusTypeInput, r.Config.Status.StatusList)
	updVacancy.Status = r.Config.Status.StatusList[number-1]
	err = r.storage.UpdateStatus(ctx, &updVacancy)
	if err != nil {
		r.Logger.Error.Printf("%v", err)
		fmt.Printf("%v", err)
	} else {
		fmt.Println("vacancy status updated")
	}
}
