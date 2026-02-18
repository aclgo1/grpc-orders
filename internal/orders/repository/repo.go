package repository

import (
	"context"
	"fmt"

	"github.com/aclgo/grpc-orders/internal/models"
	"github.com/aclgo/grpc-orders/internal/orders"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type repo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) orders.Repository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, param *models.ParamCreateOrder,
) (*models.ParamCreateOrderResult, error) {
	const query = `INSERT INTO grpc_orders (order_id, account_id, products_ids, created_at)
	VALUES ($1, $2, $3, $4) RETURNING order_id, account_id, products_ids, created_at`

	var created models.ParamCreateOrderResult

	var productsIDS []string

	if err := r.db.QueryRowx(query,
		param.OrderID,
		param.AccountID,
		pq.Array(param.ProductsIDS),
		param.CreatedAT,
	).Scan(&created.OrderID, &created.AccountID, pq.Array(&productsIDS), &created.CreatedAT); err != nil {
		return nil, fmt.Errorf("r.db.QueryRowx: %w", err)
	}

	created.ProductsIDS = productsIDS

	return &created, nil
}

func (r *repo) FindOrder(ctx context.Context, param *models.ParamFindOrder,
) (*models.ParamFindOrderResult, error) {
	const query = `SELECT * FROM grpc_orders WHERE order_id=$1`
	var result models.ParamFindOrderResult
	var productsIDS []string

	if err := r.db.QueryRowxContext(ctx, query, param.OrderID).Scan(
		&result.OrderID,
		&result.AccountID,
		pq.Array(&productsIDS),
		&result.CreatedAT,
	); err != nil {
		return nil, fmt.Errorf("r.db.QueryRow: %w", err)
	}

	result.ProductsIDS = productsIDS

	return &result, nil
}

func (r *repo) FindOrderByAccount(ctx context.Context, param *models.ParamFindOrderByAccount,
) ([]*models.ParamFindOrderByAccountResult, error) {
	const query = `SELECT * FROM grpc_orders WHERE account_id = $1`

	rows, err := r.db.QueryxContext(ctx, query, param.AccountID)
	if err != nil {
		return nil, fmt.Errorf("r.db.QueryxContext: %w", err)
	}

	var results []*models.ParamFindOrderByAccountResult

	for rows.Next() {
		var order models.ParamFindOrderByAccountResult
		var productsIDS []string

		if err := rows.Scan(
			&order.OrderID,
			&order.AccountID,
			pq.Array(&productsIDS),
			&order.CreatedAT,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		order.ProductsIDS = productsIDS

		results = append(results, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return results, nil
}

func (r *repo) FindOrderByProduct(ctx context.Context, param *models.ParamFindOrderByProduct,
) (*models.ParamFindOrderByProductResult, error) {
	const query = `SELECT * from grpc_orders WHERE $1 = ANY(products_ids);`

	row := r.db.QueryRowxContext(ctx, query, param.ProductID)

	var result models.ParamFindOrderByProductResult
	var productsIDS []string

	if err := row.Scan(
		&result.OrderID,
		&result.AccountID,
		pq.Array(&productsIDS),
		&result.CreatedAT,
	); err != nil {
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("row.Err: %w", err)
	}

	result.ProductsIDS = productsIDS

	return &result, nil
}
