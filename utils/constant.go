package vUtil

const (
	Http1     = "http/1.1"
	Http2     = "http/2"
	WebSocket = "WebSocket"
)

// 支持的协议
var Protocol = []string{Http1, Http2, WebSocket}

type Map map[string]interface{}

type VortexHeader string

func (v *VortexHeader) S() string {
	return string(*v)
}

func (v *VortexHeader) V() string {
	return "Vortex-" + v.S()
}

var VortexHeaders = struct {
	Authorization VortexHeader
	ContentType   VortexHeader
	ContentLength VortexHeader
	ContentMd5    VortexHeader
	UserAgent     VortexHeader
}{
	Authorization: "Authorization",
	ContentType:   "Content-Type",
	ContentLength: "Content-Length",
	ContentMd5:    "Content-Md5",
	UserAgent:     "User-Agent",
}

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
)
