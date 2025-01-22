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
	query := `INSERT INTO sellers (first_name, last_name, email, password, created_at) 
	VALUES ($1, $2, $3, $4, NOW())`
	_, err := s.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		log.Printf("Error creating user account: %v", err)
		return false, err
	}
	return true, nil
}

func (s *PostgresStore) UpdateUserAccount(user *models.Sellers) (bool, error) {
	query := `UPDATE sellers SET first_name = $1, last_name = $2, email = $3, password = $4 
	WHERE id = $5`
	_, err := s.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.ID)
	if err != nil {
		log.Printf("Error updating user account: %v", err)
		return false, err
	}
	return true, nil
}

func (s *PostgresStore) GetUserAccountById(id int) (*models.Sellers, error) {
	query := `SELECT id, first_name, last_name, email, password, created_at FROM sellers WHERE id = $1`
	row := s.db.QueryRow(query, id)
	var user models.Sellers
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
		}
		log.Printf("Error fetching user account by ID: %v", err)
		return nil, err
	}
	return &user, nil
}

func (s *PostgresStore) AddItem(item *models.Goods, user *models.Sellers) (bool, error) {
	query := `INSERT INTO goods (name, quantity, max_threshold, created_by) 
	VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, item.Name, item.Quantity, item.MaxThreshold, user.ID)
	if err != nil {
		log.Printf("Error adding item: %v", err)
		return false, err
	}
	return true, nil
}

func (s *PostgresStore) UpdateItem(item *models.Goods, user *models.Sellers) (bool, error) {
	query := `UPDATE goods SET name = $1, quantity = $2, max_threshold = $3 
	WHERE id = $4 AND created_by = $5`
	_, err := s.db.Exec(query, item.Name, item.Quantity, item.MaxThreshold, item.ID, user.ID)
	if err != nil {
		log.Printf("Error updating item: %v", err)
		return false, err
	}
	return true, nil
}

func (s *PostgresStore) SetItemMaxThreshold(item *models.Goods, user *models.Sellers) (bool, error) {
	query := `UPDATE goods SET max_threshold = $1 WHERE id = $2 AND created_by = $3`
	_, err := s.db.Exec(query, item.MaxThreshold, item.ID, user.ID)
	if err != nil {
		log.Printf("Error setting item max threshold: %v", err)
		return false, err
	}
	return true, nil
}

func (s *PostgresStore) GetItemById(id int) (*models.Goods, error) {
	query := `SELECT id, name, quantity, max_threshold, created_by FROM goods WHERE id = $1`
	row := s.db.QueryRow(query, id)
	var item models.Goods
	err := row.Scan(&item.ID, &item.Name, &item.Quantity, &item.MaxThreshold, &item.CreatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
		}
		log.Printf("Error fetching item by ID: %v", err)
		return nil, err
	}
	return &item, nil
}
