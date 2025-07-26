package main

import (
	"fmt"
	"log"
	"strings"

	"weather-mcp-server/internal/application/services"
	"weather-mcp-server/internal/infrastructure/weather"
)

func main() {
	// 使用提供的API key
	apiKey := "b8989edaccdc9b83fa6c4ef3915f5aef"

	// 创建天气客户端
	weatherClient := weather.NewOpenWeatherClient(apiKey)

	// 创建天气应用服务
	weatherService := services.NewWeatherApplicationService(weatherClient)

	fmt.Println("=== 天气查询测试 ===\n")

	// 测试1: 查询北京海淀区天气 (使用坐标)
	fmt.Println("测试1: 查询北京海淀区天气")
	haidianCoords := "39.9600,116.3000"
	fmt.Printf("坐标: %s\n", haidianCoords)

	weather, err := weatherService.GetWeatherByLocation(haidianCoords)
	if err != nil {
		log.Printf("查询海淀区天气失败: %v\n", err)
	} else {
		result := weatherService.FormatWeatherResponse(weather)
		fmt.Println(result)
	}

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 测试2: 查询北京市区天气 (使用城市名)
	fmt.Println("测试2: 查询北京市区天气")
	beijingWeather, err := weatherService.GetWeatherByLocation("Beijing")
	if err != nil {
		log.Printf("查询北京市区天气失败: %v\n", err)
	} else {
		beijingResult := weatherService.FormatWeatherResponse(beijingWeather)
		fmt.Println(beijingResult)
	}

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 测试3: 测试其他城市
	fmt.Println("测试3: 查询其他城市天气")
	cities := []string{"Shanghai", "Guangzhou", "Shenzhen"}

	for _, city := range cities {
		fmt.Printf("\n--- %s ---\n", city)
		cityWeather, err := weatherService.GetWeatherByLocation(city)
		if err != nil {
			fmt.Printf("查询失败: %v\n", err)
		} else {
			cityResult := weatherService.FormatWeatherResponse(cityWeather)
			fmt.Println(cityResult)
		}
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("使用方法:")
	fmt.Println("1. 直接运行: go run test_weather.go")
	fmt.Println("2. 修改城市名或坐标来测试不同地区")
	fmt.Println("3. 支持格式: 城市名(如'Beijing') 或 坐标(如'39.9600,116.3000')")
}
