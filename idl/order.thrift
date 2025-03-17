namespace go order

include "common.thrift"

// 订单请求参数（创建订单时使用）
struct OrderReq {
  1: string userId,
  2: string productId,
  3: i32 count
}

// 订单响应参数
struct OrderResp {
  1: i32 code,                // 状态码：0 成功，非 0 表示失败
  2: string message,          // 描述信息
  3: optional string orderId  // 生成的订单 ID（创建订单时返回）
}

service OrderService {
  // 创建订单：下单时调用 ReserveStock 预占库存
  OrderResp CreateOrder(1: OrderReq req),
  // 支付成功：调用 ConfirmDeductStock 确认扣减库存并更新订单状态
  OrderResp ConfirmOrder(1: string orderId),
  // 取消订单（含超时、用户取消）：调用 RollbackStock 回滚库存，并更新订单状态为取消
  OrderResp CancelOrder(1: string orderId),
  // 查询订单状态
  OrderResp QueryOrder(1: string orderId)
}
