package vortexu

import (
	"fmt"
	"time"

	"github.com/dzjyyds666/opensource/logx"
)

var vlog *logx.Logger

// 初始化日志
// 参数解析：
// logPath 日志文件路径
// logLevel 日志级别，级别低于该等级不会打印
// maxSizeMB 日志文件的最大大小，超过该值会自动生成一个新的日志文件
// consoleOut 是否输出到控制台
func InitVortexLog(logPath string, logLevel logx.LogLevel, maxSizeMB int64, consoleOut bool) error {
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
	vlog.Info(fmt.Sprintf("time=[%s] msg=%s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, args...)))
}

// 打印Debug日志
func Debugf(msg string, args ...interface{}) {
	if vlog == nil {
		panic("vortex log is not initialized")
	}
	vlog.Debug(fmt.Sprintf("time=[%s] msg=%s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, args...)))
}

// 打印Warn日志
func Warnf(msg string, args ...interface{}) {
	if vlog == nil {
		panic("vortex log is not initialized")
	}
	vlog.Warn(fmt.Sprintf("time=[%s] msg=%s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, args...)))
}

// 打印Error日志
func Errorf(msg string, args ...interface{}) {
	if vlog == nil {
		panic("vortex log is not initialized")
	}
	vlog.Error(fmt.Sprintf("time=[%s] msg=%s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, args...)))
}
