package main

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"weather-mcp-server/internal/application/services"
	"weather-mcp-server/internal/infrastructure/mcp"
	"weather-mcp-server/internal/infrastructure/weather"
)

func main() {
	// 获取OpenWeatherMap API密钥
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		apiKey = "b8989edaccdc9b83fa6c4ef3915f5aef" // 使用提供的API key
	}

	// 创建天气客户端
	weatherClient := weather.NewOpenWeatherClient(apiKey)

	// 创建天气应用服务
	weatherService := services.NewWeatherApplicationService(weatherClient)

	// 创建MCP工具
	weatherTools := mcp.NewWeatherTools(weatherService)

	// 创建MCP服务器
	mcpServer := server.NewMCPServer(
		"weather-mcp-server",
		"1.0.0",
		server.WithInstructions("这是一个天气查询MCP服务器，提供实时天气信息查询功能。"),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// 注册工具
	mcpServer.AddTools(weatherTools.GetTools()...)

	log.Println("Starting weather MCP server...")

	// 启动服务器（使用标准输入输出）
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatal("Server error:", err)
	}
}
