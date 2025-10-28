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

### 最简单的用法

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
    defer log.Sync() // 确保所有日志都被写入
    
    // 使用日志
    log.Info("这是一条信息日志")
    log.Errorf("这是一条错误日志，错误码: %d", 500)
}
```

### 自定义配置

```go
import (
    "github.com/cuisi521/zap-wrapper/logger"
)

func main() {
    // 创建自定义日志实例
    log, err := logger.New(
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
    defer log.Sync()
    
    log.Debug("调试信息")
    log.Info("普通信息")
    log.Warn("警告信息")
    log.Error("错误信息")
}
```

## 高级特性

### 日志级别分离

可以为每个日志级别配置单独的输出文件：

```go
import (
    "github.com/cuisi521/zap-wrapper/logger"
)

func main() {
    // 方式一：使用BasePath自动生成各级别文件
    levelLogger, err := logger.New(
        logger.WithLevel(logger.DebugLevel),
        logger.WithEncoding(logger.ConsoleEncoding),
        logger.WithBasePath("logs"), // 会自动生成 logs/debug.log, logs/info.log 等文件
        logger.WithConsoleOutput(true),
    )
    
    // 方式二：手动指定每个级别的文件路径
    levelLogger, err := logger.New(
        logger.WithLevel(logger.DebugLevel),
        logger.WithDebugPath("logs/custom_debug.log"),
        logger.WithInfoPath("logs/custom_info.log"),
        logger.WithWarnPath("logs/custom_warn.log"),
        logger.WithErrorLPath("logs/custom_error.log"),
        logger.WithPanicPath("logs/custom_panic.log"),
        logger.WithFatalPath("logs/custom_fatal.log"),
    )
}
```

### 异步日志

启用异步模式可以提高应用性能，特别是在高并发场景：

```go
import (
    "github.com/cuisi521/zap-wrapper/logger"
)

func main() {
    // 启用异步日志模式
    asyncLogger, err := logger.New(
        logger.WithAsyncMode(true),
        logger.WithOutputPath("logs/app.log"),
    )
    defer asyncLogger.Sync() // 重要：确保异步日志都被写入
    
    // 即使在高并发情况下也能高效处理
    for i := 0; i < 100000; i++ {
        asyncLogger.Infof("处理请求 %d", i)
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
    log, _ := logger.New()
    
    // 添加结构化字段
    userLog := log.With(
        zap.String("user_id", "123456"),
        zap.String("username", "张三"),
    )
    
    userLog.Info("用户登录成功", zap.String("ip", "192.168.1.1"))
    // 输出: {"level":"info","user_id":"123456","username":"张三","ip":"192.168.1.1","message":"用户登录成功"}
}
```

## API 概述

### 日志级别

- `Debug/Debugf`: 调试信息
- `Info/Infof`: 普通信息
- `Warn/Warnf`: 警告信息
- `Error/Errorf`: 错误信息
- `Panic/Panicf`: 导致应用中断的严重错误
- `Fatal/Fatalf`: 致命错误，记录后程序退出

### 配置选项

- `WithLevel(level)`: 设置日志级别
- `WithEncoding(encoding)`: 设置输出格式（JSON/Console）
- `WithOutputPath(path)`: 设置日志输出路径
- `WithErrorPath(path)`: 设置错误日志路径
- `WithBasePath(path)`: 设置基础路径，自动生成各级别日志文件
- `WithDebugPath/InfoPath/WarnPath/ErrorLPath/PanicPath/FatalPath`: 设置各等级日志路径
- `WithFileRotation(maxSize, maxBackups, maxAge, compress)`: 配置文件轮转
- `WithCaller(show)`: 是否显示调用者信息
- `WithStacktrace(enable)`: 是否启用堆栈跟踪
- `WithDevelopment(dev)`: 是否启用开发模式
- `WithAsyncMode(async)`: 是否启用异步日志
- `WithConsoleOutput(enable)`: 是否输出到控制台

### 其他方法

- `With(fields...)`: 添加结构化字段
- `Sync()`: 刷新日志缓冲区到磁盘
- `GetZapLogger()`: 获取原始zap logger实例（用于高级用法）

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
