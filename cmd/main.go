package main

import (
	"fmt"
	"log"

	"github.com/maSrifiyy/api"
	"github.com/maSrifiyy/db"
)

func main() {
	store, err := db.NewPostgreStore()
	if err != nil {
		log.Fatalf("Error creating database %v", err)
	}
	fmt.Printf("store mem:%+v\n", store)
	server := api.NewAPIServer(":3000", store)
	server.Run()
}
