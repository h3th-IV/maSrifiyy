package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maSrifiyy/db"
	"github.com/maSrifiyy/models"
	"github.com/stretchr/testify/require"
)

func TestCreateSellerAccount(t *testing.T) {
	//create mock
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error intializing mock db: %v", err)
	}
	defer mockDB.Close()

	store := &db.PostgresStore{DB: mockDB}

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
				mock.ExpectExec("INSERT INTO sellers").
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
				mock.ExpectExec("INSERT INTO sellers").
					WithArgs(
						sqlmock.AnyArg(), //"usr-21AxH90-360e460-00002",
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

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}
