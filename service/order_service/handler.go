package main

import (
	"context"
	order "seckill/service/order_service/kitex_gen/order"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// CreateOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *order.OrderReq) (resp *order.OrderResp, err error) {
	// TODO: Your code here...
	return
}

// QueryOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) QueryOrder(ctx context.Context, req *order.OrderQueryReq) (resp *order.OrderQueryResp, err error) {
	// TODO: Your code here...
	return
}
