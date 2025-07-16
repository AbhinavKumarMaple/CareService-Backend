package schedule

import (
	"time"

	"github.com/google/uuid"
)


type CreateScheduleRequest struct {
	ClientUserID  uuid.UUID     `json:"ClientUserID" binding:"required"`
	AssignedUserID uuid.UUID    `json:"AssignedUserID" binding:"required"`
	ServiceName   string        `json:"ServiceName" binding:"required"`   
	ScheduledSlot ScheduledSlot `json:"ScheduledSlot" binding:"required"`
	Tasks         []TaskRequest `json:"Tasks" binding:"required,min=1,dive"`
}

type TaskRequest struct {
	Title       string `json:"Title" binding:"required"`
	Description string `json:"Description"`
}

type ScheduledSlot struct {
	From time.Time `json:"From" binding:"required"`
	To   time.Time `json:"To" binding:"required"`
}

type Location struct {
	Lat  *float64 `json:"lat" binding:"required"`
	Long *float64 `json:"long" binding:"required"`
}

type Task struct {
	ID          uuid.UUID `json:"ID"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Status      string    `json:"Status"`
	Done        *bool     `json:"Done"`
	Feedback    *string   `json:"Feedback"`
}

type ClientInfo struct {
	ID            uuid.UUID `json:"ID"`
	UserName      string    `json:"UserName"`
	Email         string    `json:"Email"`
	FirstName     string    `json:"FirstName"`
	LastName      string    `json:"LastName"`
	ProfilePicture string    `json:"ProfilePicture"`
	Location      ClientLocation `json:"Location"`
}

type ClientLocation struct {
	HouseNumber string  `json:"house_number"`
	Street      string  `json:"street"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Pincode     string  `json:"pincode"`
	Lat         float64 `json:"lat"`
	Long        float64 `json:"long"`
}

type ScheduleResponse struct {
	ID               uuid.UUID      `json:"ID"`
	ClientUserID     uuid.UUID      `json:"ClientUserID"`
	ClientInfo       *ClientInfo    `json:"ClientInfo"`       
	AssignedUserID   uuid.UUID      `json:"AssignedUserID"`  
	ServiceName      string         `json:"ServiceName"`    
	ScheduledSlot    ScheduledSlot  `json:"ScheduledSlot"`
	VisitStatus      string         `json:"VisitStatus"`
	CheckinTime      *time.Time     `json:"CheckinTime"`
	CheckoutTime     *time.Time     `json:"CheckoutTime"`
	CheckinLocation  Location       `json:"CheckinLocation"`
	CheckoutLocation Location       `json:"CheckoutLocation"`
	Tasks            []Task         `json:"Tasks"`
	ServiceNote      *string        `json:"ServiceNote"`
}

type StartScheduleRequest struct {
	Timestamp time.Time `json:"timestamp" binding:"required"`
	Location  Location  `json:"location" binding:"required"`
}

type StartScheduleResponse struct {
	Message         string     `json:"Message"`
	CheckinTime     *time.Time `json:"checkin_time"`
	CheckinLocation *Location  `json:"checkin_location"`
}

type EndScheduleTaskRequest struct {
	ID          uuid.UUID `json:"ID" binding:"required"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Status      string    `json:"Status" binding:"required"`
	Done        *bool     `json:"Done" binding:"required"`
	Feedback    *string   `json:"Feedback"`
}

type EndScheduleRequest struct {
	Timestamp    time.Time `json:"timestamp" binding:"required"`
	Location     Location  `json:"location" binding:"required"`
	Tasks        []EndScheduleTaskRequest `json:"tasks"` 
}

type EndScheduleResponse struct {
	Message          string     `json:"Message"`
	CheckoutTime     *time.Time `json:"checkout_time"`
	CheckoutLocation *Location  `json:"checkout_location"`
	ServiceNote      *string    `json:"service_note"`
}

type UpdateTaskRequest struct {
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Status      string    `json:"Status" binding:"required"`
	Done        *bool     `json:"Done" binding:"required"`
	Feedback    *string   `json:"Feedback"`
}

type UpdateTaskResponse struct {
	Message string     `json:"Message"`
	Task    Task       `json:"Task"`
}

type UpdateScheduleRequest struct {
	ClientUserID     uuid.UUID     `json:"ClientUserID"`
	AssignedUserID   uuid.UUID     `json:"AssignedUserID"`
	ServiceName      string        `json:"ServiceName"`
	ScheduledSlot    *ScheduledSlot `json:"ScheduledSlot"`
	VisitStatus      string        `json:"VisitStatus"`
}

type UpdateScheduleResponse struct {
	Message  string           `json:"Message"`
	Schedule *ScheduleResponse `json:"Schedule"`
}