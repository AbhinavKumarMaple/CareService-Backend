package user

import (
	"caregiver/src/domain"
	userDomain "caregiver/src/domain/user"
	logger "caregiver/src/infrastructure/logger"
	"caregiver/src/infrastructure/repository/psql/user"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IUserUseCase interface {
	GetAll() (*[]userDomain.User, error)
	GetByID(id uuid.UUID) (*userDomain.User, error)
	GetByEmail(email string) (*userDomain.User, error)
	Create(newUser *userDomain.User) (*userDomain.User, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, userMap map[string]interface{}) (*userDomain.User, error)
	SearchPaginated(filters domain.DataFilters) (*userDomain.SearchResultUser, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
}

type UserUseCase struct {
	userRepository user.UserRepositoryInterface
	Logger         *logger.Logger
}

func NewUserUseCase(userRepository user.UserRepositoryInterface, logger *logger.Logger) IUserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
		Logger:         logger,
	}
}

func (s *UserUseCase) GetAll() (*[]userDomain.User, error) {
	s.Logger.Info("Getting all users")
	return s.userRepository.GetAll()
}

func (s *UserUseCase) GetByID(id uuid.UUID) (*userDomain.User, error) {
	s.Logger.Info("Getting user by ID", zap.String("id", id.String()))
	return s.userRepository.GetByID(id)
}

func (s *UserUseCase) GetByEmail(email string) (*userDomain.User, error) {
	s.Logger.Info("Getting user by email", zap.String("email", email))
	return s.userRepository.GetByEmail(email)
}

func (s *UserUseCase) Create(newUser *userDomain.User) (*userDomain.User, error) {
	s.Logger.Info("Creating new user", zap.String("email", newUser.Email))

	newUser.Status = true
	newUser.ID = uuid.New()

	return s.userRepository.Create(newUser)
}

func (s *UserUseCase) Delete(id uuid.UUID) error {
	s.Logger.Info("Deleting user", zap.String("id", id.String()))
	return s.userRepository.Delete(id)
}

func (s *UserUseCase) Update(id uuid.UUID, userMap map[string]interface{}) (*userDomain.User, error) {
	s.Logger.Info("Updating user", zap.String("id", id.String()))
	return s.userRepository.Update(id, userMap)
}

func (s *UserUseCase) SearchPaginated(filters domain.DataFilters) (*userDomain.SearchResultUser, error) {
	s.Logger.Info("Searching users with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.userRepository.SearchPaginated(filters)
}

func (s *UserUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching users by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.userRepository.SearchByProperty(property, searchText)
}
