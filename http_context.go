package vortex

import (
	"context"
	"fmt"
	vortexu "github.com/dzjyyds666/VortexCore/utils"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strings"
	"time"
)

type httpContext struct {
	echo.Context // Echo 上下文
}

func (h *httpContext) GetContext() context.Context {
	return h.Request().Context()
}

func (h *httpContext) GetEcho() echo.Context {
	return h.Context
}

type httpServer struct {
	ctx context.Context
	e   *echo.Echo // Echo 框架实例
}

func newHttpServer(ctx context.Context, routers []*httpRouter) *httpServer {
	e := echo.New()

	vortex := e.Group("/v1")

	for _, router := range routers {
		for _, method := range router.method {
			vortex.Add(method, router.path, func(ctx echo.Context) error {
				// 包装成自身封装的上下文
				return router.handle(&httpContext{ctx})
			}, router.ToMiddleWareList()...)
		}
	}

	return &httpServer{
		ctx: ctx,
		e:   e,
	}
}

type httpRouter struct {
	handle      func(VortexContext) error // 路由处理函数
	path        string                    // 路由路径
	method      []string                  // HTTP方法
	middleWares []VortexHttpMiddleware    // 中间件
	description string                    // 路由的描述
}

// 添加 Http 路由
func AppendHttpRouter(method []string, path string, handle func(VortexContext) error, apiDescription string, middleWares ...VortexHttpMiddleware) *httpRouter {
	// 中间件顺序调用 parseJwt -> 自定义中间件 -> verifyJwt
	mws := make([]VortexHttpMiddleware, 0)
	mws = append(mws, printRequestInfoMw(), printResponseInfoMw(), JwtParseMw())
	mws = append(mws, middleWares...)
	mws = append(mws, JwtVerifyMw())

	return &httpRouter{
		handle:      handle,
		path:        path,
		method:      method,
		middleWares: mws,
		description: apiDescription,
	}
}

// 将 VortexHttpMiddleware 转换为 Echo 中间件列表
// 这将允许 Echo 框架使用这些中间件
func (hr *httpRouter) ToMiddleWareList() []echo.MiddlewareFunc {
	middlewares := make([]echo.MiddlewareFunc, 0, len(hr.middleWares))
	for _, mw := range hr.middleWares {
		middlewares = append(middlewares, echo.MiddlewareFunc(mw))
	}
	return middlewares
}

var HttpHeaderEnums = struct {
	ContentType   string
	ContentLength string
	Authorization string

	AcceptLanguage string // 语言
}{
	ContentType:    "Content-Type",
	ContentLength:  "Content-Length",
	Authorization:  "Authorization",
	AcceptLanguage: "Accept-Language",
}

// 响应码
type VortexCode struct {
	Code    int    // http状态码
	subCode int    // 响应体状态码
	I18nKey string // 响应体状态码对应的国际化信息
}

type VortexSubCode struct {
	Code    int
	I18nKey string
}

func (vc *VortexCode) WithSubCode(subCode VortexSubCode) *VortexCode {
	vc.subCode = subCode.Code
	vc.I18nKey = subCode.I18nKey
	return vc
}

type VortexHttpResponse struct {
	Body interface{} `json:"body,omitempty"`
	Info struct {
		Url  string `json:"url,omitempty"`  // 请求地址
		Time int64  `json:"time,omitempty"` // 响应时间
		Ec   int64  `json:"ec,omitempty"`   // 响应的错误码
		Em   string `json:"em"`             // 响应的错误信息
	} `json:"info,omitempty"` // 响应的信息
}

type HttpOpt func(resp http.Header) http.Header

func WithContentType(contentType string) HttpOpt {
	return func(resp http.Header) http.Header {
		resp.Set(vortexu.VortexHeaders.ContentType.S(), contentType)
		return resp
	}
}

// HttpJsonResponse 返回json数据
func HttpJsonResponse(ctx echo.Context, vertexCode VortexCode, data interface{}, opts ...HttpOpt) error {
	// 设置响应的请求头
	for _, opt := range opts {
		opt(ctx.Response().Header())
	}

	// 从请求体重获取到国际化信息
	lang := ctx.Request().Header.Get(HttpHeaderEnums.AcceptLanguage)
	lower := strings.ToLower(fmt.Sprintf("%s.%s", vertexCode.I18nKey, lang))
	n := getI18n(lower)

	var subcode int
	if vertexCode.subCode == 0 {
		subcode = vertexCode.Code
	}

	return ctx.JSON(vertexCode.Code, VortexHttpResponse{
		Body: data,
		Info: struct {
			Url  string `json:"url,omitempty"`  // 响应的url
			Time int64  `json:"time,omitempty"` // 响应时间
			Ec   int64  `json:"ec,omitempty"`   // 响应的错误码
			Em   string `json:"em"`             // 响应的错误信息
		}{
			Url:  ctx.Request().URL.String(),
			Time: time.Now().Unix(),
			Ec:   int64(subcode),
			Em:   n,
		},
	})
}

// 流式返回数据
func HttpStreamResponse(ctx echo.Context, code int, stream io.Reader, opts ...HttpOpt) error {
	for _, opt := range opts {
		opt(ctx.Response().Header())
	}
	contentType := ctx.Response().Header().Get(vortexu.VortexHeaders.ContentType.S())
	if len(contentType) <= 0 {
		err := ctx.Stream(code, "application/octet-stream", stream)
		if nil != err {
			vortexu.Errorf("HttpStreamResponse|HttpStreamError:%v", err)
		}
	} else {
		err := ctx.Stream(code, contentType, stream)
		if nil != err {
			vortexu.Errorf("HttpStreamResponse|HttpStreamError:%v", err)
		}
	}
	ctx.Response().Flush()
	return nil
}

// http请求
func Do(ctx context.Context, hcli *http.Client, method string, reqUrl string, body io.Reader, opts ...HttpOpt) (*http.Response, error) {
	req, err := http.NewRequest(method, reqUrl, body)
	if nil != err {
		vortexu.Errorf("HttpDO|HttpRequestError:%v", err)
		return nil, err
	}

	for _, opt := range opts {
		opt(req.Header)
	}

	resp, err := hcli.Do(req)
	if nil != err {
		vortexu.Errorf("HttpDO|HttpResponseError:%v", err)
		return nil, err
	}
	return resp, nil
}

var HttpStatue = struct {
	Success        VortexCode // 请求成功 200
	InvalidParams  VortexCode // 参数错误 400
	InternalServer VortexCode // 内部服务器错误 500
	PermissionDeny VortexCode // 权限不足 403
}{
	Success:        VortexCode{Code: 200},
	InvalidParams:  VortexCode{Code: 400},
	InternalServer: VortexCode{Code: 500},
	PermissionDeny: VortexCode{Code: 403},
}
