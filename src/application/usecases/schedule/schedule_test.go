package schedule

import (
	"errors"
	"testing"
	"time"

	"caregiver/src/domain"
	domainSchedule "caregiver/src/domain/schedule"
	domainUser "caregiver/src/domain/user"
	logger "caregiver/src/infrastructure/logger"

	"github.com/google/uuid"
)

// mockScheduleRepository is a mock implementation of the IScheduleRepository interface
type mockScheduleRepository struct {
	getSchedulesFn                          func() (*[]domainSchedule.Schedule, error)
	getScheduleByIDFn                       func(id uuid.UUID) (*domainSchedule.Schedule, error)
	getTodaySchedulesFn                     func(userID uuid.UUID) (*[]domainSchedule.Schedule, error)
	updateScheduleFn                        func(id uuid.UUID, updates map[string]interface{}) (*domainSchedule.Schedule, error)
	updateTaskFn                            func(taskID uuid.UUID, updates map[string]interface{}) (*domainSchedule.Task, error)
	createFn                                func(newSchedule *domainSchedule.Schedule) (*domainSchedule.Schedule, error)
	getSchedulesByAssignedUserIDPaginatedFn func(assignedUserID uuid.UUID, filters domain.DataFilters) (*domainSchedule.SearchResultSchedule, error)
}

// Implement all methods of the IScheduleRepository interface
func (m *mockScheduleRepository) GetSchedules() (*[]domainSchedule.Schedule, error) {
	return m.getSchedulesFn()
}

func (m *mockScheduleRepository) GetScheduleByID(id uuid.UUID) (*domainSchedule.Schedule, error) {
	return m.getScheduleByIDFn(id)
}

func (m *mockScheduleRepository) GetTodaySchedules(userID uuid.UUID) (*[]domainSchedule.Schedule, error) {
	return m.getTodaySchedulesFn(userID)
}

func (m *mockScheduleRepository) UpdateSchedule(id uuid.UUID, updates map[string]interface{}) (*domainSchedule.Schedule, error) {
	return m.updateScheduleFn(id, updates)
}

func (m *mockScheduleRepository) UpdateTask(taskID uuid.UUID, updates map[string]interface{}) (*domainSchedule.Task, error) {
	return m.updateTaskFn(taskID, updates)
}

func (m *mockScheduleRepository) Create(newSchedule *domainSchedule.Schedule) (*domainSchedule.Schedule, error) {
	return m.createFn(newSchedule)
}

func (m *mockScheduleRepository) GetSchedulesByAssignedUserIDPaginated(assignedUserID uuid.UUID, filters domain.DataFilters) (*domainSchedule.SearchResultSchedule, error) {
	return m.getSchedulesByAssignedUserIDPaginatedFn(assignedUserID, filters)
}

// mockUserRepository is a mock implementation of the IUserRepository interface
type mockUserRepository struct {
	getAllFn           func() (*[]domainUser.User, error)
	createFn           func(userDomain *domainUser.User) (*domainUser.User, error)
	getByIDFn          func(id uuid.UUID) (*domainUser.User, error)
	getByEmailFn       func(email string) (*domainUser.User, error)
	updateFn           func(id uuid.UUID, userMap map[string]interface{}) (*domainUser.User, error)
	deleteFn           func(id uuid.UUID) error
	searchPaginatedFn  func(filters domain.DataFilters) (*domainUser.SearchResultUser, error)
	searchByPropertyFn func(property string, searchText string) (*[]string, error)
}

// Implement all methods of the IUserRepository interface
func (m *mockUserRepository) GetAll() (*[]domainUser.User, error) {
	return m.getAllFn()
}

func (m *mockUserRepository) Create(userDomain *domainUser.User) (*domainUser.User, error) {
	return m.createFn(userDomain)
}

func (m *mockUserRepository) GetByID(id uuid.UUID) (*domainUser.User, error) {
	return m.getByIDFn(id)
}

func (m *mockUserRepository) GetByEmail(email string) (*domainUser.User, error) {
	return m.getByEmailFn(email)
}

func (m *mockUserRepository) Update(id uuid.UUID, userMap map[string]interface{}) (*domainUser.User, error) {
	return m.updateFn(id, userMap)
}

func (m *mockUserRepository) Delete(id uuid.UUID) error {
	return m.deleteFn(id)
}

func (m *mockUserRepository) SearchPaginated(filters domain.DataFilters) (*domainUser.SearchResultUser, error) {
	return m.searchPaginatedFn(filters)
}

func (m *mockUserRepository) SearchByProperty(property string, searchText string) (*[]string, error) {
	return m.searchByPropertyFn(property, searchText)
}

// setupLogger creates a logger instance for testing
func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
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

// createTestScheduleList creates a list of test schedules
func createTestScheduleList(count int) *[]domainSchedule.Schedule {
	schedules := make([]domainSchedule.Schedule, count)
	for i := 0; i < count; i++ {
		schedules[i] = *createTestSchedule(uuid.New())
	}
	return &schedules
}

// createTestUserList creates a list of test users
func createTestUserList(count int) *[]domainUser.User {
	users := make([]domainUser.User, count)
	for i := 0; i < count; i++ {
		users[i] = *createTestUser(uuid.New())
	}
	return &users
}

// TestNewScheduleUseCase tests the creation of a new Schedule usecase
func TestNewScheduleUseCase(t *testing.T) {
	// Setup
	mockScheduleRepo := &mockScheduleRepository{}
	mockUserRepo := &mockUserRepository{}
	loggerInstance := setupLogger(t)

	// Execute
	useCase := NewScheduleUseCase(mockScheduleRepo, mockUserRepo, loggerInstance)

	// Verify
	if useCase == nil {
		t.Error("expected non-nil usecase")
	}
}

// setupTestScheduleUseCase creates a new Schedule usecase with mock repositories for testing
func setupTestScheduleUseCase(t *testing.T) (IScheduleUseCase, *mockScheduleRepository, *mockUserRepository, *logger.Logger) {
	mockScheduleRepo := &mockScheduleRepository{}
	mockUserRepo := &mockUserRepository{}
	loggerInstance := setupLogger(t)
	useCase := NewScheduleUseCase(mockScheduleRepo, mockUserRepo, loggerInstance)
	return useCase, mockScheduleRepo, mockUserRepo, loggerInstance
}

// TestGetSchedules tests the GetSchedules method
func TestGetSchedules(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, _, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		expectedSchedules := createTestScheduleList(3)
		mockScheduleRepo.getSchedulesFn = func() (*[]domainSchedule.Schedule, error) {
			return expectedSchedules, nil
		}

		// Execute
		schedules, err := useCase.GetSchedules()

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if schedules == nil {
			t.Error("expected non-nil schedules")
		}
		if len(*schedules) != 3 {
			t.Errorf("expected 3 schedules, got %d", len(*schedules))
		}
	})

	t.Run("Error", func(t *testing.T) {
		// Setup mock behavior
		mockScheduleRepo.getSchedulesFn = func() (*[]domainSchedule.Schedule, error) {
			return nil, errors.New("database error")
		}

		// Execute
		schedules, err := useCase.GetSchedules()

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if schedules != nil {
			t.Error("expected nil schedules")
		}
	})
}

// TestGetSchedulesWithClientInfo tests the GetSchedulesWithClientInfo method
func TestGetSchedulesWithClientInfo(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, mockUserRepo, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		expectedSchedules := createTestScheduleList(2)
		mockScheduleRepo.getSchedulesFn = func() (*[]domainSchedule.Schedule, error) {
			return expectedSchedules, nil
		}

		// Create a map of client IDs from the schedules
		clientIDs := make(map[uuid.UUID]bool)
		for _, schedule := range *expectedSchedules {
			clientIDs[schedule.ClientUserID] = true
		}

		// Create test users for each client ID
		expectedClients := make([]domainUser.User, 0, len(clientIDs))
		for clientID := range clientIDs {
			user := createTestUser(clientID)
			expectedClients = append(expectedClients, *user)
		}

		// Setup mock behavior for GetByID
		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			for _, client := range expectedClients {
				if client.ID == id {
					return &client, nil
				}
			}
			return nil, errors.New("user not found")
		}

		// Execute
		schedules, clients, err := useCase.GetSchedulesWithClientInfo()

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if schedules == nil {
			t.Error("expected non-nil schedules")
		}
		if len(*schedules) != 2 {
			t.Errorf("expected 2 schedules, got %d", len(*schedules))
		}
		if clients == nil {
			t.Error("expected non-nil clients")
		}
		if len(*clients) != len(clientIDs) {
			t.Errorf("expected %d clients, got %d", len(clientIDs), len(*clients))
		}
	})

	t.Run("Error getting schedules", func(t *testing.T) {
		// Setup mock behavior
		mockScheduleRepo.getSchedulesFn = func() (*[]domainSchedule.Schedule, error) {
			return nil, errors.New("database error")
		}

		// Execute
		schedules, clients, err := useCase.GetSchedulesWithClientInfo()

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if schedules != nil {
			t.Error("expected nil schedules")
		}
		if clients != nil {
			t.Error("expected nil clients")
		}
	})

	t.Run("Empty schedules", func(t *testing.T) {
		// Setup mock behavior
		emptySchedules := &[]domainSchedule.Schedule{}
		mockScheduleRepo.getSchedulesFn = func() (*[]domainSchedule.Schedule, error) {
			return emptySchedules, nil
		}

		// Execute
		schedules, clients, err := useCase.GetSchedulesWithClientInfo()

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if schedules == nil {
			t.Error("expected non-nil schedules")
		}
		if len(*schedules) != 0 {
			t.Errorf("expected 0 schedules, got %d", len(*schedules))
		}
		if clients == nil {
			t.Error("expected non-nil clients")
		}
		if len(*clients) != 0 {
			t.Errorf("expected 0 clients, got %d", len(*clients))
		}
	})
}

// TestGetScheduleByID tests the GetScheduleByID method
func TestGetScheduleByID(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, _, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()
		expectedSchedule := createTestSchedule(scheduleID)
		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return expectedSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		// Execute
		schedule, err := useCase.GetScheduleByID(scheduleID)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if schedule == nil {
			t.Error("expected non-nil schedule")
		}
		if schedule.ID != scheduleID {
			t.Errorf("expected schedule ID %s, got %s", scheduleID, schedule.ID)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		// Setup mock behavior
		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			return nil, errors.New("schedule not found")
		}

		// Execute
		schedule, err := useCase.GetScheduleByID(uuid.New())

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if schedule != nil {
			t.Error("expected nil schedule")
		}
	})
}

// TestGetScheduleWithClientInfo tests the GetScheduleWithClientInfo method
func TestGetScheduleWithClientInfo(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, mockUserRepo, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()
		clientUserID := uuid.New()

		// Create test schedule with the client user ID
		expectedSchedule := createTestSchedule(scheduleID)
		expectedSchedule.ClientUserID = clientUserID

		// Create test client user
		expectedClient := createTestUser(clientUserID)

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return expectedSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == clientUserID {
				return expectedClient, nil
			}
			return nil, errors.New("user not found")
		}

		// Execute
		schedule, client, err := useCase.GetScheduleWithClientInfo(scheduleID)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if schedule == nil {
			t.Error("expected non-nil schedule")
		}
		if schedule.ID != scheduleID {
			t.Errorf("expected schedule ID %s, got %s", scheduleID, schedule.ID)
		}
		if client == nil {
			t.Error("expected non-nil client")
		}
		if client.ID != clientUserID {
			t.Errorf("expected client ID %s, got %s", clientUserID, client.ID)
		}
	})

	t.Run("Schedule not found", func(t *testing.T) {
		// Setup mock behavior
		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			return nil, errors.New("schedule not found")
		}

		// Execute
		schedule, client, err := useCase.GetScheduleWithClientInfo(uuid.New())

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if schedule != nil {
			t.Error("expected nil schedule")
		}
		if client != nil {
			t.Error("expected nil client")
		}
	})

	t.Run("Client not found", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()
		clientUserID := uuid.New()

		// Create test schedule with the client user ID
		expectedSchedule := createTestSchedule(scheduleID)
		expectedSchedule.ClientUserID = clientUserID

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return expectedSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			return nil, errors.New("user not found")
		}

		// Execute
		schedule, client, err := useCase.GetScheduleWithClientInfo(scheduleID)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if schedule == nil {
			t.Error("expected non-nil schedule")
		}
		if schedule.ID != scheduleID {
			t.Errorf("expected schedule ID %s, got %s", scheduleID, schedule.ID)
		}
		if client != nil {
			t.Error("expected nil client when user not found")
		}
	})
}

// TestStartSchedule tests the StartSchedule method
func TestStartSchedule(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, _, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()
		timestamp := time.Now()
		lat := 12.345
		long := 67.890
		location := domainSchedule.Location{
			Lat:  &lat,
			Long: &long,
		}

		// Create test schedule
		originalSchedule := createTestSchedule(scheduleID)
		originalSchedule.VisitStatus = "upcoming"

		// Create updated schedule
		updatedSchedule := *originalSchedule
		updatedSchedule.VisitStatus = "in_progress"
		updatedSchedule.CheckinTime = &timestamp
		updatedSchedule.CheckinLocation = location

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return originalSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		mockScheduleRepo.updateScheduleFn = func(id uuid.UUID, updates map[string]interface{}) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				// Verify updates
				if updates["visit_status"] != "in_progress" {
					t.Errorf("expected visit_status to be 'in_progress', got %v", updates["visit_status"])
				}
				if updates["checkin_time"] != timestamp {
					t.Errorf("expected checkin_time to be %v, got %v", timestamp, updates["checkin_time"])
				}
				if updates["checkin_location_lat"] != location.Lat {
					t.Errorf("expected checkin_location_lat to be %v, got %v", location.Lat, updates["checkin_location_lat"])
				}
				if updates["checkin_location_long"] != location.Long {
					t.Errorf("expected checkin_location_long to be %v, got %v", location.Long, updates["checkin_location_long"])
				}

				return &updatedSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		// Execute
		result, err := useCase.StartSchedule(scheduleID, timestamp, location)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("expected non-nil result")
		}
		if result.VisitStatus != "in_progress" {
			t.Errorf("expected visit_status to be 'in_progress', got %s", result.VisitStatus)
		}
		if result.CheckinTime == nil {
			t.Error("expected non-nil checkin_time")
		} else if !result.CheckinTime.Equal(timestamp) {
			t.Errorf("expected checkin_time to be %v, got %v", timestamp, result.CheckinTime)
		}
		if result.CheckinLocation.Lat == nil {
			t.Error("expected non-nil checkin_location.Lat")
		} else if *result.CheckinLocation.Lat != lat {
			t.Errorf("expected checkin_location.Lat to be %v, got %v", lat, *result.CheckinLocation.Lat)
		}
		if result.CheckinLocation.Long == nil {
			t.Error("expected non-nil checkin_location.Long")
		} else if *result.CheckinLocation.Long != long {
			t.Errorf("expected checkin_location.Long to be %v, got %v", long, *result.CheckinLocation.Long)
		}
	})

	t.Run("Schedule not found", func(t *testing.T) {
		// Setup mock behavior
		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			return nil, errors.New("schedule not found")
		}

		// Execute
		timestamp := time.Now()
		lat := 12.345
		long := 67.890
		location := domainSchedule.Location{
			Lat:  &lat,
			Long: &long,
		}
		result, err := useCase.StartSchedule(uuid.New(), timestamp, location)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("Invalid status", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()

		// Create test schedule with invalid status
		originalSchedule := createTestSchedule(scheduleID)
		originalSchedule.VisitStatus = "in_progress" // Already in progress

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return originalSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		// Execute
		timestamp := time.Now()
		lat := 12.345
		long := 67.890
		location := domainSchedule.Location{
			Lat:  &lat,
			Long: &long,
		}
		result, err := useCase.StartSchedule(scheduleID, timestamp, location)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("Update error", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()

		// Create test schedule
		originalSchedule := createTestSchedule(scheduleID)
		originalSchedule.VisitStatus = "upcoming"

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return originalSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		mockScheduleRepo.updateScheduleFn = func(id uuid.UUID, updates map[string]interface{}) (*domainSchedule.Schedule, error) {
			return nil, errors.New("database error")
		}

		// Execute
		timestamp := time.Now()
		lat := 12.345
		long := 67.890
		location := domainSchedule.Location{
			Lat:  &lat,
			Long: &long,
		}
		result, err := useCase.StartSchedule(scheduleID, timestamp, location)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})
}

// TestEndSchedule tests the EndSchedule method
func TestEndSchedule(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, _, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()
		timestamp := time.Now()
		lat := 12.345
		long := 67.890
		location := domainSchedule.Location{
			Lat:  &lat,
			Long: &long,
		}

		// Create test schedule
		originalSchedule := createTestSchedule(scheduleID)
		originalSchedule.VisitStatus = "in_progress"

		// Create tasks for update
		done := true
		feedback := "Task completed successfully"
		tasks := []domainSchedule.Task{
			{
				ID:          uuid.New(),
				ScheduleID:  scheduleID,
				Title:       "Task 1",
				Description: "Description for task 1",
				Status:      "completed",
				Done:        &done,
				Feedback:    &feedback,
			},
		}

		// Create updated schedule
		updatedSchedule := *originalSchedule
		updatedSchedule.VisitStatus = "completed"
		updatedSchedule.CheckoutTime = &timestamp
		updatedSchedule.CheckoutLocation = location

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return originalSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		mockScheduleRepo.updateScheduleFn = func(id uuid.UUID, updates map[string]interface{}) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				// Verify updates
				if updates["visit_status"] != "completed" {
					t.Errorf("expected visit_status to be 'completed', got %v", updates["visit_status"])
				}
				if updates["checkout_time"] != timestamp {
					t.Errorf("expected checkout_time to be %v, got %v", timestamp, updates["checkout_time"])
				}
				if updates["checkout_location_lat"] != location.Lat {
					t.Errorf("expected checkout_location_lat to be %v, got %v", location.Lat, updates["checkout_location_lat"])
				}
				if updates["checkout_location_long"] != location.Long {
					t.Errorf("expected checkout_location_long to be %v, got %v", location.Long, updates["checkout_location_long"])
				}

				return &updatedSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		mockScheduleRepo.updateTaskFn = func(taskID uuid.UUID, updates map[string]interface{}) (*domainSchedule.Task, error) {
			// Find the task in the tasks list
			for _, task := range tasks {
				if task.ID == taskID {
					// Verify updates
					if updates["status"] != task.Status {
						t.Errorf("expected status to be %s, got %v", task.Status, updates["status"])
					}
					if updates["done"] != task.Done {
						t.Errorf("expected done to be %v, got %v", task.Done, updates["done"])
					}
					if updates["feedback"] != task.Feedback {
						t.Errorf("expected feedback to be %v, got %v", task.Feedback, updates["feedback"])
					}

					return &task, nil
				}
			}
			return nil, errors.New("task not found")
		}

		// Execute
		result, err := useCase.EndSchedule(scheduleID, timestamp, location, tasks)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("expected non-nil result")
		}
		if result.VisitStatus != "completed" {
			t.Errorf("expected visit_status to be 'completed', got %s", result.VisitStatus)
		}
		if result.CheckoutTime == nil {
			t.Error("expected non-nil checkout_time")
		} else if !result.CheckoutTime.Equal(timestamp) {
			t.Errorf("expected checkout_time to be %v, got %v", timestamp, result.CheckoutTime)
		}
		if result.CheckoutLocation.Lat == nil {
			t.Error("expected non-nil checkout_location.Lat")
		} else if *result.CheckoutLocation.Lat != lat {
			t.Errorf("expected checkout_location.Lat to be %v, got %v", lat, *result.CheckoutLocation.Lat)
		}
		if result.CheckoutLocation.Long == nil {
			t.Error("expected non-nil checkout_location.Long")
		} else if *result.CheckoutLocation.Long != long {
			t.Errorf("expected checkout_location.Long to be %v, got %v", long, *result.CheckoutLocation.Long)
		}
	})

	t.Run("Schedule not found", func(t *testing.T) {
		// Setup mock behavior
		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			return nil, errors.New("schedule not found")
		}

		// Execute
		timestamp := time.Now()
		lat := 12.345
		long := 67.890
		location := domainSchedule.Location{
			Lat:  &lat,
			Long: &long,
		}
		tasks := []domainSchedule.Task{}
		result, err := useCase.EndSchedule(uuid.New(), timestamp, location, tasks)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("Invalid status", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()

		// Create test schedule with invalid status
		originalSchedule := createTestSchedule(scheduleID)
		originalSchedule.VisitStatus = "upcoming" // Not in progress

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return originalSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		// Execute
		timestamp := time.Now()
		lat := 12.345
		long := 67.890
		location := domainSchedule.Location{
			Lat:  &lat,
			Long: &long,
		}
		tasks := []domainSchedule.Task{}
		result, err := useCase.EndSchedule(scheduleID, timestamp, location, tasks)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("Update error", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()

		// Create test schedule
		originalSchedule := createTestSchedule(scheduleID)
		originalSchedule.VisitStatus = "in_progress"

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return originalSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		mockScheduleRepo.updateScheduleFn = func(id uuid.UUID, updates map[string]interface{}) (*domainSchedule.Schedule, error) {
			return nil, errors.New("database error")
		}

		// Execute
		timestamp := time.Now()
		lat := 12.345
		long := 67.890
		location := domainSchedule.Location{
			Lat:  &lat,
			Long: &long,
		}
		tasks := []domainSchedule.Task{}
		result, err := useCase.EndSchedule(scheduleID, timestamp, location, tasks)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})
}

// TestUpdateTaskStatus tests the UpdateTaskStatus method
func TestUpdateTaskStatus(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, _, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		taskID := uuid.New()
		status := "completed"
		done := true
		feedback := "Task completed successfully"

		// Create updated task
		updatedTask := &domainSchedule.Task{
			ID:       taskID,
			Status:   status,
			Done:     &done,
			Feedback: &feedback,
		}

		mockScheduleRepo.updateTaskFn = func(id uuid.UUID, updates map[string]interface{}) (*domainSchedule.Task, error) {
			if id == taskID {
				// Verify updates
				if updates["Status"] != status {
					t.Errorf("expected Status to be %s, got %v", status, updates["Status"])
				}
				if updates["Done"] != done {
					t.Errorf("expected Done to be %v, got %v", done, updates["Done"])
				}
				if updates["Feedback"] != feedback {
					t.Errorf("expected Feedback to be %s, got %v", feedback, updates["Feedback"])
				}

				return updatedTask, nil
			}
			return nil, errors.New("task not found")
		}

		// Execute
		result, err := useCase.UpdateTaskStatus(taskID, status, done, feedback)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("expected non-nil result")
		}
		if result.ID != taskID {
			t.Errorf("expected task ID %s, got %s", taskID, result.ID)
		}
		if result.Status != status {
			t.Errorf("expected status %s, got %s", status, result.Status)
		}
		if result.Done == nil {
			t.Error("expected non-nil Done")
		} else if *result.Done != done {
			t.Errorf("expected Done to be %v, got %v", done, *result.Done)
		}
		if result.Feedback == nil {
			t.Error("expected non-nil Feedback")
		} else if *result.Feedback != feedback {
			t.Errorf("expected Feedback to be %s, got %s", feedback, *result.Feedback)
		}
	})

	t.Run("Update error", func(t *testing.T) {
		// Setup mock behavior
		mockScheduleRepo.updateTaskFn = func(id uuid.UUID, updates map[string]interface{}) (*domainSchedule.Task, error) {
			return nil, errors.New("database error")
		}

		// Execute
		result, err := useCase.UpdateTaskStatus(uuid.New(), "completed", true, "feedback")

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})
}

// TestCreateSchedule tests the CreateSchedule method
func TestCreateSchedule(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, mockUserRepo, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()
		clientUserID := uuid.New()
		assignedUserID := uuid.New()

		// Create test schedule
		newSchedule := createTestSchedule(uuid.Nil) // ID will be assigned by the repository
		newSchedule.ClientUserID = clientUserID
		newSchedule.AssignedUserID = assignedUserID

		// Create client and assigned users
		clientUser := createTestUser(clientUserID)
		assignedUser := createTestUser(assignedUserID)

		// Create created schedule (with ID assigned)
		createdSchedule := *newSchedule
		createdSchedule.ID = scheduleID

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == clientUserID {
				return clientUser, nil
			}
			if id == assignedUserID {
				return assignedUser, nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.createFn = func(schedule *domainSchedule.Schedule) (*domainSchedule.Schedule, error) {
			// Verify schedule properties
			if schedule.ClientUserID != clientUserID {
				t.Errorf("expected ClientUserID to be %s, got %s", clientUserID, schedule.ClientUserID)
			}
			if schedule.AssignedUserID != assignedUserID {
				t.Errorf("expected AssignedUserID to be %s, got %s", assignedUserID, schedule.AssignedUserID)
			}
			if schedule.VisitStatus != "upcoming" {
				t.Errorf("expected VisitStatus to be 'upcoming', got %s", schedule.VisitStatus)
			}

			// Verify tasks
			for _, task := range schedule.Tasks {
				if task.ID == uuid.Nil {
					t.Error("expected task ID to be assigned")
				}
				if task.Status != "pending" {
					t.Errorf("expected task Status to be 'pending', got %s", task.Status)
				}
			}

			return &createdSchedule, nil
		}

		// Execute
		result, err := useCase.CreateSchedule(newSchedule)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("expected non-nil result")
		}
		if result.ID != scheduleID {
			t.Errorf("expected schedule ID %s, got %s", scheduleID, result.ID)
		}
		if result.ClientUserID != clientUserID {
			t.Errorf("expected ClientUserID %s, got %s", clientUserID, result.ClientUserID)
		}
		if result.AssignedUserID != assignedUserID {
			t.Errorf("expected AssignedUserID %s, got %s", assignedUserID, result.AssignedUserID)
		}
		if result.VisitStatus != "upcoming" {
			t.Errorf("expected VisitStatus 'upcoming', got %s", result.VisitStatus)
		}
	})

	t.Run("Client user not found", func(t *testing.T) {
		// Setup mock behavior
		clientUserID := uuid.New()
		assignedUserID := uuid.New()

		// Create test schedule
		newSchedule := createTestSchedule(uuid.Nil)
		newSchedule.ClientUserID = clientUserID
		newSchedule.AssignedUserID = assignedUserID

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == clientUserID {
				return nil, errors.New("user not found")
			}
			return createTestUser(id), nil
		}

		// Execute
		result, err := useCase.CreateSchedule(newSchedule)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("Assigned user not found", func(t *testing.T) {
		// Setup mock behavior
		clientUserID := uuid.New()
		assignedUserID := uuid.New()

		// Create test schedule
		newSchedule := createTestSchedule(uuid.Nil)
		newSchedule.ClientUserID = clientUserID
		newSchedule.AssignedUserID = assignedUserID

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == clientUserID {
				return createTestUser(id), nil
			}
			if id == assignedUserID {
				return nil, errors.New("user not found")
			}
			return nil, errors.New("user not found")
		}

		// Execute
		result, err := useCase.CreateSchedule(newSchedule)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("Create error", func(t *testing.T) {
		// Setup mock behavior
		clientUserID := uuid.New()
		assignedUserID := uuid.New()

		// Create test schedule
		newSchedule := createTestSchedule(uuid.Nil)
		newSchedule.ClientUserID = clientUserID
		newSchedule.AssignedUserID = assignedUserID

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			return createTestUser(id), nil
		}

		mockScheduleRepo.createFn = func(schedule *domainSchedule.Schedule) (*domainSchedule.Schedule, error) {
			return nil, errors.New("database error")
		}

		// Execute
		result, err := useCase.CreateSchedule(newSchedule)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})
}

// TestGetTodaySchedules tests the GetTodaySchedules method
func TestGetTodaySchedules(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, mockUserRepo, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		userID := uuid.New()
		expectedSchedules := createTestScheduleList(2)

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == userID {
				return createTestUser(userID), nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.getTodaySchedulesFn = func(id uuid.UUID) (*[]domainSchedule.Schedule, error) {
			if id == userID {
				return expectedSchedules, nil
			}
			return nil, errors.New("schedules not found")
		}

		// Execute
		schedules, err := useCase.GetTodaySchedules(userID)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if schedules == nil {
			t.Error("expected non-nil schedules")
		}
		if len(*schedules) != 2 {
			t.Errorf("expected 2 schedules, got %d", len(*schedules))
		}
	})

	t.Run("User not found", func(t *testing.T) {
		// Setup mock behavior
		userID := uuid.New()

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			return nil, errors.New("user not found")
		}

		// Execute
		schedules, err := useCase.GetTodaySchedules(userID)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if schedules != nil {
			t.Error("expected nil schedules")
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		// Setup mock behavior
		userID := uuid.New()

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == userID {
				return createTestUser(userID), nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.getTodaySchedulesFn = func(id uuid.UUID) (*[]domainSchedule.Schedule, error) {
			return nil, errors.New("database error")
		}

		// Execute
		schedules, err := useCase.GetTodaySchedules(userID)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if schedules != nil {
			t.Error("expected nil schedules")
		}
	})
}

// TestGetTodaySchedulesWithClientInfo tests the GetTodaySchedulesWithClientInfo method
func TestGetTodaySchedulesWithClientInfo(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, mockUserRepo, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		userID := uuid.New()
		clientID1 := uuid.New()
		clientID2 := uuid.New()

		// Create test schedules with different client IDs
		schedules := make([]domainSchedule.Schedule, 2)
		schedules[0] = *createTestSchedule(uuid.New())
		schedules[0].ClientUserID = clientID1
		schedules[1] = *createTestSchedule(uuid.New())
		schedules[1].ClientUserID = clientID2

		// Create test clients
		client1 := createTestUser(clientID1)
		client2 := createTestUser(clientID2)

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == userID {
				return createTestUser(userID), nil
			}
			if id == clientID1 {
				return client1, nil
			}
			if id == clientID2 {
				return client2, nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.getTodaySchedulesFn = func(id uuid.UUID) (*[]domainSchedule.Schedule, error) {
			if id == userID {
				return &schedules, nil
			}
			return nil, errors.New("schedules not found")
		}

		// Execute
		resultSchedules, resultClients, err := useCase.GetTodaySchedulesWithClientInfo(userID)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resultSchedules == nil {
			t.Error("expected non-nil schedules")
		}
		if len(*resultSchedules) != 2 {
			t.Errorf("expected 2 schedules, got %d", len(*resultSchedules))
		}
		if resultClients == nil {
			t.Error("expected non-nil clients")
		}
		if len(*resultClients) != 2 {
			t.Errorf("expected 2 clients, got %d", len(*resultClients))
		}
	})

	t.Run("Error getting schedules", func(t *testing.T) {
		// Setup mock behavior
		userID := uuid.New()

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == userID {
				return createTestUser(userID), nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.getTodaySchedulesFn = func(id uuid.UUID) (*[]domainSchedule.Schedule, error) {
			return nil, errors.New("database error")
		}

		// Execute
		schedules, clients, err := useCase.GetTodaySchedulesWithClientInfo(userID)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if schedules != nil {
			t.Error("expected nil schedules")
		}
		if clients != nil {
			t.Error("expected nil clients")
		}
	})

	t.Run("Empty schedules", func(t *testing.T) {
		// Setup mock behavior
		userID := uuid.New()
		emptySchedules := &[]domainSchedule.Schedule{}

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == userID {
				return createTestUser(userID), nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.getTodaySchedulesFn = func(id uuid.UUID) (*[]domainSchedule.Schedule, error) {
			if id == userID {
				return emptySchedules, nil
			}
			return nil, errors.New("schedules not found")
		}

		// Execute
		schedules, clients, err := useCase.GetTodaySchedulesWithClientInfo(userID)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if schedules == nil {
			t.Error("expected non-nil schedules")
		}
		if len(*schedules) != 0 {
			t.Errorf("expected 0 schedules, got %d", len(*schedules))
		}
		if clients == nil {
			t.Error("expected non-nil clients")
		}
		if len(*clients) != 0 {
			t.Errorf("expected 0 clients, got %d", len(*clients))
		}
	})
}

// TestUpdateSchedule tests the UpdateSchedule method
func TestUpdateSchedule(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, mockUserRepo, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()
		clientUserID := uuid.New()
		assignedUserID := uuid.New()

		// Create test schedule
		originalSchedule := createTestSchedule(scheduleID)

		// Create updated schedule
		updatedSchedule := *originalSchedule
		updatedSchedule.ClientUserID = clientUserID
		updatedSchedule.AssignedUserID = assignedUserID
		updatedSchedule.ServiceName = "Updated Service"
		updatedSchedule.VisitStatus = "cancelled"

		// Create updates map
		updates := map[string]interface{}{
			"client_user_id":   clientUserID,
			"assigned_user_id": assignedUserID,
			"service_name":     "Updated Service",
			"visit_status":     "cancelled",
		}

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return originalSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == clientUserID || id == assignedUserID {
				return createTestUser(id), nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.updateScheduleFn = func(id uuid.UUID, u map[string]interface{}) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				// Verify updates
				for key, value := range updates {
					if u[key] != value {
						t.Errorf("expected %s to be %v, got %v", key, value, u[key])
					}
				}
				return &updatedSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		// Execute
		result, err := useCase.UpdateSchedule(scheduleID, updates)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("expected non-nil result")
		}
		if result.ID != scheduleID {
			t.Errorf("expected schedule ID %s, got %s", scheduleID, result.ID)
		}
		if result.ClientUserID != clientUserID {
			t.Errorf("expected ClientUserID %s, got %s", clientUserID, result.ClientUserID)
		}
		if result.AssignedUserID != assignedUserID {
			t.Errorf("expected AssignedUserID %s, got %s", assignedUserID, result.AssignedUserID)
		}
		if result.ServiceName != "Updated Service" {
			t.Errorf("expected ServiceName 'Updated Service', got %s", result.ServiceName)
		}
		if result.VisitStatus != "cancelled" {
			t.Errorf("expected VisitStatus 'cancelled', got %s", result.VisitStatus)
		}
	})

	t.Run("Schedule not found", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			return nil, errors.New("schedule not found")
		}

		// Execute
		updates := map[string]interface{}{
			"service_name": "Updated Service",
		}
		result, err := useCase.UpdateSchedule(scheduleID, updates)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("Client user not found", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()
		clientUserID := uuid.New()

		// Create test schedule
		originalSchedule := createTestSchedule(scheduleID)

		// Create updates map
		updates := map[string]interface{}{
			"client_user_id": clientUserID,
		}

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return originalSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			return nil, errors.New("user not found")
		}

		// Execute
		result, err := useCase.UpdateSchedule(scheduleID, updates)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("Invalid visit status", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()

		// Create test schedule
		originalSchedule := createTestSchedule(scheduleID)

		// Create updates map with invalid status
		updates := map[string]interface{}{
			"visit_status": "invalid_status",
		}

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return originalSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		// Execute
		result, err := useCase.UpdateSchedule(scheduleID, updates)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})

	t.Run("Cannot change from completed status", func(t *testing.T) {
		// Setup mock behavior
		scheduleID := uuid.New()

		// Create test schedule with completed status
		originalSchedule := createTestSchedule(scheduleID)
		originalSchedule.VisitStatus = "completed"

		// Create updates map trying to change status
		updates := map[string]interface{}{
			"visit_status": "cancelled",
		}

		mockScheduleRepo.getScheduleByIDFn = func(id uuid.UUID) (*domainSchedule.Schedule, error) {
			if id == scheduleID {
				return originalSchedule, nil
			}
			return nil, errors.New("schedule not found")
		}

		// Execute
		result, err := useCase.UpdateSchedule(scheduleID, updates)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if result != nil {
			t.Error("expected nil result")
		}
	})
}

// TestGetTodaySchedulesByAssignedUserID tests the GetTodaySchedulesByAssignedUserID method
func TestGetTodaySchedulesByAssignedUserID(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, mockUserRepo, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		assignedUserID := uuid.New()
		expectedSchedules := createTestScheduleList(2)

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == assignedUserID {
				return createTestUser(assignedUserID), nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.getSchedulesByAssignedUserIDPaginatedFn = func(id uuid.UUID, filters domain.DataFilters) (*domainSchedule.SearchResultSchedule, error) {
			if id == assignedUserID {
				return &domainSchedule.SearchResultSchedule{
					Data:       expectedSchedules,
					Total:      2,
					Page:       1,
					PageSize:   10,
					TotalPages: 1,
				}, nil
			}
			return nil, errors.New("schedules not found")
		}

		// Execute
		schedules, err := useCase.GetTodaySchedulesByAssignedUserID(assignedUserID)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if schedules == nil {
			t.Error("expected non-nil schedules")
		}
		if len(*schedules) != 2 {
			t.Errorf("expected 2 schedules, got %d", len(*schedules))
		}
	})

	t.Run("User not found", func(t *testing.T) {
		// Setup mock behavior
		assignedUserID := uuid.New()

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			return nil, errors.New("user not found")
		}

		// Execute
		schedules, err := useCase.GetTodaySchedulesByAssignedUserID(assignedUserID)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if schedules != nil {
			t.Error("expected nil schedules")
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		// Setup mock behavior
		assignedUserID := uuid.New()

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == assignedUserID {
				return createTestUser(assignedUserID), nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.getSchedulesByAssignedUserIDPaginatedFn = func(id uuid.UUID, filters domain.DataFilters) (*domainSchedule.SearchResultSchedule, error) {
			return nil, errors.New("database error")
		}

		// Execute
		schedules, err := useCase.GetTodaySchedulesByAssignedUserID(assignedUserID)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if schedules != nil {
			t.Error("expected nil schedules")
		}
	})
}

// TestGetTodaySchedulesByAssignedUserIDWithClientInfo tests the GetTodaySchedulesByAssignedUserIDWithClientInfo method
func TestGetTodaySchedulesByAssignedUserIDWithClientInfo(t *testing.T) {
	// Setup
	useCase, mockScheduleRepo, mockUserRepo, _ := setupTestScheduleUseCase(t)

	t.Run("Success", func(t *testing.T) {
		// Setup mock behavior
		assignedUserID := uuid.New()
		clientID1 := uuid.New()
		clientID2 := uuid.New()

		// Create test schedules with different client IDs
		schedules := make([]domainSchedule.Schedule, 2)
		schedules[0] = *createTestSchedule(uuid.New())
		schedules[0].ClientUserID = clientID1
		schedules[0].AssignedUserID = assignedUserID
		schedules[1] = *createTestSchedule(uuid.New())
		schedules[1].ClientUserID = clientID2
		schedules[1].AssignedUserID = assignedUserID

		// Create test clients
		client1 := createTestUser(clientID1)
		client2 := createTestUser(clientID2)

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == assignedUserID {
				return createTestUser(assignedUserID), nil
			}
			if id == clientID1 {
				return client1, nil
			}
			if id == clientID2 {
				return client2, nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.getSchedulesByAssignedUserIDPaginatedFn = func(id uuid.UUID, filters domain.DataFilters) (*domainSchedule.SearchResultSchedule, error) {
			if id == assignedUserID {
				return &domainSchedule.SearchResultSchedule{
					Data:       &schedules,
					Total:      2,
					Page:       1,
					PageSize:   10,
					TotalPages: 1,
				}, nil
			}
			return nil, errors.New("schedules not found")
		}

		// Execute
		resultSchedules, resultClients, err := useCase.GetTodaySchedulesByAssignedUserIDWithClientInfo(assignedUserID)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resultSchedules == nil {
			t.Error("expected non-nil schedules")
		}
		if len(*resultSchedules) != 2 {
			t.Errorf("expected 2 schedules, got %d", len(*resultSchedules))
		}
		if resultClients == nil {
			t.Error("expected non-nil clients")
		}
		if len(*resultClients) != 2 {
			t.Errorf("expected 2 clients, got %d", len(*resultClients))
		}
	})

	t.Run("Error getting schedules", func(t *testing.T) {
		// Setup mock behavior
		assignedUserID := uuid.New()

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == assignedUserID {
				return createTestUser(assignedUserID), nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.getSchedulesByAssignedUserIDPaginatedFn = func(id uuid.UUID, filters domain.DataFilters) (*domainSchedule.SearchResultSchedule, error) {
			return nil, errors.New("database error")
		}

		// Execute
		schedules, clients, err := useCase.GetTodaySchedulesByAssignedUserIDWithClientInfo(assignedUserID)

		// Verify
		if err == nil {
			t.Error("expected error, got nil")
		}
		if schedules != nil {
			t.Error("expected nil schedules")
		}
		if clients != nil {
			t.Error("expected nil clients")
		}
	})

	t.Run("Empty schedules", func(t *testing.T) {
		// Setup mock behavior
		assignedUserID := uuid.New()
		emptySchedules := &[]domainSchedule.Schedule{}

		mockUserRepo.getByIDFn = func(id uuid.UUID) (*domainUser.User, error) {
			if id == assignedUserID {
				return createTestUser(assignedUserID), nil
			}
			return nil, errors.New("user not found")
		}

		mockScheduleRepo.getSchedulesByAssignedUserIDPaginatedFn = func(id uuid.UUID, filters domain.DataFilters) (*domainSchedule.SearchResultSchedule, error) {
			if id == assignedUserID {
				return &domainSchedule.SearchResultSchedule{
					Data:       emptySchedules,
					Total:      0,
					Page:       1,
					PageSize:   10,
					TotalPages: 0,
				}, nil
			}
			return nil, errors.New("schedules not found")
		}

		// Execute
		schedules, clients, err := useCase.GetTodaySchedulesByAssignedUserIDWithClientInfo(assignedUserID)

		// Verify
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if schedules == nil {
			t.Error("expected non-nil schedules")
		}
		if len(*schedules) != 0 {
			t.Errorf("expected 0 schedules, got %d", len(*schedules))
		}
		if clients == nil {
			t.Error("expected non-nil clients")
		}
		if len(*clients) != 0 {
			t.Errorf("expected 0 clients, got %d", len(*clients))
		}
	})
}
