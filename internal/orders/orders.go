package orders

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aclgo/grpc-orders/internal/models"
	"github.com/google/uuid"
)

type UseCase interface {
	Create(context.Context, *ParamCreateOrder) (*ParamCreateOrderResult, error)
	FindOrder(context.Context, *ParamFindOrder) (*ParamFindOrderResult, error)

	FindOrderByAccount(context.Context, *ParamFindOrderByAccount,
	) ([]*ParamFindOrderByAccountResult, error)

	FindOrderByProduct(context.Context, *ParamFindOrderByProduct,
	) (*ParamFindOrderByProductResult, error)
}

type Repository interface {
	Create(context.Context, *models.ParamCreateOrder) (*models.ParamCreateOrderResult, error)
	FindOrder(context.Context, *models.ParamFindOrder) (*models.ParamFindOrderResult, error)

	FindOrderByAccount(context.Context, *models.ParamFindOrderByAccount,
	) ([]*models.ParamFindOrderByAccountResult, error)

	FindOrderByProduct(context.Context, *models.ParamFindOrderByProduct,
	) (*models.ParamFindOrderByProductResult, error)
}

type ParamCreateOrder struct {
	AccountID   string
	ProductsIDS []string
}

func (p *ParamCreateOrder) Validate() error {
	if p.AccountID == "" {
		return errors.New("account id empty")
	}

	_, err := uuid.Parse(p.AccountID)
	if err != nil {
		return fmt.Errorf("account uuid invalid: %w", err)
	}

	if len(p.ProductsIDS) <= 0 {
		return fmt.Errorf("product ids empty")
	}

	return nil
}

type ParamCreateOrderResult struct {
	OrderID     string
	AccountID   string
	ProductsIDS []string
	CreatedAT   time.Time
}

type ParamFindOrder struct {
	OrderID string
}

func (p *ParamFindOrder) Validate() error {
	if p.OrderID == "" {
		return errors.New("account id empty")
	}

	_, err := uuid.Parse(p.OrderID)
	if err != nil {
		return fmt.Errorf("account uuid invalid: %w", err)
	}
	return nil
}

type ParamFindOrderResult struct {
	OrderID     string
	AccountID   string
	ProductsIDS []string
	CreatedAT   time.Time
}

type ParamFindOrderByAccount struct {
	AccountID string
}

func (p *ParamFindOrderByAccount) Validate() error {
	if p.AccountID == "" {
		return errors.New("account id empty")
	}

	_, err := uuid.Parse(p.AccountID)
	if err != nil {
		return fmt.Errorf("account uuid invalid: %w", err)
	}

	return nil
}

type ParamFindOrderByAccountResult struct {
	OrderID     string
	AccountID   string
	ProductsIDS []string
	CreatedAT   time.Time
}

type ParamFindOrderByProduct struct {
	ProductID string
}

func (p *ParamFindOrderByProduct) Validate() error {
	if p.ProductID == "" {
		return errors.New("account id empty")
	}

	_, err := uuid.Parse(p.ProductID)
	if err != nil {
		return fmt.Errorf("account uuid invalid: %w", err)
	}
	return nil
}

type ParamFindOrderByProductResult struct {
	OrderID     string
	AccountID   string
	ProductsIDS []string
	CreatedAT   time.Time
}
