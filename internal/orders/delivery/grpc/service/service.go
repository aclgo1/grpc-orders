package service

import (
	"context"
	"encoding/json"
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

func (s *serviceGrpc) Create(ctx context.Context, req *proto.ParamCreateOrderRequest) (*proto.ParamCreateOrderResponse, error) {
	var metadata any
	if len(req.Metadata) > 0 {
		_ = json.Unmarshal(req.Metadata, &metadata)
	}

	p := orders.ParamCreateOrder{
		AccountID:     			req.AccountID,
		Type:           		req.Type.String(),
		Amount:         		req.Amount,
		PaymentMethod:  		req.PaymentMethod.String(),
		Status:         		req.Status.String(),
		Metadata:				metadata,
    	GatewayTransactionID:	req.GatewayTransactionID,
    	PixQRCode:				req.PixQRCode,
    	PixExpiration:			req.PixExpiration.AsTime(),
    	CardToken:				req.CardToken,
    	CardExpiration:			req.CardExpiration,
    	BoletoURL: 				req.BoletoURL,
    	BoletoBarcode:			req.BoletoBarcode,
    	BoletoExpiration:		req.BoletoExpiration.AsTime(),
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	created, err := s.orderUC.Create(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.orderUC.Create: %w", err)
	}

	pbOrder, err := toProtoOrderFromCreateResult(created)
	if err != nil {
		return nil, fmt.Errorf("toProtoOrderFromCreateResult: %w", err)
	}

	return &proto.ParamCreateOrderResponse{
		Order: pbOrder,
	}, nil
}

func (s *serviceGrpc) Find(ctx context.Context, req *proto.ParamFindOrderRequest) (*proto.ParamFindOrderResponse, error) {
	p := orders.ParamFindOrder{OrderID: req.OrderID}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	find, err := s.orderUC.FindOrder(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.orderUC.FindOrder: %w", err)
	}

	pbOrder, err := toProtoOrderFromFindResult(find)
	if err != nil {
		return nil, fmt.Errorf("toProtoOrderFromFindResult: %w", err)
	}

	return &proto.ParamFindOrderResponse{
		Order: pbOrder,
	}, nil
}

func (s *serviceGrpc) FindOrderByAccount(ctx context.Context, req *proto.ParamFindOrderByAccountRequest) (*proto.ParamFindOrderByAccountResponse, error) {
	p := orders.ParamFindOrderByAccount{
		AccountID: req.AccountID,
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	finds, err := s.orderUC.FindOrderByAccount(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.orderUC.FindOrderByAccount: %w", err)
	}

	var pbOrders []*proto.Orders
	for _, f := range finds {
		pbOrder, err := toProtoOrderFromAccountResult(f)
		if err != nil {
			return nil, fmt.Errorf("toProtoOrderFromAccountResult: %w", err)
		}
		pbOrders = append(pbOrders, pbOrder)
	}

	return &proto.ParamFindOrderByAccountResponse{
		Orders: pbOrders,
	}, nil
}

func (s *serviceGrpc) FindOrderByProduct(ctx context.Context, req *proto.ParamFindOrderByProductRequest) (*proto.ParamFindOrderByProductResponse, error) {
	p := orders.ParamFindOrderByProduct{
		ProductID: req.ProductID,
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("p.Validate: %w", err)
	}

	finds, err := s.orderUC.FindOrderByProduct(ctx, &p)
	if err != nil {
		return nil, fmt.Errorf("s.orderUC.FindOrderByProduct: %w", err)
	}

	var pbOrders []*proto.Orders
	for _, f := range finds {
		pbOrder, err := toProtoOrderFromProductResult(f)
		if err != nil {
			return nil, fmt.Errorf("toProtoOrderFromProductResult: %w", err)
		}
		pbOrders = append(pbOrders, pbOrder)
	}

	return &proto.ParamFindOrderByProductResponse{
		Orders: pbOrders,
	}, nil
}

func toProtoOrderFromCreateResult(res *orders.ParamCreateOrderResult) (*proto.Orders, error) {
	metaBytes, _ := json.Marshal(res.Metadata)

	return &proto.Orders{
		OrderID:              res.OrderID,
		AccountID:            res.AccountID,
		Type:                 proto.OrderType(proto.OrderType_value[res.Type]),
		Amount:               res.Amount,
		PaymentMethod:        proto.PaymentMethod(proto.PaymentMethod_value[res.PaymentMethod]),
		Status:               proto.OrderStatus(proto.OrderStatus_value[res.Status]),
		Metadata:             metaBytes,
			GatewayTransactionID: res.GatewayTransactionID,
			PixQRCode:           res.PixQRCode,	
			PixExpiration:        timestamppb.New(res.	PixExpiration),
			CardToken:            res.CardToken,	
			CardExpiration:      res.CardExpiration,	
			BoletoURL:           res.BoletoURL,	
			BoletoBarcode:       res.BoletoBarcode,	
			BoletoExpiration:     timestamppb.New(res.	BoletoExpiration),
		CreatedAT:            timestamppb.New(res.CreatedAt),
		UpdatedAT:            timestamppb.New(res.UpdatedAt),
	}, nil
}

func toProtoOrderFromFindResult(res *orders.ParamFindOrderResult) (*proto.Orders, error) {
	metaBytes, _ := json.Marshal(res.Metadata)

	return &proto.Orders{
		OrderID:              res.OrderID,
		AccountID:            res.AccountID,
		Type:                 proto.OrderType(proto.OrderType_value[res.Type]),
		Amount:               res.Amount,
		PaymentMethod:        proto.PaymentMethod(proto.PaymentMethod_value[res.PaymentMethod]),
		Status:               proto.OrderStatus(proto.OrderStatus_value[res.Status]),
		Metadata:             metaBytes,
		GatewayTransactionID: res.GatewayTransactionID,
		PixQRCode:            res.PixQRCode,	
		PixExpiration:        timestamppb.New(res.	PixExpiration),
		CardToken:            res.CardToken,	
		CardExpiration:       res.CardExpiration,	
		BoletoURL:            res.BoletoURL,	
		BoletoBarcode:        res.BoletoBarcode,	
		BoletoExpiration:     timestamppb.New(res.	BoletoExpiration),
		CreatedAT:            timestamppb.New(res.CreatedAt),
		UpdatedAT:            timestamppb.New(res.UpdatedAt),
	}, nil
}

func toProtoOrderFromAccountResult(res *orders.ParamFindOrderByAccountResult) (*proto.Orders, error) {
	metaBytes, _ := json.Marshal(res.Metadata)

	return &proto.Orders{
		OrderID:              res.OrderID,
		AccountID:            res.AccountID,
		Type:                 proto.OrderType(proto.OrderType_value[res.Type]),
		Amount:               res.Amount,
		PaymentMethod:        proto.PaymentMethod(proto.PaymentMethod_value[res.PaymentMethod]),
		Status:               proto.OrderStatus(proto.OrderStatus_value[res.Status]),
		Metadata:             metaBytes,
		GatewayTransactionID: res.GatewayTransactionID,
		PixQRCode:            res.PixQRCode,	
		PixExpiration:        timestamppb.New(res.	PixExpiration),
		CardToken:            res.CardToken,	
		CardExpiration:       res.CardExpiration,	
		BoletoURL:            res.BoletoURL,	
		BoletoBarcode:        res.BoletoBarcode,	
		BoletoExpiration:     timestamppb.New(res.	BoletoExpiration),
		CreatedAT:            timestamppb.New(res.CreatedAt),
		UpdatedAT:            timestamppb.New(res.UpdatedAt),
	}, nil
}

func toProtoOrderFromProductResult(res *orders.ParamFindOrderByProductResult) (*proto.Orders, error) {
	metaBytes, _ := json.Marshal(res.Metadata)

	return &proto.Orders{
		OrderID:              res.OrderID,
		AccountID:            res.AccountID,
		Type:                 proto.OrderType(proto.OrderType_value[res.Type]),
		Amount:               res.Amount,
		PaymentMethod:        proto.PaymentMethod(proto.PaymentMethod_value[res.PaymentMethod]),
		Status:               proto.OrderStatus(proto.OrderStatus_value[res.Status]),
		Metadata:             metaBytes,
		GatewayTransactionID: res.GatewayTransactionID,
		PixQRCode:            res.PixQRCode,	
		PixExpiration:        timestamppb.New(res.	PixExpiration),
		CardToken:            res.CardToken,	
		CardExpiration:       res.CardExpiration,	
		BoletoURL:            res.BoletoURL,	
		BoletoBarcode:        res.BoletoBarcode,	
		BoletoExpiration:     timestamppb.New(res.	BoletoExpiration),
		CreatedAT:            timestamppb.New(res.CreatedAt),
		UpdatedAT:            timestamppb.New(res.UpdatedAt),
	}, nil
}
