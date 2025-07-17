package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAll() (*[]domainUser.User, error) {
	args := m.Called()
	return args.Get(0).(*[]domainUser.User), args.Error(1)
}

func (m *MockUserService) GetByID(id uuid.UUID) (*domainUser.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserService) Create(user *domainUser.User) (*domainUser.User, error) {
	args := m.Called(user)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserService) GetOneByMap(userMap map[string]interface{}) (*domainUser.User, error) {
	args := m.Called(userMap)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserService) Update(id uuid.UUID, userMap map[string]interface{}) (*domainUser.User, error) {
	args := m.Called(id, userMap)
	return args.Get(0).(*domainUser.User), args.Error(1)
}

func (m *MockUserService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) SearchPaginated(filters domain.DataFilters) (*domainUser.SearchResultUser, error) {
	args := m.Called(filters)
	return args.Get(0).(*domainUser.SearchResultUser), args.Error(1)
}

func (m *MockUserService) SearchByProperty(property string, searchText string) (*[]string, error) {
	args := m.Called(property, searchText)
	return args.Get(0).(*[]string), args.Error(1)
}

func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
}

func TestNewUserController(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	assert.NotNil(t, controller)
	assert.Equal(t, mockService, controller.(*UserController).userService)
	assert.Equal(t, loggerInstance, controller.(*UserController).Logger)
}

func TestDomainToResponseMapper(t *testing.T) {
	now := time.Now()
	domainUser := &domainUser.User{
		ID:           uuid.New(),
		UserName:     "testuser",
		Email:        "test@example.com",
		FirstName:    "Test",
		LastName:     "User",
		Status:       true,
		HashPassword: "hashedpassword",
		Role:         "caregiver",
		Location:     domainUser.Location{HouseNumber: "1", Street: "Main St"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	response := domainToResponseMapper(domainUser)

	assert.Equal(t, domainUser.ID, response.ID)
	assert.Equal(t, domainUser.UserName, response.UserName)
	assert.Equal(t, domainUser.Email, response.Email)
	assert.Equal(t, domainUser.FirstName, response.FirstName)
	assert.Equal(t, domainUser.LastName, response.LastName)
	assert.Equal(t, domainUser.Status, response.Status)
	assert.Equal(t, domainUser.CreatedAt, response.CreatedAt)
	assert.Equal(t, domainUser.UpdatedAt, response.UpdatedAt)
}

func TestArrayDomainToResponseMapper(t *testing.T) {
	now := time.Now()
	users := []domainUser.User{
		{
			ID:           uuid.New(),
			UserName:     "user1",
			Email:        "user1@example.com",
			FirstName:    "User",
			LastName:     "One",
			Status:       true,
			HashPassword: "hash1",
			Role:         "caregiver",
			Location:     domainUser.Location{HouseNumber: "1", Street: "Main St"},
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New(),
			UserName:     "user2",
			Email:        "user2@example.com",
			FirstName:    "User",
			LastName:     "Two",
			Status:       false,
			HashPassword: "hash2",
			Role:         "client",
			Location:     domainUser.Location{HouseNumber: "2", Street: "Second St"},
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	responses := arrayDomainToResponseMapper(&users)

	assert.Len(t, *responses, 2)
	assert.Equal(t, users[0].ID, (*responses)[0].ID)
	assert.Equal(t, users[1].ID, (*responses)[1].ID)
}

func TestToUsecaseMapper(t *testing.T) {
	request := &NewUserRequest{
		UserName:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role: "user",
		Location: LocationRequest{
			HouseNumber: "1",
			Street:      "Main St",
			City:        "Anytown",
			State:       "CA",
			Pincode:     "12345",
			Lat:         34.0522,
			Long:        -118.2437,
		},
	}

	domainUser := toUsecaseMapper(request)

	assert.Equal(t, request.UserName, domainUser.UserName)
	assert.Equal(t, request.Email, domainUser.Email)
	assert.Equal(t, request.FirstName, domainUser.FirstName)
	assert.Equal(t, request.LastName, domainUser.LastName)
	assert.Equal(t, request.Role, domainUser.Role)
	assert.Equal(t, request.Location.HouseNumber, domainUser.Location.HouseNumber)
	assert.Equal(t, request.Location.Street, domainUser.Location.Street)
	assert.Equal(t, request.Location.City, domainUser.Location.City)
	assert.Equal(t, request.Location.State, domainUser.Location.State)
	assert.Equal(t, request.Location.Pincode, domainUser.Location.Pincode)
	assert.Equal(t, request.Location.Lat, domainUser.Location.Lat)
	assert.Equal(t, request.Location.Long, domainUser.Location.Long)
}

func TestUpdateValidation(t *testing.T) {
	validRequest := map[string]any{
		"user_name": "validuser",
		"email":     "valid@example.com",
		"firstName": "Valid",
		"lastName":  "User",
	}

	err := updateValidation(validRequest)
	assert.NoError(t, err)

	emptyRequest := map[string]any{
		"user_name": "",
		"email":     "",
	}

	err = updateValidation(emptyRequest)
	assert.Error(t, err)

	invalidEmailRequest := map[string]any{
		"email": "invalid-email",
	}

	err = updateValidation(invalidEmailRequest)
	assert.Error(t, err)

	shortUserNameRequest := map[string]any{
		"user_name": "ab",
	}

	err = updateValidation(shortUserNameRequest)
	assert.Error(t, err)

	longUserNameRequest := map[string]any{
		"user_name": "verylongusernameverylongusernameverylongusernameverylongusernameverylongusernameverylongusernameverylongusernameverylongusernameverylongusernameverylongusername",
	}

	err = updateValidation(longUserNameRequest)
	assert.Error(t, err)

	shortFirstNameRequest := map[string]any{
		"firstName": "a",
	}

	err = updateValidation(shortFirstNameRequest)
	assert.Error(t, err)

	longFirstNameRequest := map[string]any{
		"firstName": "verylongfirstnameverylongfirstnameverylongfirstnameverylongfirstnameverylongfirstnameverylongfirstnameverylongfirstnameverylongfirstnameverylongfirstname",
	}

	err = updateValidation(longFirstNameRequest)
	assert.Error(t, err)
}

func setupGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	router := gin.New()
	// Add error handling middleware
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			c.JSON(500, gin.H{"error": err.Error()})
		}
	})
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	return c, w
}

func TestUserController_NewUser(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	t.Run("Success", func(t *testing.T) {
		c, w := setupGinContext()
		request := NewUserRequest{
			UserName:  "testuser",
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Role:      "user",
			Location:  LocationRequest{HouseNumber: "1", Street: "Main St"},
		}
		jsonData, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		expectedUser := &domainUser.User{
			ID:           uuid.New(),
			UserName:     "testuser",
			Email:        "test@example.com",
			FirstName:    "Test",
			LastName:     "User",
			Status:       true,
			HashPassword: "hashedpassword",
			Role:         "user",
			Location:     domainUser.Location{HouseNumber: "1", Street: "Main St"},
		}

		mockService.On("Create", mock.Anything).Return(expectedUser, nil)

		controller.NewUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("POST", "/users", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.NewUser(c)

		assert.Equal(t, http.StatusOK, w.Code) 
	})

	t.Run("Service Error", func(t *testing.T) {
		c, w := setupGinContext()
		request := NewUserRequest{
			UserName:  "testuser",
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Role: "user",
			Location: LocationRequest{HouseNumber: "1", Street: "Main St"},
		}
		jsonData, _ := json.Marshal(request)
		c.Request = httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		mockService.On("Create", mock.Anything).Return((*domainUser.User)(nil), errors.New("service error"))

		controller.NewUser(c)

		assert.Equal(t, http.StatusOK, w.Code) 
		mockService.AssertExpectations(t)
	})
}

func TestUserController_GetAllUsers(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	t.Run("Success", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("GET", "/users", nil)

		expectedUsers := &[]domainUser.User{
			{ID: uuid.New(), UserName: "user1", Email: "user1@example.com"},
			{ID: uuid.New(), UserName: "user2", Email: "user2@example.com"},
		}

		mockService.On("GetAll").Return(expectedUsers, nil)

		controller.GetAllUsers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("GET", "/users", nil)

		mockService.On("GetAll").Return((*[]domainUser.User)(nil), errors.New("service error"))

		controller.GetAllUsers(c)

		assert.Equal(t, http.StatusOK, w.Code) 
		mockService.AssertExpectations(t)
	})
}

func TestUserController_GetUsersByID(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	t.Run("Success", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("GET", "/users/1", nil)
		id := uuid.New()
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		expectedUser := &domainUser.User{
			ID:       id,
			UserName: "user1",
			Email:    "user1@example.com",
		}

		mockService.On("GetByID", id).Return(expectedUser, nil)

		controller.GetUsersByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("GET", "/users/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		controller.GetUsersByID(c)

		assert.Equal(t, http.StatusOK, w.Code) 
	})

	t.Run("Service Error", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("GET", "/users/1", nil)
		id := uuid.New()
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		mockService.On("GetByID", id).Return((*domainUser.User)(nil), errors.New("service error"))

		controller.GetUsersByID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserController_UpdateUser(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	t.Run("Success", func(t *testing.T) {
		c, w := setupGinContext()
		updateData := map[string]any{
			"user_name": "updateduser",
			"email":     "updated@example.com",
		}
		jsonData, _ := json.Marshal(updateData)
		c.Request = httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")
		id := uuid.New()
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		expectedUser := &domainUser.User{
			ID:       id,
			UserName: "updateduser",
			Email:    "updated@example.com",
		}

		mockService.On("Update", id, updateData).Return(expectedUser, nil)

		controller.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("PUT", "/users/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		controller.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code) 
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("PUT", "/users/1", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")
		id := uuid.New()
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		controller.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code) 
	})

	t.Run("Service Error", func(t *testing.T) {
		c, w := setupGinContext()
		updateData := map[string]any{"user_name": "updateduser"}
		jsonData, _ := json.Marshal(updateData)
		c.Request = httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")
		id := uuid.New()
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		mockService.On("Update", id, updateData).Return((*domainUser.User)(nil), errors.New("service error"))

		controller.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code) 
		mockService.AssertExpectations(t)
	})
}

func TestUserController_DeleteUser(t *testing.T) {
	mockService := &MockUserService{}
	loggerInstance := setupLogger(t)
	controller := NewUserController(mockService, loggerInstance)

	t.Run("Success", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("DELETE", "/users/1", nil)
		id := uuid.New()
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		mockService.On("Delete", id).Return(nil)

		controller.DeleteUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("DELETE", "/users/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		controller.DeleteUser(c)

		assert.Equal(t, http.StatusOK, w.Code) 
	})

	t.Run("Service Error", func(t *testing.T) {
		c, w := setupGinContext()
		c.Request = httptest.NewRequest("DELETE", "/users/1", nil)
		id := uuid.New()
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		mockService.On("Delete", id).Return(errors.New("service error"))

		controller.DeleteUser(c)

		assert.Equal(t, http.StatusOK, w.Code) 
		mockService.AssertExpectations(t)
	})
}
