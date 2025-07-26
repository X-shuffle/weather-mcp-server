package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"weather-mcp-server/internal/application/services"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// WeatherTools MCP天气工具
type WeatherTools struct {
	weatherService *services.WeatherApplicationService
}

// NewWeatherTools 创建新的天气工具
func NewWeatherTools(weatherService *services.WeatherApplicationService) *WeatherTools {
	return &WeatherTools{
		weatherService: weatherService,
	}
}

// GetTools 获取所有天气工具
func (wt *WeatherTools) GetTools() []server.ServerTool {
	return []server.ServerTool{
		{
			Tool: mcp.Tool{
				Name:        "get_weather",
				Description: "获取指定位置的当前天气信息",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]any{
						"location": map[string]any{
							"type":        "string",
							"description": "位置信息，可以是城市名（如：北京）或坐标（如：39.9042,116.4074）",
						},
					},
					Required: []string{"location"},
				},
			},
			Handler: wt.handleGetWeather,
		},
	}
}

// handleGetWeather 处理天气查询请求
func (wt *WeatherTools) handleGetWeather(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Location string `json:"location"`
	}

	// 将arguments转换为JSON字节
	argsBytes, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %w", err)
	}

	if err := json.Unmarshal(argsBytes, &args); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if args.Location == "" {
		return nil, fmt.Errorf("location parameter is required")
	}

	// 获取天气信息
	weather, err := wt.weatherService.GetWeatherByLocation(args.Location)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("❌ 获取天气信息失败: %s", err.Error()),
				},
			},
		}, nil
	}

	// 格式化响应
	formattedResponse := wt.weatherService.FormatWeatherResponse(weather)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: formattedResponse,
			},
		},
	}, nil
}
