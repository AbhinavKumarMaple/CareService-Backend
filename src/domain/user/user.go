package user

import (
	"time"

	"caregiver/src/domain"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"primaryKey"`
	UserName       string    `gorm:"column:user_name;unique"`
	Email          string    `gorm:"unique"`
	FirstName      string    `gorm:"column:first_name"`
	LastName       string    `gorm:"column:last_name"`
	Status         bool      `gorm:"column:status"`
	HashPassword   string    `gorm:"column:hash_password"`
	Role           string    `gorm:"column:role"`
	ProfilePicture string    `gorm:"column:profile_picture"`
	Location       Location  `gorm:"embedded;embeddedPrefix:location_"`
	CreatedAt      time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime:milli"`
}

type Location struct {
	HouseNumber string  `json:"house_number"`
	Street      string  `json:"street"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Pincode     string  `json:"pincode"`
	Lat         float64 `json:"lat"`
	Long        float64 `json:"long"`
}

type SearchResultUser struct {
	Data       *[]User
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

type IUserService interface {
	GetAll() (*[]User, error)
	GetByID(id uuid.UUID) (*User, error)
	Create(newUser *User) (*User, error)
	Delete(id uuid.UUID) error
	Update(id uuid.UUID, userMap map[string]interface{}) (*User, error)
	SearchPaginated(filters domain.DataFilters) (*SearchResultUser, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
}

type IUserRepository interface {
	GetAll() (*[]User, error)
	Create(userDomain *User) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(id uuid.UUID, userMap map[string]interface{}) (*User, error)
	Delete(id uuid.UUID) error
	SearchPaginated(filters domain.DataFilters) (*SearchResultUser, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
}
