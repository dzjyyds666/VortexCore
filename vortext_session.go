package vortex

import (
	"context"

	"github.com/labstack/echo/v4"
)

type VortexContext interface {
	GetContext() context.Context // 解析协议
}

type HttpContext struct {
	ctx echo.Context // Echo 上下文
}

func (h *HttpContext) GetContext() context.Context {
	return h.ctx.Request().Context()
}
