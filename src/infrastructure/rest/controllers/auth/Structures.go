package auth

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"Email" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

type AccessTokenRequest struct {
	RefreshToken string `json:"RefreshToken" binding:"required"`
}

type UserData struct {
	UserName  string    `json:"UserName"`
	Email     string    `json:"Email"`
	FirstName string    `json:"FirstName"`
	LastName  string    `json:"LastName"`
	Status    bool      `json:"Status"`
	ID        uuid.UUID `json:"ID"`
}

type SecurityData struct {
	JWTAccessToken            string    `json:"JWTAccessToken"`
	JWTRefreshToken           string    `json:"JWTRefreshToken"`
	ExpirationAccessDateTime  time.Time `json:"ExpirationAccessDateTime"`
	ExpirationRefreshDateTime time.Time `json:"ExpirationRefreshDateTime"`
}

type LoginResponse struct {
	Data     UserData     `json:"Data"`
	Security SecurityData `json:"Security"`
}
