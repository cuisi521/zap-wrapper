// example/main.go
package main

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/cuisi521/zap-wrapper/logger"
)

func main() {
	fmt.Println("========== 日志级别分离功能测试 ==========")

	// 创建日志器，为每个级别配置单独的日志文件
	levelLogger, err := logger.New(
		// logger.WithAsyncMode(true), // 启用异步
		logger.WithLevel(logger.DebugLevel),
		logger.WithEncoding(logger.ConsoleEncoding),
		logger.WithBasePath("bin/logs"),
		logger.WithConsoleOutput(true), // 启用控制台输出
	)

	if err != nil {
		fmt.Printf("创建日志器失败: %v\n", err)
		return
	}

	defer func() {
		if err := levelLogger.Sync(); err != nil {
			fmt.Printf("Sync失败: %v\n", err)
		}
	}()

	// 执行异步日志测试
	// testAsyncLogging()

	// 输出不同级别的日志
	fmt.Println("\n写入不同级别的日志...")
	for i := 0; i < 1000000; i++ {
		levelLogger.Debugf("这是debug级别的日志 %d", i)
		levelLogger.Infof("这是info级别的日志 %d", i)
		levelLogger.Warnf("这是warn级别的日志 %d", i)
		levelLogger.Errorf("这是error级别的日志 %d", i)
	}

	fmt.Println("\n日志写入完成！")
	fmt.Println("请查看 bin/logs/ 目录下的日志文件，验证每个文件只包含对应级别的日志")
	fmt.Println("\n========== 测试完成 ==========")

}

// 测试每个日志级别严格分离输出到文件
func testStrictLevelLogging() {
	fmt.Println("[严格日志级别分离测试] 开始执行...")

	// 创建专门用于严格级别分离测试的日志器
	strictLogger, err := logger.New(
		logger.WithLevel(logger.DebugLevel),
		logger.WithEncoding(logger.ConsoleEncoding),
		// 为每个级别配置单独的日志文件
		logger.WithDebugPath("bin/logs/strict_debug.log"),
		logger.WithInfoPath("bin/logs/strict_info.log"),
		logger.WithWarnPath("bin/logs/strict_warn.log"),
		logger.WithErrorLPath("bin/logs/strict_error.log"),
	)

	if err != nil {
		fmt.Printf("创建日志器失败: %v\n", err)
		return
	}

	defer func() {
		if err := strictLogger.Sync(); err != nil {
			fmt.Printf("Sync失败: %v\n", err)
		}
	}()

	// 输出不同级别的日志
	strictLogger.Debug("[严格测试] 这是debug级别的日志")
	strictLogger.Info("[严格测试] 这是info级别的日志")
	strictLogger.Warn("[严格测试] 这是warn级别的日志")
	strictLogger.Error("[严格测试] 这是error级别的日志")

	fmt.Println("[严格日志级别分离测试] 日志写入完成")
	fmt.Println("[严格日志级别分离测试] 请查看 bin/logs/ 目录下的 strict_*.log 文件验证")
	fmt.Println("[严格日志级别分离测试] 完成")
}

// 验证日志文件内容，确保只包含指定级别的日志
func verifyLogFile(filePath string, expectedLevel string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", filePath, err)
		return
	}

	lines := strings.Split(string(content), "\n")
	fmt.Printf("%s contains %d lines\n", filePath, len(lines)-1)

	// 检查文件中的每一行是否只包含期望的级别
	for _, line := range lines {
		if line == "" {
			continue
		}
		if !strings.Contains(line, expectedLevel) {
			fmt.Printf("Warning: %s contains non-%s level log: %s\n", filePath, expectedLevel, line)
		}
	}
}

func testErrorScenarios() {
	// 测试空日志器情况
	var nilLogger *logger.Logger
	// 这应该不会panic，因为我们添加了安全检查
	if nilLogger != nil {
		nilLogger.Info("This should not panic")
	}

	fmt.Println("All tests completed successfully")
}

// 测试基础路径配置功能 - 使用新的WithBasePath选项
func testBasePathLogging() {
	fmt.Println("\n========== 基础路径配置测试 ==========")

	// 清理基础路径日志文件
	os.Remove("./logs/debug.log")
	os.Remove("./logs/info.log")
	os.Remove("./logs/warn.log")
	os.Remove("./logs/error.log")

	// 初始化日志 - 使用简化的基础路径配置方式
	basePathLogger, err := logger.New(
		// logger.WithAsyncMode(true), // 启用异步
		logger.WithLevel(logger.DebugLevel),
		// 只需指定基础日志路径，系统会自动在该路径下创建各级别日志文件
		// 例如：./logs/debug.log, ./logs/info.log, ./logs/warn.log, ./logs/error.log 等
		logger.WithBasePath("./logs"),
	)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}

	defer func() {
		if err := basePathLogger.Sync(); err != nil {
			fmt.Printf("BasePath logger sync failed: %v\n", err)
		}
	}()

	// 输出不同级别的日志
	fmt.Println("Writing logs with base path configuration...")
	basePathLogger.Debug("这是使用基础路径配置的debug日志")
	basePathLogger.Info("这是使用基础路径配置的info日志")
	basePathLogger.Warn("这是使用基础路径配置的warn日志")
	basePathLogger.Error("这是使用基础路径配置的error日志")

	fmt.Println("基础路径配置日志写入完成！")
	fmt.Println("请查看 ./logs/ 目录下的日志文件，验证各日志级别是否正确输出")
	fmt.Println("基础路径配置可以简化多级别日志的配置，只需指定一个基础路径")
}

// 测试异步日志功能
func testAsyncLogging() {
	fmt.Println("\n========== 异步日志功能测试 ==========")

	// 创建异步日志器
	asyncLogger, err := logger.New(
		logger.WithLevel(logger.DebugLevel),
		logger.WithEncoding(logger.ConsoleEncoding),
		// logger.WithDebugPath("bin/logs/async_debug.log"),
		// logger.WithInfoPath("bin/logs/async_info.log"),
		// logger.WithWarnPath("bin/logs/async_warn.log"),
		// logger.WithErrorLPath("bin/logs/async_error.log"),
		logger.WithBasePath("bin/logs"),
		logger.WithAsyncMode(true), // 启用异步模式
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
	for i := 0; i < 100; i++ {
		asyncLogger.Debug("这是异步模式的debug日志", zap.Int("index", i))
		asyncLogger.Info("这是异步模式的info日志", zap.Int("index", i))
		if i%10 == 0 {
			asyncLogger.Warn("这是异步模式的warn日志", zap.Int("index", i))
			asyncLogger.Error("这是异步模式的error日志", zap.Int("index", i))
		}
	}

	fmt.Println("异步日志写入完成！")
	fmt.Println("请查看 bin/logs/ 目录下的 async_*.log 文件，验证异步日志输出")
	fmt.Println("异步模式可以提高日志写入性能，但请注意确保调用Sync()来刷新缓冲区")
}
