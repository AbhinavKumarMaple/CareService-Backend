package controllers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestBindJSON(t *testing.T) {
	validJSON := `{"name": "test", "email": "test@example.com"}`

	c, _ := setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(validJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	var request map[string]string
	err := BindJSON(c, &request)

	assert.NoError(t, err)
	assert.Equal(t, "test", request["name"])
	assert.Equal(t, "test@example.com", request["email"])

	invalidJSON := `{"name": "test", "email": "test@example.com"`

	c, _ = setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	err = BindJSON(c, &request)

	assert.Error(t, err)

	c, _ = setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(""))
	c.Request.Header.Set("Content-Type", "application/json")

	err = BindJSON(c, &request)

	assert.Error(t, err)
}

func TestBindJSONMap(t *testing.T) {
	validJSON := `{"name": "test", "email": "test@example.com", "age": 25}`

	c, _ := setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(validJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	var request map[string]any
	err := BindJSONMap(c, &request)

	assert.NoError(t, err)
	assert.Equal(t, "test", request["name"])
	assert.Equal(t, "test@example.com", request["email"])
	assert.Equal(t, float64(25), request["age"]) 

	invalidJSON := `{"name": "test", "email": "test@example.com"`

	c, _ = setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	err = BindJSONMap(c, &request)

	assert.Error(t, err)

	c, _ = setupGinContext()
	c.Request = httptest.NewRequest("POST", "/test", bytes.NewBufferString(""))
	c.Request.Header.Set("Content-Type", "application/json")

	err = BindJSONMap(c, &request)

	assert.Error(t, err)
}

func TestPaginationValues(t *testing.T) {
	numPages, nextCursor, prevCursor := PaginationValues(10, 2, 25)
	assert.Equal(t, int64(3), numPages)   
	assert.Equal(t, int64(3), nextCursor) 
	assert.Equal(t, int64(1), prevCursor) 

	numPages, nextCursor, prevCursor = PaginationValues(10, 1, 25)
	assert.Equal(t, int64(3), numPages)   
	assert.Equal(t, int64(2), nextCursor) 
	assert.Equal(t, int64(0), prevCursor) 

	// Test case 3: Last page
	numPages, nextCursor, prevCursor = PaginationValues(10, 3, 25)
	assert.Equal(t, int64(3), numPages)   
	assert.Equal(t, int64(0), nextCursor) 
	assert.Equal(t, int64(2), prevCursor) 

	// Test case 4: Single page
	numPages, nextCursor, prevCursor = PaginationValues(10, 1, 5)
	assert.Equal(t, int64(1), numPages)  
	assert.Equal(t, int64(0), nextCursor) 
	assert.Equal(t, int64(0), prevCursor)

	// Test case 5: Empty result
	numPages, nextCursor, prevCursor = PaginationValues(10, 1, 0)
	assert.Equal(t, int64(0), numPages)   
	assert.Equal(t, int64(0), nextCursor)
	assert.Equal(t, int64(0), prevCursor) 

	// Test case 6: Large numbers
	numPages, nextCursor, prevCursor = PaginationValues(100, 5, 1000)
	assert.Equal(t, int64(10), numPages)  
	assert.Equal(t, int64(6), nextCursor) 
	assert.Equal(t, int64(4), prevCursor) 
}

func TestMessageResponse(t *testing.T) {
	message := MessageResponse{
		Message: "Test message",
	}

	jsonData, err := json.Marshal(message)
	assert.NoError(t, err)

	var unmarshaled MessageResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, message.Message, unmarshaled.Message)
}

func TestSortByDataRequest(t *testing.T) {
	sortRequest := SortByDataRequest{
		Field:     "name",
		Direction: "asc",
	}

	jsonData, err := json.Marshal(sortRequest)
	assert.NoError(t, err)

	var unmarshaled SortByDataRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, sortRequest.Field, unmarshaled.Field)
	assert.Equal(t, sortRequest.Direction, unmarshaled.Direction)
}

func TestFieldDateRangeDataRequest(t *testing.T) {
	dateRangeRequest := FieldDateRangeDataRequest{
		Field:     "created_at",
		StartDate: "2023-01-01",
		EndDate:   "2023-12-31",
	}

	jsonData, err := json.Marshal(dateRangeRequest)
	assert.NoError(t, err)

	var unmarshaled FieldDateRangeDataRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, dateRangeRequest.Field, unmarshaled.Field)
	assert.Equal(t, dateRangeRequest.StartDate, unmarshaled.StartDate)
	assert.Equal(t, dateRangeRequest.EndDate, unmarshaled.EndDate)
}
