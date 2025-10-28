package main

import (
	"fmt"
	"testing"

	"github.com/cuisi521/zap-wrapper/logger"
)

// 测试输出每个日志级别是否严格分离到不同文件
func TestTrictLogging(test *testing.T) {
	// 创建异步日志器
	asyncLogger, err := logger.New(
		logger.WithLevel(logger.DebugLevel),
		logger.WithEncoding(logger.ConsoleEncoding),
		logger.WithBasePath("../bin/logs"),
		logger.WithConsoleOutput(true), // 启用控制台输出
		logger.WithAsyncMode(true),     // 启用异步模式

	)

	if err != nil {
		fmt.Printf("创建异步日志器失败: %v\n", err)
		return
	}

	defer func() {
		if err := asyncLogger.Sync(); err != nil {
			fmt.Printf("异步日志Sync失败: %v\n", err)
		}
	}()
	// 输出不同级别的异步日志
	fmt.Println("写入异步日志...")
	for i := 0; i < 10; i++ {
		asyncLogger.Debug(fmt.Sprintf("这是debug级别的日志 %d", i))
		asyncLogger.Info(fmt.Sprintf("这是info级别的日志 %d", i))
		asyncLogger.Warn(fmt.Sprintf("这是warn级别的日志 %d", i))
		asyncLogger.Error(fmt.Sprintf("这是error级别的日志 %d", i))
	}
	fmt.Println("异步日志写入完成！")
	fmt.Println("请查看 bin/logs/ 目录下的 async_*.log 文件，验证异步日志输出")
	fmt.Println("异步模式可以提高日志写入性能，但请注意确保调用Sync()来刷新缓冲区")
}
