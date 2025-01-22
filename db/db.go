package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/maSrifiyy/models"
)

type Storage interface {
	CreateUserAccount(*models.Sellers) (bool, error)
	UpdateUserAccount(*models.Sellers) (bool, error)
	GetUserAccountById(int) (*models.Sellers, error)
	AddItem(*models.Goods, *models.Sellers) (bool, error)
	UpdateItem(*models.Goods, *models.Sellers) (bool, error)
	SetItemMaxThreshold(*models.Goods, *models.Sellers) (bool, error)
	GetItemById(int) (*models.Goods, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgreStore() (*PostgresStore, error) {
	//created new sqlDB, then decide to port that to gormDB
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
		db: db,
	}, nil
}

func (s *PostgresStore) CreateUserAccount(user *models.Sellers) (bool, error) {
	return true, nil
}

func (s *PostgresStore) UpdateUserAccount(*models.Sellers) (bool, error) {
	return true, nil
}

func (s *PostgresStore) GetUserAccountById(id int) (*models.Sellers, error) {
	return nil, nil
}

func (s *PostgresStore) AddItem(*models.Goods, *models.Sellers) (bool, error) {
	return true, nil
}

func (s *PostgresStore) UpdateItem(*models.Goods, *models.Sellers) (bool, error) {
	return true, nil
}

func (s *PostgresStore) SetItemMaxThreshold(*models.Goods, *models.Sellers) (bool, error) {
	return true, nil
}

func (s *PostgresStore) GetItemById(int) (*models.Goods, error) {
	return nil, nil
}
