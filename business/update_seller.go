package business

import (
	"github.com/maSrifiyy/models"
)

func updateSellerAccount(seller *models.Sellers) (*models.Sellers, error) {
	updatedseller, err := pgStore.UpdateUserAccount(seller.FirstName, seller.LastName, seller.Email, seller.Password, seller.ID)
	if err != nil {
		return nil, err
	}
	return updatedseller, nil
}
