package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"weather-mcp-server/internal/domain/weather"
)

// OpenWeatherClient OpenWeatherMap API客户端
type OpenWeatherClient struct {
	apiKey      string
	client      *http.Client
	baseURL     string
	cityMapping *CityMapping
}

// NewOpenWeatherClient 创建新的OpenWeatherMap客户端
func NewOpenWeatherClient(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey:      apiKey,
		client:      &http.Client{Timeout: 10 * time.Second},
		baseURL:     "https://api.openweathermap.org/data/2.5",
		cityMapping: NewCityMapping(),
	}
}

// GetCurrentWeather 获取当前天气
func (c *OpenWeatherClient) GetCurrentWeather(lat, lon float64) (*weather.Weather, error) {
	params := url.Values{}
	params.Add("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	params.Add("lon", strconv.FormatFloat(lon, 'f', -1, 64))
	params.Add("appid", c.apiKey)
	params.Add("units", "metric")
	params.Add("lang", "zh_cn")

	resp, err := c.client.Get(fmt.Sprintf("%s/weather?%s", c.baseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var apiResp OpenWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return c.convertToWeather(&apiResp), nil
}

// GetWeatherByCity 根据城市名获取天气
func (c *OpenWeatherClient) GetWeatherByCity(city string) (*weather.Weather, error) {
	// 检查是否为中文城市名，如果是则转换为英文
	queryCity := city
	if c.cityMapping.IsChineseCity(city) {
		if englishName, exists := c.cityMapping.GetEnglishName(city); exists {
			queryCity = englishName
		}
	}

	params := url.Values{}
	params.Add("q", queryCity)
	params.Add("appid", c.apiKey)
	params.Add("units", "metric")
	params.Add("lang", "zh_cn")

	resp, err := c.client.Get(fmt.Sprintf("%s/weather?%s", c.baseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var apiResp OpenWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return c.convertToWeather(&apiResp), nil
}

// OpenWeatherResponse OpenWeatherMap API响应结构
type OpenWeatherResponse struct {
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
	Name string `json:"name"`
	Dt   int64  `json:"dt"`
}

// ForecastAPIResponse OpenWeatherMap 预报API响应结构
type ForecastAPIResponse struct {
	City struct {
		Name    string `json:"name"`
		Country string `json:"country"`
		Coord   struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
	} `json:"city"`
	List []struct {
		Dt   int64 `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		} `json:"wind"`
	} `json:"list"`
}

// convertToWeather 将API响应转换为领域模型
func (c *OpenWeatherClient) convertToWeather(resp *OpenWeatherResponse) *weather.Weather {
	var description, icon string
	if len(resp.Weather) > 0 {
		description = resp.Weather[0].Description
		icon = resp.Weather[0].Icon
	}

	windDir := getWindDirection(resp.Wind.Deg)

	return &weather.Weather{
		Location: weather.Location{
			City:    resp.Name,
			Country: resp.Sys.Country,
			Lat:     resp.Coord.Lat,
			Lon:     resp.Coord.Lon,
		},
		Current: weather.CurrentWeather{
			Temperature: resp.Main.Temp,
			FeelsLike:   resp.Main.FeelsLike,
			Humidity:    resp.Main.Humidity,
			Pressure:    resp.Main.Pressure,
			WindSpeed:   resp.Wind.Speed,
			WindDir:     windDir,
			Description: description,
			Icon:        icon,
		},
		LastUpdated: time.Unix(resp.Dt, 0),
	}
}

// GetHourlyWeatherByCoords 获取未来小时天气预报（经纬度）
func (c *OpenWeatherClient) GetHourlyWeatherByCoords(lat, lon float64, hours int) (*weather.HourlyWeatherResult, error) {
	params := url.Values{}
	params.Add("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	params.Add("lon", strconv.FormatFloat(lon, 'f', -1, 64))
	params.Add("appid", c.apiKey)
	params.Add("units", "metric")
	params.Add("lang", "zh_cn")

	resp, err := c.client.Get(fmt.Sprintf("%s/forecast?%s", c.baseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var apiResp ForecastAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	hourly := make([]weather.HourlyWeather, 0, hours)
	for i, item := range apiResp.List {
		if i >= hours {
			break
		}
		var desc, icon string
		if len(item.Weather) > 0 {
			desc = item.Weather[0].Description
			icon = item.Weather[0].Icon
		}
		windDir := getWindDirection(item.Wind.Deg)
		hourly = append(hourly, weather.HourlyWeather{
			Date:        time.Unix(item.Dt, 0),
			Temperature: item.Main.Temp,
			FeelsLike:   item.Main.FeelsLike,
			Humidity:    item.Main.Humidity,
			Pressure:    item.Main.Pressure,
			WindSpeed:   item.Wind.Speed,
			WindDir:     windDir,
			Description: desc,
			Icon:        icon,
		})
	}

	return &weather.HourlyWeatherResult{
		Location: weather.Location{
			City:    apiResp.City.Name,
			Country: apiResp.City.Country,
			Lat:     apiResp.City.Coord.Lat,
			Lon:     apiResp.City.Coord.Lon,
		},
		Hourly:      hourly,
		LastUpdated: time.Now(),
	}, nil
}

// GetHourlyWeatherByCity 获取未来小时天气预报（城市名）
func (c *OpenWeatherClient) GetHourlyWeatherByCity(city string, hours int) (*weather.HourlyWeatherResult, error) {
	// 检查是否为中文城市名，如果是则转换为英文
	queryCity := city
	if c.cityMapping.IsChineseCity(city) {
		if englishName, exists := c.cityMapping.GetEnglishName(city); exists {
			queryCity = englishName
		}
	}
	params := url.Values{}
	params.Add("q", queryCity)
	params.Add("appid", c.apiKey)
	params.Add("units", "metric")
	params.Add("lang", "zh_cn")

	resp, err := c.client.Get(fmt.Sprintf("%s/forecast?%s", c.baseURL, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var apiResp ForecastAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	hourly := make([]weather.HourlyWeather, 0, hours)
	for i, item := range apiResp.List {
		if i >= hours {
			break
		}
		var desc, icon string
		if len(item.Weather) > 0 {
			desc = item.Weather[0].Description
			icon = item.Weather[0].Icon
		}
		windDir := getWindDirection(item.Wind.Deg)
		hourly = append(hourly, weather.HourlyWeather{
			Date:        time.Unix(item.Dt, 0),
			Temperature: item.Main.Temp,
			FeelsLike:   item.Main.FeelsLike,
			Humidity:    item.Main.Humidity,
			Pressure:    item.Main.Pressure,
			WindSpeed:   item.Wind.Speed,
			WindDir:     windDir,
			Description: desc,
			Icon:        icon,
		})
	}

	return &weather.HourlyWeatherResult{
		Location: weather.Location{
			City:    apiResp.City.Name,
			Country: apiResp.City.Country,
			Lat:     apiResp.City.Coord.Lat,
			Lon:     apiResp.City.Coord.Lon,
		},
		Hourly:      hourly,
		LastUpdated: time.Now(),
	}, nil
}

// getWindDirection 根据角度获取风向
func getWindDirection(deg int) string {
	directions := []string{"北", "东北", "东", "东南", "南", "西南", "西", "西北"}
	index := (deg + 22) / 45 % 8
	return directions[index]
}
