package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maSrifiyy/db"
	"github.com/maSrifiyy/models"
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

	router.HandleFunc("/acct", makeHTTPHandleFunc(s.handleAcct))
	router.HandleFunc("/create", makeHTTPHandleFunc(s.handleGetAcct))
	log.Println("json http server running on port :3000")
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAcct(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAcct(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAcct(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAcct(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAcct(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, &models.User{})
}

func (s *APIServer) handleCreateAcct(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAcct(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type APIFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
			//handle eror here
		}
	}
}
