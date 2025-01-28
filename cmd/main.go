package main

import (
	"fmt"
	"log"

	"github.com/maSrifiyy/api"
	"github.com/maSrifiyy/business"
	"github.com/maSrifiyy/db"
	"github.com/robfig/cron/v3"
)

func main() {
	store, err := db.NewPostgreStore()
	if err != nil {
		log.Fatalf("Error creating database %v", err)
	}
	croner := cron.New()
	croner.AddFunc("0 8 * * *", func() {
		err := business.SendThresholdNotification()
		if err != nil {
			log.Printf("Error running threshold notification: %v", err)
		}
	})
	fmt.Printf("store mem: %+v\n", store)

	//init db tables
	store.CreateSellersTable()
	store.CreateGoodsTable()
	//start server
	server := api.NewAPIServer(":3000", store)
	log.Println("Cron job scheduler started.")
	go croner.Start()
	server.Run()
}
