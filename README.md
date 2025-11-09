# Zap Wrapper

一个基于 [Uber Zap](https://github.com/uber-go/zap) 的 Go 日志库封装，提供更友好的 API 和丰富的配置选项。本库在保留 Zap 高性能特性的同时，简化了使用流程，并增强了日志分级、异步处理等功能。

## 特性

- ✅ 基于高性能的 Zap 日志库，保持极高的吞吐量和极低的延迟
- ✅ 支持 JSON 和 Console 两种输出格式
- ✅ 支持日志文件自动轮转（大小限制、保留时间、压缩归档）
- ✅ 支持日志级别分离，每个级别可输出到不同文件
- ✅ 支持异步日志模式，提高应用性能
- ✅ 结构化日志和格式化日志双模式支持
- ✅ 开发模式与生产模式快速切换
- ✅ 线程安全设计，适用于高并发场景
- ✅ 自动创建日志目录，无需手动管理
- ✅ 支持调用者信息和堆栈跟踪
- ✅ 支持多目标输出（同时输出到控制台和文件）
- ✅ 友好的API设计，使用简单直观

## 安装

```bash
go get github.com/cuisi521/zap-wrapper
```

## 基本用法

### 使用全局日志器（推荐）

全局日志器是最简便的使用方式，只需初始化一次，然后在代码的任何地方都可以直接使用。

```go
import (
    "github.com/cuisi521/zap-wrapper/logger"
)

func main() {
    // 初始化全局日志器（只需调用一次）
    _, err := logger.NewDefault()
    if err != nil {
        panic(err)
    }
    defer logger.Sync() // 确保所有日志都被写入
    
    // 在任何地方直接使用全局日志方法
    logger.Info("这是一条信息日志")
    logger.Errorf("这是一条错误日志，错误码: %d", 500)
    logger.Warn("警告信息")
    logger.Debug("调试信息")
}
```

### 自定义配置的全局日志器

```go
import (
    "github.com/cuisi521/zap-wrapper/logger"
)

func main() {
    // 初始化自定义配置的全局日志器
    _, err := logger.New(
        logger.WithLevel(logger.DebugLevel),
        logger.WithEncoding(logger.ConsoleEncoding),
        logger.WithOutputPath("logs/app.log"),
        logger.WithErrorPath("logs/error.log"),
        logger.WithConsoleOutput(true),
        logger.WithAsyncMode(true),
        logger.WithFileRotation(100, 10, 7, true), // 100MB, 保留10个备份，7天，启用压缩
    )
    if err != nil {
        panic(err)
    }
    defer logger.Sync() // 全局同步函数
    
    // 直接使用全局日志方法
    logger.Debug("调试信息")
    logger.Info("普通信息")
    logger.Warn("警告信息")
    logger.Error("错误信息")
    logger.Infof("带参数的信息: %s, %d", "测试", 123)
}
```

### 传统实例方式（备用）

如果你需要在不同地方使用不同配置的日志器，可以使用传统的实例方式：

```go
import (
    "github.com/cuisi521/zap-wrapper/logger"
)

func main() {
    // 创建默认日志实例
    log, err := logger.NewDefault()
    if err != nil {
        panic(err)
    }
    defer log.Sync()
    
    // 使用日志实例
    log.Info("这是一条信息日志")
    log.Errorf("这是一条错误日志，错误码: %d", 500)
}
```
```

## 高级特性

### 日志级别分离

可以为每个日志级别配置单独的输出文件，初始化后直接使用全局日志方法：

```go
import (
    "github.com/cuisi521/zap-wrapper/logger"
)

func main() {
    // 初始化时配置日志级别分离
    _, err := logger.New(
        logger.WithLevel(logger.DebugLevel),
        logger.WithEncoding(logger.ConsoleEncoding),
        logger.WithBasePath("logs"), // 会自动生成 logs/debug.log, logs/info.log 等文件
        logger.WithConsoleOutput(true),
    )
    if err != nil {
        panic(err)
    }
    defer logger.Sync()
    
    // 直接使用全局日志方法，日志会自动写入到对应的级别文件
    logger.Debug("这条日志会写入debug.log")
    logger.Info("这条日志会写入info.log")
    logger.Warn("这条日志会写入warn.log")
    logger.Error("这条日志会写入error_l.log")
}

// 另一个函数中直接使用
func anotherFunction() {
    // 无需再次初始化，直接使用
    logger.Info("在另一个函数中使用全局日志")
}
```

### 异步日志

启用异步模式可以提高应用性能，特别是在高并发场景：

```go
import (
    "github.com/cuisi521/zap-wrapper/logger"
)

func main() {
    // 初始化时启用异步日志模式
    _, err := logger.New(
        logger.WithAsyncMode(true),
        logger.WithOutputPath("logs/app.log"),
        logger.WithLevel(logger.InfoLevel),
    )
    if err != nil {
        panic(err)
    }
    defer logger.Sync() // 重要：确保异步日志都被写入
    
    // 即使在高并发情况下也能高效处理
    for i := 0; i < 100000; i++ {
        logger.Infof("处理请求 %d", i) // 直接使用全局方法
    }
}
```

### 结构化日志

支持添加结构化字段，便于日志分析和查询：

```go
import (
    "github.com/cuisi521/zap-wrapper/logger"
    "go.uber.org/zap"
)

func main() {
    // 初始化全局日志器
    _, _ = logger.New(
        logger.WithLevel(logger.InfoLevel),
        logger.WithEncoding(logger.JSONEncoding), // JSON格式更适合结构化日志
    )
    
    // 全局添加结构化字段
    logger.Info("用户登录成功", 
        zap.String("user_id", "123456"),
        zap.String("username", "张三"),
        zap.String("ip", "192.168.1.1"),
        zap.Int("status_code", 200),
    )
    
    // 输出: {"level":"info","user_id":"123456","username":"张三","ip":"192.168.1.1","status_code":200,"message":"用户登录成功"}
}

// 创建带上下文的日志实例
func processUserRequest(userID string, username string) {
    // 创建带有用户信息的日志实例
    userLogger := logger.With(
        zap.String("user_id", userID),
        zap.String("username", username),
    )
    
    // 使用带上下文的日志实例
    userLogger.Info("开始处理用户请求")
    // 处理逻辑...
    userLogger.Info("用户请求处理完成")
}
```

## API 概述

### 全局日志函数（推荐使用）

#### 日志级别函数

- **logger.Debug(msg string, fields ...zap.Field)** - 全局Debug级别日志
- **logger.Debugf(format string, args ...interface{})** - 格式化的全局Debug级别日志
- **logger.Info(msg string, fields ...zap.Field)** - 全局Info级别日志
- **logger.Infof(format string, args ...interface{})** - 格式化的全局Info级别日志
- **logger.Warn(msg string, fields ...zap.Field)** - 全局Warn级别日志
- **logger.Warnf(format string, args ...interface{})** - 格式化的全局Warn级别日志
- **logger.Error(msg string, fields ...zap.Field)** - 全局Error级别日志
- **logger.Errorf(format string, args ...interface{})** - 格式化的全局Error级别日志
- **logger.Panic(msg string, fields ...zap.Field)** - 全局Panic级别日志（会导致程序中断）
- **logger.Panicf(format string, args ...interface{})** - 格式化的全局Panic级别日志
- **logger.Fatal(msg string, fields ...zap.Field)** - 全局Fatal级别日志（记录后程序退出）
- **logger.Fatalf(format string, args ...interface{})** - 格式化的全局Fatal级别日志

#### 全局辅助函数

- **logger.With(fields ...zap.Field)** - 创建带有结构化字段的日志实例
- **logger.Sync()** - 同步全局日志器，将缓冲区内容写入磁盘
- **logger.L()** - 获取全局日志器实例

### 日志器初始化函数

- **logger.New(options ...Option)** - 创建日志器实例并设置为全局日志器
- **logger.NewDefault()** - 创建默认配置的日志器实例并设置为全局日志器
- **logger.InitGlobal(options ...Option)** - 仅初始化全局日志器，不返回实例

### 配置选项

- **logger.WithLevel(level)** - 设置日志级别（DebugLevel, InfoLevel, WarnLevel, ErrorLevel等）
- **logger.WithEncoding(encoding)** - 设置输出格式（JSONEncoding 或 ConsoleEncoding）
- **logger.WithOutputPath(path)** - 设置主日志输出路径
- **logger.WithErrorPath(path)** - 设置错误日志路径
- **logger.WithBasePath(path)** - 设置基础路径，自动生成各级别日志文件
- **logger.WithDebugPath/InfoPath/WarnPath/ErrorLPath/PanicPath/FatalPath(path)** - 设置特定级别日志路径
- **logger.WithFileRotation(maxSize, maxBackups, maxAge, compress)** - 配置文件轮转参数
- **logger.WithCaller(show)** - 是否显示调用者信息
- **logger.WithStacktrace(enable)** - 是否启用堆栈跟踪
- **logger.WithDevelopment(dev)** - 是否启用开发模式
- **logger.WithAsyncMode(async)** - 是否启用异步日志模式
- **logger.WithConsoleOutput(enable)** - 是否同时输出到控制台

### 实例方法（传统方式）

如果使用日志实例而不是全局函数，以下是可用的实例方法：

- **log.Debug/Info/Warn/Error/Panic/Fatal** - 实例日志方法
- **log.Debugf/Infof/Warnf/Errorf/Panicf/Fatalf** - 实例格式化日志方法
- **log.With(fields...)** - 为实例添加结构化字段
- **log.Sync()** - 同步实例日志缓冲区
- **log.GetZapLogger()** - 获取原始zap logger实例（高级用法）

## 注意事项

1. **始终调用Sync()**: 使用defer确保所有日志都被写入，尤其是异步模式
2. **目录权限**: 确保应用有足够权限创建和写入日志目录
3. **异步模式考虑**: 异步模式下日志可能会有延迟，重要日志考虑使用同步模式
4. **Panic/Fatal级别**: 这些级别会中断程序执行，请谨慎使用
5. **性能调优**: 对于极高并发场景，可调整异步工作协程数和通道大小

## 性能特点

- 基于Uber Zap，保持高性能特性
- 异步模式可进一步提高吞吐量
- 支持配置文件轮转，防止日志文件过大
- 优化的空指针检查和错误处理

## 开发与生产环境

开发环境推荐配置：
```go
logger.New(
    logger.WithLevel(logger.DebugLevel),
    logger.WithEncoding(logger.ConsoleEncoding),
    logger.WithDevelopment(true),
    logger.WithConsoleOutput(true),
)
```

生产环境推荐配置：
```go
logger.New(
    logger.WithLevel(logger.InfoLevel),
    logger.WithEncoding(logger.JSONEncoding),
    logger.WithOutputPath("/var/log/app.log"),
    logger.WithErrorPath("/var/log/error.log"),
    logger.WithAsyncMode(true),
    logger.WithFileRotation(100, 30, 15, true),
)
```
