package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/maSrifiyy/models"
)

type Storage interface {
	CreateUserAccount(*models.User) (bool, error)
	UpdateUserAccount(*models.User) (bool, error)
	GetUserAccountById(int) (*models.User, error)
	AddItem(*models.Good, *models.User) (bool, error)
	UpdateItem(*models.Good, *models.User) (bool, error)
	SetItemMaxThreshold(*models.Good, *models.User) (bool, error)
	GetItemById(int) (*models.Good, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgreStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=h3th sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("An error occured when connecting to postgres db: %v", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Printf("UNable to test database connection: %v", err)
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateUserAccount(*models.User) (bool, error) {
	return true, nil
}

func (s *PostgresStore) UpdateUserAccount(*models.User) (bool, error) {
	return true, nil
}

func (s *PostgresStore) GetUserAccountById(id int) (*models.User, error) {
	return nil, nil
}

func (s *PostgresStore) AddItem(*models.Good, *models.User) (bool, error) {
	return true, nil
}

func (s *PostgresStore) UpdateItem(*models.Good, *models.User) (bool, error) {
	return true, nil
}

func (s *PostgresStore) SetItemMaxThreshold(*models.Good, *models.User) (bool, error) {
	return true, nil
}

func (s *PostgresStore) GetItemById(int) (*models.Good, error) {
	return nil, nil
}
