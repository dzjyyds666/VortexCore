package httpx

var HttpStatus = struct {
	StatusOk                  int
	StatusBadRequest          int
	StatusInternalServerError int
	StatusNotFound            int
	StatusUnauthorized        int
	StatusForbidden           int
	StatusMethodNotAllowed    int
	StatusServiceUnavailable  int
	StatusGatewayTimeout      int
	StatusTooManyRequests     int
}{
	StatusOk:                  200,
	StatusBadRequest:          400,
	StatusUnauthorized:        401,
	StatusForbidden:           403,
	StatusNotFound:            404,
	StatusMethodNotAllowed:    405,
	StatusTooManyRequests:     429, // 请求过多
	StatusInternalServerError: 500,
	StatusServiceUnavailable:  503, // 服务不可用
	StatusGatewayTimeout:      504, // 网关超时

}
