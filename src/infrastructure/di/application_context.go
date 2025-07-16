package di

import (
	"sync"

	authUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/auth"
	scheduleUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/schedule"
	userUseCase "github.com/gbrayhan/microservices-go/src/application/usecases/user"
	domainSchedule "github.com/gbrayhan/microservices-go/src/domain/schedule"
	scheduleRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/schedule"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql"
	userRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	scheduleController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/schedule"
	userController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	"gorm.io/gorm"
)

type ApplicationContext struct {
	DB                 *gorm.DB
	Logger             *logger.Logger
	AuthController     authController.IAuthController
	UserController     userController.IUserController
	ScheduleController scheduleController.IScheduleController
	JWTService         security.IJWTService
	UserRepository     userRepo.UserRepositoryInterface
	ScheduleRepository domainSchedule.IScheduleRepository
	AuthUseCase        authUseCase.IAuthUseCase
	UserUseCase        userUseCase.IUserUseCase
	ScheduleUseCase    scheduleUseCase.IScheduleUseCase
}

var (
	loggerInstance *logger.Logger
	loggerOnce     sync.Once
)

func GetLogger() *logger.Logger {
	loggerOnce.Do(func() {
		loggerInstance, _ = logger.NewLogger()
	})
	return loggerInstance
}

func SetupDependencies(loggerInstance *logger.Logger) (*ApplicationContext, error) {
	db, err := psql.InitPSQLDB(loggerInstance)
	if err != nil {
		return nil, err
	}

	jwtService := security.NewJWTService()

	userRepo := userRepo.NewUserRepository(db, loggerInstance)
	scheduleRepo := scheduleRepo.NewScheduleRepository(db, loggerInstance)

	authUC := authUseCase.NewAuthUseCase(userRepo, jwtService, loggerInstance)
	userUC := userUseCase.NewUserUseCase(userRepo, loggerInstance)
	scheduleUC := scheduleUseCase.NewScheduleUseCase(scheduleRepo, userRepo, loggerInstance)

	authController := authController.NewAuthController(authUC, loggerInstance)
	userController := userController.NewUserController(userUC, loggerInstance)
	scheduleController := scheduleController.NewScheduleController(scheduleUC, loggerInstance)

	return &ApplicationContext{
		DB:                 db,
		Logger:             loggerInstance,
		AuthController:     authController,
		UserController:     userController,
		ScheduleController: scheduleController,
		JWTService:         jwtService,
		UserRepository:     userRepo,
		ScheduleRepository: scheduleRepo,
		AuthUseCase:        authUC,
		UserUseCase:        userUC,
		ScheduleUseCase:    scheduleUC,
	}, nil
}

func NewTestApplicationContext(
	mockUserRepo userRepo.UserRepositoryInterface,
	mockScheduleRepo domainSchedule.IScheduleRepository,
	mockJWTService security.IJWTService,
	loggerInstance *logger.Logger,
) *ApplicationContext {
	authUC := authUseCase.NewAuthUseCase(mockUserRepo, mockJWTService, loggerInstance)
	userUC := userUseCase.NewUserUseCase(mockUserRepo, loggerInstance)
	scheduleUC := scheduleUseCase.NewScheduleUseCase(mockScheduleRepo, mockUserRepo, loggerInstance)

	authController := authController.NewAuthController(authUC, loggerInstance)
	userController := userController.NewUserController(userUC, loggerInstance)
	scheduleController := scheduleController.NewScheduleController(scheduleUC, loggerInstance)

	return &ApplicationContext{
		Logger:             loggerInstance,
		AuthController:     authController,
		UserController:     userController,
		ScheduleController: scheduleController,
		JWTService:         mockJWTService,
		UserRepository:     mockUserRepo,
		ScheduleRepository: mockScheduleRepo,
		AuthUseCase:        authUC,
		UserUseCase:        userUC,
		ScheduleUseCase:    scheduleUC,
	}
}
