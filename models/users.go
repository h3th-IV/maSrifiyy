package models

import (
	"time"
)

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement;unique" json:"id"`
	FirstName string    `gorm:"size:256;not null" json:"first_name"`
	LastName  string    `gorm:"size:256;not null" json:"last_name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"size:256;not null" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func NewUser(firstName, lastName, email, password string) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
}
