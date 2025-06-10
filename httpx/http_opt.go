package httpx

import (
	"context"
	"io"
	"net/http"
	"time"

	vortexu "github.com/dzjyyds666/VortexCore/utils"
	"github.com/labstack/echo/v4"
)

type VortexHttpResponse struct {
	Body interface{} `json:"body,omitempty"`
	Code int         `json:"code,omitempty"`
	Info struct {
		Url  string `json:"url,omitempty"`  // 请求地址
		Time int64  `json:"time,omitempty"` // 响应时间
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
func HttpJsonResponse(ctx echo.Context, code int, data interface{}, opts ...HttpOpt) error {
	// 设置响应的请求头
	for _, opt := range opts {
		opt(ctx.Response().Header())
	}

	return ctx.JSON(code, VortexHttpResponse{
		Code: code,
		Body: data,
		Info: struct {
			Url  string `json:"url,omitempty"`  // 响应的url
			Time int64  `json:"time,omitempty"` // 响应时间
		}{
			Url:  ctx.Request().URL.String(),
			Time: time.Now().Unix(),
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
		return ctx.Stream(code, "application/octet-stream", stream)
	} else {
		return ctx.Stream(code, contentType, stream)
	}
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
