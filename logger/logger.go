package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 包装的日志器
type Logger struct {
	zapLogger *zap.Logger
	config    *Config
}

// New 创建新的日志实例，并设置为全局日志器
func New(options ...Option) (*Logger, error) {
	config := &Config{
		Level:         DefaultLevel,
		Encoding:      DefaultEncoding,
		MaxSize:       DefaultMaxSize,
		MaxBackups:    DefaultMaxBackups,
		MaxAge:        DefaultMaxAge,
		Compress:      DefaultCompress,
		ShowCaller:    DefaultShowCaller,
		Stacktrace:    DefaultStacktrace,
		ConsoleOutput: DefaultConsoleOutput,
	}

	for _, opt := range options {
		opt(config)
	}

	// 如果设置了BasePath且未设置OutputPath，则自动设置OutputPath
	if config.BasePath != "" && config.OutputPath == "" && config.OutputPath != "stdout" {
		config.OutputPath = filepath.Join(config.BasePath, "app.log")
	}

	// 如果设置了BasePath且未设置ErrorPath，则自动设置ErrorPath
	if config.BasePath != "" && config.ErrorPath == "" {
		config.ErrorPath = filepath.Join(config.BasePath, "error.log")
	}

	// 如果设置了BasePath且未设置各级别路径，则自动生成各级别路径
	if config.BasePath != "" {
		if config.DebugPath == "" {
			config.DebugPath = filepath.Join(config.BasePath, "debug.log")
		}
		if config.InfoPath == "" {
			config.InfoPath = filepath.Join(config.BasePath, "info.log")
		}
		if config.WarnPath == "" {
			config.WarnPath = filepath.Join(config.BasePath, "warn.log")
		}
		if config.ErrorLPath == "" {
			config.ErrorLPath = filepath.Join(config.BasePath, "error_l.log")
		}
		if config.PanicPath == "" {
			config.PanicPath = filepath.Join(config.BasePath, "panic.log")
		}
		if config.FatalPath == "" {
			config.FatalPath = filepath.Join(config.BasePath, "fatal.log")
		}
	}

	// 创建 zap 配置
	zapConfig := zap.NewProductionConfig()

	// 设置日志级别
	level, err := parseLevel(config.Level)
	if err != nil {
		return nil, err
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// 设置编码格式
	zapConfig.Encoding = string(config.Encoding)

	// 开发模式配置
	if config.Development {
		zapConfig = zap.NewDevelopmentConfig()
	}

	// 创建 encoder
	encoderConfig := zapConfig.EncoderConfig
	encoder := getEncoder(encoderConfig, config.Encoding)

	// 创建 core
	cores := []zapcore.Core{}

	// 控制台输出
	if (config.OutputPath == "" || config.OutputPath == "stdout") || config.ConsoleOutput {
		consoleCore := zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= level
			}),
		)
		cores = append(cores, consoleCore)
	}

	// 文件输出 - 捕获该级别及以上的所有日志
	if config.OutputPath != "" && config.OutputPath != "stdout" {
		fileCore := createFileCore(config.OutputPath, encoder, level, config, false)
		cores = append(cores, fileCore)
	}

	// 错误日志单独输出（兼容旧配置） - 捕获Error及以上级别日志
	if config.ErrorPath != "" {
		errorLevel := zap.ErrorLevel
		errorCore := createFileCore(config.ErrorPath, encoder, errorLevel, config, false)
		cores = append(cores, errorCore)
	}

	// 为各个级别单独创建日志文件输出 - 只捕获特定级别的日志
	if config.DebugPath != "" {
		debugCore := createFileCore(config.DebugPath, encoder, zap.DebugLevel, config, true)
		cores = append(cores, debugCore)
	}

	if config.InfoPath != "" {
		infoCore := createFileCore(config.InfoPath, encoder, zap.InfoLevel, config, true)
		cores = append(cores, infoCore)
	}

	if config.WarnPath != "" {
		warnCore := createFileCore(config.WarnPath, encoder, zap.WarnLevel, config, true)
		cores = append(cores, warnCore)
	}

	if config.ErrorLPath != "" {
		errorLCore := createFileCore(config.ErrorLPath, encoder, zap.ErrorLevel, config, true)
		cores = append(cores, errorLCore)
	}

	if config.PanicPath != "" {
		panicCore := createFileCore(config.PanicPath, encoder, zap.PanicLevel, config, true)
		cores = append(cores, panicCore)
	}

	if config.FatalPath != "" {
		fatalCore := createFileCore(config.FatalPath, encoder, zap.FatalLevel, config, true)
		cores = append(cores, fatalCore)
	}

	// 创建 logger
	core := zapcore.NewTee(cores...)

	// 添加 caller 信息
	opts := []zap.Option{}
	if config.ShowCaller {
		opts = append(opts, zap.AddCaller())
	}

	// 添加堆栈跟踪
	// 注意：这里只在Panic级别添加堆栈跟踪，避免所有Error日志都包含堆栈
	if config.Stacktrace {
		// 可以根据需要调整级别，这里使用PanicLevel避免普通错误日志包含堆栈
		opts = append(opts, zap.AddStacktrace(zap.PanicLevel))
	}

	// 创建基础logger
	baseLogger := zap.New(core, opts...)

	var zapLogger *zap.Logger
	if config.AsyncMode {
		// 保留对基础logger的引用，用于异步日志处理
		zapLogger = baseLogger
	} else {
		// 同步模式直接使用基础logger
		zapLogger = baseLogger
	}

	logger := &Logger{
		zapLogger: zapLogger,
		config:    config,
	}
	
	// 设置为全局日志器
	globalMutex.Lock()
	globalLogger = logger
	globalMutex.Unlock()
	
	return logger, nil
}

// NewDefault 创建默认日志实例
func NewDefault() (*Logger, error) {
	return New()
}

// 在 createFileCore 函数中添加更好的错误处理
func createFileCore(filePath string, encoder zapcore.Encoder, level zapcore.Level, config *Config, isLevelSpecific bool) zapcore.Core {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		// 如果创建目录失败，回退到控制台输出并记录警告
		fmt.Printf("WARN: Failed to create log directory %s: %v. Falling back to console output.\n", dir, err)
		return zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				if isLevelSpecific {
					// 只捕获特定级别的日志
					return lvl == level
				}
				// 捕获该级别及以上的所有日志
				return lvl >= level
			}),
		)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	writer := zapcore.AddSync(lumberJackLogger)

	return zapcore.NewCore(
		encoder,
		writer,
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			if isLevelSpecific {
				// 只捕获特定级别的日志
				return lvl == level
			}
			// 捕获该级别及以上的所有日志
			return lvl >= level
		}),
	)
}

// getEncoder 获取编码器
func getEncoder(encoderConfig zapcore.EncoderConfig, encoding Encoding) zapcore.Encoder {
	switch encoding {
	case ConsoleEncoding:
		return zapcore.NewConsoleEncoder(encoderConfig)
	case JSONEncoding:
		fallthrough
	default:
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

// parseLevel 解析日志级别
func parseLevel(level Level) (zapcore.Level, error) {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel, nil
	case InfoLevel:
		return zapcore.InfoLevel, nil
	case WarnLevel:
		return zapcore.WarnLevel, nil
	case ErrorLevel:
		return zapcore.ErrorLevel, nil
	case PanicLevel:
		return zapcore.PanicLevel, nil
	case FatalLevel:
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown level: %s", level)
	}
}

// 异步日志相关变量
var (
	// 用于异步日志处理的工作池
	logWorkers = 4
	// 日志任务队列
	logChan = make(chan func(), 1000)
)

func init() {
	// 启动异步日志工作协程
	for i := 0; i < logWorkers; i++ {
		go func() {
			for task := range logChan {
				task()
			}
		}()
	}
}

// 基础日志方法
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	if l == nil || l.zapLogger == nil {
		fmt.Printf("DEBUG: %s\n", msg)
		return
	}

	if l.config.AsyncMode {
		// 异步模式：将日志写入任务发送到通道
		logMsg := msg
		logFields := make([]zap.Field, len(fields))
		copy(logFields, fields)

		select {
		case logChan <- func() {
			l.zapLogger.Debug(logMsg, logFields...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行以避免日志丢失
			l.zapLogger.Debug(msg, fields...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Debug(msg, fields...)
	}
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	if l == nil || l.zapLogger == nil {
		fmt.Printf("INFO: %s\n", msg)
		return
	}

	if l.config.AsyncMode {
		// 异步模式：将日志写入任务发送到通道
		logMsg := msg
		logFields := make([]zap.Field, len(fields))
		copy(logFields, fields)

		select {
		case logChan <- func() {
			l.zapLogger.Info(logMsg, logFields...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行以避免日志丢失
			l.zapLogger.Info(msg, fields...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Info(msg, fields...)
	}
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	if l == nil || l.zapLogger == nil {
		fmt.Printf("WARN: %s\n", msg)
		return
	}

	if l.config.AsyncMode {
		// 异步模式：将日志写入任务发送到通道
		logMsg := msg
		logFields := make([]zap.Field, len(fields))
		copy(logFields, fields)

		select {
		case logChan <- func() {
			l.zapLogger.Warn(logMsg, logFields...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行以避免日志丢失
			l.zapLogger.Warn(msg, fields...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Warn(msg, fields...)
	}
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	if l == nil || l.zapLogger == nil {
		fmt.Printf("ERROR: %s\n", msg)
		return
	}

	if l.config.AsyncMode {
		// 异步模式：将日志写入任务发送到通道
		logMsg := msg
		logFields := make([]zap.Field, len(fields))
		copy(logFields, fields)

		select {
		case logChan <- func() {
			l.zapLogger.Error(logMsg, logFields...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行以避免日志丢失
			l.zapLogger.Error(msg, fields...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Error(msg, fields...)
	}
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	if l == nil || l.zapLogger == nil {
		panic(fmt.Sprintf("PANIC: %s", msg))
	}

	// Panic级别的日志通常需要立即执行，不建议异步处理
	// 但我们仍然提供异步选项，同时保留原始行为
	if l.config.AsyncMode {
		logMsg := msg
		logFields := make([]zap.Field, len(fields))
		copy(logFields, fields)

		select {
		case logChan <- func() {
			l.zapLogger.Panic(logMsg, logFields...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行
			l.zapLogger.Panic(msg, fields...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Panic(msg, fields...)
	}
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	if l == nil || l.zapLogger == nil {
		fmt.Printf("FATAL: %s\n", msg)
		os.Exit(1)
	}

	// Fatal级别的日志通常需要立即执行，不建议异步处理
	// 但我们仍然提供异步选项，同时保留原始行为
	if l.config.AsyncMode {
		logMsg := msg
		logFields := make([]zap.Field, len(fields))
		copy(logFields, fields)

		select {
		case logChan <- func() {
			l.zapLogger.Fatal(logMsg, logFields...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行
			l.zapLogger.Fatal(msg, fields...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Fatal(msg, fields...)
	}
}

// 格式化日志方法
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l == nil || l.zapLogger == nil {
		fmt.Printf("DEBUG: "+format+"\n", args...)
		return
	}

	if l.config.AsyncMode {
		// 异步模式：将日志写入任务发送到通道
		logFormat := format
		logArgs := make([]interface{}, len(args))
		copy(logArgs, args)

		select {
		case logChan <- func() {
			l.zapLogger.Sugar().Debugf(logFormat, logArgs...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行以避免日志丢失
			l.zapLogger.Sugar().Debugf(format, args...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Sugar().Debugf(format, args...)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l == nil || l.zapLogger == nil {
		fmt.Printf("INFO: "+format+"\n", args...)
		return
	}

	if l.config.AsyncMode {
		// 异步模式：将日志写入任务发送到通道
		logFormat := format
		logArgs := make([]interface{}, len(args))
		copy(logArgs, args)

		select {
		case logChan <- func() {
			l.zapLogger.Sugar().Infof(logFormat, logArgs...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行以避免日志丢失
			l.zapLogger.Sugar().Infof(format, args...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Sugar().Infof(format, args...)
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l == nil || l.zapLogger == nil {
		fmt.Printf("WARN: "+format+"\n", args...)
		return
	}

	if l.config.AsyncMode {
		// 异步模式：将日志写入任务发送到通道
		logFormat := format
		logArgs := make([]interface{}, len(args))
		copy(logArgs, args)

		select {
		case logChan <- func() {
			l.zapLogger.Sugar().Warnf(logFormat, logArgs...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行以避免日志丢失
			l.zapLogger.Sugar().Warnf(format, args...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Sugar().Warnf(format, args...)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l == nil || l.zapLogger == nil {
		fmt.Printf("ERROR: "+format+"\n", args...)
		return
	}

	if l.config.AsyncMode {
		// 异步模式：将日志写入任务发送到通道
		logFormat := format
		logArgs := make([]interface{}, len(args))
		copy(logArgs, args)

		select {
		case logChan <- func() {
			l.zapLogger.Sugar().Errorf(logFormat, logArgs...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行以避免日志丢失
			l.zapLogger.Sugar().Errorf(format, args...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Sugar().Errorf(format, args...)
	}
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	if l == nil || l.zapLogger == nil {
		panic(fmt.Sprintf(format, args...))
	}

	// Panic级别的日志通常需要立即执行，不建议异步处理
	if l.config.AsyncMode {
		logFormat := format
		logArgs := make([]interface{}, len(args))
		copy(logArgs, args)

		select {
		case logChan <- func() {
			l.zapLogger.Sugar().Panicf(logFormat, logArgs...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行
			l.zapLogger.Sugar().Panicf(format, args...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Sugar().Panicf(format, args...)
	}
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l == nil || l.zapLogger == nil {
		fmt.Printf(format+"\n", args...)
		os.Exit(1)
	}

	// Fatal级别的日志通常需要立即执行，不建议异步处理
	if l.config.AsyncMode {
		logFormat := format
		logArgs := make([]interface{}, len(args))
		copy(logArgs, args)

		select {
		case logChan <- func() {
			l.zapLogger.Sugar().Fatalf(logFormat, logArgs...)
		}:
			// 任务已发送到通道
		default:
			// 通道已满，同步执行
			l.zapLogger.Sugar().Fatalf(format, args...)
		}
	} else {
		// 同步模式：直接写入日志
		l.zapLogger.Sugar().Fatalf(format, args...)
	}
}

// With 添加字段
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		zapLogger: l.zapLogger.With(fields...),
		config:    l.config,
	}
}

// Sync 将缓冲区刷新到磁盘
func (l *Logger) Sync() error {
	if l == nil || l.zapLogger == nil {
		return nil
	}
	// 尝试sync，但优雅地处理错误
	// 特别是当输出到stdout时，sync操作会失败，这是正常的
	if err := l.zapLogger.Sync(); err != nil {
		// 检查是否是因为stdout导致的错误
		if err.Error() == "sync /dev/stdout: inappropriate ioctl for device" ||
			err.Error() == "sync /dev/stderr: inappropriate ioctl for device" {
			// 这些错误是预期的，我们可以忽略它们
			return nil
		}
		// 对于其他错误，仍然返回
		return err
	}
	return nil
}

// GetZapLogger 获取原始的 zap logger（用于高级用法）
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.zapLogger
}
