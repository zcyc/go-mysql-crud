package result

type ResponseResult struct {
	Result bool        `json:"result"` // 是否成功
	Msg    string      `json:"msg"`    // 错误描述
	Code   int         `json:"code"`   // 错误码
	Data   interface{} `json:"data"`   // 返回数据
}

// 成功响应,带返回值
func SuccessDate(data interface{}) ResponseResult {
	return ResponseResult{
		Result: true,
		Msg:    "success",
		Code:   SUCCESS,
		Data:   data,
	}
}
