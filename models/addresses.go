package models

type Address struct {
	ID         int64  `json:"id" db:"id"`
	Address    string `json:"address" db:"address"`
	City       string `json:"city" db:"city"`
	PostalCode string `json:"postal_code" db:"postal_code"`
	Country    string `json:"country" db:"country"`
}
