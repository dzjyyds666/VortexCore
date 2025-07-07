package vortex

import (
	"errors"
	"github.com/dzjyyds666/VortexCore/utils"
	"time"

	"github.com/dzjyyds666/opensource/sdk"
	"github.com/labstack/echo/v4"
)

const (
	JwtVerifySuccess = "Jwt-Verify-Success" // JWT 验证成功
	JwtOption        = "Jwt-Option"         // JWT 验证选项
	JwtVerifySkip    = "Jwt-Verify-Skip"    // 跳过 JWT 验证

)

type VortexHttpMiddleware echo.MiddlewareFunc // Vortex HTTP 中间件类型

func JwtParseMw() VortexHttpMiddleware {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			token := ctx.Request().Header.Get(vortexu.VortexHeaders.Authorization.S())
			jwtToken, err := sdk.ParseJwtToken("", token)
			if nil != err {
				ctx.Set(JwtVerifySuccess, false)
			} else {
				ctx.Set(JwtVerifySuccess, true)
				ctx.Set(JwtOption, jwtToken)
			}
			return next(ctx)
		}
	}
}

func JwtSkipMw() VortexHttpMiddleware {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set(JwtVerifySkip, true)
			return next(ctx)
		}
	}
}

func JwtVerifyMw() VortexHttpMiddleware {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			skip := ctx.Get(JwtVerifySkip)
			succ := ctx.Get(JwtVerifySuccess)
			if skip != nil || succ.(bool) {
				return next(ctx)
			} else {
				return errors.New("jwt verify failed")
			}
		}
	}
}

// 打印请求信息
func printRequestInfoMw() VortexHttpMiddleware {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			vortexu.Infof("\" START ==> %s ==> %s ==> UserAgent=%s\"", ctx.Request().Method, ctx.Request().Host+ctx.Request().URL.Path, ctx.Request().Header.Get(vortexu.VortexHeaders.UserAgent.S()))
			ctx.Set("BeginTime", time.Now().UnixMilli())
			return next(ctx)
		}
	}
}

// 打印响应信息
func printResponseInfoMw() VortexHttpMiddleware {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			err := next(ctx)
			if err != nil {
				vortexu.Errorf("\" %s  %s UserAgent=%s \"", ctx.Request().Method, ctx.Request().Host+ctx.Request().URL.Path, ctx.Request().Header.Get(vortexu.VortexHeaders.UserAgent.S()))
			}
			beginTime := ctx.Get("BeginTime")
			vortexu.Infof("\" END   ==> %s ==> %s ==> time=%vms %v \"", ctx.Request().Method, ctx.Request().Host+ctx.Request().URL.Path, time.Now().UnixMilli()-beginTime.(int64), ctx.Response().Status)
			return nil
		}
	}
}
