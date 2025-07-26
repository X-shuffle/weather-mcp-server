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
				Description: "获取指定位置的天气信息，支持实时天气和未来小时预报",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]any{
						"location": map[string]any{
							"type":        "string",
							"description": "位置信息，可以是城市名（如：北京）或坐标（如：39.9042,116.4074）",
						},
						"hours": map[string]any{
							"type":        "integer",
							"description": "需要查询的小时数，0或不传表示查询实时天气，1-12表示查询未来小时预报",
							"minimum":     0,
							"maximum":     12,
							"default":     0,
						},
					},
					Required: []string{"location"},
				},
			},
			Handler: wt.handleGetWeather,
		},
	}
}

// handleGetWeather 处理天气查询请求（支持实时和小时预报）
func (wt *WeatherTools) handleGetWeather(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Location string `json:"location"`
		Hours    int    `json:"hours"`
	}

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
	if args.Hours < 0 || args.Hours > 12 {
		return nil, fmt.Errorf("hours parameter must be between 0 and 12")
	}

	// 根据 hours 参数决定查询类型
	if args.Hours == 0 {
		// 查询实时天气
		weather, err := wt.weatherService.GetWeatherByLocation(args.Location)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("❌ 获取实时天气信息失败: %s", err.Error()),
					},
				},
			}, nil
		}
		formattedResponse := wt.weatherService.FormatWeatherResponse(weather)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: formattedResponse,
				},
			},
		}, nil
	} else {
		// 查询小时级天气预报
		hourly, err := wt.weatherService.GetHourlyWeatherByLocation(args.Location, args.Hours)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("❌ 获取小时级天气预报失败: %s", err.Error()),
					},
				},
			}, nil
		}
		formattedResponse := wt.weatherService.FormatHourlyWeatherResponse(hourly)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: formattedResponse,
				},
			},
		}, nil
	}
}
