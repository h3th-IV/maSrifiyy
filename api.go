package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type APIServer struct {
	listenAddr string
}
type APIError struct {
	Error string
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

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/acct", makeHTTPHandleFunc(s.handleAcct))
	router.HandleFunc("/get-acct/{id}", makeHTTPHandleFunc(s.handleGetAcct))
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
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	return writeJSON(w, http.StatusOK, &account{})
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
