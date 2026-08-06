package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aclgo/grpc-orders/internal/models"
	"github.com/aclgo/grpc-orders/internal/orders"
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) orders.Repository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, param *models.ParamCreateOrder) (*models.ParamCreateOrderResult, error) {
	const query = `
		INSERT INTO grpc_orders (
			order_id,
			account_id,
			type,
			amount,
			payment_method,
			status,
			metadata,
			gateway_transaction_id,
			pix_qr_code,
			pix_expiration,
			card_token,
			card_expiration,
			boleto_url,
			boleto_barcode,
			boleto_expiration,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW())
		RETURNING 
			order_id, account_id, type, amount, payment_method, status, metadata,
			gateway_transaction_id, pix_qr_code, pix_expiration, card_token, card_expiration,
			boleto_url, boleto_barcode, boleto_expiration, created_at, updated_at
	`
	metaBytes, err := json.Marshal(param.Metadata)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal metadata: %w", err)
	}

	var created models.ParamCreateOrderResult
	var metaRaw json.RawMessage

	err = r.db.QueryRowxContext(ctx, query,
		param.OrderID,
		param.AccountID,
		param.Type,
		param.Amount,
		param.PaymentMethod,
		param.Status,
		metaBytes,
		param.GatewayTransactionID,
		param.PixQRCode,
		param.PixExpiration,
		param.CardToken,
		param.CardExpiration,
		param.BoletoURL,
		param.BoletoBarcode,
		param.BoletoExpiration,
	).Scan(
		&created.OrderID,
		&created.AccountID,
		&created.Type,
		&created.Amount,
		&created.PaymentMethod,
		&created.Status,
		&metaRaw,
		&created.GatewayTransactionID,
		&created.PixQRCode,
		&created.PixExpiration,
		&created.CardToken,
		&created.CardExpiration,
		&created.BoletoURL,
		&created.BoletoBarcode,
		&created.BoletoExpiration,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("r.db.QueryRowxContext: %w", err)
	}

	created.Metadata = metaRaw

	return &created, nil
}

func (r *repo) FindOrder(ctx context.Context, param *models.ParamFindOrder) (*models.ParamFindOrderResult, error) {
	const query = `
		SELECT 
			order_id, account_id, type, amount, payment_method, status, metadata,
			gateway_transaction_id, pix_qr_code, pix_expiration, card_token, card_expiration,
			boleto_url, boleto_barcode, boleto_expiration, created_at, updated_at
		FROM grpc_orders 
		WHERE order_id = $1
	`

	var result models.ParamFindOrderResult
	var metaRaw json.RawMessage

	if err := r.db.QueryRowxContext(ctx, query, param.OrderID).Scan(
		&result.OrderID,
		&result.AccountID,
		&result.Type,
		&result.Amount,
		&result.PaymentMethod,
		&result.Status,
		&metaRaw,
		&result.GatewayTransactionID,
		&result.PixQRCode,
		&result.PixExpiration,
		&result.CardToken,
		&result.CardExpiration,
		&result.BoletoURL,
		&result.BoletoBarcode,
		&result.BoletoExpiration,
		&result.CreatedAt,
		&result.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("r.db.QueryRowxContext: %w", err)
	}

	result.Metadata = metaRaw

	return &result, nil
}

func (r *repo) FindOrderByAccount(ctx context.Context, param *models.ParamFindOrderByAccount) ([]*models.ParamFindOrderByAccountResult, error) {
	const query = `
		SELECT 
			order_id, account_id, type, amount, payment_method, status, metadata,
			gateway_transaction_id, pix_qr_code, pix_expiration, card_token, card_expiration,
			boleto_url, boleto_barcode, boleto_expiration, created_at, updated_at
		FROM grpc_orders 
		WHERE account_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryxContext(ctx, query, param.AccountID)
	if err != nil {
		return nil, fmt.Errorf("r.db.QueryxContext: %w", err)
	}
	defer rows.Close()

	var results []*models.ParamFindOrderByAccountResult

	for rows.Next() {
		var order models.ParamFindOrderByAccountResult
		var metaRaw json.RawMessage

		if err := rows.Scan(
			&order.OrderID,
			&order.AccountID,
			&order.Type,
			&order.Amount,
			&order.PaymentMethod,
			&order.Status,
			&metaRaw,
			&order.GatewayTransactionID,
			&order.PixQRCode,
			&order.PixExpiration,
			&order.CardToken,
			&order.CardExpiration,
			&order.BoletoURL,
			&order.BoletoBarcode,
			&order.BoletoExpiration,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		order.Metadata = metaRaw
		results = append(results, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return results, nil
}

func (r *repo) FindOrderByProduct(ctx context.Context, param *models.ParamFindOrderByProduct) ([]*models.ParamFindOrderByProductResult, error) {
	// Busca dentro do JSONB no campo metadata -> products
	const query = `
		SELECT 
			order_id, account_id, type, amount, payment_method, status, metadata,
			gateway_transaction_id, pix_qr_code, pix_expiration, card_token, card_expiration,
			boleto_url, boleto_barcode, boleto_expiration, created_at, updated_at
		FROM grpc_orders 
		WHERE metadata @> jsonb_build_object('products', jsonb_build_array(jsonb_build_object('product_id', $1::text)))
	`

	rows, err := r.db.QueryxContext(ctx, query, param.ProductID)
	if err != nil {
		return nil, fmt.Errorf("r.db.QueryxContext: %w", err)
	}
	defer rows.Close()

	var results []*models.ParamFindOrderByProductResult

	for rows.Next() {
		var result models.ParamFindOrderByProductResult
		var metaRaw json.RawMessage

		if err := rows.Scan(
			&result.OrderID,
			&result.AccountID,
			&result.Type,
			&result.Amount,
			&result.PaymentMethod,
			&result.Status,
			&metaRaw,
			&result.GatewayTransactionID,
			&result.PixQRCode,
			&result.PixExpiration,
			&result.CardToken,
			&result.CardExpiration,
			&result.BoletoURL,
			&result.BoletoBarcode,
			&result.BoletoExpiration,
			&result.CreatedAt,
			&result.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		result.Metadata = metaRaw
		results = append(results, &result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return results, nil
}