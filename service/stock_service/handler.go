package main

import (
	"context"
	stock "seckill/service/stock_service/kitex_gen/stock"
)

// StockServiceImpl implements the last service interface defined in the IDL.
type StockServiceImpl struct{}

// QueryStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) QueryStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// TODO: Your code here...
	return
}

// ProDeductStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) ProDeductStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// TODO: Your code here...
	return
}

// RollCallBackStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) RollCallBackStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// TODO: Your code here...
	return
}
