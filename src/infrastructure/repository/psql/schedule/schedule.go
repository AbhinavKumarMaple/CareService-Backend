package schedule

import (
	"encoding/json"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainSchedule "github.com/gbrayhan/microservices-go/src/domain/schedule"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Schedule struct {
	ID               uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ClientUserID     uuid.UUID `gorm:"column:client_user_id;type:uuid"`
	AssignedUserID   uuid.UUID `gorm:"column:assigned_user_id;type:uuid"` 
	ServiceName      string    `gorm:"column:service_name"`               
	ScheduledSlotFrom time.Time `gorm:"column:scheduled_slot_from"`
	ScheduledSlotTo   time.Time `gorm:"column:scheduled_slot_to"`
	VisitStatus      string    `gorm:"column:visit_status"`
	CheckinTime      *time.Time `gorm:"column:checkin_time"`
	CheckoutTime     *time.Time `gorm:"column:checkout_time"`
	CheckinLocationLat  *float64 `gorm:"column:checkin_location_lat"`
	CheckinLocationLong *float64 `gorm:"column:checkin_location_long"`
	CheckoutLocationLat  *float64 `gorm:"column:checkout_location_lat"`
	CheckoutLocationLong *float64 `gorm:"column:checkout_location_long"`
	Tasks            []Task    `gorm:"foreignKey:ScheduleID"`
	ServiceNote      *string   `gorm:"column:service_note"`
	CreatedAt        time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime:milli"`
}

type Task struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ScheduleID  uuid.UUID  `gorm:"column:schedule_id;type:uuid"`
	Title       string     `gorm:"column:title"`
	Description string     `gorm:"column:description"`
	Status      string     `gorm:"column:status"`
	Done        *bool      `gorm:"column:done"`
	Feedback    *string    `gorm:"column:feedback"`
	CreatedAt   time.Time  `gorm:"autoCreateTime:milli"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime:milli"`
}

func (Schedule) TableName() string {
	return "schedules"
}

func (Task) TableName() string {
	return "tasks"
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewScheduleRepository(db *gorm.DB, loggerInstance *logger.Logger) domainSchedule.IScheduleRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetSchedules() (*[]domainSchedule.Schedule, error) {
	var schedules []Schedule
	if err := r.DB.Preload("Tasks").Find(&schedules).Error; err != nil {
		r.Logger.Error("Error getting all schedules", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainMapper(&schedules), nil
}

func (r *Repository) GetScheduleByID(id uuid.UUID) (*domainSchedule.Schedule, error) {
	var schedule Schedule
	err := r.DB.Preload("Tasks").Where("id = ?", id).First(&schedule).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Schedule not found", zap.String("id", id.String()))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting schedule by ID", zap.Error(err), zap.String("id", id.String()))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return nil, err
	}
	return schedule.toDomainMapper(), nil
}

func (r *Repository) GetTodaySchedules(userID uuid.UUID) (*[]domainSchedule.Schedule, error) {
	var schedules []Schedule
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	if err := r.DB.Preload("Tasks").
		Where("client_user_id = ?", userID).
		Where("scheduled_slot_from >= ? AND scheduled_slot_from < ?", today, tomorrow).
		Find(&schedules).Error; err != nil {
		r.Logger.Error("Error getting today's schedules", zap.Error(err), zap.String("userID", userID.String()))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainMapper(&schedules), nil
}

func (r *Repository) UpdateSchedule(id uuid.UUID, updates map[string]interface{}) (*domainSchedule.Schedule, error) {
	var scheduleObj Schedule
	if err := r.DB.Preload("Tasks").Where("id = ?", id).First(&scheduleObj).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Schedule not found for update", zap.String("id", id.String()))
			return nil, domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		}
		r.Logger.Error("Error retrieving schedule for update", zap.Error(err), zap.String("id", id.String()))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}


	err := r.DB.Model(&scheduleObj).Updates(updates).Error
	if err != nil {
		r.Logger.Error("Error updating schedule", zap.Error(err), zap.String("id", id.String()))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return nil, errUnmarshal
		}
		switch newError.Number {
		case 1062: 
			return nil, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}

	if err := r.DB.Preload("Tasks").Where("id = ?", id).First(&scheduleObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated schedule", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	return scheduleObj.toDomainMapper(), nil
}

func (r *Repository) UpdateTask(taskID uuid.UUID, updates map[string]interface{}) (*domainSchedule.Task, error) {
	var taskObj Task
	taskObj.ID = taskID

	err := r.DB.Model(&taskObj).Updates(updates).Error
	if err != nil {
		r.Logger.Error("Error updating task", zap.Error(err), zap.String("taskID", taskID.String()))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	if err := r.DB.Where("id = ?", taskID).First(&taskObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated task", zap.Error(err), zap.String("taskID", taskID.String()))
		return nil, err
	}

	return taskObj.toDomainMapper(), nil
}

func (s *Schedule) toDomainMapper() *domainSchedule.Schedule {
	tasksDomain := make([]domainSchedule.Task, len(s.Tasks))
	for i, task := range s.Tasks {
		tasksDomain[i] = *task.toDomainMapper()
	}

	return &domainSchedule.Schedule{
		ID:               s.ID,
		ClientUserID:     s.ClientUserID,
		AssignedUserID:   s.AssignedUserID, 
		ServiceName:      s.ServiceName,     
		ScheduledSlot: domainSchedule.ScheduledSlot{
			From: s.ScheduledSlotFrom,
			To:   s.ScheduledSlotTo,
		},
		VisitStatus: s.VisitStatus,
		CheckinTime: s.CheckinTime,
		CheckoutTime: s.CheckoutTime,
		CheckinLocation: domainSchedule.Location{
			Lat:  s.CheckinLocationLat,
			Long: s.CheckinLocationLong,
		},
		CheckoutLocation: domainSchedule.Location{
			Lat:  s.CheckoutLocationLat,
			Long: s.CheckoutLocationLong,
		},
		Tasks:       tasksDomain,
		ServiceNote: s.ServiceNote,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func (t *Task) toDomainMapper() *domainSchedule.Task {
	return &domainSchedule.Task{
		ID:          t.ID,
		ScheduleID:  t.ScheduleID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Done:        t.Done,
		Feedback:    t.Feedback,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func arrayToDomainMapper(schedules *[]Schedule) *[]domainSchedule.Schedule {
	schedulesDomain := make([]domainSchedule.Schedule, len(*schedules))
	for i, schedule := range *schedules {
		schedulesDomain[i] = *schedule.toDomainMapper()
	}
	return &schedulesDomain
}

func (r *Repository) Create(newSchedule *domainSchedule.Schedule) (*domainSchedule.Schedule, error) {
	r.Logger.Info("Creating new schedule in repository", zap.String("clientUserID", newSchedule.ClientUserID.String()))

	scheduleModel := fromDomainMapper(newSchedule)

	err := r.DB.Create(scheduleModel).Error
	if err != nil {
		r.Logger.Error("Error creating schedule", zap.Error(err), zap.String("clientUserID", newSchedule.ClientUserID.String()))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return nil, errUnmarshal
		}
		switch newError.Number {
		case 1062: 
			return nil, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}

	r.Logger.Info("Schedule created successfully in repository", zap.String("scheduleID", scheduleModel.ID.String()))
	return scheduleModel.toDomainMapper(), nil
}

func fromDomainMapper(s *domainSchedule.Schedule) *Schedule {
	tasksModel := make([]Task, len(s.Tasks))
	for i, task := range s.Tasks {
		tasksModel[i] = Task{
			ID:          task.ID,
			ScheduleID:  task.ScheduleID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Done:        task.Done,
			Feedback:    task.Feedback,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}

	return &Schedule{
		ID:               s.ID,
		ClientUserID:     s.ClientUserID,
		AssignedUserID:   s.AssignedUserID, 
		ServiceName:      s.ServiceName,     
		ScheduledSlotFrom: s.ScheduledSlot.From,
		ScheduledSlotTo:   s.ScheduledSlot.To,
		VisitStatus:      s.VisitStatus,
		CheckinTime:      s.CheckinTime,
		CheckoutTime:     s.CheckoutTime,
		CheckinLocationLat:  s.CheckinLocation.Lat,
		CheckinLocationLong: s.CheckinLocation.Long,
		CheckoutLocationLat:  s.CheckoutLocation.Lat,
		CheckoutLocationLong: s.CheckoutLocation.Long,
		Tasks:            tasksModel,
		ServiceNote:      s.ServiceNote,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}

func (r *Repository) GetSchedulesByAssignedUserIDPaginated(assignedUserID uuid.UUID, filters domain.DataFilters) (*domainSchedule.SearchResultSchedule, error) {

	query := r.DB.Session(&gorm.Session{PrepareStmt: false}).Model(&Schedule{}).Preload("Tasks").Where("assigned_user_id = ?", assignedUserID)

	for _, dateFilter := range filters.DateRangeFilters {
		if dateFilter.Field == "scheduled_slot_from" { // Assuming filtering on scheduled_slot_from
			if dateFilter.Start != nil {
				query = query.Where("scheduled_slot_from >= ?", dateFilter.Start)
			}
			if dateFilter.End != nil {
				query = query.Where("scheduled_slot_from <= ?", dateFilter.End)
			}
		}
	}

	if len(filters.SortBy) > 0 && filters.SortDirection.IsValid() {
		for _, sortField := range filters.SortBy {
			query = query.Order(sortField + " " + string(filters.SortDirection))
		}
	}

	var total int64
	clonedQuery := query
	clonedQuery.Count(&total)

	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 10
	}
	offset := (filters.Page - 1) * filters.PageSize

	var schedules []Schedule
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&schedules).Error; err != nil {
		r.Logger.Error("Error searching schedules by assigned user ID", zap.Error(err), zap.String("assignedUserID", assignedUserID.String()))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domainSchedule.SearchResultSchedule{
		Data:       arrayToDomainMapper(&schedules),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	r.Logger.Info("Successfully searched schedules by assigned user ID",
		zap.String("assignedUserID", assignedUserID.String()),
		zap.Int64("total", total),
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))

	return result, nil
}

func (r *Repository) GetSchedulesInProgressByAssignedUserID(assignedUserID uuid.UUID) (*[]domainSchedule.Schedule, error) {
	var schedules []Schedule
	if err := r.DB.Preload("Tasks").
		Where("assigned_user_id = ? AND visit_status = ?", assignedUserID, "in_progress").
		Find(&schedules).Error; err != nil {
		r.Logger.Error("Error getting schedules in progress by assigned user ID", zap.Error(err), zap.String("assignedUserID", assignedUserID.String()))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainMapper(&schedules), nil
}