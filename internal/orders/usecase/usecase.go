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
			OrderID  :uuid.NewString(),
			AccountID      :param.AccountID,
			Type            :param.Type,
			Amount         :    param.Amount,
			PaymentMethod :param.PaymentMethod,
			Status         :param.Status,
			Metadata        :param.Metadata,
			GatewayTransactionID :param.GatewayTransactionID,
			PixQRCode         :param.PixQRCode,
			PixExpiration     :param.PixExpiration,
			CardToken          :param.CardToken,
			CardExpiration    :param.CardExpiration,
			BoletoURL          :param.BoletoURL,
			BoletoBarcode       :param.BoletoBarcode,
			BoletoExpiration    :param.BoletoExpiration,
	}

	create, err := o.repo.Create(ctx, &mo)
	if err != nil {
		return nil, fmt.Errorf("o.repo.Create: %w", err)
	}

	result := orders.ParamCreateOrderResult{
			OrderID  :create.OrderID,
			AccountID      :create.AccountID,
			Type            :create.Type,
			Amount         :    create.Amount,
			PaymentMethod :create.PaymentMethod,
			Status         :create.Status,
			Metadata        :create.Metadata,
			GatewayTransactionID: strVal(create.GatewayTransactionID),
    	    PixQRCode:            strVal(create.PixQRCode),
    	    PixExpiration:        timeVal(create.PixExpiration),
    	    CardToken:            strVal(create.CardToken),
    	    CardExpiration:       strVal(create.CardExpiration),
    	    BoletoURL:            strVal(create.BoletoURL),
    	    BoletoBarcode:        strVal(create.BoletoBarcode),
    	    BoletoExpiration:     timeVal(create.BoletoExpiration),
			CreatedAt: create.CreatedAt,
			UpdatedAt: create.UpdatedAt,
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
			OrderID  :find.OrderID,
			AccountID      :find.AccountID,
			Type            :find.Type,
			Amount         :    find.Amount,
			PaymentMethod :find.PaymentMethod,
			Status         :find.Status,
			Metadata        :find.Metadata,
			GatewayTransactionID: strVal(find.GatewayTransactionID),
    	    PixQRCode:            strVal(find.PixQRCode),
    	    PixExpiration:        timeVal(find.PixExpiration),
    	    CardToken:            strVal(find.CardToken),
    	    CardExpiration:       strVal(find.CardExpiration),
    	    BoletoURL:            strVal(find.BoletoURL),
    	    BoletoBarcode:        strVal(find.BoletoBarcode),
    	    BoletoExpiration:     timeVal(find.BoletoExpiration),
			CreatedAt: find.CreatedAt,
			UpdatedAt: find.UpdatedAt,
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

	for  i := range find {
		result := orders.ParamFindOrderByAccountResult{
			OrderID  :find[i].OrderID,
			AccountID      :find[i].AccountID,
			Type            :find[i].Type,
			Amount         :    find[i].Amount,
			PaymentMethod :find[i].PaymentMethod,
			Status         :find[i].Status,
			Metadata        :find[i].Metadata,
			GatewayTransactionID: strVal(find[i].GatewayTransactionID),
    	    PixQRCode:            strVal(find[i].PixQRCode),
    	    PixExpiration:        timeVal(find[i].PixExpiration),
    	    CardToken:            strVal(find[i].CardToken),
    	    CardExpiration:       strVal(find[i].CardExpiration),
    	    BoletoURL:            strVal(find[i].BoletoURL),
    	    BoletoBarcode:        strVal(find[i].BoletoBarcode),
    	    BoletoExpiration:     timeVal(find[i].BoletoExpiration),
			CreatedAt: find[i].CreatedAt,
			UpdatedAt: find[i].UpdatedAt,
		}

		results = append(results, &result)
	}

	return results, nil
}

func (o *orderUseCase) FindOrderByProduct(ctx context.Context, param *orders.ParamFindOrderByProduct,
) ([]*orders.ParamFindOrderByProductResult, error) {

	mo := models.ParamFindOrderByProduct{
		ProductID: param.ProductID,
	}

	find, err := o.repo.FindOrderByProduct(ctx, &mo)
	if err != nil {
		return nil, fmt.Errorf("o.repo.FindOrderByProduct: %w", err)
	}

	var results []*orders.ParamFindOrderByProductResult

	for i := range find {
		product := orders.ParamFindOrderByProductResult{
			OrderID  :find[i].OrderID,
			AccountID      :find[i].AccountID,
			Type            :find[i].Type,
			Amount         :    find[i].Amount,
			PaymentMethod :find[i].PaymentMethod,
			Status         :find[i].Status,
			Metadata        :find[i].Metadata,
			GatewayTransactionID: strVal(find[i].GatewayTransactionID),
    	    PixQRCode:            strVal(find[i].PixQRCode),
    	    PixExpiration:        timeVal(find[i].PixExpiration),
    	    CardToken:            strVal(find[i].CardToken),
    	    CardExpiration:       strVal(find[i].CardExpiration),
    	    BoletoURL:            strVal(find[i].BoletoURL),
    	    BoletoBarcode:        strVal(find[i].BoletoBarcode),
    	    BoletoExpiration:     timeVal(find[i].BoletoExpiration),
			CreatedAt: find[i].CreatedAt,
			UpdatedAt: find[i].UpdatedAt,
		}

		results = append(results, &product)
	}

	return results, nil
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func timeVal(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}