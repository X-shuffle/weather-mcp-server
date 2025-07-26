package weather

import (
	"time"
)

// Weather 天气实体
type Weather struct {
	Location    Location
	Current     CurrentWeather
	Forecast    []ForecastWeather
	LastUpdated time.Time
}

// Location 位置值对象
type Location struct {
	City    string
	Country string
	Lat     float64
	Lon     float64
}

// CurrentWeather 当前天气值对象
type CurrentWeather struct {
	Temperature float64
	FeelsLike   float64
	Humidity    int
	Pressure    int
	WindSpeed   float64
	WindDir     string
	Description string
	Icon        string
}

// ForecastWeather 预报天气值对象
type ForecastWeather struct {
	Date        time.Time
	Temperature struct {
		Min float64
		Max float64
	}
	Humidity    int
	Description string
	Icon        string
}

// WeatherRepository 天气仓储接口
type WeatherRepository interface {
	GetCurrentWeather(lat, lon float64) (*Weather, error)
	GetWeatherByCity(city string) (*Weather, error)
}

// WeatherService 天气服务接口
type WeatherService interface {
	GetCurrentWeather(lat, lon float64) (*Weather, error)
	GetWeatherByCity(city string) (*Weather, error)
}
