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

// 全局日志方法 - 基础日志

// Debug 全局Debug级别日志
func Debug(msg string, fields ...zap.Field) {
	L().Debug(msg, fields...)
}

// Info 全局Info级别日志
func Info(msg string, fields ...zap.Field) {
	L().Info(msg, fields...)
}

// Warn 全局Warn级别日志
func Warn(msg string, fields ...zap.Field) {
	L().Warn(msg, fields...)
}

// Error 全局Error级别日志
func Error(msg string, fields ...zap.Field) {
	L().Error(msg, fields...)
}

// Panic 全局Panic级别日志
func Panic(msg string, fields ...zap.Field) {
	L().Panic(msg, fields...)
}

// Fatal 全局Fatal级别日志
func Fatal(msg string, fields ...zap.Field) {
	L().Fatal(msg, fields...)
}

// 全局日志方法 - 格式化日志

// Debugf 全局格式化Debug级别日志
func Debugf(format string, args ...interface{}) {
	L().Debugf(format, args...)
}

// Infof 全局格式化Info级别日志
func Infof(format string, args ...interface{}) {
	L().Infof(format, args...)
}

// Warnf 全局格式化Warn级别日志
func Warnf(format string, args ...interface{}) {
	L().Warnf(format, args...)
}

// Errorf 全局格式化Error级别日志
func Errorf(format string, args ...interface{}) {
	L().Errorf(format, args...)
}

// Panicf 全局格式化Panic级别日志
func Panicf(format string, args ...interface{}) {
	L().Panicf(format, args...)
}

// Fatalf 全局格式化Fatal级别日志
func Fatalf(format string, args ...interface{}) {
	L().Fatalf(format, args...)
}

// With 全局添加结构化字段
func With(fields ...zap.Field) *Logger {
	return L().With(fields...)
}

// Sync 同步全局日志器（与SyncGlobal功能相同）
func Sync() error {
	return L().Sync()
}
