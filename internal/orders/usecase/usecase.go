package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/aclgo/grpc-orders/internal/models"
	"github.com/aclgo/grpc-orders/internal/orders"
	"github.com/google/uuid"
)

type orderUseCase struct {
	repo orders.Repository
}

func NewOrderUseCase(repo orders.Repository) orders.UseCase {
	return &orderUseCase{repo: repo}
}

func (o *orderUseCase) Create(ctx context.Context, param *orders.ParamCreateOrder,
) (*orders.ParamCreateOrderResult, error) {
	mo := models.ParamCreateOrder{
		OrderID:     uuid.NewString(),
		AccountID:   param.AccountID,
		ProductsIDS: param.ProductsIDS,
		CreatedAT:   time.Now(),
	}

	create, err := o.repo.Create(ctx, &mo)
	if err != nil {
		return nil, fmt.Errorf("o.repo.Create: %w", err)
	}

	result := orders.ParamCreateOrderResult{
		OrderID:     create.OrderID,
		AccountID:   create.AccountID,
		ProductsIDS: create.ProductsIDS,
		CreatedAT:   create.CreatedAT,
	}

	return &result, nil
}

func (o *orderUseCase) FindOrder(ctx context.Context, param *orders.ParamFindOrder,
) (*orders.ParamFindOrderResult, error) {

	mo := models.ParamFindOrder{OrderID: param.OrderID}

	find, err := o.repo.FindOrder(ctx, &mo)
	if err != nil {
		return nil, fmt.Errorf("o.repo.FindOrder: %w", err)
	}

	result := orders.ParamFindOrderResult{
		OrderID:     find.OrderID,
		AccountID:   find.AccountID,
		ProductsIDS: find.ProductsIDS,
		CreatedAT:   find.CreatedAT,
	}

	return &result, nil

}

func (o *orderUseCase) FindOrderByAccount(ctx context.Context, param *orders.ParamFindOrderByAccount,
) ([]*orders.ParamFindOrderByAccountResult, error) {

	mo := models.ParamFindOrderByAccount{AccountID: param.AccountID}

	find, err := o.repo.FindOrderByAccount(ctx, &mo)
	if err != nil {
		return nil, fmt.Errorf("o.repo.FindOrderByAccount: %w", err)
	}

	var results []*orders.ParamFindOrderByAccountResult

	for _, f := range find {
		result := orders.ParamFindOrderByAccountResult{
			OrderID:     f.OrderID,
			AccountID:   f.AccountID,
			ProductsIDS: f.ProductsIDS,
			CreatedAT:   f.CreatedAT,
		}

		results = append(results, &result)
	}

	return results, nil
}

func (o *orderUseCase) FindOrderByProduct(ctx context.Context, param *orders.ParamFindOrderByProduct,
) (*orders.ParamFindOrderByProductResult, error) {

	mo := models.ParamFindOrderByProduct{
		ProductID: param.ProductID,
	}

	find, err := o.repo.FindOrderByProduct(ctx, &mo)
	if err != nil {
		return nil, fmt.Errorf("o.repo.FindOrderByProduct: %w", err)
	}

	result := orders.ParamFindOrderByProductResult{
		OrderID:     find.OrderID,
		AccountID:   find.AccountID,
		ProductsIDS: find.ProductsIDS,
		CreatedAT:   find.CreatedAT,
	}

	return &result, nil
}
