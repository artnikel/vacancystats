package routes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/artnikel/vacancystats/internal/config"
	"github.com/artnikel/vacancystats/internal/logging"
	"github.com/artnikel/vacancystats/internal/model"
	"github.com/artnikel/vacancystats/internal/routes/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	code := m.Run()
	_ = os.RemoveAll("testlog")
	os.Exit(code)
}

func setup(t *testing.T) (*Routes, *mocks.MockStorage) {
	mockStorage := mocks.NewMockStorage(t)
	cfg := &config.Config{
		Resource: config.ResourceConfig{ResourceList: []string{"headhunter", "linkedin"}},
		Status:   config.StatusConfig{StatusList: []string{"waiting", "interview"}},
	}
	logger, _ := logging.NewLogger("testlog")
	return NewRoutes(mockStorage, logger, cfg), mockStorage
}

func withInput(input string) func() {
	old := os.Stdin
	r, w, _ := os.Pipe()
	w.WriteString(input)
	w.Close()
	os.Stdin = r
	return func() { os.Stdin = old }
}

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	return buf.String()
}

func TestGetStats(t *testing.T) {
	r, mockStorage := setup(t)

	mockStorage.EXPECT().GetAll(mock.Anything).Return([]*model.Vacancy{
		{Status: "waiting"},
		{Status: "waiting"},
		{Status: "interview"},
	}, nil)

	output := captureOutput(func() {
		r.GetStats(context.Background())
	})

	require.Contains(t, output, "waiting")
	require.Contains(t, output, "interview")
}

func TestDelete(t *testing.T) {
	r, mockStorage := setup(t)

	id := uuid.New()
	defer withInput(id.String() + "\n")()

	mockStorage.EXPECT().Delete(mock.Anything, id).Return(nil)

	r.Delete(context.Background())
}

func TestUpdateStatus(t *testing.T) {
	r, mockStorage := setup(t)

	id := uuid.New()
	input := fmt.Sprintf("%s\n1\n", id.String())
	defer withInput(input)()

	mockStorage.EXPECT().UpdateStatus(mock.Anything, mock.MatchedBy(func(v *model.Vacancy) bool {
		return v.VacancyID == id && v.Status == "waiting"
	})).Return(nil)

	r.UpdateStatus(context.Background())
}
