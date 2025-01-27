package api_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/maSrifiyy/api"
	"github.com/maSrifiyy/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// mock postgres store
type MockStorage struct {
	mock.Mock
}

var (
	mockStorage = new(MockStorage)
	server      = api.NewAPIServer(":3000", mockStorage)
	testRouter  = mux.NewRouter()
)

func TestHandleCreateAcct(t *testing.T) {
	testRouter.HandleFunc("/create", api.MakeHTTPHandleFunc(server.HandleCreateAcct))

	t.Run("Seller Account creation", func(t *testing.T) {

		reqBody := models.NewUser("Thread", "Miller", "threadmiller@wooler.com", "!p@ssw0rd")
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 8)
		reqBody.Password = string(hashedPassword)
		//mock
		mockStorage.On("CreateUserAccount", mock.Anything).Return(true, nil)
		respBytes, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatal(err)
		}

		//mock request call
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(respBytes))
		req.Header.Set("Content-Type", "application/json")

		//fake response writer
		recorder := httptest.NewRecorder()
		err = server.HandleCreateAcct(recorder, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, recorder.Code)
	})
}
func TestGetItemByProductId(t *testing.T) {
	testRouter.HandleFunc("/get-product/{productId}", api.MakeHTTPHandleFunc(server.GetItemByProductID))
	t.Run("Get product by Product ID", func(t *testing.T) {
		//request param
		reqParam := mock.Anything
		mockStorage.On("GetItemByProductID", reqParam).Return(models.Goods{
			ProductID: "12345",
			Name:      "Sample Product",
			Quantity:  10,
		}, nil) //mock dbcall
		request := httptest.NewRequest(http.MethodGet, "/get-product/prddd00aa67-ba85-4f58-971a-6fb4b8a10c57", nil)
		request.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		err := server.HandleCreateAcct(recorder, request)
		log.Printf("%v", err)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, recorder.Code)
	})
}

func (m *MockStorage) CreateUserAccount(user *models.Sellers) (bool, error) {
	args := m.Called(user)
	log.Println(args...)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) CreateSellersTable() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) CreateGoodsTable() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) UpdateUserAccount(user *models.Sellers) (bool, error) {
	args := m.Called(user)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) GetUserAccountById(id int) (*models.Sellers, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Sellers), args.Error(1)
}

func (m *MockStorage) GetUserAccountByUserID(userID string) (*models.Sellers, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.Sellers), args.Error(1)
}

func (m *MockStorage) GetUserAccountByEmail(email string) (*models.Sellers, error) {
	args := m.Called(email)
	return args.Get(0).(*models.Sellers), args.Error(1)
}

func (m *MockStorage) AddItem(item *models.Goods, seller *models.Sellers) (bool, error) {
	args := m.Called(item, seller)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) UpdateItem(item *models.Goods, seller *models.Sellers) (bool, error) {
	args := m.Called(item, seller)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) SetItemMaxThreshold(item *models.Goods, seller *models.Sellers) (bool, error) {
	args := m.Called(item, seller)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) GetItemById(id int) (*models.Goods, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Goods), args.Error(1)
}

func (m *MockStorage) GetItemByProductID(productID string) (*models.Goods, error) {
	args := m.Called(productID)
	return args.Get(0).(*models.Goods), args.Error(1)
}

func (m *MockStorage) AddItemToInventory(productID string, quantity int) (bool, error) {
	args := m.Called(productID, quantity)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) RemoveItemFromInventory(productID string, quantity int) (bool, error) {
	args := m.Called(productID, quantity)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) GetAllItem() ([]*models.Goods, error) {
	args := m.Called()
	return args.Get(0).([]*models.Goods), args.Error(1)
}

func (m *MockStorage) SentThresholdNotification() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) GetLowStockProducts() ([]*models.ItemUser, error) {
	args := m.Called()
	return args.Get(0).([]*models.ItemUser), args.Error(1)
}
