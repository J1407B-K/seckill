namespace go order

include "common.thrift"

// 下单请求
struct OrderReq {
  1: string userId,       // 用户 ID
  2: string productId,    // 商品 ID
  3: i32 quantity         // 商品数量
}

// 下单响应
struct OrderResp {
  1: i32 code,
  2: string message,
  3: optional string orderId   // 成功返回订单号
}

// 订单查询请求（可扩展）
struct OrderQueryRequest {
  1: string orderId
}

// 订单查询响应
struct OrderQueryResponse {
  1: i32 code,
  2: string message,
  3: optional Order order   // 订单详细信息
}

// 订单详细信息结构体
struct Order {
  1: string orderId,
  2: string userId,
  3: string productId,
  4: i32 quantity,
  5: i64 timestamp       // 订单生成时间
}

service OrderService {
  // 创建订单接口
  OrderResp CreateOrder(1: OrderReq req),
  // 查询订单接口
  OrderQueryResponse QueryOrder(1: OrderQueryRequest req)
}
