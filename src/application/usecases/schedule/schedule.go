package schedule

import (
	"errors"
	"time"

	"caregiver/src/domain"
	domainErrors "caregiver/src/domain/errors"
	domainSchedule "caregiver/src/domain/schedule"
	domainUser "caregiver/src/domain/user"
	logger "caregiver/src/infrastructure/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IScheduleUseCase interface {
	GetSchedules() (*[]domainSchedule.Schedule, error)
	GetSchedulesWithClientInfo() (*[]domainSchedule.Schedule, *[]domainUser.User, error)
	GetScheduleByID(id uuid.UUID) (*domainSchedule.Schedule, error)
	GetScheduleWithClientInfo(id uuid.UUID) (*domainSchedule.Schedule, *domainUser.User, error)
	GetTodaySchedules(userID uuid.UUID) (*[]domainSchedule.Schedule, error)
	GetTodaySchedulesWithClientInfo(userID uuid.UUID) (*[]domainSchedule.Schedule, *[]domainUser.User, error)
	StartSchedule(scheduleID uuid.UUID, timestamp time.Time, location domainSchedule.Location) (*domainSchedule.Schedule, error)
	EndSchedule(scheduleID uuid.UUID, timestamp time.Time, location domainSchedule.Location, tasks []domainSchedule.Task) (*domainSchedule.Schedule, error)
	UpdateTaskStatus(taskID uuid.UUID, status string, done bool, feedback string) (*domainSchedule.Task, error)
	UpdateSchedule(scheduleID uuid.UUID, updates map[string]interface{}) (*domainSchedule.Schedule, error)
	CreateSchedule(newSchedule *domainSchedule.Schedule) (*domainSchedule.Schedule, error)
	GetTodaySchedulesByAssignedUserID(assignedUserID uuid.UUID) (*[]domainSchedule.Schedule, error)
	GetTodaySchedulesByAssignedUserIDWithClientInfo(assignedUserID uuid.UUID) (*[]domainSchedule.Schedule, *[]domainUser.User, error)
	GetSchedulesInProgressByAssignedUserID(assignedUserID uuid.UUID) (*[]domainSchedule.Schedule, error)
}

type ScheduleUseCase struct {
	scheduleRepository domainSchedule.IScheduleRepository
	userRepository     domainUser.IUserRepository
	Logger             *logger.Logger
}

func NewScheduleUseCase(scheduleRepository domainSchedule.IScheduleRepository, userRepository domainUser.IUserRepository, logger *logger.Logger) IScheduleUseCase {
	return &ScheduleUseCase{
		scheduleRepository: scheduleRepository,
		userRepository:     userRepository,
		Logger:             logger,
	}
}

func (s *ScheduleUseCase) GetSchedules() (*[]domainSchedule.Schedule, error) {
	s.Logger.Info("Getting all schedules")
	return s.scheduleRepository.GetSchedules()
}

func (s *ScheduleUseCase) GetScheduleByID(id uuid.UUID) (*domainSchedule.Schedule, error) {
	s.Logger.Info("Getting schedule by ID", zap.String("id", id.String()))
	return s.scheduleRepository.GetScheduleByID(id)
}

func (s *ScheduleUseCase) GetTodaySchedules(userID uuid.UUID) (*[]domainSchedule.Schedule, error) {
	s.Logger.Info("Getting today's schedules for user", zap.String("userID", userID.String()))
	// Optional: Validate if userID exists in the system
	_, err := s.userRepository.GetByID(userID)
	if err != nil {
		s.Logger.Error("User not found for today's schedules", zap.Error(err), zap.String("userID", userID.String()))
		return nil, domainErrors.NewAppError(errors.New("user not found"), domainErrors.NotFound)
	}
	return s.scheduleRepository.GetTodaySchedules(userID)
}

func (s *ScheduleUseCase) StartSchedule(scheduleID uuid.UUID, timestamp time.Time, location domainSchedule.Location) (*domainSchedule.Schedule, error) {
	s.Logger.Info("Starting schedule", zap.String("scheduleID", scheduleID.String()))

	schedule, err := s.scheduleRepository.GetScheduleByID(scheduleID)
	if err != nil {
		s.Logger.Error("Schedule not found for start", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		return nil, err
	}

	if schedule.VisitStatus != "upcoming" {
		s.Logger.Warn("Cannot start schedule, invalid status", zap.String("scheduleID", scheduleID.String()), zap.String("status", schedule.VisitStatus))
		return nil, domainErrors.NewAppError(errors.New("schedule is not in 'upcoming' status"), domainErrors.ValidationError)
	}

	// Check if the current time is before the scheduled start time
	if timestamp.Before(schedule.ScheduledSlot.From) {
		s.Logger.Warn("Cannot start schedule before scheduled time",
			zap.String("scheduleID", scheduleID.String()),
			zap.Time("currentTime", timestamp),
			zap.Time("scheduledStartTime", schedule.ScheduledSlot.From))
		return nil, domainErrors.NewAppError(errors.New("cannot start schedule before the scheduled start time"), domainErrors.ValidationError)
	}

	// Check if there are any other schedules in progress for the same assigned user
	schedulesInProgress, err := s.scheduleRepository.GetSchedulesInProgressByAssignedUserID(schedule.AssignedUserID)
	if err != nil {
		s.Logger.Error("Error checking for schedules in progress", zap.Error(err), zap.String("assignedUserID", schedule.AssignedUserID.String()))
		return nil, err
	}

	if schedulesInProgress != nil && len(*schedulesInProgress) > 0 {
		s.Logger.Warn("Cannot start schedule, another schedule is already in progress",
			zap.String("scheduleID", scheduleID.String()),
			zap.String("assignedUserID", schedule.AssignedUserID.String()),
			zap.Int("inProgressCount", len(*schedulesInProgress)))
		return nil, domainErrors.NewAppError(errors.New("cannot start schedule: another schedule is already in progress for this user"), domainErrors.ValidationError)
	}

	updates := map[string]interface{}{
		"visit_status":          "in_progress",
		"checkin_time":          timestamp,
		"checkin_location_lat":  location.Lat,
		"checkin_location_long": location.Long,
	}

	updatedSchedule, err := s.scheduleRepository.UpdateSchedule(scheduleID, updates)
	if err != nil {
		s.Logger.Error("Error updating schedule for start", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		return nil, err
	}
	s.Logger.Info("Schedule started successfully", zap.String("scheduleID", scheduleID.String()))
	return updatedSchedule, nil
}

func (s *ScheduleUseCase) EndSchedule(scheduleID uuid.UUID, timestamp time.Time, location domainSchedule.Location, tasks []domainSchedule.Task) (*domainSchedule.Schedule, error) {
	s.Logger.Info("Ending schedule", zap.String("scheduleID", scheduleID.String()))

	schedule, err := s.scheduleRepository.GetScheduleByID(scheduleID)
	if err != nil {
		s.Logger.Error("Schedule not found for end", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		return nil, err
	}

	if schedule.VisitStatus != "in_progress" {
		s.Logger.Warn("Cannot end schedule, invalid status", zap.String("scheduleID", scheduleID.String()), zap.String("status", schedule.VisitStatus))
		return nil, domainErrors.NewAppError(errors.New("schedule is not in 'in_progress' status"), domainErrors.ValidationError)
	}

	updates := map[string]interface{}{
		"visit_status":           "completed",
		"checkout_time":          timestamp,
		"checkout_location_lat":  location.Lat,
		"checkout_location_long": location.Long,
	}

	updatedSchedule, err := s.scheduleRepository.UpdateSchedule(scheduleID, updates)
	if err != nil {
		s.Logger.Error("Error updating schedule for end", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		return nil, err
	}

	for _, task := range tasks {
		_, err := s.scheduleRepository.UpdateTask(task.ID, map[string]interface{}{
			"status":   task.Status,
			"done":     task.Done,
			"feedback": task.Feedback,
		})
		if err != nil {
			s.Logger.Error("Error updating task during EndSchedule", zap.Error(err), zap.String("taskID", task.ID.String()))

		}
	}
	if err != nil {
		s.Logger.Error("Error updating schedule for end", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		return nil, err
	}
	s.Logger.Info("Schedule ended successfully", zap.String("scheduleID", scheduleID.String()))
	return updatedSchedule, nil
}

func (s *ScheduleUseCase) UpdateTaskStatus(taskID uuid.UUID, status string, done bool, feedback string) (*domainSchedule.Task, error) {
	s.Logger.Info("Updating task status", zap.String("taskID", taskID.String()))

	updates := map[string]interface{}{
		"Status":   status,
		"Done":     done,
		"Feedback": feedback,
	}

	updatedTask, err := s.scheduleRepository.UpdateTask(taskID, updates)
	if err != nil {
		s.Logger.Error("Error updating task status", zap.Error(err), zap.String("taskID", taskID.String()))
		return nil, err
	}
	s.Logger.Info("Task status updated successfully", zap.String("taskID", taskID.String()))
	return updatedTask, nil
}

func (s *ScheduleUseCase) CreateSchedule(newSchedule *domainSchedule.Schedule) (*domainSchedule.Schedule, error) {
	s.Logger.Info("Creating new schedule", zap.String("clientUserID", newSchedule.ClientUserID.String()), zap.String("assignedUserID", newSchedule.AssignedUserID.String()))

	_, err := s.userRepository.GetByID(newSchedule.ClientUserID)
	if err != nil {
		s.Logger.Error("Client user not found for schedule creation", zap.Error(err), zap.String("clientUserID", newSchedule.ClientUserID.String()))
		return nil, domainErrors.NewAppError(errors.New("client user not found"), domainErrors.NotFound)
	}

	_, err = s.userRepository.GetByID(newSchedule.AssignedUserID)
	if err != nil {
		s.Logger.Error("Assigned user not found for schedule creation", zap.Error(err), zap.String("assignedUserID", newSchedule.AssignedUserID.String()))
		return nil, domainErrors.NewAppError(errors.New("assigned user not found"), domainErrors.NotFound)
	}

	newSchedule.VisitStatus = "upcoming"

	for i := range newSchedule.Tasks {
		if newSchedule.Tasks[i].ID == uuid.Nil {
			newSchedule.Tasks[i].ID = uuid.New()
		}
		newSchedule.Tasks[i].Status = "pending"
	}

	createdSchedule, err := s.scheduleRepository.Create(newSchedule)
	if err != nil {
		s.Logger.Error("Error creating schedule in repository", zap.Error(err), zap.String("clientUserID", newSchedule.ClientUserID.String()))
		return nil, err
	}

	s.Logger.Info("Schedule created successfully in use case", zap.String("scheduleID", createdSchedule.ID.String()))
	return createdSchedule, nil
}

func (s *ScheduleUseCase) GetTodaySchedulesByAssignedUserID(assignedUserID uuid.UUID) (*[]domainSchedule.Schedule, error) {
	s.Logger.Info("Getting today's schedules by assigned user ID", zap.String("assignedUserID", assignedUserID.String()))

	_, err := s.userRepository.GetByID(assignedUserID)
	if err != nil {
		s.Logger.Error("Error getting today's schedules by assigned user ID", zap.Error(err), zap.String("assignedUserID", assignedUserID.String()))
		return nil, domainErrors.NewAppError(errors.New("assigned user not found"), domainErrors.NotFound)
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour).Add(-time.Nanosecond) // End of today

	filters := domain.DataFilters{
		DateRangeFilters: []domain.DateRangeFilter{
			{
				Field: "scheduled_slot_from",
				Start: &todayStart,
				End:   &todayEnd,
			},
		},
	}

	schedulesResult, err := s.scheduleRepository.GetSchedulesByAssignedUserIDPaginated(assignedUserID, filters)
	if err != nil {
		s.Logger.Error("Error retrieving today's schedules by assigned user ID from repository", zap.Error(err), zap.String("assignedUserID", assignedUserID.String()))
		return nil, err
	}

	return schedulesResult.Data, nil
}

func (s *ScheduleUseCase) GetScheduleWithClientInfo(id uuid.UUID) (*domainSchedule.Schedule, *domainUser.User, error) {
	s.Logger.Info("Getting schedule with client info by ID", zap.String("id", id.String()))

	schedule, err := s.scheduleRepository.GetScheduleByID(id)
	if err != nil {
		s.Logger.Error("Schedule not found", zap.Error(err), zap.String("id", id.String()))
		return nil, nil, err
	}

	client, err := s.userRepository.GetByID(schedule.ClientUserID)
	if err != nil {
		s.Logger.Error("Client user not found", zap.Error(err), zap.String("clientUserID", schedule.ClientUserID.String()))
		return schedule, nil, nil
	}

	return schedule, client, nil
}

func (s *ScheduleUseCase) GetTodaySchedulesWithClientInfo(userID uuid.UUID) (*[]domainSchedule.Schedule, *[]domainUser.User, error) {
	s.Logger.Info("Getting today's schedules with client info for user", zap.String("userID", userID.String()))

	schedules, err := s.GetTodaySchedules(userID)
	if err != nil {
		return nil, nil, err
	}

	if schedules == nil || len(*schedules) == 0 {
		return schedules, &[]domainUser.User{}, nil
	}

	clientIDs := make(map[uuid.UUID]bool)
	for _, schedule := range *schedules {
		clientIDs[schedule.ClientUserID] = true
	}

	clients := make([]domainUser.User, 0, len(clientIDs))
	for clientID := range clientIDs {
		client, err := s.userRepository.GetByID(clientID)
		if err != nil {
			s.Logger.Warn("Client user not found", zap.Error(err), zap.String("clientUserID", clientID.String()))
			continue
		}
		clients = append(clients, *client)
	}

	return schedules, &clients, nil
}

func (s *ScheduleUseCase) GetTodaySchedulesByAssignedUserIDWithClientInfo(assignedUserID uuid.UUID) (*[]domainSchedule.Schedule, *[]domainUser.User, error) {
	s.Logger.Info("Getting today's schedules with client info by assigned user ID", zap.String("assignedUserID", assignedUserID.String()))

	schedules, err := s.GetTodaySchedulesByAssignedUserID(assignedUserID)
	if err != nil {
		return nil, nil, err
	}

	if schedules == nil || len(*schedules) == 0 {
		return schedules, &[]domainUser.User{}, nil
	}

	clientIDs := make(map[uuid.UUID]bool)
	for _, schedule := range *schedules {
		clientIDs[schedule.ClientUserID] = true
	}

	clients := make([]domainUser.User, 0, len(clientIDs))
	for clientID := range clientIDs {
		client, err := s.userRepository.GetByID(clientID)
		if err != nil {
			s.Logger.Warn("Client user not found", zap.Error(err), zap.String("clientUserID", clientID.String()))
			continue
		}
		clients = append(clients, *client)
	}

	return schedules, &clients, nil
}

func (s *ScheduleUseCase) GetSchedulesWithClientInfo() (*[]domainSchedule.Schedule, *[]domainUser.User, error) {
	s.Logger.Info("Getting all schedules with client info")

	schedules, err := s.scheduleRepository.GetSchedules()
	if err != nil {
		s.Logger.Error("Error getting all schedules", zap.Error(err))
		return nil, nil, err
	}

	if schedules == nil || len(*schedules) == 0 {
		return schedules, &[]domainUser.User{}, nil
	}

	clientIDs := make(map[uuid.UUID]bool)
	for _, schedule := range *schedules {
		clientIDs[schedule.ClientUserID] = true
	}

	clients := make([]domainUser.User, 0, len(clientIDs))
	for clientID := range clientIDs {
		client, err := s.userRepository.GetByID(clientID)
		if err != nil {
			s.Logger.Warn("Client user not found", zap.Error(err), zap.String("clientUserID", clientID.String()))
			continue
		}
		clients = append(clients, *client)
	}

	return schedules, &clients, nil
}

func (s *ScheduleUseCase) UpdateSchedule(scheduleID uuid.UUID, updates map[string]interface{}) (*domainSchedule.Schedule, error) {
	s.Logger.Info("Updating schedule", zap.String("scheduleID", scheduleID.String()))

	existingSchedule, err := s.scheduleRepository.GetScheduleByID(scheduleID)
	if err != nil {
		s.Logger.Error("Schedule not found for update", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		return nil, domainErrors.NewAppError(errors.New("schedule not found"), domainErrors.NotFound)
	}

	if clientUserID, ok := updates["client_user_id"].(uuid.UUID); ok {
		_, err := s.userRepository.GetByID(clientUserID)
		if err != nil {
			s.Logger.Error("New client user not found", zap.Error(err), zap.String("clientUserID", clientUserID.String()))
			return nil, domainErrors.NewAppError(errors.New("new client user not found"), domainErrors.NotFound)
		}
	}

	if assignedUserID, ok := updates["assigned_user_id"].(uuid.UUID); ok {
		_, err := s.userRepository.GetByID(assignedUserID)
		if err != nil {
			s.Logger.Error("New assigned user not found", zap.Error(err), zap.String("assignedUserID", assignedUserID.String()))
			return nil, domainErrors.NewAppError(errors.New("new assigned user not found"), domainErrors.NotFound)
		}
	}

	if status, ok := updates["visit_status"].(string); ok {
		validStatuses := map[string]bool{
			"upcoming":    true,
			"in_progress": true,
			"completed":   true,
			"cancelled":   true,
		}

		if !validStatuses[status] {
			s.Logger.Error("Invalid visit status", zap.String("status", status))
			return nil, domainErrors.NewAppError(errors.New("invalid visit status"), domainErrors.ValidationError)
		}

		currentStatus := existingSchedule.VisitStatus

		if currentStatus == "completed" && status != "completed" {
			s.Logger.Error("Cannot change status from completed", zap.String("currentStatus", currentStatus), zap.String("newStatus", status))
			return nil, domainErrors.NewAppError(errors.New("cannot change status from completed"), domainErrors.ValidationError)
		}

		if currentStatus == "cancelled" && status != "cancelled" {
			s.Logger.Error("Cannot change status from cancelled", zap.String("currentStatus", currentStatus), zap.String("newStatus", status))
			return nil, domainErrors.NewAppError(errors.New("cannot change status from cancelled"), domainErrors.ValidationError)
		}
	}

	updatedSchedule, err := s.scheduleRepository.UpdateSchedule(scheduleID, updates)
	if err != nil {
		s.Logger.Error("Error updating schedule", zap.Error(err), zap.String("scheduleID", scheduleID.String()))
		return nil, err
	}

	s.Logger.Info("Schedule updated successfully", zap.String("scheduleID", scheduleID.String()))
	return updatedSchedule, nil
}

func (s *ScheduleUseCase) GetSchedulesInProgressByAssignedUserID(assignedUserID uuid.UUID) (*[]domainSchedule.Schedule, error) {
	s.Logger.Info("Getting schedules in progress by assigned user ID", zap.String("assignedUserID", assignedUserID.String()))
	return s.scheduleRepository.GetSchedulesInProgressByAssignedUserID(assignedUserID)
}
