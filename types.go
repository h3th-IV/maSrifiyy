package main

import "math/rand"

type account struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	BankNumber int64  `json:"bank_number"`
	Balance    int64  `json:"balance"`
}

func NewAccount(firstName, lastName string) *account {
	return &account{
		ID:         rand.Intn(10000),
		FirstName:  firstName,
		LastName:   lastName,
		BankNumber: int64(rand.Intn(1000000000)),
	}
}
