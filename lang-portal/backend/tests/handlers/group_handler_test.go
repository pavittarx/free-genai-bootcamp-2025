package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pavittarx/lang-portal/backend/pkg/handlers"
	"github.com/pavittarx/lang-portal/backend/pkg/models"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
	"github.com/pavittarx/lang-portal/backend/pkg/services"
	"github.com/stretchr/testify/assert"
	"github.com/pavittarx/lang-portal/backend/tests/testutils"
)

func setupGroupTest(t *testing.T) (*echo.Echo, *handlers.GroupHandler, func()) {
	e := echo.New()
	db, cleanup, err := testutils.CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}

	repo := repository.NewSQLiteGroupRepository(db)
	service := services.NewGroupService(repo)
	handler := handlers.NewGroupHandler(service, repo)

	// Predefined groups for testing
	predefinedGroups := []models.Group{
		{Name: "Travel Words"},
		{Name: "Food Vocabulary"},
	}

	ctx := context.Background()
	for _, group := range predefinedGroups {
		err := repo.Create(ctx, &group)
		assert.NoError(t, err)
	}

	return e, handler, cleanup
}

func TestGroupHandler_CreateGroup(t *testing.T) {
	e, handler, cleanup := setupGroupTest(t)
	defer cleanup()

	tests := []struct {
		name       string
		group      models.Group
		wantStatus int
		wantGroup  bool
	}{
		{
			name: "valid group",
			group: models.Group{
				Name: "Test Group",
			},
			wantStatus: http.StatusCreated,
			wantGroup:  true,
		},
		{
			name: "invalid group - empty name",
			group: models.Group{
				Name: "",
			},
			wantStatus: http.StatusBadRequest,
			wantGroup:  false,
		},
		{
			name: "invalid group - too short name",
			group: models.Group{
				Name: "A",
			},
			wantStatus: http.StatusBadRequest,
			wantGroup:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, _ := json.Marshal(tt.group)
			req := httptest.NewRequest(http.MethodPost, "/groups", bytes.NewReader(jsonBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.CreateGroup(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantGroup {
				var createdGroup models.Group
				err := json.Unmarshal(rec.Body.Bytes(), &createdGroup)
				assert.NoError(t, err)
				assert.NotZero(t, createdGroup.ID)
				assert.Equal(t, tt.group.Name, createdGroup.Name)
			}
		})
	}
}

func TestGroupHandler_GetGroupByID(t *testing.T) {
	e, handler, cleanup := setupGroupTest(t)
	defer cleanup()

	// First create a group to test with
	group := &models.Group{
		Name: "Test Group",
	}
	jsonBytes, _ := json.Marshal(group)
	req := httptest.NewRequest(http.MethodPost, "/groups", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	err := handler.CreateGroup(c)
	assert.NoError(t, err)

	var createdGroup models.Group
	err = json.Unmarshal(rec.Body.Bytes(), &createdGroup)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "existing group",
			id:         strconv.FormatInt(createdGroup.ID, 10),
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing group",
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
			req := httptest.NewRequest(http.MethodGet, "/groups/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			err := handler.GetGroupByID(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				var retrievedGroup models.Group
				err := json.Unmarshal(rec.Body.Bytes(), &retrievedGroup)
				assert.NoError(t, err)
				assert.Equal(t, createdGroup.ID, retrievedGroup.ID)
				assert.Equal(t, createdGroup.Name, retrievedGroup.Name)
			}
		})
	}
}

func TestGroupHandler_UpdateGroup(t *testing.T) {
	e, handler, cleanup := setupGroupTest(t)
	defer cleanup()

	// First create a group to test with
	group := &models.Group{
		Name: "Test Group",
	}
	jsonBytes, _ := json.Marshal(group)
	req := httptest.NewRequest(http.MethodPost, "/groups", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	err := handler.CreateGroup(c)
	assert.NoError(t, err)

	var createdGroup models.Group
	err = json.Unmarshal(rec.Body.Bytes(), &createdGroup)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		groupID    string
		updateData models.Group
		wantStatus int
	}{
		{
			name:    "valid update",
			groupID: strconv.FormatInt(createdGroup.ID, 10),
			updateData: models.Group{
				ID:   createdGroup.ID,
				Name: "Updated Group",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "invalid update - empty name",
			groupID: strconv.FormatInt(createdGroup.ID, 10),
			updateData: models.Group{
				ID:   createdGroup.ID,
				Name: "",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "non-existing group",
			groupID: "999",
			updateData: models.Group{
				ID:   999,
				Name: "Updated Group",
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:    "non-existing group",
			groupID: "999",
			updateData: models.Group{
				ID:   0,
				Name: "Updated Group",
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, _ := json.Marshal(tt.updateData)
			req := httptest.NewRequest(http.MethodPut, "/groups/:id", bytes.NewReader(jsonBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.groupID)

			err := handler.UpdateGroup(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantStatus == http.StatusOK {
				var updatedGroup models.Group
				err := json.Unmarshal(rec.Body.Bytes(), &updatedGroup)
				assert.NoError(t, err)
				assert.Equal(t, tt.updateData.Name, updatedGroup.Name)
			}
		})
	}
}

func TestGroupHandler_DeleteGroup(t *testing.T) {
	e, handler, cleanup := setupGroupTest(t)
	defer cleanup()

	// First create a group to test with
	group := &models.Group{
		Name: "Test Group",
	}
	jsonBytes, _ := json.Marshal(group)
	req := httptest.NewRequest(http.MethodPost, "/groups", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	
	err := handler.CreateGroup(c)
	assert.NoError(t, err)

	var createdGroup models.Group
	err = json.Unmarshal(rec.Body.Bytes(), &createdGroup)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		groupID    string
		wantStatus int
	}{
		{
			name:       "existing group",
			groupID:    strconv.FormatInt(createdGroup.ID, 10),
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing group",
			groupID:    "999",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid id",
			groupID:    "abc",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/groups/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.groupID)

			err := handler.DeleteGroup(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestGroupHandler_ListGroups(t *testing.T) {
	e, handler, cleanup := setupGroupTest(t)
	defer cleanup()

	tests := []struct {
		name       string
		page       int
		pageSize   int
		search     string
		wantStatus int
		wantTotal  int
	}{
		{
			name:       "list all groups",
			page:       1,
			pageSize:   10,
			search:     "",
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name:       "list with pagination",
			page:       1,
			pageSize:   1,
			search:     "",
			wantStatus: http.StatusOK,
			wantTotal:  2,
		},
		{
			name:       "search groups",
			page:       1,
			pageSize:   10,
			search:     "Travel",
			wantStatus: http.StatusOK,
			wantTotal:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/groups", nil)
			q := req.URL.Query()
			q.Add("page", strconv.Itoa(tt.page))
			q.Add("pageSize", strconv.Itoa(tt.pageSize))
			q.Add("search", tt.search)
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.ListGroups(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatus, rec.Code)

			var response map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			_ = response["groups"].([]interface{})
			total, ok := response["total"].(float64)
			assert.True(t, ok, "total should be a float64")

			assert.Equal(t, tt.wantTotal, int(total))
		})
	}
}
