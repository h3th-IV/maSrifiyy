package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/maSrifiyy/models"
	"github.com/maSrifiyy/utils"
)

type Storage interface {
	CreateSellersTable() error
	CreateGoodsTable() error
	CreateUserAccount(*models.Sellers) (bool, error)
	UpdateUserAccount(*models.Sellers) (bool, error)
	GetUserAccountById(int) (*models.Sellers, error)
	GetUserAccountByUserID(userID string) (*models.Sellers, error)
	GetUserAccountByEmail(email string) (*models.Sellers, error)
	AddItem(*models.Goods, *models.Sellers) (bool, error)
	UpdateItem(*models.Goods, *models.Sellers) (bool, error)
	SetItemMaxThreshold(*models.Goods, *models.Sellers) (bool, error)
	GetItemById(int) (*models.Goods, error)
	GetItemByProductID(productID string) (*models.Goods, error)
	AddItemToInventory(productID string, quantity int) (bool, error)
	RemoveItemFromInventory(productID string, quantity int) (bool, error)
	GetAllItem() ([]*models.Goods, error)
	SentThresholdNotification() error
	GetLowStockProducts() ([]*models.ItemUser, error)
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

func (s *PostgresStore) CreateSellersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS sellers (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := s.db.Exec(query)
	if err != nil {
		log.Printf("Error creating sellers table: %v", err)
		return err
	}
	log.Println("Sellers table created successfully.")
	return nil
}

func (s *PostgresStore) CreateGoodsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS goods (
		id SERIAL PRIMARY KEY,
		product_id VARCHAR(255) NOT NULL,
		name VARCHAR(100) NOT NULL,
		quantity INT NOT NULL,
		max_threshold INT NOT NULL,
		min_threshold INT NOT NULL,
		created_by INT REFERENCES sellers(id) ON DELETE CASCADE
	);`
	_, err := s.db.Exec(query)
	if err != nil {
		log.Printf("Error creating goods table: %v", err)
		return err
	}
	log.Println("Goods table created successfully.")
	return nil
}

func (s *PostgresStore) DropSellersTable() error {
	query := `DROP TABLE IF EXISTS sellers;`
	_, err := s.db.Exec(query)
	if err != nil {
		log.Printf("Error dropping sellers table: %v", err)
		return err
	}
	log.Println("Sellers table dropped successfully.")
	return nil
}

func (s *PostgresStore) DropGoodsTable() error {
	query := `DROP TABLE IF EXISTS goods;`
	_, err := s.db.Exec(query)
	if err != nil {
		log.Printf("Error dropping goods table: %v", err)
		return err
	}
	log.Println("Goods table dropped successfully.")
	return nil
}

func (s *PostgresStore) CreateUserAccount(user *models.Sellers) (bool, error) {
	query := `INSERT INTO sellers (user_id, first_name, last_name, email, password, created_at)
	VALUES ($1, $2, $3, $4, $5, NOW())`
	createUser := models.NewUser(user.FirstName, user.LastName, user.Email, user.Password)
	result, err := s.db.Exec(query, createUser.UserID, createUser.FirstName, createUser.LastName, createUser.Email, createUser.Password)
	if err != nil {
		log.Printf("Error creating user account: %v", err)
		if strings.Contains(err.Error(), "sellers_email_key") {
			return false, fmt.Errorf("user with email already exist")
		}
		return false, err
	}
	row, err := result.RowsAffected()
	if err != nil {
		log.Printf("err getting affected row: %v", err)
		return false, err
	}
	if row <= 0 {
		log.Printf("err checking row affected, rowId: %d", row)
		return false, err
	}
	return true, nil
}

func (s *PostgresStore) UpdateUserAccount(user *models.Sellers) (bool, error) {
	query := `UPDATE sellers SET first_name = $1, last_name = $2, email = $3, password = $4
	WHERE id = $5`
	result, err := s.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.ID)
	if err != nil {
		log.Printf("err updating user: %v", err)
		return false, err
	}
	row, err := result.RowsAffected()
	if row <= 0 {
		log.Printf("affected row: %v", row)
		return false, err
	}
	if err != nil {
		log.Printf("Error updating user account: %v", err)
		return false, err
	}
	return true, nil
}

func (s *PostgresStore) GetUserAccountByEmail(email string) (*models.Sellers, error) {
	query := `SELECT id, user_id, first_name, last_name, email, password, created_at FROM sellers WHERE email = $1`
	row := s.db.QueryRow(query, email)
	var user models.Sellers
	err := row.Scan(&user.ID, &user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Printf("Error fetching user account by email: %v", err)
		return nil, err
	}
	return &user, nil
}

func (s *PostgresStore) GetUserAccountById(id int) (*models.Sellers, error) {
	query := `SELECT id, user_id, first_name, last_name, email, password, created_at FROM sellers WHERE id = $1`
	row := s.db.QueryRow(query, id)
	var user models.Sellers
	err := row.Scan(&user.ID, &user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Printf("Error fetching user account by ID: %v", err)
		return nil, err
	}
	return &user, nil
}

func (s *PostgresStore) GetUserAccountByUserID(userID string) (*models.Sellers, error) {
	query := `SELECT id, user_id, first_name, last_name, email, password, created_at FROM sellers WHERE user_id = $1`
	row := s.db.QueryRow(query, userID)
	var user models.Sellers
	err := row.Scan(&user.ID, &user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Printf("Error fetching user account by user_id: %v", err)
		return nil, err
	}
	return &user, nil
}

func (s *PostgresStore) AddItem(item *models.Goods, user *models.Sellers) (bool, error) {
	query := `
	INSERT INTO goods (product_id, name, quantity, max_threshold, min_threshold, created_by)
	VALUES ($1, $2, $3, $4, $5, $6)`

	result, err := s.db.Exec(
		query,
		item.ProductID,
		item.Name,
		item.Quantity,
		item.MaxThreshold,
		item.MinThreshold,
		user.ID,
	)
	if err != nil {
		log.Printf("Error adding item: %v", err)
		return false, err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if row <= 0 {
		log.Printf("No rows were affected, rowId: %d", row)
		return false, err
	}
	return true, nil
}

func (s *PostgresStore) GetItemByProductID(productID string) (*models.Goods, error) {
	query := `SELECT id, product_id, name, quantity, max_threshold, min_threshold, created_by
	          FROM goods WHERE product_id = $1`
	row := s.db.QueryRow(query, productID)

	var item models.Goods
	err := row.Scan(
		&item.ID,
		&item.ProductID,
		&item.Name,
		&item.Quantity,
		&item.MaxThreshold,
		&item.MinThreshold,
		&item.CreatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("item not found")
		}
		log.Printf("Error retrieving item: %v", err)
		return nil, err
	}
	return &item, nil
}

func (s *PostgresStore) GetAllItem() ([]*models.Goods, error) {
	var (
		items []*models.Goods
	)
	item := new(models.Goods)
	query := `SELECT * from goods`
	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("err: %v", err)
		return nil, err
	}
	fmt.Println("a ...any")
	for rows.Next() {
		fmt.Println("a ...anies")

		err = rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.Quantity, &item.MaxThreshold, &item.MinThreshold, &item.CreatedBy)
		if err != nil {
			log.Printf("err: %v", err)
			return nil, err
		}
		fmt.Println("a ...aniex")
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		log.Printf("err: %v", err)
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no item found in the inventory")
	}
	return items, nil
}
func (s *PostgresStore) UpdateItem(item *models.Goods, user *models.Sellers) (bool, error) {
	query := `UPDATE goods SET name = $1, quantity = $2, max_threshold = $3
	WHERE id = $4 AND created_by = $5`
	result, err := s.db.Exec(query, item.Name, item.Quantity, item.MaxThreshold, item.ID, user.ID)
	if err != nil {
		log.Printf("Error updating item: %v", err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		return false, err
	}

	if rowsAffected == 0 {
		log.Println("No rows were affected during the update.")
		return false, errors.New("no rows updated")
	}

	return true, nil
}

func (s *PostgresStore) SetItemMaxThreshold(item *models.Goods, user *models.Sellers) (bool, error) {
	query := `UPDATE goods SET max_threshold = $1 WHERE product_id = $2 AND created_by = $3`
	result, err := s.db.Exec(query, item.MaxThreshold, item.ID, user.ID)
	if err != nil {
		log.Printf("Error setting item max threshold: %v", err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		return false, err
	}

	if rowsAffected == 0 {
		log.Println("No rows were affected during the update.")
		return false, errors.New("no rows updated")
	}
	return true, nil
}

func (s *PostgresStore) GetItemById(id int) (*models.Goods, error) {
	query := `SELECT id, product_id, name, quantity, max_threshold, min_threshold, created_by FROM goods WHERE id = $1`
	row := s.db.QueryRow(query, id)
	var item models.Goods
	err := row.Scan(&item.ID, &item.Name, &item.Quantity, &item.MaxThreshold, &item.CreatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("item not found")
		}
		log.Printf("Error fetching item by ID: %v", err)
		return nil, err
	}
	return &item, nil
}

func (s *PostgresStore) AddItemToInventory(productID string, quantity int) (bool, error) {
	query := `UPDATE goods SET quantity = quantity + $1 WHERE product_id = $2 AND quantity + $1 <= max_threshold`
	result, err := s.db.Exec(query, quantity, productID)
	if err != nil {
		log.Printf("Error adding quantity: %v", err)
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, fmt.Errorf("update failed: quantity exceeds max threshold")
	}

	return true, nil
}

func (s *PostgresStore) RemoveItemFromInventory(productID string, quantity int) (bool, error) {
	query := `UPDATE goods SET quantity = quantity - $1 WHERE product_id = $2 AND quantity - $1 >= min_threshold`
	result, err := s.db.Exec(query, quantity, productID)
	if err != nil {
		log.Printf("Error removing quantity: %v", err)
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, fmt.Errorf("update failed: quantity below min threshold")
	}

	return true, nil
}

func (s *PostgresStore) GetLowStockProducts() ([]*models.ItemUser, error) {
	query := `
	SELECT g.id, g.product_id, g.name, g.quantity, g.max_threshold, g.min_threshold, g.created_by, s.first_name, s.email FROM goods g INNER JOIN sellers s ON g.created_by = s.id WHERE g.quantity <= g.min_threshold OR g.quantity <= g.min_threshold + 0.1 * g.min_threshold;`

	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("Error querying low stock products: %v", err)
		return nil, err
	}
	defer rows.Close()

	lowStockItems := []*models.ItemUser{}
	for rows.Next() {
		item := new(models.ItemUser)
		err := rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.Quantity, &item.MaxThreshold, &item.MinThreshold, &item.CreatedBy, &item.FirstName, &item.Email)
		if err != nil {
			log.Printf("Error scanning product: %v", err)
			continue
		}
		lowStockItems = append(lowStockItems, item)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, err
	}

	return lowStockItems, nil
}

func (s *PostgresStore) SentThresholdNotification() error {
	thresholdProducts, err := s.GetLowStockProducts()
	if err != nil {
		log.Printf("Error fetching low stock products: %v", err)
		return err
	}

	if len(thresholdProducts) == 0 {
		log.Println("No low-stock products found. No notifications sent.")
		return nil
	}

	for _, tped := range thresholdProducts {
		log.Printf("Processing low stock alert: ProductID=%s, Name=%s, SellerEmail=%s", tped.ProductID, tped.Name, tped.Email)
		if err := utils.SendEmail(tped.Email, tped.FirstName, tped.Name, tped.ProductID); err != nil {
			log.Printf("Error sending email to %s for product %s: %v", tped.Email, tped.ProductID, err)
			continue
		}
		log.Printf("Email successfully sent to %s for product %s", tped.Email, tped.ProductID)
	}

	log.Println("All low stock notifications processed successfully.")
	return nil
}
