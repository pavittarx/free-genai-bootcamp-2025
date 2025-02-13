package handlers

import (
	"net/http"
	"strconv"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
	"github.com/pavittarx/lang-portal/backend/pkg/services"
)

// WordHandler handles HTTP requests related to words
type WordHandler struct {
	wordService *services.WordService
	wordRepo    *repository.SQLiteWordRepository
}

// NewWordHandler creates a new instance of WordHandler
func NewWordHandler(service *services.WordService, repo *repository.SQLiteWordRepository) *WordHandler {
	return &WordHandler{wordService: service, wordRepo: repo}
}

// CreateWord handles the creation of a new word
// @Summary Create a new word
// @Description Create a new word with details like Hindi, English, Hinglish, etc.
// @Tags words
// @Accept json
// @Produce json
// @Param word body models.Word true "Word details"
// @Success 201 {object} models.Word "Word created successfully"
// @Failure 400 {object} map[string]string "Invalid request payload or validation error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/words [post]
func (h *WordHandler) CreateWord(c echo.Context) error {
	// Create a new word instance
	word := &models.Word{}

	// Bind the request body to the word
	if err := c.Bind(word); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate the word
	if err := word.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Create the word
	if err := h.wordService.CreateWord(c.Request().Context(), word); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Return the created word
	return c.JSON(http.StatusCreated, word)
}

// GetWordByID retrieves a word by its ID
// @Summary Get a word by ID
// @Description Retrieve a word's details using its unique identifier
// @Tags words
// @Produce json
// @Param id path int true "Word ID"
// @Success 200 {object} models.Word "Word retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid word ID"
// @Failure 404 {object} map[string]string "Word not found"
// @Router /api/v1/words/{id} [get]
func (h *WordHandler) GetWordByID(c echo.Context) error {
	// Parse the ID from the URL parameter
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid word ID",
		})
	}

	// Retrieve the word
	word, err := h.wordService.GetWordByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, word)
}

// UpdateWord updates an existing word
// @Summary Update a word
// @Description Update an existing word's details by its ID
// @Tags words
// @Accept json
// @Produce json
// @Param id path int true "Word ID"
// @Param word body models.Word true "Updated word details"
// @Success 200 {object} models.Word "Word updated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload or word ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/words/{id} [put]
func (h *WordHandler) UpdateWord(c echo.Context) error {
	// Parse the ID from the URL parameter
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid word ID",
		})
	}

	// Create a new word instance with the ID
	word := &models.Word{ID: id}

	// Bind the request body to the word
	if err := c.Bind(word); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Update the word
	if err := h.wordService.UpdateWord(c.Request().Context(), word); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Return the updated word
	return c.JSON(http.StatusOK, word)
}

// DeleteWord removes a word by its ID
// @Summary Delete a word
// @Description Delete an existing word by its unique identifier
// @Tags words
// @Produce json
// @Param id path int true "Word ID"
// @Success 200 {object} map[string]string "Word deleted successfully"
// @Failure 400 {object} map[string]string "Invalid word ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/words/{id} [delete]
func (h *WordHandler) DeleteWord(c echo.Context) error {
	// Parse the ID from the URL parameter
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid word ID",
		})
	}

	// Delete the word
	if err := h.wordService.DeleteWord(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Word deleted successfully",
	})
}

// ListWords retrieves a list of words
// @Summary List words
// @Description Retrieve a paginated list of words with optional search and filtering
// @Tags words
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Param search query string false "Search term"
// @Param language query string false "Filter by language"
// @Success 200 {object} map[string]interface{} "Words retrieved successfully"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/words [get]
func (h *WordHandler) ListWords(c echo.Context) error {
	// Parse pagination parameters with defaults
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // Default page size
	}

	// Prepare list parameters
	params := repository.ListWordsParams{
		Page:     page,
		PageSize: pageSize,
		Search:   c.QueryParam("search"),
		Language: c.QueryParam("language"),
	}

	// Call service to list words
	words, totalCount, err := h.wordService.ListWords(c.Request().Context(), params)
	if err != nil {
		// Log the error for server-side debugging
		log.Printf("Error in ListWords handler: %v", err)

		// Return a more detailed error response
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to list words",
			"details": err.Error(),
		})
	}

	// Prepare response
	response := map[string]interface{}{
		"words":      words,
		"totalCount": totalCount,
		"page":       page,
		"pageSize":   pageSize,
	}

	return c.JSON(http.StatusOK, response)
}

// SearchWords provides a search endpoint for words
// @Summary Search words
// @Description Search for words using a query term and optional language filter
// @Tags words
// @Produce json
// @Param query query string true "Search query"
// @Param language query string false "Language filter"
// @Success 200 {object} map[string]interface{} "Search results retrieved successfully"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/words/search [get]
func (h *WordHandler) SearchWords(c echo.Context) error {
	// Get search query and language
	query := c.QueryParam("query")
	language := c.QueryParam("language")

	// Perform search
	words, totalCount, err := h.wordService.SearchWords(c.Request().Context(), query, language)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Return search results
	return c.JSON(http.StatusOK, map[string]interface{}{
		"words":      words,
		"totalCount": totalCount,
	})
}

// GetRandomWord handles the request to get a random word
// @Summary Get a random word
// @Description Retrieve a random word from the database
// @Tags words
// @Produce json
// @Success 200 {object} models.Word "Random word retrieved successfully"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/words/random [get]
func (h *WordHandler) GetRandomWord(c echo.Context) error {
	ctx := c.Request().Context()

	log.Printf("DEBUG: GetRandomWord method called")
	log.Printf("DEBUG: WordHandler repo: %+v", h.wordRepo)

	if h.wordRepo == nil {
		log.Printf("ERROR: wordRepo is nil")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Repository not initialized",
		})
	}

	log.Printf("DEBUG: Attempting to retrieve random word")

	word, err := h.wordRepo.GetRandomWord(ctx)
	if err != nil {
		log.Printf("ERROR retrieving random word: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve random word",
			"details": err.Error(),
		})
	}

	log.Printf("Retrieved random word: %+v", word)

	return c.JSON(http.StatusOK, word)
}

// GetWordsByGroupID handles the request to get all words in a specific group
func (h *WordHandler) GetWordsByGroupID(c echo.Context) error {
	// Parse group ID from URL parameter
	groupIDStr := c.Param("group-id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid group ID format",
		})
	}

	// Get words by group ID
	words, err := h.wordRepo.GetWordsByGroupID(c.Request().Context(), groupID)
	if err != nil {
		log.Printf("Error getting words by group ID: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve words for the group",
		})
	}

	// Return empty array if no words found
	if len(words) == 0 {
		return c.JSON(http.StatusOK, []models.Word{})
	}

	return c.JSON(http.StatusOK, words)
}
