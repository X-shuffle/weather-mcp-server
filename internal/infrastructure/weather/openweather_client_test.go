package weather

import (
	"testing"
	"time"
)

func TestNewOpenWeatherClient(t *testing.T) {
	apiKey := "test_api_key"
	client := NewOpenWeatherClient(apiKey)

	if client.apiKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.apiKey)
	}

	if client.baseURL != "https://api.openweathermap.org/data/2.5" {
		t.Errorf("Expected base URL %s, got %s", "https://api.openweathermap.org/data/2.5", client.baseURL)
	}

	if client.client.Timeout != 10*time.Second {
		t.Errorf("Expected timeout %v, got %v", 10*time.Second, client.client.Timeout)
	}
}

func TestGetWindDirection(t *testing.T) {
	tests := []struct {
		deg      int
		expected string
	}{
		{0, "北"},
		{45, "东北"},
		{90, "东"},
		{135, "东南"},
		{180, "南"},
		{225, "西南"},
		{270, "西"},
		{315, "西北"},
		{360, "北"},
	}

	for _, test := range tests {
		result := getWindDirection(test.deg)
		if result != test.expected {
			t.Errorf("For degree %d, expected %s, got %s", test.deg, test.expected, result)
		}
	}
}

func TestConvertToWeather(t *testing.T) {
	client := NewOpenWeatherClient("test_key")

	resp := &OpenWeatherResponse{
		Coord: struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		}{
			Lat: 39.9042,
			Lon: 116.4074,
		},
		Weather: []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		}{
			{
				ID:          800,
				Main:        "Clear",
				Description: "晴天",
				Icon:        "01d",
			},
		},
		Main: struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		}{
			Temp:      25.3,
			FeelsLike: 26.1,
			Pressure:  1013,
			Humidity:  65,
		},
		Wind: struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		}{
			Speed: 3.2,
			Deg:   45,
		},
		Sys: struct {
			Country string `json:"country"`
		}{
			Country: "CN",
		},
		Name: "北京",
		Dt:   1642248600, // 2022-01-15 14:30:00 UTC
	}

	weather := client.convertToWeather(resp)

	// 验证位置信息
	if weather.Location.City != "北京" {
		t.Errorf("Expected city %s, got %s", "北京", weather.Location.City)
	}
	if weather.Location.Country != "CN" {
		t.Errorf("Expected country %s, got %s", "CN", weather.Location.Country)
	}
	if weather.Location.Lat != 39.9042 {
		t.Errorf("Expected lat %f, got %f", 39.9042, weather.Location.Lat)
	}
	if weather.Location.Lon != 116.4074 {
		t.Errorf("Expected lon %f, got %f", 116.4074, weather.Location.Lon)
	}

	// 验证天气信息
	if weather.Current.Temperature != 25.3 {
		t.Errorf("Expected temperature %f, got %f", 25.3, weather.Current.Temperature)
	}
	if weather.Current.FeelsLike != 26.1 {
		t.Errorf("Expected feels like %f, got %f", 26.1, weather.Current.FeelsLike)
	}
	if weather.Current.Humidity != 65 {
		t.Errorf("Expected humidity %d, got %d", 65, weather.Current.Humidity)
	}
	if weather.Current.Pressure != 1013 {
		t.Errorf("Expected pressure %d, got %d", 1013, weather.Current.Pressure)
	}
	if weather.Current.WindSpeed != 3.2 {
		t.Errorf("Expected wind speed %f, got %f", 3.2, weather.Current.WindSpeed)
	}
	if weather.Current.WindDir != "东北" {
		t.Errorf("Expected wind direction %s, got %s", "东北", weather.Current.WindDir)
	}
	if weather.Current.Description != "晴天" {
		t.Errorf("Expected description %s, got %s", "晴天", weather.Current.Description)
	}
	if weather.Current.Icon != "01d" {
		t.Errorf("Expected icon %s, got %s", "01d", weather.Current.Icon)
	}

	// 验证时间
	expectedTime := time.Unix(1642248600, 0)
	if !weather.LastUpdated.Equal(expectedTime) {
		t.Errorf("Expected time %v, got %v", expectedTime, weather.LastUpdated)
	}
}
