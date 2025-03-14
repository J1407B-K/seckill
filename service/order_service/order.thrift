namespace go order

include "/home/kq/GolandProjects/seckill/common/thrift/common.thrift"

//下单请求
struct OrderReq{
    1: string userId,
    2: string productId,
    3: i32 quantity,
}

//下单响应
struct OrderResp{
    1: common.Resp resp,
    2: optional string orderId,
}

//订单查询
struct OrderQueryReq{
    1: string orderId,
}

//订单查询响应
struct OrderQueryResp{
    1: common.Resp resp,
    2: optional Order order,
}

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
  OrderQueryResp QueryOrder(1: OrderQueryReq req)
}