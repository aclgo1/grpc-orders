package models

import (
	"encoding/json"
	"time"
)

type ParamCreateOrder struct {
	OrderID              string          `db:"order_id"`
	AccountID            string          `db:"account_id"`
	Type                 string          `db:"type"`
	Amount               int64           `db:"amount"`
	PaymentMethod        string          `db:"payment_method"`
	Status               string          `db:"status"`
	Metadata             any     `db:"metadata"`
	GatewayTransactionID string         `db:"gateway_transaction_id"`
	PixQRCode            string         `db:"pix_qr_code"`
	PixExpiration        time.Time      `db:"pix_expiration"`
	CardToken            string         `db:"card_token"`
	CardExpiration       string         `db:"card_expiration"`
	BoletoURL            string         `db:"boleto_url"`
	BoletoBarcode        string         `db:"boleto_barcode"`
	BoletoExpiration     time.Time      `db:"boleto_expiration"`
}

type ParamCreateOrderResult struct {
	OrderID              string          `db:"order_id"`
	AccountID            string          `db:"account_id"`
	Type                 string          `db:"type"`
	Amount               int64           `db:"amount"`
	PaymentMethod        string          `db:"payment_method"`
	Status               string          `db:"status"`
	Metadata             json.RawMessage `db:"metadata"`
	GatewayTransactionID *string         `db:"gateway_transaction_id"`
	PixQRCode            *string         `db:"pix_qr_code"`
	PixExpiration        *time.Time      `db:"pix_expiration"`
	CardToken            *string         `db:"card_token"`
	CardExpiration       *string         `db:"card_expiration"`
	BoletoURL            *string         `db:"boleto_url"`
	BoletoBarcode        *string         `db:"boleto_barcode"`
	BoletoExpiration     *time.Time      `db:"boleto_expiration"`
	CreatedAt            time.Time       `db:"created_at"`
	UpdatedAt            time.Time       `db:"updated_at"`
}

type ParamFindOrder struct {
	OrderID string
}


type ParamFindOrderResult = ParamCreateOrderResult

type ParamFindOrderByAccount struct {
	AccountID string
}

type ParamFindOrderByAccountResult = ParamCreateOrderResult

type ParamFindOrderByProduct struct {
	ProductID string
}

type ParamFindOrderByProductResult = ParamCreateOrderResult

