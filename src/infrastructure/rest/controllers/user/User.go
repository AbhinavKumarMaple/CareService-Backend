package user

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"caregiver/src/domain"
	domainErrors "caregiver/src/domain/errors"
	domainUser "caregiver/src/domain/user"
	logger "caregiver/src/infrastructure/logger"
	"caregiver/src/infrastructure/repository/psql/user"
	"caregiver/src/infrastructure/rest/controllers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Structures
type LocationRequest struct {
	HouseNumber string  `json:"HouseNumber"`
	Street      string  `json:"Street"`
	City        string  `json:"City"`
	State       string  `json:"State"`
	Pincode     string  `json:"Pincode"`
	Lat         float64 `json:"Lat"`
	Long        float64 `json:"Long"`
}

type NewUserRequest struct {
	UserName  string          `json:"UserName" binding:"required"`
	Email     string          `json:"Email" binding:"required"`
	FirstName string          `json:"FirstName" binding:"required"`
	LastName  string          `json:"LastName" binding:"required"`
	Role      string          `json:"Role" binding:"required"`
	Location  LocationRequest `json:"Location"`
}

type ResponseUser struct {
	ID        uuid.UUID       `json:"ID"`
	UserName  string          `json:"UserName"`
	Email     string          `json:"Email"`
	FirstName string          `json:"FirstName"`
	LastName  string          `json:"LastName"`
	Status    bool            `json:"Status"`
	Role      string          `json:"Role"`
	Location  LocationRequest `json:"Location"`
	CreatedAt time.Time       `json:"CreatedAt,omitempty"`
	UpdatedAt time.Time       `json:"UpdatedAt,omitempty"`
}

type IUserController interface {
	NewUser(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetUsersByID(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
	SearchByProperty(ctx *gin.Context)
}

type UserController struct {
	userService domainUser.IUserService
	Logger      *logger.Logger
}

func NewUserController(userService domainUser.IUserService, loggerInstance *logger.Logger) IUserController {
	return &UserController{userService: userService, Logger: loggerInstance}
}

func (c *UserController) NewUser(ctx *gin.Context) {
	c.Logger.Info("Creating new user")
	var request NewUserRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new user", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	userModel, err := c.userService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating user", zap.Error(err), zap.String("email", request.Email))
		_ = ctx.Error(err)
		return
	}
	userResponse := domainToResponseMapper(userModel)
	c.Logger.Info("User created successfully", zap.String("email", request.Email), zap.String("id", userModel.ID.String()))
	ctx.JSON(http.StatusOK, userResponse)
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	c.Logger.Info("Getting all users")
	users, err := c.userService.GetAll()
	if err != nil {
		c.Logger.Error("Error getting all users", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all users", zap.Int("count", len(*users)))
	ctx.JSON(http.StatusOK, arrayDomainToResponseMapper(users))
}

func (c *UserController) GetUsersByID(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("user id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting user by ID", zap.String("id", userID.String()))
	user, err := c.userService.GetByID(userID)
	if err != nil {
		c.Logger.Error("Error getting user by ID", zap.Error(err), zap.String("id", userID.String()))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved user by ID", zap.String("id", userID.String()))
	ctx.JSON(http.StatusOK, domainToResponseMapper(user))
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating user", zap.String("id", userID.String()))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for user update", zap.Error(err), zap.String("id", userID.String()))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for user update", zap.Error(err), zap.String("id", userID.String()))
		_ = ctx.Error(err)
		return
	}
	userUpdated, err := c.userService.Update(userID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating user", zap.Error(err), zap.String("id", userID.String()))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("User updated successfully", zap.String("id", userID.String()))
	ctx.JSON(http.StatusOK, domainToResponseMapper(userUpdated))
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting user", zap.String("id", userID.String()))
	err = c.userService.Delete(userID)
	if err != nil {
		c.Logger.Error("Error deleting user", zap.Error(err), zap.String("id", userID.String()))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("User deleted successfully", zap.String("id", userID.String()))
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

func (c *UserController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching users with pagination")

	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	// Build filters
	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	// Parse like filters
	likeFilters := make(map[string][]string)
	for field := range user.ColumnsUserMapping {
		if values := ctx.QueryArray(field + "_Like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range user.ColumnsUserMapping {
		if values := ctx.QueryArray(field + "_Match"); len(values) > 0 {
			matches[field] = values
		}
	}
	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range user.ColumnsUserMapping {
		startStr := ctx.Query(field + "_Start")
		endStr := ctx.Query(field + "_End")

		if startStr != "" || endStr != "" {
			dateRange := domain.DateRangeFilter{Field: field}

			if startStr != "" {
				if startTime, err := time.Parse(time.RFC3339, startStr); err == nil {
					dateRange.Start = &startTime
				}
			}

			if endStr != "" {
				if endTime, err := time.Parse(time.RFC3339, endStr); err == nil {
					dateRange.End = &endTime
				}
			}

			dateRanges = append(dateRanges, dateRange)
		}
	}
	filters.DateRangeFilters = dateRanges

	// Parse sorting
	sortBy := ctx.QueryArray("sortBy")
	if len(sortBy) > 0 {
		filters.SortBy = sortBy
	}

	sortDirection := domain.SortDirection(ctx.DefaultQuery("sortDirection", "asc"))
	if sortDirection.IsValid() {
		filters.SortDirection = sortDirection
	}

	result, err := c.userService.SearchPaginated(filters)
	if err != nil {
		c.Logger.Error("Error searching users", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	response := gin.H{
		"Data":       arrayDomainToResponseMapper(result.Data),
		"Total":      result.Total,
		"Page":       result.Page,
		"PageSize":   result.PageSize,
		"TotalPages": result.TotalPages,
		"Filters":    filters,
	}

	c.Logger.Info("Successfully searched users",
		zap.Int64("total", result.Total),
		zap.Int("page", result.Page))
	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) SearchByProperty(ctx *gin.Context) {
	property := ctx.Query("property")
	searchText := ctx.Query("searchText")

	if property == "" || searchText == "" {
		c.Logger.Error("Missing property or searchText parameter")
		appError := domainErrors.NewAppError(errors.New("missing property or searchText parameter"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	// Validate property
	allowed := map[string]bool{
		"UserName":    true,
		"Email":       true,
		"FirstName":   true,
		"LastName":    true,
		"Status":      true,
		"Role":        true,
		"HouseNumber": true,
		"Street":      true,
		"City":        true,
		"State":       true,
		"Pincode":     true,
		"Lat":         true,
		"Long":        true,
	}
	if !allowed[property] {
		c.Logger.Error("Invalid property for search", zap.String("property", property))
		appError := domainErrors.NewAppError(errors.New("invalid property"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	coincidences, err := c.userService.SearchByProperty(property, searchText)
	if err != nil {
		c.Logger.Error("Error searching by property", zap.Error(err), zap.String("property", property))
		_ = ctx.Error(err)
		return
	}

	c.Logger.Info("Successfully searched by property",
		zap.String("property", property),
		zap.Int("results", len(*coincidences)))
	ctx.JSON(http.StatusOK, coincidences)
}

// Mappers
func domainToResponseMapper(domainUser *domainUser.User) *ResponseUser {
	return &ResponseUser{
		ID:        domainUser.ID,
		UserName:  domainUser.UserName,
		Email:     domainUser.Email,
		FirstName: domainUser.FirstName,
		LastName:  domainUser.LastName,
		Status:    domainUser.Status,
		Role:      domainUser.Role,
		Location: LocationRequest{
			HouseNumber: domainUser.Location.HouseNumber,
			Street:      domainUser.Location.Street,
			City:        domainUser.Location.City,
			State:       domainUser.Location.State,
			Pincode:     domainUser.Location.Pincode,
			Lat:         domainUser.Location.Lat,
			Long:        domainUser.Location.Long,
		},
		CreatedAt: domainUser.CreatedAt,
		UpdatedAt: domainUser.UpdatedAt,
	}
}

func arrayDomainToResponseMapper(users *[]domainUser.User) *[]ResponseUser {
	res := make([]ResponseUser, len(*users))
	for i, u := range *users {
		res[i] = *domainToResponseMapper(&u)
	}
	return &res
}

func toUsecaseMapper(req *NewUserRequest) *domainUser.User {
	return &domainUser.User{
		UserName:  req.UserName,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		Location: domainUser.Location{
			HouseNumber: req.Location.HouseNumber,
			Street:      req.Location.Street,
			City:        req.Location.City,
			State:       req.Location.State,
			Pincode:     req.Location.Pincode,
			Lat:         req.Location.Lat,
			Long:        req.Location.Long,
		},
	}
}
