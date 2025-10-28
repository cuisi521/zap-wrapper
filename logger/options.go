package logger

// Option 配置选项
type Option func(*Config)

// Config 日志配置
type Config struct {
	Level      Level    `json:"level" yaml:"level"`
	Encoding   Encoding `json:"encoding" yaml:"encoding"`
	OutputPath string   `json:"output_path" yaml:"output_path"`
	ErrorPath  string   `json:"error_path" yaml:"error_path"`
	// 基础日志路径，如果设置了此路径且未明确指定各个级别路径，将自动生成各级别日志文件
	BasePath string `json:"base_path" yaml:"base_path"`
	// 各个级别的日志文件路径
	DebugPath string `json:"debug_path" yaml:"debug_path"`
	InfoPath  string `json:"info_path" yaml:"info_path"`
	WarnPath  string `json:"warn_path" yaml:"warn_path"`
	ErrorLPath string `json:"error_l_path" yaml:"error_l_path"`
	PanicPath string `json:"panic_path" yaml:"panic_path"`
	FatalPath string `json:"fatal_path" yaml:"fatal_path"`

	// 文件轮转配置
	MaxSize    int  `json:"max_size" yaml:"max_size"`
	MaxBackups int  `json:"max_backups" yaml:"max_backups"`
	MaxAge     int  `json:"max_age" yaml:"max_age"`
	Compress   bool `json:"compress" yaml:"compress"`

	// 其他配置
	ShowCaller    bool `json:"show_caller" yaml:"show_caller"`
	Stacktrace    bool `json:"stacktrace" yaml:"stacktrace"`
	Development   bool `json:"development" yaml:"development"`
	AsyncMode     bool `json:"async_mode" yaml:"async_mode"` // 异步日志模式
	ConsoleOutput bool `json:"console_output" yaml:"console_output"` // 是否输出到控制台
}

// WithLevel 设置日志级别
func WithLevel(level Level) Option {
	return func(c *Config) {
		c.Level = level
	}
}

// WithEncoding 设置编码格式
func WithEncoding(encoding Encoding) Option {
	return func(c *Config) {
		c.Encoding = encoding
	}
}

// WithOutputPath 设置输出路径
func WithOutputPath(path string) Option {
	return func(c *Config) {
		c.OutputPath = path
	}
}

// WithErrorPath 设置错误日志路径
func WithErrorPath(path string) Option {
	return func(c *Config) {
		c.ErrorPath = path
	}
}

// WithDebugPath 设置Debug级别日志路径
func WithDebugPath(path string) Option {
	return func(c *Config) {
		c.DebugPath = path
	}
}

// WithInfoPath 设置Info级别日志路径
func WithInfoPath(path string) Option {
	return func(c *Config) {
		c.InfoPath = path
	}
}

// WithWarnPath 设置Warn级别日志路径
func WithWarnPath(path string) Option {
	return func(c *Config) {
		c.WarnPath = path
	}
}

// WithErrorLPath 设置Error级别日志路径
func WithErrorLPath(path string) Option {
	return func(c *Config) {
		c.ErrorLPath = path
	}
}

// WithPanicPath 设置Panic级别日志路径
func WithPanicPath(path string) Option {
	return func(c *Config) {
		c.PanicPath = path
	}
}

// WithFatalPath 设置Fatal级别日志路径
func WithFatalPath(path string) Option {
	return func(c *Config) {
		c.FatalPath = path
	}
}

// WithFileRotation 设置文件轮转配置
func WithFileRotation(maxSize, maxBackups, maxAge int, compress bool) Option {
	return func(c *Config) {
		c.MaxSize = maxSize
		c.MaxBackups = maxBackups
		c.MaxAge = maxAge
		c.Compress = compress
	}
}

// WithCaller 设置是否显示调用者信息
func WithCaller(show bool) Option {
	return func(c *Config) {
		c.ShowCaller = show
	}
}

// WithStacktrace 设置是否记录堆栈跟踪
func WithStacktrace(enable bool) Option {
	return func(c *Config) {
		c.Stacktrace = enable
	}
}

// WithDevelopment 设置开发模式
func WithDevelopment(dev bool) Option {
	return func(c *Config) {
		c.Development = dev
	}
}

// WithAsyncMode 设置是否启用异步日志模式
func WithAsyncMode(async bool) Option {
	return func(c *Config) {
		c.AsyncMode = async
	}
}

// WithBasePath 设置基础日志路径，系统会自动在该路径下为各个级别生成对应的日志文件
// 例如：如果设置为"./logs"，则会自动生成"./logs/debug.log", "./logs/info.log"等
func WithBasePath(path string) Option {
	return func(c *Config) {
		c.BasePath = path
	}
}

// WithConsoleOutput 设置是否输出日志到控制台
func WithConsoleOutput(enable bool) Option {
	return func(c *Config) {
		c.ConsoleOutput = enable
	}
}
