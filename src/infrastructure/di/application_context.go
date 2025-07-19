package di

import (
	"sync"

	authUseCase "caregiver/src/application/usecases/auth"
	scheduleUseCase "caregiver/src/application/usecases/schedule"
	userUseCase "caregiver/src/application/usecases/user"
	domainSchedule "caregiver/src/domain/schedule"
	scheduleRepo "caregiver/src/infrastructure/repository/psql/schedule"

	logger "caregiver/src/infrastructure/logger"
	"caregiver/src/infrastructure/repository/psql"
	userRepo "caregiver/src/infrastructure/repository/psql/user"
	authController "caregiver/src/infrastructure/rest/controllers/auth"
	scheduleController "caregiver/src/infrastructure/rest/controllers/schedule"
	userController "caregiver/src/infrastructure/rest/controllers/user"
	"caregiver/src/infrastructure/security"

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
