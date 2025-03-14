namespace go common

//返回码
const i32 Success = 200
const i32 Wrong = 500

//返回前端结构体
struct Resp {
    1: i32 code
    2: string msg
    3: string data
}