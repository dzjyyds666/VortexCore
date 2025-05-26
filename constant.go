package vortex

const (
	http1     = "http/1.1"
	http2     = "http/2"
	webSocket = "webSocket"
)

// 支持的协议
var Protocol = []string{http1, http2, webSocket}
