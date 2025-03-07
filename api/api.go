package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/maSrifiyy/business"
	"github.com/maSrifiyy/db"
	"github.com/maSrifiyy/models"
	"github.com/maSrifiyy/utils"
	"golang.org/x/crypto/bcrypt"
)

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIServer struct {
	listenAddr string
	Storage    db.Storage
}
type APIError struct {
	Error string
}

func NewAPIServer(listenAddr string, store db.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		Storage:    store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/acct", MakeHTTPHandleFunc(s.handleAcct))
	router.HandleFunc("/create", MakeHTTPHandleFunc(s.HandleCreateAcct))
	router.HandleFunc("/login", MakeHTTPHandleFunc(s.Login))
	router.HandleFunc("/add-item", MakeHTTPHandleFunc(s.handleAddItemToInventory))
	router.HandleFunc("/update-inventory", MakeHTTPHandleFunc(s.updateItem))
	router.HandleFunc("/get-items", MakeHTTPHandleFunc(s.GetAllItems))
	router.HandleFunc("/get-product/{productId}", MakeHTTPHandleFunc(s.GetItemByProductID))

	log.Println("json http server running on port :3000")
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAcct(w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAcct(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, &models.Sellers{})
}

func (s *APIServer) HandleCreateAcct(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("Method not allowed")
	}
	acctReq := new(models.CreateAccount)
	if err := json.NewDecoder(r.Body).Decode(acctReq); err != nil {

	}
	pass_hash, err := bcrypt.GenerateFromPassword([]byte(acctReq.Password), 8)
	if err != nil {
		log.Printf("err hashing password: %v", err)
		return err
	}
	acctReq.Password = string(pass_hash)
	newAcct, err := business.CreateSellerAccount(acctReq)
	return writeJSON(w, http.StatusOK, newAcct)
}

func (s *APIServer) Login(w http.ResponseWriter, r *http.Request) error {
	login := new(models.Login) //return login pointer
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(login)
	if err != nil {
		return err
	}
	seller, err := s.Storage.GetUserAccountByEmail(login.Email)
	if err != nil {
		return err
	}
	sellerIDd, err := s.Storage.GetUserAccountById(seller.ID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(seller.Password), []byte(login.Password))
	if err != nil {
		log.Println(err)
		return fmt.Errorf("incorrect password")
	}
	JWToken, tokenErr := utils.GenerateJWT(*sellerIDd, 2*time.Hour, utils.ISSUER, utils.SECRET)
	if tokenErr != nil {
		return fmt.Errorf("err generating jwt token")
	}
	res := map[string]interface{}{}
	res["jwtToken"] = JWToken
	res["user"] = sellerIDd
	return writeJSON(w, http.StatusOK, res)
}

func (s *APIServer) handleAddItemToInventory(w http.ResponseWriter, r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid authorization header format"})
	}

	tokenString := parts[1]
	claims, err := utils.DecodeJWT(tokenString, utils.SECRET)
	if err != nil {
		return writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	//extract user_id from claims
	user_id, ok := claims["user"]
	if !ok {
		return writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
	}
	fmt.Printf("%T", user_id)

	userID := user_id.(string)
	//fetch user account using user_id
	user, err := s.Storage.GetUserAccountByUserID(userID)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "user not found"})
	}

	item := new(models.CreateGood)
	if err := json.NewDecoder(r.Body).Decode(item); err != nil {
		return writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
	}
	defer r.Body.Close()

	newItem, err := business.AddNewItemToInventory(item, user.ID)
	if err != nil {
		log.Printf("Error creating new item: %v", err)
		return writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to retrieve added item"})
	}

	res := map[string]interface{}{}
	res["message"] = "Item added successfully"
	res["item"] = newItem
	return writeJSON(w, http.StatusOK, res)
}

func (s *APIServer) updateItem(w http.ResponseWriter, r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid authorization header format"})
	}

	tokenString := parts[1]
	claims, err := utils.DecodeJWT(tokenString, utils.SECRET)
	if err != nil {
		return writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	user_id, ok := claims["user"]
	if !ok {
		return writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
	}

	userID := user_id.(string)
	//fetch user account using d user_id
	user, err := s.Storage.GetUserAccountByUserID(userID)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "user not found"})
	}

	var payload struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
	}

	defer r.Body.Close()

	item, err := s.Storage.GetItemByProductID(payload.ProductID)
	if err != nil {
		return writeJSON(w, http.StatusNotFound, map[string]string{"error": "item not found"})
	}

	if item.CreatedBy != user.ID {
		return writeJSON(w, http.StatusForbidden, map[string]string{"error": "user not authorized to update this item"})
	}

	switch r.Method {
	case http.MethodPost:
		success, err := s.Storage.AddItemToInventory(payload.ProductID, payload.Quantity)
		if err != nil || !success {
			return writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

	case http.MethodDelete:
		success, err := s.Storage.RemoveItemFromInventory(payload.ProductID, payload.Quantity)
		if err != nil || !success {
			return writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

	default:
		return writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}

	return writeJSON(w, http.StatusOK, map[string]string{"message": "inventory updated successfully"})
}

// no auth for my reason
func (s *APIServer) GetAllItems(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("Method not allowed")
	}
	items, err := s.Storage.GetAllItem()
	if err != nil {
		return err
	}
	res := map[string]interface{}{}
	res["message"] = "All items retunred"
	res["items"] = items
	return writeJSON(w, http.StatusOK, res)
}

func (s *APIServer) GetItemByProductID(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("method not allowed")
	}
	vars := mux.Vars(r)
	productId := vars["productId"]
	item, err := s.Storage.GetItemByProductID(productId)
	if err != nil {
		log.Printf("err getting product: %v", err)
		return err
	}
	return writeJSON(w, http.StatusOK, item)
}

type APIFunc func(http.ResponseWriter, *http.Request) error

func MakeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
			//handle eror here
		}
	}
}
