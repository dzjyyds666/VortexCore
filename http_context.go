package vortex

import (
	"context"

	"github.com/labstack/echo/v4"
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

func NewHttpServer(ctx context.Context, routers []*HttpRouter) *httpServer {
	e := echo.New()

	vortex := e.Group("/v1")

	for _, router := range routers {
		for _, method := range router.Method {
			vortex.Add(method, router.path, func(ctx echo.Context) error {
				return nil
			}, router.ToMiddleWareList()...)
		}
	}

	return &httpServer{
		ctx: ctx,
		e:   e,
	}
}

type VortexHttpMiddleware echo.MiddlewareFunc // Vortex HTTP 中间件类型

type HttpRouter struct {
	handle      func(VortexContext) error // 路由处理函数
	path        string                    // 路由路径
	Method      []string                  // HTTP方法
	middleWares []VortexHttpMiddleware    // 中间件
}

func AppendHttpRouter(method []string, path string, handle func(VortexContext) error, middlwWares ...VortexHttpMiddleware) *HttpRouter {
	return &HttpRouter{
		handle:      handle,
		path:        path,
		Method:      method,
		middleWares: middlwWares,
	}
}

func (hr *HttpRouter) ToMiddleWareList() []echo.MiddlewareFunc {
	middlewares := make([]echo.MiddlewareFunc, 0, len(hr.middleWares))
	for _, mw := range hr.middleWares {
		middlewares = append(middlewares, echo.MiddlewareFunc(mw))
	}
	return middlewares
}
