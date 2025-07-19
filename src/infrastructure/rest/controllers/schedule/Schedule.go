package schedule

import (
	"errors"
	"net/http"

	scheduleUseCase "caregiver/src/application/usecases/schedule"
	domainErrors "caregiver/src/domain/errors"
	domainSchedule "caregiver/src/domain/schedule"
	domainUser "caregiver/src/domain/user"
	logger "caregiver/src/infrastructure/logger"
	"caregiver/src/infrastructure/rest/controllers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IScheduleController interface {
	GetSchedules(ctx *gin.Context)
	GetTodaySchedules(ctx *gin.Context)
	GetScheduleByID(ctx *gin.Context)
	StartSchedule(ctx *gin.Context)
	EndSchedule(ctx *gin.Context)
	UpdateTask(ctx *gin.Context)
	UpdateSchedule(ctx *gin.Context)
	CreateSchedule(ctx *gin.Context)
	GetTodaySchedulesByAssignedUserID(ctx *gin.Context)
}

type Controller struct {
	scheduleUseCase scheduleUseCase.IScheduleUseCase
	Logger          *logger.Logger
}

func NewScheduleController(scheduleUseCase scheduleUseCase.IScheduleUseCase, loggerInstance *logger.Logger) IScheduleController {
	return &Controller{scheduleUseCase: scheduleUseCase, Logger: loggerInstance}
}

func (c *Controller) GetSchedules(ctx *gin.Context) {
	c.Logger.Info("Getting all schedules")
	schedules, clients, err := c.scheduleUseCase.GetSchedulesWithClientInfo()
	if err != nil {
		c.Logger.Error("Error getting all schedules", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved all schedules", zap.Int("count", len(*schedules)))
	ctx.JSON(http.StatusOK, arrayDomainToResponseMapperWithClients(*schedules, *clients))
}

func (c *Controller) CreateSchedule(ctx *gin.Context) {
	c.Logger.Info("Creating new schedule")
	var request CreateScheduleRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new schedule", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	if request.ClientUserID == uuid.Nil {
		c.Logger.Error("ClientUserID is required for new schedule")
		appError := domainErrors.NewAppError(errors.New("ClientUserID is required"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	if request.ScheduledSlot.From.IsZero() || request.ScheduledSlot.To.IsZero() {
		c.Logger.Error("ScheduledSlot (From, To) is required for new schedule")
		appError := domainErrors.NewAppError(errors.New("ScheduledSlot (From, To) is required"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	if request.ScheduledSlot.From.After(request.ScheduledSlot.To) {
		c.Logger.Error("ScheduledSlot 'From' cannot be after 'To'")
		appError := domainErrors.NewAppError(errors.New("ScheduledSlot 'From' cannot be after 'To'"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	if len(request.Tasks) == 0 {
		c.Logger.Error("At least one task is required for new schedule")
		appError := domainErrors.NewAppError(errors.New("at least one task is required"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	domainTasks := make([]domainSchedule.Task, len(request.Tasks))
	for i, taskReq := range request.Tasks {
		domainTasks[i] = domainSchedule.Task{
			Title:       taskReq.Title,
			Description: taskReq.Description,
			Status:      "pending",
			Done:        nil,
			Feedback:    nil,
		}
	}

	newSchedule := &domainSchedule.Schedule{
		ClientUserID:   request.ClientUserID,
		AssignedUserID: request.AssignedUserID,
		ServiceName:    request.ServiceName,
		ScheduledSlot:  domainSchedule.ScheduledSlot{From: request.ScheduledSlot.From, To: request.ScheduledSlot.To},
		Tasks:          domainTasks,
		VisitStatus:    "upcoming",
	}

	createdSchedule, err := c.scheduleUseCase.CreateSchedule(newSchedule)
	if err != nil {
		c.Logger.Error("Error creating schedule", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	c.Logger.Info("Schedule created successfully", zap.String("scheduleID", createdSchedule.ID.String()))
	ctx.JSON(http.StatusOK, domainToResponseMapper(createdSchedule))
}

func clientToResponseMapper(u *domainUser.User) *ClientInfo {
	if u == nil {
		return nil
	}
	return &ClientInfo{
		ID:             u.ID,
		UserName:       u.UserName,
		Email:          u.Email,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		ProfilePicture: u.ProfilePicture,
		Location: ClientLocation{
			HouseNumber: u.Location.HouseNumber,
			Street:      u.Location.Street,
			City:        u.Location.City,
			State:       u.Location.State,
			Pincode:     u.Location.Pincode,
			Lat:         u.Location.Lat,
			Long:        u.Location.Long,
		},
	}
}

func domainToResponseMapper(s *domainSchedule.Schedule) *ScheduleResponse {
	tasksResponse := make([]Task, len(s.Tasks))
	for i, task := range s.Tasks {
		tasksResponse[i] = Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Done:        task.Done,
			Feedback:    task.Feedback,
		}
	}

	return &ScheduleResponse{
		ID:             s.ID,
		ClientUserID:   s.ClientUserID,
		ClientInfo:     nil,
		AssignedUserID: s.AssignedUserID,
		ServiceName:    s.ServiceName,
		ScheduledSlot: ScheduledSlot{
			From: s.ScheduledSlot.From,
			To:   s.ScheduledSlot.To,
		},
		VisitStatus:  s.VisitStatus,
		CheckinTime:  s.CheckinTime,
		CheckoutTime: s.CheckoutTime,
		CheckinLocation: Location{
			Lat:  s.CheckinLocation.Lat,
			Long: s.CheckinLocation.Long,
		},
		CheckoutLocation: Location{
			Lat:  s.CheckoutLocation.Lat,
			Long: s.CheckoutLocation.Long,
		},
		Tasks:       tasksResponse,
		ServiceNote: s.ServiceNote,
	}
}

func arrayDomainToResponseMapper(schedules []domainSchedule.Schedule) []ScheduleResponse {
	res := make([]ScheduleResponse, len(schedules))
	for i, s := range schedules {
		res[i] = *domainToResponseMapper(&s)
	}
	return res
}

func arrayDomainToResponseMapperWithClients(schedules []domainSchedule.Schedule, clients []domainUser.User) []ScheduleResponse {
	res := make([]ScheduleResponse, len(schedules))

	// Create a map for quick client lookup
	clientMap := make(map[uuid.UUID]*domainUser.User)
	for i := range clients {
		clientMap[clients[i].ID] = &clients[i]
	}

	for i, s := range schedules {
		response := domainToResponseMapper(&s)
		if client, exists := clientMap[s.ClientUserID]; exists {
			response.ClientInfo = clientToResponseMapper(client)
		}
		res[i] = *response
	}
	return res
}

func (c *Controller) GetTodaySchedules(ctx *gin.Context) {
	c.Logger.Info("Getting today's schedules")

	userIDStr := ctx.Query("ClientUserID")
	if userIDStr == "" {
		c.Logger.Error("Missing ClientUserID query parameter for today's schedules")
		appError := domainErrors.NewAppError(errors.New("ClientUserID query parameter is required"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Logger.Error("Invalid ClientUserID format", zap.Error(err), zap.String("ClientUserID", userIDStr))
		appError := domainErrors.NewAppError(errors.New("Invalid ClientUserID format"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	schedules, clients, err := c.scheduleUseCase.GetTodaySchedulesWithClientInfo(userID)
	if err != nil {
		c.Logger.Error("Error getting today's schedules", zap.Error(err), zap.String("userID", userID.String()))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved today's schedules", zap.Int("count", len(*schedules)), zap.String("userID", userID.String()))
	ctx.JSON(http.StatusOK, arrayDomainToResponseMapperWithClients(*schedules, *clients))
}

func (c *Controller) GetScheduleByID(ctx *gin.Context) {
	scheduleIDStr := ctx.Param("id")
	scheduleID, err := uuid.Parse(scheduleIDStr)
	if err != nil {
		c.Logger.Error("Invalid schedule ID parameter", zap.Error(err), zap.String("id", scheduleIDStr))
		appError := domainErrors.NewAppError(errors.New("schedule id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting schedule by ID", zap.String("id", scheduleID.String()))
	schedule, client, err := c.scheduleUseCase.GetScheduleWithClientInfo(scheduleID)
	if err != nil {
		c.Logger.Error("Error getting schedule by ID", zap.Error(err), zap.String("id", scheduleID.String()))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved schedule by ID", zap.String("id", scheduleID.String()))

	response := domainToResponseMapper(schedule)
	response.ClientInfo = clientToResponseMapper(client)
	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) StartSchedule(ctx *gin.Context) {
	scheduleIDStr := ctx.Param("id")
	scheduleID, err := uuid.Parse(scheduleIDStr)
	if err != nil {
		c.Logger.Error("Invalid schedule ID parameter for start", zap.Error(err), zap.String("id", scheduleIDStr))
		appError := domainErrors.NewAppError(errors.New("schedule id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	var request StartScheduleRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for start schedule", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	if request.Timestamp.IsZero() {
		c.Logger.Error("Timestamp is required for start schedule", zap.String("ScheduleID", scheduleID.String()))
		appError := domainErrors.NewAppError(errors.New("Timestamp is required"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	if request.Location.Lat == nil || request.Location.Long == nil {
		c.Logger.Error("Location (Lat, Long) is required for start schedule", zap.String("ScheduleID", scheduleID.String()))
		appError := domainErrors.NewAppError(errors.New("Location (Lat, Long) is required"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	schedule, err := c.scheduleUseCase.StartSchedule(scheduleID, request.Timestamp, domainSchedule.Location{Lat: request.Location.Lat, Long: request.Location.Long})
	if err != nil {
		c.Logger.Error("Error starting schedule", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		_ = ctx.Error(err)
		return
	}

	c.Logger.Info("Schedule started successfully", zap.String("scheduleID", scheduleID.String()))
	ctx.JSON(http.StatusOK, StartScheduleResponse{
		Message:         "Check-in recorded successfully",
		CheckinTime:     schedule.CheckinTime,
		CheckinLocation: &Location{Lat: schedule.CheckinLocation.Lat, Long: schedule.CheckinLocation.Long},
	})
}

func (c *Controller) EndSchedule(ctx *gin.Context) {
	scheduleIDStr := ctx.Param("id")
	scheduleID, err := uuid.Parse(scheduleIDStr)
	if err != nil {
		c.Logger.Error("Invalid schedule ID parameter for end", zap.Error(err), zap.String("id", scheduleIDStr))
		appError := domainErrors.NewAppError(errors.New("schedule id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	var request EndScheduleRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for end schedule", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	if request.Timestamp.IsZero() {
		c.Logger.Error("Timestamp is required for end schedule", zap.String("ScheduleID", scheduleID.String()))
		appError := domainErrors.NewAppError(errors.New("Timestamp is required"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	if request.Location.Lat == nil || request.Location.Long == nil {
		c.Logger.Error("Location (Lat, Long) is required for end schedule", zap.String("ScheduleID", scheduleID.String()))
		appError := domainErrors.NewAppError(errors.New("Location (Lat, Long) is required"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	domainTasks := make([]domainSchedule.Task, len(request.Tasks))
	for i, taskReq := range request.Tasks {
		domainTasks[i] = domainSchedule.Task{
			ID:          taskReq.ID,
			Title:       taskReq.Title,
			Description: taskReq.Description,
			Status:      taskReq.Status,
			Done:        taskReq.Done,
			Feedback:    taskReq.Feedback,
		}
	}

	schedule, err := c.scheduleUseCase.EndSchedule(scheduleID, request.Timestamp, domainSchedule.Location{Lat: request.Location.Lat, Long: request.Location.Long}, domainTasks)
	if err != nil {
		c.Logger.Error("Error ending schedule", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		_ = ctx.Error(err)
		return
	}

	c.Logger.Info("Schedule ended successfully", zap.String("scheduleID", scheduleID.String()))
	ctx.JSON(http.StatusOK, EndScheduleResponse{
		Message:          "Check-out recorded successfully",
		CheckoutTime:     schedule.CheckoutTime,
		CheckoutLocation: &Location{Lat: schedule.CheckoutLocation.Lat, Long: schedule.CheckoutLocation.Long},
	})
}

func (c *Controller) UpdateTask(ctx *gin.Context) {
	taskIDStr := ctx.Param("taskId") // Corrected to match route parameter case

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.Logger.Error("Invalid task ID parameter for update ", zap.Error(err), zap.String("taskID", taskIDStr))
		appError := domainErrors.NewAppError(errors.New("task id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	var request UpdateTaskRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for update task", zap.Error(err), zap.String("taskID", taskID.String()))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	if request.Status == "" {
		c.Logger.Error("Status is required for task update", zap.String("TaskID", taskID.String()))
		appError := domainErrors.NewAppError(errors.New("Status is required"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	if request.Done == nil {
		c.Logger.Error("Done field is required for task update", zap.String("TaskID", taskID.String()))
		appError := domainErrors.NewAppError(errors.New("Done field is required"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	var feedback string
	if request.Feedback != nil {
		feedback = *request.Feedback
	}

	updatedTask, err := c.scheduleUseCase.UpdateTaskStatus(taskID, request.Status, *request.Done, feedback)
	if err != nil {
		c.Logger.Error("Error updating task status", zap.Error(err), zap.String("taskID", taskID.String()))
		_ = ctx.Error(err)
		return
	}

	c.Logger.Info("Task updated successfully", zap.String("taskID", taskID.String()))
	ctx.JSON(http.StatusOK, UpdateTaskResponse{
		Message: "Task updated successfully",
		Task:    Task{ID: updatedTask.ID, Status: updatedTask.Status, Done: updatedTask.Done, Feedback: updatedTask.Feedback},
	})
}

func (c *Controller) GetTodaySchedulesByAssignedUserID(ctx *gin.Context) {
	assignedUserIDStr := ctx.Param("assignedUserID")
	assignedUserID, err := uuid.Parse(assignedUserIDStr)
	if err != nil {
		c.Logger.Error("Invalid assigned user ID parameter", zap.Error(err), zap.String("assignedUserID", assignedUserIDStr))
		appError := domainErrors.NewAppError(errors.New("assigned user ID is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	c.Logger.Info("Getting today's schedules by assigned user ID", zap.String("assignedUserID", assignedUserID.String()))

	schedules, clients, err := c.scheduleUseCase.GetTodaySchedulesByAssignedUserIDWithClientInfo(assignedUserID)
	if err != nil {
		c.Logger.Error("Error getting today's schedules by assigned user ID", zap.Error(err), zap.String("assignedUserID", assignedUserID.String()))
		_ = ctx.Error(err)
		return
	}

	c.Logger.Info("Successfully retrieved today's schedules by assigned user ID", zap.Int("count", len(*schedules)), zap.String("assignedUserID", assignedUserID.String()))
	ctx.JSON(http.StatusOK, arrayDomainToResponseMapperWithClients(*schedules, *clients))
}

func (c *Controller) UpdateSchedule(ctx *gin.Context) {
	scheduleIDStr := ctx.Param("id")
	scheduleID, err := uuid.Parse(scheduleIDStr)
	if err != nil {
		c.Logger.Error("Invalid schedule ID parameter for update", zap.Error(err), zap.String("id", scheduleIDStr))
		appError := domainErrors.NewAppError(errors.New("schedule id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	_, _, err = c.scheduleUseCase.GetScheduleWithClientInfo(scheduleID)
	if err != nil {
		c.Logger.Error("Error getting schedule for update", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		_ = ctx.Error(err)
		return
	}

	var request UpdateScheduleRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for update schedule", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	updates := make(map[string]interface{})

	if request.ClientUserID != uuid.Nil {
		updates["client_user_id"] = request.ClientUserID
	}

	if request.AssignedUserID != uuid.Nil {
		updates["assigned_user_id"] = request.AssignedUserID
	}

	if request.ServiceName != "" {
		updates["service_name"] = request.ServiceName
	}

	if request.VisitStatus != "" {
		updates["visit_status"] = request.VisitStatus
	}

	if request.ScheduledSlot != nil {
		if request.ScheduledSlot.From.IsZero() || request.ScheduledSlot.To.IsZero() {
			c.Logger.Error("Both From and To dates must be provided for ScheduledSlot", zap.String("scheduleID", scheduleID.String()))
			appError := domainErrors.NewAppError(errors.New("Both From and To dates must be provided for ScheduledSlot"), domainErrors.ValidationError)
			_ = ctx.Error(appError)
			return
		}

		if request.ScheduledSlot.From.After(request.ScheduledSlot.To) {
			c.Logger.Error("ScheduledSlot 'From' cannot be after 'To'", zap.String("scheduleID", scheduleID.String()))
			appError := domainErrors.NewAppError(errors.New("ScheduledSlot 'From' cannot be after 'To'"), domainErrors.ValidationError)
			_ = ctx.Error(appError)
			return
		}

		updates["scheduled_slot_from"] = request.ScheduledSlot.From
		updates["scheduled_slot_to"] = request.ScheduledSlot.To
	}

	if len(updates) == 0 {
		c.Logger.Warn("No valid fields to update", zap.String("scheduleID", scheduleID.String()))
		appError := domainErrors.NewAppError(errors.New("No valid fields to update"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	updatedSchedule, err := c.scheduleUseCase.UpdateSchedule(scheduleID, updates)
	if err != nil {
		c.Logger.Error("Error updating schedule", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		_ = ctx.Error(err)
		return
	}

	_, client, _ := c.scheduleUseCase.GetScheduleWithClientInfo(scheduleID)

	response := domainToResponseMapper(updatedSchedule)
	response.ClientInfo = clientToResponseMapper(client)

	c.Logger.Info("Schedule updated successfully", zap.String("scheduleID", scheduleID.String()))
	ctx.JSON(http.StatusOK, UpdateScheduleResponse{
		Message:  "Schedule updated successfully",
		Schedule: response,
	})
}
