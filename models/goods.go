package models

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
