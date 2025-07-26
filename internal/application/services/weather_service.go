package services

import (
	"fmt"
	"strconv"
	"strings"

	"weather-mcp-server/internal/domain/weather"
)

// WeatherApplicationService 天气应用服务
type WeatherApplicationService struct {
	weatherRepo weather.WeatherRepository
}

// NewWeatherApplicationService 创建新的天气应用服务
func NewWeatherApplicationService(weatherRepo weather.WeatherRepository) *WeatherApplicationService {
	return &WeatherApplicationService{
		weatherRepo: weatherRepo,
	}
}

// GetWeatherByLocation 根据位置获取天气
func (s *WeatherApplicationService) GetWeatherByLocation(location string) (*weather.Weather, error) {
	// 检查是否是坐标格式 (lat,lon)
	if strings.Contains(location, ",") {
		coords := strings.Split(location, ",")
		if len(coords) == 2 {
			lat, err := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid latitude: %w", err)
			}
			lon, err := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid longitude: %w", err)
			}
			return s.weatherRepo.GetCurrentWeather(lat, lon)
		}
	}

	// 否则按城市名处理
	return s.weatherRepo.GetWeatherByCity(location)
}

// GetHourlyWeatherByLocation 获取未来小时天气预报
func (s *WeatherApplicationService) GetHourlyWeatherByLocation(location string, hours int) (*weather.HourlyWeatherResult, error) {
	if strings.Contains(location, ",") {
		coords := strings.Split(location, ",")
		if len(coords) == 2 {
			lat, err := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid latitude: %w", err)
			}
			lon, err := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid longitude: %w", err)
			}
			return s.weatherRepo.GetHourlyWeatherByCoords(lat, lon, hours)
		}
	}
	return s.weatherRepo.GetHourlyWeatherByCity(location, hours)
}

// FormatWeatherResponse 格式化天气响应
func (s *WeatherApplicationService) FormatWeatherResponse(w *weather.Weather) string {
	if w == nil {
		return "无法获取天气信息"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📍 %s, %s\n", w.Location.City, w.Location.Country))
	sb.WriteString(fmt.Sprintf("🌡️  温度: %.1f°C (体感: %.1f°C)\n", w.Current.Temperature, w.Current.FeelsLike))
	sb.WriteString(fmt.Sprintf("💧 湿度: %d%%\n", w.Current.Humidity))
	sb.WriteString(fmt.Sprintf("🌪️  风速: %.1f m/s (%s)\n", w.Current.WindSpeed, w.Current.WindDir))
	sb.WriteString(fmt.Sprintf("🌡️  气压: %d hPa\n", w.Current.Pressure))
	sb.WriteString(fmt.Sprintf("☁️  天气: %s\n", w.Current.Description))
	sb.WriteString(fmt.Sprintf("🕐 更新时间: %s", w.LastUpdated.Format("2006-01-02 15:04:05")))

	return sb.String()
}

// FormatHourlyWeatherResponse 格式化小时级天气响应
func (s *WeatherApplicationService) FormatHourlyWeatherResponse(hw *weather.HourlyWeatherResult) string {
	if hw == nil || len(hw.Hourly) == 0 {
		return "无法获取小时级天气预报信息"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📍 %s, %s\n", hw.Location.City, hw.Location.Country))
	for i, h := range hw.Hourly {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, h.Date.Format("2006-01-02 15:04")))
		sb.WriteString(fmt.Sprintf("  🌡️ %.1f°C (体感: %.1f°C), 💧%d%%, 🌪️ %.1fm/s(%s), ☁️ %s\n",
			h.Temperature, h.FeelsLike, h.Humidity, h.WindSpeed, h.WindDir, h.Description))
	}
	sb.WriteString(fmt.Sprintf("🕐 更新时间: %s", hw.LastUpdated.Format("2006-01-02 15:04:05")))
	return sb.String()
}
