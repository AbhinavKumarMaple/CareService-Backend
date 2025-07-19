package user

import (
	"errors"
	"reflect"
	"testing"

	"caregiver/src/domain"
	userDomain "caregiver/src/domain/user"
	logger "caregiver/src/infrastructure/logger"

	"github.com/google/uuid"
)

type mockUserService struct {
	getAllFn     func() (*[]userDomain.User, error)
	getByIDFn    func(id uuid.UUID) (*userDomain.User, error)
	getByEmailFn func(email string) (*userDomain.User, error)
	createFn     func(u *userDomain.User) (*userDomain.User, error)
	deleteFn     func(id uuid.UUID) error
	updateFn     func(id uuid.UUID, m map[string]interface{}) (*userDomain.User, error)
}

func (m *mockUserService) GetAll() (*[]userDomain.User, error) {
	return m.getAllFn()
}
func (m *mockUserService) GetByID(id uuid.UUID) (*userDomain.User, error) {
	return m.getByIDFn(id)
}
func (m *mockUserService) GetByEmail(email string) (*userDomain.User, error) {
	return m.getByEmailFn(email)
}
func (m *mockUserService) Create(newUser *userDomain.User) (*userDomain.User, error) {
	return m.createFn(newUser)
}
func (m *mockUserService) Delete(id uuid.UUID) error {
	return m.deleteFn(id)
}
func (m *mockUserService) Update(id uuid.UUID, userMap map[string]interface{}) (*userDomain.User, error) {
	return m.updateFn(id, userMap)
}
func (m *mockUserService) SearchPaginated(filters domain.DataFilters) (*userDomain.SearchResultUser, error) {
	return nil, nil
}
func (m *mockUserService) SearchByProperty(property string, searchText string) (*[]string, error) {
	return nil, nil
}

func setupLogger(t *testing.T) *logger.Logger {
	loggerInstance, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return loggerInstance
}

func TestUserUseCase(t *testing.T) {

	mockRepo := &mockUserService{}
	logger := setupLogger(t)
	useCase := NewUserUseCase(mockRepo, logger)

	t.Run("Test GetAll", func(t *testing.T) {
		mockRepo.getAllFn = func() (*[]userDomain.User, error) {
			return &[]userDomain.User{{ID: uuid.New()}}, nil
		}
		us, err := useCase.GetAll()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(*us) != 1 {
			t.Error("expected 1 user from GetAll")
		}
	})

	t.Run("Test GetByID", func(t *testing.T) {
		mockRepo.getByIDFn = func(id uuid.UUID) (*userDomain.User, error) {
			if id == uuid.Nil { // Using uuid.Nil for "not found" equivalent
				return nil, errors.New("not found")
			}
			return &userDomain.User{ID: id}, nil
		}
		_, err := useCase.GetByID(uuid.Nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		u, err := useCase.GetByID(uuid.New())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if u == nil || u.ID == uuid.Nil {
			t.Errorf("expected a valid user, got %v", u)
		}
	})

	t.Run("Test GetByEmail", func(t *testing.T) {
		mockRepo.getByEmailFn = func(email string) (*userDomain.User, error) {
			if email == "notfound@example.com" {
				return nil, errors.New("not found")
			}
			return &userDomain.User{ID: uuid.New(), Email: email}, nil
		}
		_, err := useCase.GetByEmail("notfound@example.com")
		if err == nil {
			t.Error("expected error, got nil")
		}
		u, err := useCase.GetByEmail("test@example.com")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if u.Email != "test@example.com" {
			t.Errorf("expected email=test@example.com, got %s", u.Email)
		}
	})

	t.Run("Test Create (OK)", func(t *testing.T) {
		mockRepo.createFn = func(newU *userDomain.User) (*userDomain.User, error) {
			if !newU.Status {
				t.Error("expected user.Status to be true")
			}
			if newU.Email == "" {
				return nil, errors.New("bad data")
			}
			newU.ID = uuid.New()
			return newU, nil
		}
		created, err := useCase.Create(&userDomain.User{Email: "test@mail.com"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if created.ID == uuid.Nil {
			t.Error("expected a non-nil ID after create")
		}
	})

	t.Run("Test Create (Error empty email)", func(t *testing.T) {
		_, err := useCase.Create(&userDomain.User{Email: ""})
		if err == nil {
			t.Error("expected error on create user with empty email")
		}
	})

	t.Run("Test Delete", func(t *testing.T) {
		mockRepo.deleteFn = func(id uuid.UUID) error {
			if id != uuid.Nil {
				return nil
			}
			return errors.New("cannot delete")
		}
		err := useCase.Delete(uuid.Nil)
		if err == nil {
			t.Error("expected error for cannot delete")
		}
		err = useCase.Delete(uuid.New())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Test Update", func(t *testing.T) {
		mockRepo.updateFn = func(id uuid.UUID, m map[string]interface{}) (*userDomain.User, error) {
			if id == uuid.Nil {
				return nil, errors.New("not found")
			}
			return &userDomain.User{ID: id, UserName: "Updated"}, nil
		}
		_, err := useCase.Update(uuid.Nil, map[string]interface{}{"userName": "any"})
		if err == nil {
			t.Error("expected error, got nil")
		}
		updated, err := useCase.Update(uuid.New(), map[string]interface{}{"userName": "whatever"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if updated.UserName != "Updated" {
			t.Error("expected userName=Updated")
		}
	})
}

func TestNewUserUseCase(t *testing.T) {
	mockRepo := &mockUserService{}
	loggerInstance := setupLogger(t)
	useCase := NewUserUseCase(mockRepo, loggerInstance)
	if reflect.TypeOf(useCase).String() != "*user.UserUseCase" {
		t.Error("expected *user.UserUseCase type")
	}
}
