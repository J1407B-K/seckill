package main

import (
	"context"
	stock "seckill/idl/kitex_gen/stock"
)

// StockServiceImpl implements the last service interface defined in the IDL.
type StockServiceImpl struct{}

// QueryStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) QueryStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// TODO: Your code here...
	return
}

// PreDeductStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) PreDeductStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// TODO: Your code here...
	return
}

// RollbackStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) RollbackStock(ctx context.Context, req *stock.StockReq) (resp *stock.StockResp, err error) {
	// TODO: Your code here...
	return
}
