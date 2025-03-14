namespace go user

include "/home/kq/GolandProjects/seckill/common/thrift/common.thrift"

//注册请求
struct RegisterReq{
    1: string username
    2: string password
    3: optional string email
}

//注册响应
struct RegisterResp{
    1: common.Resp resp
}

//登录请求
struct LoginReq{
    1: string username
    2: string password
}

//登录响应
struct LoginResp{
    1: common.Resp resp
}

service UserService{
    //注册
    RegisterResp Register(1: RegisterReq req),
    //登录
    LoginResp Login(1: LoginReq req),
}