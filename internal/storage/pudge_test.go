package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/artnikel/vacancystats/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (db *Pudge, cleanup func()) {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "testdb")

	p, err := NewPudge(dbPath)
	require.NoError(t, err)

	cleanup = func() {
		_ = os.RemoveAll(tmpDir)
	}

	return p, cleanup
}

func TestCreateAndGetAll(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	vac := &model.Vacancy{
		VacancyID:    uuid.New(),
		Resource:     "linkedin.com",
		Company:      "OpenAI",
		Status:       "viewed",
		ResponseTime: time.Now(),
	}

	err := store.Create(ctx, vac)
	require.NoError(t, err)

	vacancies, err := store.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, vacancies, 1)
	require.Equal(t, vac.VacancyID, vacancies[0].VacancyID)
	require.Equal(t, vac.Company, vacancies[0].Company)
}

func TestDelete(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	vac := &model.Vacancy{
		VacancyID:    uuid.New(),
		Resource:     "rabota.by",
		Company:      "Google",
		Status:       "interview",
		ResponseTime: time.Now(),
	}

	err := store.Create(ctx, vac)
	require.NoError(t, err)

	err = store.Delete(ctx, vac.VacancyID)
	require.NoError(t, err)

	vacancies, err := store.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, vacancies, 0)
}

func TestUpdateStatus(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	vac := &model.Vacancy{
		VacancyID:    uuid.New(),
		Resource:     "headhunter.ru",
		Company:      "Netflix",
		Status:       "rejection",
		ResponseTime: time.Now(),
	}

	err := store.Create(ctx, vac)
	require.NoError(t, err)

	vac.Status = "Interview"
	err = store.UpdateStatus(ctx, vac)
	require.NoError(t, err)

	vacancies, err := store.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, vacancies, 1)
	require.Equal(t, "Interview", vacancies[0].Status)
}
