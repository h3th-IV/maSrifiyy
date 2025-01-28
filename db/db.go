package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/maSrifiyy/models"
)

func (s *PostgresStore) CreateUserAccount(userID, firstName, lastName, email, password string) (bool, error) {
	result, err := s.DB.Exec(createSellerAcct, userID, firstName, lastName, email, password)
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

func (s *PostgresStore) UpdateUserAccount(firstName, lastName, email, password string, id int) (*models.Sellers, error) {
	result, err := s.DB.Exec(updateSellerAcct, firstName, lastName, email, password, id)
	if err != nil {
		log.Printf("err updating user: %v", err)
		return nil, err
	}
	row, err := result.RowsAffected()
	if row <= 0 {
		log.Printf("affected row: %v", row)
		return nil, err
	}
	if err != nil {
		log.Printf("Error updating user account: %v", err)
		return nil, err
	}
	seller, err := s.GetUserAccountById(id)
	if err != nil {
		log.Printf("err fetching user after update: %v", err)
		return nil, err
	}

	return seller, nil
}

func (s *PostgresStore) GetUserAccountById(id int) (*models.Sellers, error) {
	row := s.DB.QueryRow(getSellerAcctbyID, id)
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

func (s *PostgresStore) GetUserAccountByEmail(email string) (*models.Sellers, error) {
	row := s.DB.QueryRow(getSellerAcctbyEmail, email)
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

func (s *PostgresStore) GetUserAccountByUserID(userID string) (*models.Sellers, error) {
	row := s.DB.QueryRow(getSellerAcctbyUserId, userID)
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

func (s *PostgresStore) AddItem(productID, name string, quantity, maxThreshold, minThreshold, id int) (bool, error) {
	result, err := s.DB.Exec(
		createItemInInventory,
		productID,
		name,
		quantity,
		maxThreshold,
		minThreshold,
		id,
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
	row := s.DB.QueryRow(getItembyProductID, productID)

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
	rows, err := s.DB.Query(getAllItem)
	if err != nil {
		log.Printf("err: %v", err)
		return nil, err
	}
	for rows.Next() {

		err = rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.Quantity, &item.MaxThreshold, &item.MinThreshold, &item.CreatedBy)
		if err != nil {
			log.Printf("err: %v", err)
			return nil, err
		}
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
	result, err := s.DB.Exec(updateItemInInventory, item.Name, item.Quantity, item.MaxThreshold, item.ID, user.ID)
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
	result, err := s.DB.Exec(setMaxThreshold, item.MaxThreshold, item.ID, user.ID)
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
	row := s.DB.QueryRow(getItembyID, id)
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
	result, err := s.DB.Exec(incrementItem, quantity, productID)
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
	result, err := s.DB.Exec(decrementItem, quantity, productID)
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
	rows, err := s.DB.Query(getLowStockPorducts)
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
