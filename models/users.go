package models

import (
	"time"

	"github.com/google/uuid"
)

type Sellers struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAccount struct {
	FirstName string `json:"first_name"`
	LAstName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ItemUser struct {
	Goods
	FirstName string
	Email     string
}

func NewUser(firstName, lastName, email, password string) *Sellers {
	return &Sellers{
		UserID:    "usr" + uuid.NewString(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
}
