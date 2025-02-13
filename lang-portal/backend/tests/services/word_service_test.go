package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
	"github.com/pavittarx/lang-portal/backend/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWordRepository is a mock implementation of WordRepository
type MockWordRepository struct {
	mock.Mock
}

func (m *MockWordRepository) Create(ctx context.Context, word *models.Word) error {
	args := m.Called(ctx, word)
	return args.Error(0)
}

func (m *MockWordRepository) GetByID(ctx context.Context, id int64) (*models.Word, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Word), args.Error(1)
}

func (m *MockWordRepository) Update(ctx context.Context, word *models.Word) error {
	args := m.Called(ctx, word)
	return args.Error(0)
}

func (m *MockWordRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWordRepository) List(ctx context.Context, params repository.ListWordsParams) ([]models.Word, int, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]models.Word), args.Int(1), args.Error(2)
}

func (m *MockWordRepository) GetRandomWord(ctx context.Context) (*models.Word, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Word), args.Error(1)
}

func (m *MockWordRepository) GetWordsByGroupID(ctx context.Context, groupID int64) ([]models.Word, error) {
	args := m.Called(ctx, groupID)
	return args.Get(0).([]models.Word), args.Error(1)
}

func createTestWord() *models.Word {
	return &models.Word{
		ID:        1,
		Hindi:     "नमस्ते",
		Hinglish:  "Namaste",
		English:   "Hello",
		CreatedAt: time.Now(),
	}
}

func TestWordService_CreateWord(t *testing.T) {
	mockRepo := new(MockWordRepository)
	service := services.NewWordService(mockRepo)

	ctx := context.Background()
	word := createTestWord()

	// Set up expectations
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Word")).Return(nil)

	// Call the method
	err := service.CreateWord(ctx, word)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	assert.NotEmpty(t, word.Scrambled)
}

func TestWordService_CreateWord_Validation(t *testing.T) {
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
			mockRepo := new(MockWordRepository)
			ctx := context.Background()

			if !tt.wantErr {
				mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Word")).Return(nil)
			}

			err := services.NewWordService(mockRepo).CreateWord(ctx, tt.word)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWordService_GetWordByID(t *testing.T) {
	mockRepo := new(MockWordRepository)
	service := services.NewWordService(mockRepo)

	ctx := context.Background()
	word := createTestWord()

	// Set up expectations
	mockRepo.On("GetByID", ctx, word.ID).Return(word, nil)

	// Call the method
	retrievedWord, err := service.GetWordByID(ctx, word.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, word, retrievedWord)
	mockRepo.AssertExpectations(t)
}

func TestWordService_UpdateWord(t *testing.T) {
	mockRepo := new(MockWordRepository)
	service := services.NewWordService(mockRepo)

	ctx := context.Background()
	existingWord := createTestWord()
	updatedWord := &models.Word{
		ID:       existingWord.ID,
		Hindi:    "अलविदा",
		Hinglish: "Alvida",
		English:  "Goodbye",
	}

	// Set up expectations
	mockRepo.On("GetByID", ctx, existingWord.ID).Return(existingWord, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*models.Word")).Return(nil)

	// Call the method
	err := service.UpdateWord(ctx, updatedWord)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWordService_UpdateWord_Validation(t *testing.T) {
	tests := []struct {
		name    string
		word    *models.Word
		wantErr bool
	}{
		{
			name: "valid update",
			word: &models.Word{
				ID:       1,
				Hindi:    "नमस्ते",
				English:  "Hello",
				Hinglish: "Namaste",
			},
			wantErr: false,
		},
		{
			name: "empty word ID",
			word: &models.Word{
				ID: 0,
			},
			wantErr: true,
		},
		{
			name: "empty hindi",
			word: &models.Word{
				ID:       1,
				Hinglish: "Alvida",
				English:  "Goodbye",
			},
			wantErr: true,
		},
		{
			name: "empty english",
			word: &models.Word{
				ID:       1,
				Hindi:    "अलविदा",
				Hinglish: "Alvida",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockWordRepository)
			ctx := context.Background()

			if !tt.wantErr {
				existingWord := createTestWord()
				mockRepo.On("GetByID", ctx, tt.word.ID).Return(existingWord, nil)
				mockRepo.On("Update", ctx, mock.AnythingOfType("*models.Word")).Return(nil)
			}

			err := services.NewWordService(mockRepo).UpdateWord(ctx, tt.word)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWordService_DeleteWord(t *testing.T) {
	mockRepo := new(MockWordRepository)
	service := services.NewWordService(mockRepo)

	ctx := context.Background()
	word := createTestWord()

	// Set up expectations
	mockRepo.On("GetByID", ctx, word.ID).Return(word, nil)
	mockRepo.On("Delete", ctx, word.ID).Return(nil)

	// Call the method
	err := service.DeleteWord(ctx, word.ID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWordService_ListWords(t *testing.T) {
	mockRepo := new(MockWordRepository)
	service := services.NewWordService(mockRepo)

	ctx := context.Background()
	words := []models.Word{
		*createTestWord(),
		{
			ID:        2,
			Hindi:     "अलविदा",
			Hinglish:  "Alvida",
			English:   "Goodbye",
			CreatedAt: time.Now(),
		},
	}

	params := repository.ListWordsParams{
		Page:     1,
		PageSize: 10,
	}

	// Set up expectations
	mockRepo.On("List", ctx, params).Return(words, len(words), nil)

	// Call the method
	listedWords, total, err := service.ListWords(ctx, params)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, len(words), total)
	assert.Equal(t, words, listedWords)
	mockRepo.AssertExpectations(t)
}

func TestWordService_SearchWords(t *testing.T) {
	mockRepo := new(MockWordRepository)
	service := services.NewWordService(mockRepo)

	ctx := context.Background()
	words := []models.Word{
		*createTestWord(),
	}

	// Set up expectations
	mockRepo.On("List", ctx, repository.ListWordsParams{
		Search:   "नमस्ते",
		Language: "hindi",
		Page:     1,
		PageSize: 50,
	}).Return(words, len(words), nil)

	// Call the method
	listedWords, total, err := service.SearchWords(ctx, "नमस्ते", "hindi")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, len(words), total)
	assert.Equal(t, words, listedWords)
	mockRepo.AssertExpectations(t)
}

func TestWordService_GetWordsByGroupID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockWordRepository)
	wordService := services.NewWordService(mockRepo)

	groupID := int64(1)
	expectedWords := []models.Word{
		{
			ID:       1,
			English:  "test1",
			Hindi:    "परीक्षण1",
			Hinglish: "test1",
		},
		{
			ID:       2,
			English:  "test2",
			Hindi:    "परीक्षण2",
			Hinglish: "test2",
		},
	}

	// Mock the repository method
	mockRepo.On("GetWordsByGroupID", ctx, groupID).Return(expectedWords, nil)

	// Call the service method
	words, err := wordService.GetWordsByGroupID(ctx, groupID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedWords, words)
	assert.Len(t, words, 2)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}
