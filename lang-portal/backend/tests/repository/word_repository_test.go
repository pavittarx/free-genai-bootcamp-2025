package repository_test

import (
"github.com/pavittarx/lang-portal/backend/tests/testutils"
	"context"
	"testing"
	"time"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
	"github.com/stretchr/testify/assert"
)

func setupWordRepositoryTest(t *testing.T) (*repository.SQLiteWordRepository, func()) {
	db, cleanup, err := testutils.CreateTestDB()
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create words table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hindi TEXT NOT NULL,
			scrambled TEXT NOT NULL,
			hinglish TEXT NOT NULL,
			english TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		cleanup()
		t.Fatalf("Failed to create words table: %v", err)
	}

	// Clear existing data
	_, err = db.Exec(`DELETE FROM words`)
	if err != nil {
		cleanup()
		t.Fatalf("Failed to clear words table: %v", err)
	}

	repo := repository.NewSQLiteWordRepository(db)
	return repo, cleanup
}

func createTestWord() *models.Word {
	return &models.Word{
		Hindi:     "नमस्ते",
		Hinglish:  "Namaste",
		English:   "Hello",
		CreatedAt: time.Now(),
	}
}

func TestWordRepository_Create(t *testing.T) {
	repo, cleanup := setupWordRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	word := createTestWord()
	err := repo.Create(ctx, word)
	assert.NoError(t, err)
	assert.NotZero(t, word.ID)
	assert.NotEmpty(t, word.Scrambled)
}

func TestWordRepository_CreateValidation(t *testing.T) {
	repo, cleanup := setupWordRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	tests := []struct {
		name    string
		word    *models.Word
		wantErr bool
	}{
		{
			name: "valid word",
			word: createTestWord(),
			wantErr: false,
		},
		{
			name: "empty hindi",
			word: &models.Word{
				Hindi:    "",
				Hinglish: "Namaste",
				English:  "Hello",
			},
			wantErr: true,
		},
		{
			name: "empty english",
			word: &models.Word{
				Hindi:    "नमस्ते",
				Hinglish: "Namaste",
				English:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.word)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWordRepository_GetByID(t *testing.T) {
	repo, cleanup := setupWordRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	// First create a word
	word := createTestWord()
	err := repo.Create(ctx, word)
	assert.NoError(t, err)

	// Then retrieve it
	retrievedWord, err := repo.GetByID(ctx, word.ID)
	assert.NoError(t, err)
	assert.Equal(t, word.Hindi, retrievedWord.Hindi)
	assert.Equal(t, word.English, retrievedWord.English)
	assert.Equal(t, word.Hinglish, retrievedWord.Hinglish)
	assert.NotZero(t, retrievedWord.CreatedAt)

	// Try to get non-existing word
	_, err = repo.GetByID(ctx, 999)
	assert.Error(t, err)
}

func TestWordRepository_Update(t *testing.T) {
	repo, cleanup := setupWordRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	// First create a word
	word := createTestWord()
	err := repo.Create(ctx, word)
	assert.NoError(t, err)

	// Update the word
	word.Hindi = "अलविदा"
	word.English = "Goodbye"
	err = repo.Update(ctx, word)
	assert.NoError(t, err)

	// Retrieve and verify
	updatedWord, err := repo.GetByID(ctx, word.ID)
	assert.NoError(t, err)
	assert.Equal(t, "अलविदा", updatedWord.Hindi)
	assert.Equal(t, "Goodbye", updatedWord.English)

	// Try to update non-existing word
	nonExistingWord := &models.Word{
		ID:       999,
		Hindi:    "Test",
		English:  "Test",
		Hinglish: "Test",
	}
	err = repo.Update(ctx, nonExistingWord)
	assert.Error(t, err)
}

func TestWordRepository_Delete(t *testing.T) {
	repo, cleanup := setupWordRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	// First create a word
	word := createTestWord()
	err := repo.Create(ctx, word)
	assert.NoError(t, err)

	// Delete the word
	err = repo.Delete(ctx, word.ID)
	assert.NoError(t, err)

	// Try to retrieve deleted word
	_, err = repo.GetByID(ctx, word.ID)
	assert.Error(t, err)

	// Try to delete non-existing word
	err = repo.Delete(ctx, 999)
	assert.Error(t, err)
}

func TestWordRepository_List(t *testing.T) {
	repo, cleanup := setupWordRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple words
	words := []models.Word{
		{
			Hindi:     "नमस्ते",
			Hinglish:  "Namaste",
			English:   "Hello",
			CreatedAt: time.Now(),
		},
		{
			Hindi:     "अलविदा",
			Hinglish:  "Alvida",
			English:   "Goodbye",
			CreatedAt: time.Now(),
		},
		{
			Hindi:     "धन्यवाद",
			Hinglish:  "Dhanyavaad",
			English:   "Thank you",
			CreatedAt: time.Now(),
		},
	}

	for i := range words {
		err := repo.Create(ctx, &words[i])
		assert.NoError(t, err)
	}

	// Test listing all words
	params := repository.ListWordsParams{
		Page:     1,
		PageSize: 10,
	}
	listedWords, total, err := repo.List(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, len(words), total)
	assert.Equal(t, len(words), len(listedWords))

	// Test pagination
	params.PageSize = 2
	listedWords, total, err = repo.List(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, len(words), total)
	assert.Equal(t, 2, len(listedWords))

	// Test search
	params.Search = "नमस्ते"
	listedWords, total, err = repo.List(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(listedWords))
	assert.Equal(t, "नमस्ते", listedWords[0].Hindi)
}

func TestWordRepository_GetRandomWord(t *testing.T) {
	repo, cleanup := setupWordRepositoryTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple words
	words := []models.Word{
		{
			Hindi:     "नमस्ते",
			Hinglish:  "Namaste",
			English:   "Hello",
			CreatedAt: time.Now(),
		},
		{
			Hindi:     "अलविदा",
			Hinglish:  "Alvida",
			English:   "Goodbye",
			CreatedAt: time.Now(),
		},
	}

	for i := range words {
		err := repo.Create(ctx, &words[i])
		assert.NoError(t, err)
	}

	// Get random word
	randomWord, err := repo.GetRandomWord(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, randomWord)
	assert.NotEmpty(t, randomWord.Hindi)
	assert.NotEmpty(t, randomWord.English)
}
