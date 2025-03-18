package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"gorm.io/gorm"
	"log"
	order "seckill/idl/kitex_gen/order"
	"seckill/idl/kitex_gen/stock"
	"seckill/idl/kitex_gen/stock/stockservice"
	"seckill/rpc/order/model"
	"strconv"
	"time"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct {
	db       *gorm.DB
	stockCli stockservice.Client
}

func NewStockClient() (stockservice.Client, error) {
	// 使用时请传入真实 etcd 的服务地址，本例中为 127.0.0.1:2379
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}
	return stockservice.NewClient("stockservice", client.WithResolver(r)) // 指定 Resolver
}

// CreateOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *order.OrderReq) (resp *order.OrderResp, err error) {
	stockReq := &stock.StockReq{
		ProductId: req.ProductId,
		Count:     req.Count,
	}

	//预占库存
	stockResp, err := s.stockCli.ReserveStock(ctx, stockReq)
	if err != nil || stockResp.Code != 0 {
		return &order.OrderResp{Code: 1, Message: "库存预占失败"}, err
	}

	// 生成唯一订单ID
	orderID := fmt.Sprintf("%d-%s", time.Now().UnixNano(), req.UserId)

	// 创建订单记录（状态：0 表示待支付）
	newOrder := model.Order{
		OrderID:   orderID,
		UserID:    req.UserId,
		ProductID: req.ProductId,
		Count:     int(req.Count),
		Status:    0,
		CreatedAt: time.Now(),
	}
	if err := s.db.Create(&newOrder).Error; err != nil {
		// 订单创建失败，释放预占库存
		_, err := s.stockCli.ReleaseStock(ctx, stockReq)
		if err != nil {
			return nil, err
		}
		return &order.OrderResp{Code: 2, Message: "订单创建失败"}, err
	}

	return &order.OrderResp{Code: 0, Message: "订单创建成功", OrderId: &orderID}, nil
}

// ConfirmOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ConfirmOrder(ctx context.Context, orderId string) (resp *order.OrderResp, err error) {
	//查找订单
	var o model.Order
	if err := s.db.Where("order_id = ?", orderId).First(&o).Error; err != nil {
		return &order.OrderResp{Code: 1, Message: "订单不存在"}, err
	}
	if o.Status != 0 {
		return &order.OrderResp{Code: 2, Message: "订单状态错误，无法确认支付"}, nil
	}

	//扣减库存
	stockReq := &stock.StockReq{
		ProductId: o.ProductID,
		Count:     int32(o.Count),
	}

	log.Printf("Attempting to pre-deduct stock for ProductID: %s, Count: %d", o.ProductID, o.Count)
	stockResp, err := s.stockCli.PreDeductStock(ctx, stockReq)
	if err != nil || stockResp.Code != 0 {
		log.Printf("PreDeductStock failed for ProductID: %s, Count: %d, Error: %v", o.ProductID, o.Count, err)

		// 释放库存
		log.Printf("Attempting to release reserved stock for ProductID: %s, Count: %d", o.ProductID, o.Count)
		releaseResp, releaseErr := s.stockCli.ReleaseStock(ctx, stockReq)
		if releaseErr != nil || releaseResp.Code != 0 {
			log.Printf("ReleaseStock failed for ProductID: %s, Count: %d, Error: %v", o.ProductID, o.Count, releaseErr)
			return &order.OrderResp{Code: 3, Message: "库存扣减失败，且释放库存失败"}, releaseErr
		}
		return &order.OrderResp{Code: 3, Message: "库存扣减失败，订单无法支付"}, err
	}

	//更新订单状态
	o.Status = 1 // 1代表已支付
	if err := s.db.Save(&o).Error; err != nil {
		// 如果更新订单状态失败，则释放库存
		releaseResp, releaseErr := s.stockCli.ReleaseStock(ctx, stockReq)
		if releaseErr != nil || releaseResp.Code != 0 {
			return &order.OrderResp{Code: 4, Message: "订单支付成功，但更新失败，库存已释放"}, releaseErr
		}
		return &order.OrderResp{Code: 4, Message: "订单支付成功，但订单更新失败"}, err
	}

	//返回成功响应
	return &order.OrderResp{Code: 0, Message: "订单支付成功", OrderId: &o.OrderID}, nil
}

// CancelOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, orderId string) (resp *order.OrderResp, err error) {
	// 查找订单
	var o model.Order
	if err := s.db.Where("order_id = ?", orderId).First(&o).Error; err != nil {
		return &order.OrderResp{Code: 1, Message: "订单不存在"}, err
	}

	// 如果订单已支付（状态为1），则需要回滚库存
	if o.Status == 1 {
		// 执行回滚库存操作
		stockReq := &stock.StockReq{
			ProductId: o.ProductID,
			Count:     int32(o.Count),
		}

		// 调用库存回滚接口
		stockResp, err := s.stockCli.RollbackStock(ctx, stockReq)
		if err != nil || stockResp.Code != 0 {
			return &order.OrderResp{Code: 2, Message: "库存回滚失败，订单无法取消"}, err
		}

		// 更新订单状态为已取消
		o.Status = 2 // 2代表已取消
		if err := s.db.Save(&o).Error; err != nil {
			return &order.OrderResp{Code: 3, Message: "订单取消失败，无法更新订单状态"}, err
		}

		return &order.OrderResp{Code: 0, Message: "订单取消成功，库存已回滚", OrderId: &o.OrderID}, nil
	}

	// 如果订单未支付，直接释放预占库存
	if o.Status == 0 {
		stockReq := &stock.StockReq{
			ProductId: o.ProductID,
			Count:     int32(o.Count),
		}

		// 调用释放预占库存接口
		stockResp, err := s.stockCli.ReleaseStock(ctx, stockReq)
		if err != nil || stockResp.Code != 0 {
			return &order.OrderResp{Code: 4, Message: "库存释放失败，订单无法取消"}, err
		}

		// 更新订单状态为已取消
		o.Status = 2 // 2代表已取消
		if err := s.db.Save(&o).Error; err != nil {
			return &order.OrderResp{Code: 5, Message: "订单取消失败，无法更新订单状态"}, err
		}

		return &order.OrderResp{Code: 0, Message: "订单取消成功，库存已释放", OrderId: &o.OrderID}, nil
	}

	// 如果订单已被取消或已处理，则返回错误
	return &order.OrderResp{Code: 6, Message: "订单已被取消或处理过，无法再次取消"}, nil
}

// QueryOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) QueryOrder(ctx context.Context, orderId string) (resp *order.OrderResp, err error) {
	var o model.Order
	if err := s.db.Where("order_id = ?", orderId).First(&o).Error; err != nil {
		return &order.OrderResp{Code: 1, Message: "订单不存在"}, err
	}
	return &order.OrderResp{Code: 0, Message: o.ProductID + "		" + strconv.Itoa(o.Count) + "	 " + strconv.Itoa(o.Status)}, nil
}
