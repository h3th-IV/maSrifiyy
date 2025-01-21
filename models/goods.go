package models

import (
	"math/rand"
)

type Good struct {
	ID           int    `gorm:"primaryKey;AutoIncrement;unique" json:"id"`
	Name         string `gorm:"size:255;not null" json:"name"`
	Quantity     int    `gorm:"not null" json:"quantity"`
	MaxThreshold int    `gorm:"not null" json:"max_threshold"`
	CreatedBy    int    `gorm:"nut null" json:"created_by"`
}

func NewGood(name string, quantity int) *Good {
	return &Good{
		ID:   rand.Intn(1000000),
		Name: name,
	}
}
