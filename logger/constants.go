package logger

// 日志级别
type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	PanicLevel Level = "panic"
	FatalLevel Level = "fatal"
)

// 编码格式
type Encoding string

const (
	JSONEncoding    Encoding = "json"
	ConsoleEncoding Encoding = "console"
)

// 默认配置
const (
	DefaultLevel      = InfoLevel
	DefaultEncoding   = JSONEncoding
	DefaultMaxSize    = 100 // MB
	DefaultMaxBackups = 3
	DefaultMaxAge     = 30 // days
	DefaultCompress   = true
	DefaultShowCaller = true
	DefaultStacktrace = false
	DefaultAsyncMode  = false // 默认不使用异步模式，保持向后兼容
	DefaultConsoleOutput = true // 默认输出到控制台
)
