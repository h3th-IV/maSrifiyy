package models

import (
	"math/rand"
)

type Goods struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Quantity     int    `json:"quantity"`
	MaxThreshold int    `json:"max_threshold"`
	CreatedBy    int    `json:"created_by"`
}

func NewGood(name string, quantity int) *Goods {
	return &Goods{
		ID:   rand.Intn(1000000),
		Name: name,
	}
}
