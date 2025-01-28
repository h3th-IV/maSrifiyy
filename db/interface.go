package db

import (
	"database/sql"
	"log"

	"github.com/maSrifiyy/models"
)

type Storage interface {
	CreateSellersTable() error
	CreateGoodsTable() error
	CreateUserAccount(userID, firstName, lastName, email, password string) (bool, error)
	UpdateUserAccount(firstName, lastName, email, password string, id int) (*models.Sellers, error)
	GetUserAccountById(int) (*models.Sellers, error)
	GetUserAccountByUserID(userID string) (*models.Sellers, error)
	GetUserAccountByEmail(email string) (*models.Sellers, error)
	AddItem(productID, name string, quantity, maxThreshold, minThreshold, id int) (bool, error)
	UpdateItem(*models.Goods, *models.Sellers) (bool, error)
	SetItemMaxThreshold(*models.Goods, *models.Sellers) (bool, error)
	GetItemById(int) (*models.Goods, error)
	GetItemByProductID(productID string) (*models.Goods, error)
	AddItemToInventory(productID string, quantity int) (bool, error)
	RemoveItemFromInventory(productID string, quantity int) (bool, error)
	GetAllItem() ([]*models.Goods, error)
	GetLowStockProducts() ([]*models.ItemUser, error)
}

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgreStore() (*PostgresStore, error) {

	connStr := "user=heth dbname=masrifiyy password=yourpassword host=172.17.0.2 port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("An error occurred when connecting to postgres db: %v", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Printf("Unable to test database connection: %v", err)
		return nil, err
	}
	return &PostgresStore{
		DB: db,
	}, nil
}

func (s *PostgresStore) CreateSellersTable() error {
	_, err := s.DB.Exec(createSellerTableStmt)
	if err != nil {
		log.Printf("Error creating sellers table: %v", err)
		return err
	}
	log.Println("Sellers table created successfully.")
	return nil
}

func (s *PostgresStore) CreateGoodsTable() error {
	_, err := s.DB.Exec(createGoodsTableStmt)
	if err != nil {
		log.Printf("Error creating goods table: %v", err)
		return err
	}
	log.Println("Goods table created successfully.")
	return nil
}

func (s *PostgresStore) DropSellersTable() error {
	_, err := s.DB.Exec(dropSellerTableStmt)
	if err != nil {
		log.Printf("Error dropping sellers table: %v", err)
		return err
	}
	log.Println("Sellers table dropped successfully.")
	return nil
}

func (s *PostgresStore) DropGoodsTable() error {
	_, err := s.DB.Exec(dropGoodsTableStmt)
	if err != nil {
		log.Printf("Error dropping goods table: %v", err)
		return err
	}
	log.Println("Goods table dropped successfully.")
	return nil
}
