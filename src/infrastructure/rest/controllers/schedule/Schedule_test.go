package schedule

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domainSchedule "caregiver/src/domain/schedule"
	domainUser "caregiver/src/domain/user"
	logger "caregiver/src/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// mockScheduleUseCase is a mock implementation of the IScheduleUseCase interface
type mockScheduleUseCase struct {
	getSchedulesFn                                    func() (*[]domainSchedule.Schedule, error)
	getSchedulesWithClientInfoFn                      func() (*[]domainSchedule.Schedule, *[]domainUser.User, error)
	getScheduleByIDFn                                 func(id uuid.UUID) (*domainSchedule.Schedule, error)
	getScheduleWithClientInfoFn                       func(id uuid.UUID) (*domainSchedule.Schedule, *domainUser.User, error)
	getTodaySchedulesFn                               func(userID uuid.UUID) (*[]domainSchedule.Schedule, error)
	getTodaySchedulesWithClientInfoFn                 func(userID uuid.UUID) (*[]domainSchedule.Schedule, *[]domainUser.User, error)
	startScheduleFn                                   func(scheduleID uuid.UUID, timestamp time.Time, location domainSchedule.Location) (*domainSchedule.Schedule, error)
	endScheduleFn                                     func(scheduleID uuid.UUID, timestamp time.Time, location domainSchedule.Location, tasks []domainSchedule.Task) (*domainSchedule.Schedule, error)
	updateTaskStatusFn                                func(taskID uuid.UUID, status string, done bool, feedback string) (*domainSchedule.Task, error)
	updateScheduleFn                                  func(scheduleID uuid.UUID, updates map[string]interface{}) (*domainSchedule.Schedule, error)
	createScheduleFn                                  func(newSchedule *domainSchedule.Schedule) (*domainSchedule.Schedule, error)
	getTodaySchedulesByAssignedUserIDFn               func(assignedUserID uuid.UUID) (*[]domainSchedule.Schedule, error)
	getTodaySchedulesByAssignedUserIDWithClientInfoFn func(assignedUserID uuid.UUID) (*[]domainSchedule.Schedule, *[]domainUser.User, error)
}

// Implement all methods of the IScheduleUseCase interface
func (m *mockScheduleUseCase) GetSchedules() (*[]domainSchedule.Schedule, error) {
	return m.getSchedulesFn()
}

func (m *mockScheduleUseCase) GetSchedulesWithClientInfo() (*[]domainSchedule.Schedule, *[]domainUser.User, error) {
	return m.getSchedulesWithClientInfoFn()
}

func (m *mockScheduleUseCase) GetScheduleByID(id uuid.UUID) (*domainSchedule.Schedule, error) {
	return m.getScheduleByIDFn(id)
}

func (m *mockScheduleUseCase) GetScheduleWithClientInfo(id uuid.UUID) (*domainSchedule.Schedule, *domainUser.User, error) {
	return m.getScheduleWithClientInfoFn(id)
}

func (m *mockScheduleUseCase) GetTodaySchedules(userID uuid.UUID) (*[]domainSchedule.Schedule, error) {
	return m.getTodaySchedulesFn(userID)
}

func (m *mockScheduleUseCase) GetTodaySchedulesWithClientInfo(userID uuid.UUID) (*[]domainSchedule.Schedule, *[]domainUser.User, error) {
	return m.getTodaySchedulesWithClientInfoFn(userID)
}

func (m *mockScheduleUseCase) StartSchedule(scheduleID uuid.UUID, timestamp time.Time, location domainSchedule.Location) (*domainSchedule.Schedule, error) {
	return m.startScheduleFn(scheduleID, timestamp, location)
}

func (m *mockScheduleUseCase) EndSchedule(scheduleID uuid.UUID, timestamp time.Time, location domainSchedule.Location, tasks []domainSchedule.Task) (*domainSchedule.Schedule, error) {
	return m.endScheduleFn(scheduleID, timestamp, location, tasks)
}

func (m *mockScheduleUseCase) UpdateTaskStatus(taskID uuid.UUID, status string, done bool, feedback string) (*domainSchedule.Task, error) {
	return m.updateTaskStatusFn(taskID, status, done, feedback)
}

func (m *mockScheduleUseCase) UpdateSchedule(scheduleID uuid.UUID, updates map[string]interface{}) (*domainSchedule.Schedule, error) {
	return m.updateScheduleFn(scheduleID, updates)
}

func (m *mockScheduleUseCase) CreateSchedule(newSchedule *domainSchedule.Schedule) (*domainSchedule.Schedule, error) {
	return m.createScheduleFn(newSchedule)
}

func (m *mockScheduleUseCase) GetTodaySchedulesByAssignedUserID(assignedUserID uuid.UUID) (*[]domainSchedule.Schedule, error) {
	return m.getTodaySchedulesByAssignedUserIDFn(assignedUserID)
}

func (m *mockScheduleUseCase) GetTodaySchedulesByAssignedUserIDWithClientInfo(assignedUserID uuid.UUID) (*[]domainSchedule.Schedule, *[]domainUser.User, error) {
	return m.getTodaySchedulesByAssignedUserIDWithClientInfoFn(assignedUserID)
}

// setupLogger creates a logger instance for testing
func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
}

// setupTestController creates a new Schedule controller with mock usecase for testing
func setupTestController(t *testing.T) (*Controller, *mockScheduleUseCase, *gin.Engine) {
	gin.SetMode(gin.TestMode)

	mockUseCase := &mockScheduleUseCase{}
	loggerInstance := setupLogger(t)
	controller := &Controller{
		scheduleUseCase: mockUseCase,
		Logger:          loggerInstance,
	}

	router := gin.New()
	router.Use(gin.Recovery())
	// Add error handling middleware
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			c.JSON(500, gin.H{"error": err.Error()})
		}
	})

	return controller, mockUseCase, router
}

// createTestSchedule creates a test schedule with the given ID
func createTestSchedule(id uuid.UUID) *domainSchedule.Schedule {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	clientUserID := uuid.New()
	assignedUserID := uuid.New()

	task1ID := uuid.New()
	task2ID := uuid.New()

	done := false

	return &domainSchedule.Schedule{
		ID:             id,
		ClientUserID:   clientUserID,
		AssignedUserID: assignedUserID,
		ServiceName:    "Test Service",
		ScheduledSlot: domainSchedule.ScheduledSlot{
			From: tomorrow,
			To:   tomorrow.Add(2 * time.Hour),
		},
		VisitStatus: "upcoming",
		Tasks: []domainSchedule.Task{
			{
				ID:          task1ID,
				ScheduleID:  id,
				Title:       "Task 1",
				Description: "Description for task 1",
				Status:      "pending",
				Done:        &done,
			},
			{
				ID:          task2ID,
				ScheduleID:  id,
				Title:       "Task 2",
				Description: "Description for task 2",
				Status:      "pending",
				Done:        &done,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// createTestUser creates a test user with the given ID
func createTestUser(id uuid.UUID) *domainUser.User {
	return &domainUser.User{
		ID:        id,
		UserName:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Status:    true,
		Location: domainUser.Location{
			HouseNumber: "123",
			Street:      "Test Street",
			City:        "Test City",
			State:       "Test State",
			Pincode:     "12345",
			Lat:         12.345,
			Long:        67.890,
		},
	}
}

// TestGetSchedules tests the GetSchedules controller method
func TestGetSchedules(t *testing.T) {
	// Setup
	controller, mockUseCase, router := setupTestController(t)

	// Setup route
	router.GET("/schedules", controller.GetSchedules)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		scheduleID1 := uuid.New()
		scheduleID2 := uuid.New()

		schedule1 := createTestSchedule(scheduleID1)
		schedule2 := createTestSchedule(scheduleID2)

		schedules := []domainSchedule.Schedule{*schedule1, *schedule2}
		clients := []domainUser.User{*createTestUser(schedule1.ClientUserID), *createTestUser(schedule2.ClientUserID)}

		mockUseCase.getSchedulesWithClientInfoFn = func() (*[]domainSchedule.Schedule, *[]domainUser.User, error) {
			return &schedules, &clients, nil
		}

		// Execute request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/schedules", nil)
		router.ServeHTTP(w, req)

		// Verify
		assert.Equal(t, http.StatusOK, w.Code)

		var response []ScheduleResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, scheduleID1, response[0].ID)
		assert.Equal(t, scheduleID2, response[1].ID)
	})

	t.Run("Error", func(t *testing.T) {
		// Setup mock behavior
		mockUseCase.getSchedulesWithClientInfoFn = func() (*[]domainSchedule.Schedule, *[]domainUser.User, error) {
			return nil, nil, errors.New("database error")
		}

		// Execute request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/schedules", nil)
		router.ServeHTTP(w, req)

		// Verify
		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

// TestCreateSchedule tests the CreateSchedule controller method
func TestCreateSchedule(t *testing.T) {
	// Setup
	controller, mockUseCase, router := setupTestController(t)

	// Setup route
	router.POST("/schedules", controller.CreateSchedule)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()
		clientUserID := uuid.New()
		assignedUserID := uuid.New()

		now := time.Now()
		tomorrow := now.Add(24 * time.Hour)

		// Create request body
		requestBody := CreateScheduleRequest{
			ClientUserID:   clientUserID,
			AssignedUserID: assignedUserID,
			ServiceName:    "Test Service",
			ScheduledSlot: ScheduledSlot{
				From: tomorrow,
				To:   tomorrow.Add(2 * time.Hour),
			},
			Tasks: []TaskRequest{
				{
					Title:       "Task 1",
					Description: "Description for task 1",
				},
				{
					Title:       "Task 2",
					Description: "Description for task 2",
				},
			},
		}

		// Create expected schedule
		createdSchedule := createTestSchedule(scheduleID)
		createdSchedule.ClientUserID = clientUserID
		createdSchedule.AssignedUserID = assignedUserID

		mockUseCase.createScheduleFn = func(newSchedule *domainSchedule.Schedule) (*domainSchedule.Schedule, error) {
			// Verify schedule properties
			assert.Equal(t, clientUserID, newSchedule.ClientUserID)
			assert.Equal(t, assignedUserID, newSchedule.AssignedUserID)
			assert.Equal(t, "Test Service", newSchedule.ServiceName)
			assert.Equal(t, "upcoming", newSchedule.VisitStatus)
			assert.Len(t, newSchedule.Tasks, 2)

			return createdSchedule, nil
		}

		// Execute request
		w := httptest.NewRecorder()
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/schedules", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Verify
		assert.Equal(t, http.StatusOK, w.Code)

		var response ScheduleResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, scheduleID, response.ID)
		assert.Equal(t, clientUserID, response.ClientUserID)
		assert.Equal(t, assignedUserID, response.AssignedUserID)
		assert.Equal(t, "Test Service", response.ServiceName)
		assert.Equal(t, "upcoming", response.VisitStatus)
		assert.Len(t, response.Tasks, 2)
	})

	t.Run("Missing ClientUserID", func(t *testing.T) {
		// Create request body with missing ClientUserID
		requestBody := CreateScheduleRequest{
			AssignedUserID: uuid.New(),
			ServiceName:    "Test Service",
			ScheduledSlot: ScheduledSlot{
				From: time.Now().Add(24 * time.Hour),
				To:   time.Now().Add(26 * time.Hour),
			},
			Tasks: []TaskRequest{
				{
					Title:       "Task 1",
					Description: "Description for task 1",
				},
			},
		}

		// Execute request
		w := httptest.NewRecorder()
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/schedules", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Verify
		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("Missing ScheduledSlot", func(t *testing.T) {
		// Create request body with missing ScheduledSlot
		requestBody := CreateScheduleRequest{
			ClientUserID:   uuid.New(),
			AssignedUserID: uuid.New(),
			ServiceName:    "Test Service",
			Tasks: []TaskRequest{
				{
					Title:       "Task 1",
					Description: "Description for task 1",
				},
			},
		}

		// Execute request
		w := httptest.NewRecorder()
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/schedules", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Verify
		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid ScheduledSlot", func(t *testing.T) {
		// Create request body with invalid ScheduledSlot (From after To)
		now := time.Now()
		requestBody := CreateScheduleRequest{
			ClientUserID:   uuid.New(),
			AssignedUserID: uuid.New(),
			ServiceName:    "Test Service",
			ScheduledSlot: ScheduledSlot{
				From: now.Add(26 * time.Hour), // From is after To
				To:   now.Add(24 * time.Hour),
			},
			Tasks: []TaskRequest{
				{
					Title:       "Task 1",
					Description: "Description for task 1",
				},
			},
		}

		// Execute request
		w := httptest.NewRecorder()
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/schedules", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Verify
		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("No Tasks", func(t *testing.T) {
		// Create request body with no tasks
		requestBody := CreateScheduleRequest{
			ClientUserID:   uuid.New(),
			AssignedUserID: uuid.New(),
			ServiceName:    "Test Service",
			ScheduledSlot: ScheduledSlot{
				From: time.Now().Add(24 * time.Hour),
				To:   time.Now().Add(26 * time.Hour),
			},
			Tasks: []TaskRequest{}, // Empty tasks
		}

		// Execute request
		w := httptest.NewRecorder()
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/schedules", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Verify
		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("UseCase Error", func(t *testing.T) {
		// Setup mock behavior
		clientUserID := uuid.New()
		assignedUserID := uuid.New()

		now := time.Now()
		tomorrow := now.Add(24 * time.Hour)

		// Create request body
		requestBody := CreateScheduleRequest{
			ClientUserID:   clientUserID,
			AssignedUserID: assignedUserID,
			ServiceName:    "Test Service",
			ScheduledSlot: ScheduledSlot{
				From: tomorrow,
				To:   tomorrow.Add(2 * time.Hour),
			},
			Tasks: []TaskRequest{
				{
					Title:       "Task 1",
					Description: "Description for task 1",
				},
			},
		}

		mockUseCase.createScheduleFn = func(newSchedule *domainSchedule.Schedule) (*domainSchedule.Schedule, error) {
			return nil, errors.New("database error")
		}

		// Execute request
		w := httptest.NewRecorder()
		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/schedules", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Verify
		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}
