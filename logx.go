package vortex

import (
	"fmt"

	"github.com/dzjyyds666/opensource/logx"
)

var vlog *logx.Logger

func initVortexLog(logPath string, logLevel logx.LogLevel, maxSizeMB int64, consoleOut bool) error {
	logger, err := logx.NewLogger(logPath, logLevel, maxSizeMB, consoleOut)
	if nil != err {
		return err
	}
	vlog = logger
	vlog.StartWorker()
	return nil
}

// 打印Info日志
func Infof(msg string, args ...interface{}) {
	if vlog == nil {
		panic("vortex log is not initialized")
	}
	vlog.Info(fmt.Sprintf(msg, args...))
}

// 打印Debug日志
func Debugf(msg string, args ...interface{}) {
	if vlog == nil {
		panic("vortex log is not initialized")
	}
	vlog.Debug(fmt.Sprintf(msg, args...))
}

// 打印Warn日志
func Warnf(msg string, args ...interface{}) {
	if vlog == nil {
		panic("vortex log is not initialized")
	}
	vlog.Warn(fmt.Sprintf(msg, args...))
}

// 打印Error日志
func Errorf(msg string, args ...interface{}) {
	if vlog == nil {
		panic("vortex log is not initialized")
	}
	vlog.Error(fmt.Sprintf(msg, args...))
}
