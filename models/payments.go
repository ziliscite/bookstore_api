package models

import "time"

type PaymentStatus int

const (
	Pending PaymentStatus = iota + 1
	Completed
	Canceled
)

func (ps PaymentStatus) String() string {
	return [...]string{"pending", "completed", "canceled"}[ps-1]
}

func (ps PaymentStatus) EnumIndex() int {
	return int(ps)
}

type PaymentResult struct {
	ID           int64         `json:"id"`
	Status       PaymentStatus `json:"status" db:"status" enums:"pending,completed,canceled" default:"pending"`
	EmailAddress string        `json:"email_address" db:"email_address"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
}
