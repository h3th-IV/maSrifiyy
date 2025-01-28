package business

import (
	"log"

	"github.com/google/uuid"
	"github.com/maSrifiyy/models"
)

func NewGood(name string, quantity, max_threshold int) *models.Goods {
	minThreshold := max_threshold / 10 //10% of max
	return &models.Goods{
		ProductID:    "prd" + uuid.NewString(),
		Name:         name,
		Quantity:     quantity,
		MaxThreshold: max_threshold,
		MinThreshold: minThreshold,
	}
}

func AddNewItemToInventory(good *models.CreateGood, userID int) (*models.Goods, error) {
	newGood := NewGood(good.Name, good.Quantity, good.MaxThreshold)
	success, err := pgStore.AddItem(newGood.ProductID, newGood.Name, newGood.Quantity, newGood.MaxThreshold, newGood.MinThreshold, userID)
	if err != nil {
		return nil, err
	}
	if !success {
		log.Printf("err creating new product: %v", err)
	}
	return newGood, nil
}
