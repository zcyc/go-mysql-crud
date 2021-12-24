package result

const (
	SUCCESS    = 200  // 正常响应
	ERROR_USER = 2001 // 服务器正常，用户数据错误
	ERROR      = 500  // 服务器异常
	ERROR_AUTH = 401  // 拒绝访问
)
