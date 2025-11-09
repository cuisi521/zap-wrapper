// example/main.go
package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/cuisi521/zap-wrapper/logger"
)

func main() {
	fmt.Println("========== 全局日志功能测试 ==========")

	// 创建日志器，同时会自动设置为全局日志器
	_, err := logger.New(
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

	// 使用全局Sync函数确保所有日志都被写入
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("Sync失败: %v\n", err)
		}
	}()

	// 现在可以直接使用全局日志方法
	fmt.Println("\n使用全局日志方法...")
	logger.Debug("这是全局debug级别日志")
	logger.Info("这是全局info级别日志")
	logger.Warn("这是全局warn级别日志")
	logger.Error("这是全局error级别日志")

	// 使用格式化日志方法
	logger.Infof("这是带参数的全局info日志，ID: %s, 名称: %s", "12345", "测试用户")
	logger.Errorf("这是带参数的全局error日志，错误码: %d, 错误信息: %s", 500, "服务器内部错误")

	// 结构化日志示例
	logger.Info("用户登录成功",
		zap.String("user_id", "1001"),
		zap.String("username", "张三"),
		zap.String("ip", "192.168.1.100"),
	)

	// 示例: 输出少量测试日志
	fmt.Println("\n写入少量测试日志...")
	for i := 0; i < 10; i++ {
		logger.Debugf("这是debug级别的日志 %d", i)
		logger.Infof("这是info级别的日志 %d", i)
		logger.Warnf("这是warn级别的日志 %d", i)
		logger.Errorf("这是error级别的日志 %d", i)
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
