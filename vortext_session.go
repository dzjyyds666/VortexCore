package vortex

import (
	"context"

	"github.com/labstack/echo/v4"
)

type VortexContext interface {
	GetContext() context.Context // 解析协议
	GetEcho() echo.Context       // 获取 Echo 上下文
}

type VortexHttpResponse struct {
	Body interface{} `json:"body,omitempty"`
	Code int         `json:"code,omitempty"`
	Info struct {
		Url  string `json:"url,omitempty"`  // 请求地址
		Time int64  `json:"time,omitempty"` // 响应时间
	} `json:"info,omitempty"` // 响应的信息
}
