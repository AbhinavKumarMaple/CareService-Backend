package schedule

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	"github.com/google/uuid"
)

type Schedule struct {
	ID               uuid.UUID     `gorm:"primaryKey"`
	ClientUserID     uuid.UUID     `gorm:"column:client_user_id"`
	AssignedUserID   uuid.UUID     `gorm:"column:assigned_user_id"` 
	ServiceName      string        `gorm:"column:service_name"`     
	ScheduledSlot    ScheduledSlot `gorm:"embedded;embeddedPrefix:scheduled_slot_"`
	VisitStatus      string        `gorm:"column:visit_status"` 
	CheckinTime      *time.Time    `gorm:"column:checkin_time"`
	CheckoutTime     *time.Time    `gorm:"column:checkout_time"`
	CheckinLocation  Location      `gorm:"embedded;embeddedPrefix:checkin_location_"`
	CheckoutLocation Location      `gorm:"embedded;embeddedPrefix:checkout_location_"`
	Tasks            []Task        `gorm:"foreignKey:ScheduleID"` 
	ServiceNote      *string       `gorm:"column:service_note"`
	CreatedAt        time.Time     `gorm:"autoCreateTime:milli"`
	UpdatedAt        time.Time     `gorm:"autoUpdateTime:milli"`
}

type ScheduledSlot struct {
	From time.Time `gorm:"column:from"`
	To   time.Time `gorm:"column:to"`
}

type Location struct {
	Lat  *float64 `gorm:"column:lat"`
	Long *float64 `gorm:"column:long"`
}

type Task struct {
	ID          uuid.UUID  `gorm:"primaryKey"`
	ScheduleID  uuid.UUID  `gorm:"column:schedule_id"`
	Title       string     `gorm:"column:title"`
	Description string     `gorm:"column:description"`
	Status      string     `gorm:"column:status"` 
	Done        *bool      `gorm:"column:done"`
	Feedback    *string    `gorm:"column:feedback"`
	CreatedAt   time.Time  `gorm:"autoCreateTime:milli"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime:milli"`
}

type SearchResultSchedule struct {
	Data       *[]Schedule
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type IScheduleRepository interface {
	GetSchedules() (*[]Schedule, error)
	GetScheduleByID(id uuid.UUID) (*Schedule, error)
	GetTodaySchedules(userID uuid.UUID) (*[]Schedule, error)
	UpdateSchedule(id uuid.UUID, updates map[string]interface{}) (*Schedule, error)
	UpdateTask(taskID uuid.UUID, updates map[string]interface{}) (*Task, error)
	Create(newSchedule *Schedule) (*Schedule, error) 
	GetSchedulesByAssignedUserIDPaginated(assignedUserID uuid.UUID, filters domain.DataFilters) (*SearchResultSchedule, error)
}