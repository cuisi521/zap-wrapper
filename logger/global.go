package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	globalLogger *Logger
	globalMutex  sync.RWMutex
)

// InitGlobal 初始化全局日志器
func InitGlobal(options ...Option) error {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	logger, err := New(options...)
	if err != nil {
		return err
	}

	globalLogger = logger
	return nil
}

// L 获取全局日志器（如果未初始化则返回控制台日志器）
func L() *Logger {
	globalMutex.RLock()
	defer globalMutex.RUnlock()

	if globalLogger == nil {
		// 返回一个安全的控制台日志器
		return &Logger{
			zapLogger: zap.NewExample(), // 使用示例日志器作为回退
		}
	}

	return globalLogger
}

// SyncGlobal 同步全局日志器
func SyncGlobal() {
	globalMutex.RLock()
	defer globalMutex.RUnlock()

	if globalLogger != nil {
		globalLogger.Sync()
	}
}
