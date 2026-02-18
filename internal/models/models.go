package models

import "time"

type ParamCreateOrder struct {
	OrderID     string
	AccountID   string
	ProductsIDS []string
	CreatedAT   time.Time
}

type ParamCreateOrderResult struct {
	OrderID     string    `db:"order_id"`
	AccountID   string    `db:"account_id"`
	ProductsIDS []string  `db:"products_ids"`
	CreatedAT   time.Time `db:"created_at"`
}

type ParamFindOrder struct {
	OrderID string
}

type ParamFindOrderResult struct {
	OrderID     string    `db:"order_id"`
	AccountID   string    `db:"account_id"`
	ProductsIDS []string  `db:"products_ids"`
	CreatedAT   time.Time `db:"created_at"`
}

type ParamFindOrderByAccount struct {
	AccountID string
}

type ParamFindOrderByAccountResult struct {
	OrderID     string    `db:"order_id"`
	AccountID   string    `db:"account_id"`
	ProductsIDS []string  `db:"products_ids"`
	CreatedAT   time.Time `db:"created_at"`
}

type ParamFindOrderByProduct struct {
	ProductID string
}

type ParamFindOrderByProductResult struct {
	OrderID     string    `db:"order_id"`
	AccountID   string    `db:"account_id"`
	ProductsIDS []string  `db:"product_ids"`
	CreatedAT   time.Time `db:"created_at"`
}
