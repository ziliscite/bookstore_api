package models

import "time"

type PaymentMethod int

const (
	PayPal PaymentMethod = iota + 1
	Bank
	QRIS
)

func (pm PaymentMethod) String() string {
	return [...]string{"PayPal", "Bank", "QRIS"}[pm-1]
}

func (pm PaymentMethod) EnumIndex() int {
	return int(pm)
}

type Order struct {
	ID        int64 `json:"id" db:"id"`
	UserID    int   `json:"user_id" db:"user_id"`
	AddressID int   `json:"address_id" db:"address_id"`

	ProductPrice float64 `json:"product_price" db:"product_price"`
	TaxFee       float64 `json:"tax_fee" db:"tax_fee"`
	TotalPrice   float64 `json:"total_price" db:"total_price"`

	PaymentMethod   PaymentMethod `json:"payment_method" db:"payment_method"`
	PaymentResultId int           `json:"payment_result_id" db:"payment_result_id"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type OrderBook struct {
	ID       int64 `json:"id" db:"id"` // Should I just constrain both bookid and orderid together instead of using an id?
	Quantity int   `json:"quantity" db:"quantity"`

	BookID  int `json:"book_id" db:"book_id"`
	OrderID int `json:"order_id" db:"order_id"`
}

/* Will look something like this btw
{
    "id": 12,
    "user_id": 1,
    "address_id": 13,

    "product_price": 151.41 -- calculated from aggregating book price from book_id * quantity from the orderbook

    "tax_fee": 10.8 -- calculated off a percentage from product_price (or the aggregation result)

    "total_price": 162.21 -- total

    "books": [
        {
            "id": 41,
            "quantity": 2
            "product_id": 3
        },
        {
            "id": 42,
            "quantity": 1
            "product_id": 5
        },
        {
            "id": 43,
            "quantity": 8
            "product_id": 1
        },
    ]
}
*/
