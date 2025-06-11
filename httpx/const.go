package httpx

var HttpStatus = struct {
	Ok                  int
	BadRequest          int
	InternalServerError int
	NotFound            int
	Unauthorized        int
	Forbidden           int
	MethodNotAllowed    int
	ServiceUnavailable  int
	GatewayTimeout      int
	TooManyRequests     int
	ParamError          int
}{
	Ok:                  200,
	BadRequest:          400,
	ParamError:          400,
	Unauthorized:        401,
	Forbidden:           403,
	NotFound:            404,
	MethodNotAllowed:    405,
	TooManyRequests:     429, // 请求过多
	InternalServerError: 500,
	ServiceUnavailable:  503, // 服务不可用
	GatewayTimeout:      504, // 网关超时
}
