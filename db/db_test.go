package db_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maSrifiyy/db"
	"github.com/maSrifiyy/models"
	"github.com/stretchr/testify/require"
)

var (
	mockDB sqlmock.Sqlmock
	store  *db.PostgresStore
)

// setupTestDB initializes the mock database and store.
func setupTestDB(t *testing.T) {
	var dbMock *sql.DB
	var err error

	dbMock, mockDB, err = sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock db: %v", err)
	}

	store = &db.PostgresStore{DB: dbMock}
}

func cleanupTestDB() {
	if store != nil {
		_ = store.DB.Close()
	}
}

// TestMain runs before any tests are executed.
func TestMain(m *testing.M) {
	// Run setup before all tests
	setupTestDB(nil)
	defer cleanupTestDB()

	m.Run()
}

func TestCreateSellerAccount(t *testing.T) {
	testCases := []struct {
		desc           string
		user           *models.Sellers
		mockSetup      func()
		expectedResult bool
		expectedErr    error
	}{
		{
			desc: "User account created successfully",
			user: &models.Sellers{
				UserID:    "usr-10sdr90-190d393-00001",
				FirstName: "Thread",
				LastName:  "Miller",
				Email:     "threadMiller@wool.com",
				Password:  "!pAssw0rd",
			},
			mockSetup: func() {
				mockDB.ExpectExec("INSERT INTO sellers").
					WithArgs(
						sqlmock.AnyArg(),
						"Thread",
						"Miller",
						"threadMiller@wool.com",
						"!pAssw0rd",
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedResult: true,
			expectedErr:    nil,
		},
		{
			desc: "Duplicate Email Error",
			user: &models.Sellers{
				UserID:    "usr-21AxH90-360e460-00002",
				FirstName: "Sew",
				LastName:  "Wearer",
				Email:     "sewWearer@wool.com",
				Password:  "securepassword",
			},
			mockSetup: func() {
				mockDB.ExpectExec("INSERT INTO sellers").
					WithArgs(
						sqlmock.AnyArg(),
						"Sew",
						"Wearer",
						"sewWearer@wool.com",
						"securepassword",
					).
					WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"sellers_email_key\""))
			},
			expectedResult: false,
			expectedErr:    errors.New("user with email already exist"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.mockSetup()

			res, err := store.CreateUserAccount(tc.user)

			require.Equal(t, tc.expectedResult, res)
			require.Equal(t, tc.expectedErr, err)

			err = mockDB.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestAddItem(t *testing.T) {
	t.Run("Add Item Success", func(t *testing.T) {
		mockDB.ExpectExec("INSERT INTO goods").
			WithArgs("prd-ew783rrow", "hp Omen 16pro", 10, 100, 10, 3).
			WillReturnResult(sqlmock.NewResult(1, 1))

		res, err := store.AddItem(&models.Goods{
			ProductID:    "prd-ew783rrow",
			Name:         "hp Omen 16pro",
			Quantity:     10,
			MaxThreshold: 100,
			MinThreshold: 10,
		}, &models.Sellers{
			ID: 3,
		})

		require.True(t, res)
		require.NoError(t, err)

		err = mockDB.ExpectationsWereMet()
		require.NoError(t, err)
	})
}

// this try to fetch item from the db good/user
func TestGetItemFromDB(t *testing.T) {
	testCases := []struct {
		desc           string
		UserID         string
		mockSetup      func()
		expectedResult *models.Sellers
		expectedErr    error
	}{
		{
			desc:   "user fetched successfully",
			UserID: "usr-12ee34c-e8hf9023-29092h2e",
			mockSetup: func() {
				createdAt := time.Date(2025, 1, 25, 10, 0, 0, 0, time.UTC)
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "first_name", "last_name", "email", "password", "created_at"}).AddRow(1, "usr-1903ds0390-290de02e00-028e0020", "Thread", "Miller", "threadmiller@clother.com", "epwe9032jiwdj0i2je10e18e01ewjidq30eiwdq0212`w1e2wij03d23", createdAt)
				mockDB.ExpectQuery("SELECT id, user_id, first_name, last_name, email, password, created_at FROM sellers WHERE user_id = \\$1").WithArgs("usr-12ee34c-e8hf9023-29092h2e").WillReturnRows(rows)
			},
			expectedResult: &models.Sellers{
				CreatedAt: time.Now(),
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			//call mock function
			tc.mockSetup()

			seller, err := store.GetUserAccountByUserID(tc.UserID)
			require.IsType(t, tc.expectedResult, seller)
			require.NoError(t, err)
		})
	}
}

// db update ops test
func TestUpdateUser(t *testing.T) {
	t.Run("Update db op", func(t *testing.T) {
		mockDB.ExpectExec("UPDATE goods SET").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(0, 1))

		res, err := store.UpdateItem(&models.Goods{ID: 2}, &models.Sellers{ID: 2})

		require.True(t, res)
		require.NoError(t, err)
	})
}
