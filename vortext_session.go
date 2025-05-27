package vortex

import (
	"context"

	"github.com/labstack/echo/v4"
)

type VortexContext interface {
	GetContext() context.Context // 解析协议
	GetEcho() echo.Context       // 获取 Echo 上下文
}
