package repository

import (
	"time"

	"github.com/lib/pq"
)

type ParamCreateOrderResult struct {
	OrderID     string         `db:"order_id"`
	AccountID   string         `db:"account_id"`
	ProductsIDS pq.StringArray `db:"products_ids"`
	CreatedAT   time.Time      `db:"created_at"`
}
