package business

import (
	"log"

	"github.com/google/uuid"
	"github.com/maSrifiyy/db"
	"github.com/maSrifiyy/models"
)

func NewUser(firstName, lastName, email, password string) *models.Sellers {
	return &models.Sellers{
		UserID:    "usr" + uuid.NewString(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}
}

var (
	pgStore, _ = db.NewPostgreStore()
)

func CreateSellerAccount(seller *models.CreateAccount) (*models.Sellers, error) {
	new_seller := NewUser(seller.FirstName, seller.LastName, seller.Email, seller.Password)
	success, err := pgStore.CreateUserAccount(new_seller.UserID, new_seller.FirstName, new_seller.LastName, new_seller.Email, new_seller.Password)
	if err != nil {
		return nil, err
	}
	if !success {
		log.Println("creating user failed without an error")
		return nil, err
	}
	return new_seller, nil
}
