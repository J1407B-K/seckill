namespace go stock

include "/home/kq/GolandProjects/seckill/common/thrift/common.thrift"

//库存操作请求
struct StockReq{
    1: string productId
}

//库存操作响应
struct StockResp{
    1: common.Resp resp,
}

service StockService{
    //查询库存接口
    StockResp QueryStock(1: StockReq req),
    //减少库存
    StockResp ProDeductStock(1: StockReq req),
    //回滚库存(取消订单)
    StockResp RollCallBackStock(1: StockReq req),
}