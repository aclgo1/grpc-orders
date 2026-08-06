package orders

import (
	"context"
	"encoding/json"
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
	) ([]*ParamFindOrderByProductResult, error)
}

type Repository interface {
	Create(context.Context, *models.ParamCreateOrder) (*models.ParamCreateOrderResult, error)
	FindOrder(context.Context, *models.ParamFindOrder) (*models.ParamFindOrderResult, error)

	FindOrderByAccount(context.Context, *models.ParamFindOrderByAccount,
	) ([]*models.ParamFindOrderByAccountResult, error)

	FindOrderByProduct(context.Context, *models.ParamFindOrderByProduct,
	) ([]*models.ParamFindOrderByProductResult, error)
}

type ParamCreateOrder struct {
	AccountID            string          
	Type                 string         
	Amount               int64           
	PaymentMethod        string          
	Status               string        
	Metadata             any            
	GatewayTransactionID string        
	PixQRCode            string        
	PixExpiration        time.Time     
	CardToken            string        
	CardExpiration       string         
	BoletoURL            string        
	BoletoBarcode        string         
	BoletoExpiration     time.Time      
}

type ParamCreateOrderResult struct {
	OrderID              string         
	AccountID            string         
	Type                 string          
	Amount               int64           
	PaymentMethod        string          
	Status               string          
	Metadata             json.RawMessage 
	GatewayTransactionID string         
	PixQRCode            string         
	PixExpiration        time.Time      
	CardToken            string        
	CardExpiration       string         
	BoletoURL            string         
	BoletoBarcode        string         
	BoletoExpiration     time.Time     
	CreatedAt            time.Time       
	UpdatedAt            time.Time      
}

type ParamFindOrderResult = ParamCreateOrderResult
type ParamFindOrderByAccountResult = ParamCreateOrderResult
type ParamFindOrderByProductResult = ParamCreateOrderResult

func (p *ParamCreateOrder) Validate() error {
	if p.AccountID == "" {
		return errors.New("account id empty")
	}

	_, err := uuid.Parse(p.AccountID)
	if err != nil {
		return fmt.Errorf("account uuid invalid: %w", err)
	}

	// if len(p.ProductsIDS) <= 0 {
	// 	return fmt.Errorf("product ids empty")
	// }

	return nil
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

