package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/handlers"
	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
	"github.com/pavittarx/lang-portal/backend/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/pavittarx/lang-portal/backend/tests/testutils"
)

func setupTest(t *testing.T) (*echo.Echo, *handlers.WordHandler, func()) {
	e := echo.New()
	db, cleanup, err := testutils.CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}

	repo := repository.NewSQLiteWordRepository(db)
	service := services.NewWordService(repo)
	handler := handlers.NewWordHandler(service, repo)

	return e, handler, cleanup
}

func TestWordHandler_CreateWord(t *testing.T) {
	e, handler, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name       string
		word       models.Word
		wantStatus int
	}{
		{
			name: "valid word",
			word: models.Word{
				Hindi:    "नमस्ते",
				English:  "Hello",
				Hinglish: "Namaste",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid word",
			word: models.Word{
				Hindi:    "",
				English:  "Hello",
				Hinglish: "Namaste",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, _ := json.Marshal(tt.word)
			req := httptest.NewRequest(http.MethodPost, "/words", bytes.NewReader(jsonBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.CreateWord(c)
			if err != nil {
				t.Errorf("WordHandler.CreateWord() error = %v", err)
				return
			}

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestWordHandler_GetWordByID(t *testing.T) {
	e, handler, cleanup := setupTest(t)
	defer cleanup()

	// First create a word to test with
	word := &models.Word{
		Hindi:    "नमस्ते",
		English:  "Hello",
		Hinglish: "Namaste",
	}
	jsonBytes, _ := json.Marshal(word)
	req := httptest.NewRequest(http.MethodPost, "/words", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler.CreateWord(c)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "existing word",
			id:         "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing word",
			id:         "999",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid id",
			id:         "abc",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/words/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			err := handler.GetWordByID(c)
			if err != nil {
				t.Errorf("WordHandler.GetWordByID() error = %v", err)
				return
			}

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestWordHandler_ListWords(t *testing.T) {
	e, handler, cleanup := setupTest(t)
	defer cleanup()

	// Add some test words
	words := []models.Word{
		{Hindi: "नमस्ते", English: "Hello", Hinglish: "Namaste"},
		{Hindi: "धन्यवाद", English: "Thank you", Hinglish: "Dhanyavaad"},
	}

	for _, w := range words {
		jsonBytes, _ := json.Marshal(w)
		req := httptest.NewRequest(http.MethodPost, "/words", bytes.NewReader(jsonBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler.CreateWord(c)
	}

	// Test listing words
	req := httptest.NewRequest(http.MethodGet, "/words", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.ListWords(c)
	if err != nil {
		t.Errorf("WordHandler.ListWords() error = %v", err)
		return
	}

	assert.Equal(t, http.StatusOK, rec.Code)

	var response struct {
		Words []models.Word `json:"words"`
		Total int          `json:"total"`
	}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, len(words), len(response.Words))
}

func TestWordHandler_GetRandomWord(t *testing.T) {
	e, handler, cleanup := setupTest(t)
	defer cleanup()

	// Add some test words
	words := []models.Word{
		{Hindi: "नमस्ते", English: "Hello", Hinglish: "Namaste"},
		{Hindi: "धन्यवाद", English: "Thank you", Hinglish: "Dhanyavaad"},
	}

	for _, w := range words {
		jsonBytes, _ := json.Marshal(w)
		req := httptest.NewRequest(http.MethodPost, "/words", bytes.NewReader(jsonBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler.CreateWord(c)
	}

	// Test getting random word
	req := httptest.NewRequest(http.MethodGet, "/words/random", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetRandomWord(c)
	if err != nil {
		t.Errorf("WordHandler.GetRandomWord() error = %v", err)
		return
	}

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Word
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NotEmpty(t, response)
}
