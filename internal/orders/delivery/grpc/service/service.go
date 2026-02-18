package service

import (
	"context"
	"fmt"

	"github.com/aclgo/grpc-orders/internal/orders"
	"github.com/aclgo/grpc-orders/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type serviceGrpc struct {
	orderUC orders.UseCase
	proto.UnimplementedServiceOrderServer
}

func NewServiceGprc(orderUC orders.UseCase) *serviceGrpc {
	return &serviceGrpc{orderUC: orderUC}
}

func (s *serviceGrpc) Create(ctx context.Context, req *proto.ParamCreateOrderRequest,
) (*proto.ParamCreateOrderResponse, error) {
	p := orders.ParamCreateOrder{
		AccountID:   req.AccountID,
		ProductsIDS: req.ProductsIDS,
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	created, err := s.orderUC.Create(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.orderUC.Create: %w", err)
	}

	r := proto.ParamCreateOrderResponse{
		Order: &proto.Orders{
			OrderID:     created.OrderID,
			AccountID:   created.AccountID,
			ProductsIDS: created.ProductsIDS,
			CreatedAT:   timestamppb.New(created.CreatedAT),
		},
	}

	return &r, nil
}

func (s *serviceGrpc) Find(ctx context.Context, req *proto.ParamFindOrderRequest,
) (*proto.ParamFindOrderResponse, error) {
	p := orders.ParamFindOrder{OrderID: req.OrderID}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	find, err := s.orderUC.FindOrder(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.orderUC.FindOrder: %w", err)
	}

	r := proto.ParamFindOrderResponse{
		Order: &proto.Orders{
			OrderID:     find.OrderID,
			AccountID:   find.AccountID,
			ProductsIDS: find.ProductsIDS,
			CreatedAT:   timestamppb.New(find.CreatedAT),
		},
	}

	return &r, nil
}

func (s *serviceGrpc) FindOrderByAccount(ctx context.Context, req *proto.ParamFindOrderByAccountRequest,
) (*proto.ParamFindOrderByAccountResponse, error) {

	p := orders.ParamFindOrderByAccount{
		AccountID: req.AccountID,
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	find, err := s.orderUC.FindOrderByAccount(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.orderUC.FindOrderByAccount: %w", err)
	}

	var orders []*proto.Orders

	for _, f := range find {
		order := proto.Orders{
			OrderID:     f.OrderID,
			AccountID:   f.AccountID,
			ProductsIDS: f.ProductsIDS,
			CreatedAT:   timestamppb.New(f.CreatedAT),
		}

		orders = append(orders, &order)
	}

	var results proto.ParamFindOrderByAccountResponse
	results.Orders = append(results.Orders, orders...)

	return &results, nil
}

func (s *serviceGrpc) FindOrderByProduct(ctx context.Context, req *proto.ParamFindOrderByProductRequest,
) (*proto.ParamFindOrderByProductResponse, error) {

	p := orders.ParamFindOrderByProduct{
		ProductID: req.ProductID,
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	find, err := s.orderUC.FindOrderByProduct(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.orderUC.FindOrderByProduct: %w", err)
	}

	var result proto.ParamFindOrderByProductResponse

	order := proto.Orders{
		OrderID:     find.OrderID,
		AccountID:   find.AccountID,
		ProductsIDS: find.ProductsIDS,
		CreatedAT:   timestamppb.New(find.CreatedAT),
	}

	result.Order = &order

	return &result, nil
}
