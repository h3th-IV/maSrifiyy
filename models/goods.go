package models

import (
	"github.com/google/uuid"
)

type Goods struct {
	ID           int    `json:"id"`
	ProductID    string `json:"product_id"`
	Name         string `json:"name"`
	Quantity     int    `json:"quantity"`
	MaxThreshold int    `json:"max_threshold"`
	MinThreshold int    `json:"min_threshold"`
	CreatedBy    int    `json:"created_by"`
}

type CreateGood struct {
	Name         string `json:"name"`
	Quantity     int    `json:"quantity"`
	MaxThreshold int    `json:"max_threshold"`
}

func NewGood(name string, quantity, max_threshold int) *Goods {
	minThreshold := max_threshold / 10 //10% of max
	return &Goods{
		ProductID:    "prd" + uuid.NewString(),
		Name:         name,
		Quantity:     quantity,
		MaxThreshold: max_threshold,
		MinThreshold: minThreshold,
	}
}
