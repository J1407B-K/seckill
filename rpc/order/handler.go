package main

import (
	"context"
	"gorm.io/gorm"
	order "seckill/idl/kitex_gen/order"
	"seckill/idl/kitex_gen/stock/stockservice"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct {
	db       *gorm.DB
	stockCli stockservice.Client
}

// CreateOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *order.OrderReq) (resp *order.OrderResp, err error) {
	// TODO: Your code here...
	return
}

// ConfirmOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ConfirmOrder(ctx context.Context, orderId string) (resp *order.OrderResp, err error) {
	// TODO: Your code here...
	return
}

// CancelOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, orderId string) (resp *order.OrderResp, err error) {
	// TODO: Your code here...
	return
}

// QueryOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) QueryOrder(ctx context.Context, orderId string) (resp *order.OrderResp, err error) {
	// TODO: Your code here...
	return
}
