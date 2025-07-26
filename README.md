# Weather MCP Server

一个基于Go语言开发的天气查询MCP（Model Context Protocol）服务器，提供实时天气信息查询功能。

## 功能特性

- 🌤️ 实时天气查询
- 📍 支持城市名和坐标查询
- 🌍 多语言支持（中文）
- 🇨🇳 支持中文城市名和区级地名查询
- 🔧 基于MCP协议，易于集成
- 🏗️ 遵循领域驱动设计（DDD）原则

## 技术栈

- **语言**: Go 1.21+
- **MCP库**: github.com/mark3labs/mcp-go
- **天气API**: OpenWeatherMap
- **架构**: 领域驱动设计（DDD）

## 快速开始

### 1. 获取OpenWeatherMap API密钥

1. 访问 [OpenWeatherMap](https://openweathermap.org/api)
2. 注册免费账户
3. 获取API密钥

### 2. 克隆项目

```bash
git clone <repository-url>
cd weather-mcp-server
```

### 3. 设置环境变量

```bash
export OPENWEATHER_API_KEY=your_api_key_here
```

### 4. 运行服务器

#### 方式一：直接运行（推荐开发使用）

```bash
go run cmd/server/main.go
```

#### 方式二：构建后运行

```bash
# 构建项目
go build -o bin/weather-mcp-server cmd/server/main.go

# 运行服务器
./bin/weather-mcp-server
```

### 5. 测试功能

```bash
# 运行测试脚本
go run test_weather.go
```

## 项目结构

```
weather-mcp-server/
├── cmd/server/              # 应用入口
│   └── main.go
├── internal/                # 内部包
│   ├── domain/weather/      # 天气领域模型
│   ├── application/services/ # 应用服务
│   ├── infrastructure/      # 基础设施
│   │   ├── weather/         # 天气API客户端
│   │   └── mcp/             # MCP协议实现
│   └── interfaces/          # 接口层
├── configs/                 # MCP配置文件
├── test_weather.go          # 测试脚本
└── README.md
```

## MCP工具

### get_weather

获取指定位置的天气信息，支持实时天气和未来小时预报。

**参数:**
- `location` (string, 必需): 位置信息，可以是城市名（如：北京、Beijing）或坐标（如：39.9042,116.4074）
- `hours` (integer, 可选): 需要查询的小时数，0或不传表示查询实时天气，1-12表示查询未来小时预报

**注意**: OpenWeatherMap 的预报 API 返回的是3小时间隔的数据。例如：
- 请求3小时会返回 [当前+3h, 当前+6h, 当前+9h] 的数据
- 请求6小时会返回 [当前+3h, 当前+6h] 的数据
- 请求9小时会返回 [当前+3h, 当前+6h, 当前+9h] 的数据

**示例:**
```json
// 查询实时天气
{
  "location": "北京"
}

// 查询未来3小时预报
{
  "location": "北京",
  "hours": 3
}

// 查询未来6小时预报
{
  "location": "39.9600,116.3000",
  "hours": 6
}
```

**支持的城市名格式:**
- 中文城市名：北京、上海、广州、深圳等
- 英文城市名：Beijing、Shanghai、Guangzhou、Shenzhen等
- 区级地名：北京海淀、上海浦东、广州天河等
- 坐标格式：39.9042,116.4074

**响应示例:**
```
📍 北京, CN
🌡️  温度: 25.3°C (体感: 26.1°C)
💧 湿度: 65%
🌪️  风速: 3.2 m/s (东北)
🌡️  气压: 1013 hPa
☁️  天气: 多云
🕐 更新时间: 2024-01-15 14:30:00
```

## 开发

### 运行测试

```bash
go test ./...
```

### 代码格式化

```bash
go fmt ./...
```

### 代码检查

```bash
golangci-lint run
```

### 清理构建文件

```bash
rm -rf bin/
```

## 配置

### 服务器配置

通过环境变量进行配置：

#### 环境变量

- `OPENWEATHER_API_KEY`: OpenWeatherMap API密钥（必需）

### MCP客户端配置

要使用天气MCP服务器，需要在MCP客户端中配置 `mcp_settings.json` 文件。

#### 1. 使用go run命令（开发模式）

```json
{
  "mcpServers": {
    "weather-mcp-server": {
      "command": "go",
      "args": [
        "run",
        "cmd/server/main.go"
      ],
      "env": {
        "OPENWEATHER_API_KEY": "your_openweather_api_key_here"
      },
      "cwd": "/path/to/weather-mcp-server"
    }
  }
}
```

#### 2. 使用构建的二进制文件

```json
{
  "mcpServers": {
    "weather-mcp-server": {
      "command": "./bin/weather-mcp-server",
      "args": [],
      "env": {
        "OPENWEATHER_API_KEY": "your_openweather_api_key_here"
      },
      "cwd": "/path/to/weather-mcp-server"
    }
  }
}
```

#### 3. 使用绝对路径的二进制文件

```json
{
  "mcpServers": {
    "weather-mcp-server": {
      "command": "/absolute/path/to/weather-mcp-server/bin/weather-mcp-server",
      "args": [],
      "env": {
        "OPENWEATHER_API_KEY": "your_openweather_api_key_here"
      }
    }
  }
}
```

#### 配置说明

- `command`: 要执行的命令
- `args`: 命令行参数数组
- `env`: 环境变量对象
- `cwd`: 工作目录路径（可选）

#### 配置文件位置

根据不同的MCP客户端，配置文件位置可能不同：

- **Claude Desktop**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Cursor**: `~/.cursor/mcp_settings.json`
- **其他客户端**: 请参考客户端文档

#### 示例配置文件

项目提供了配置文件示例：

- `configs/mcp_settings.json.example` - 基本配置
- `configs/mcp_settings_relative.json.example` - 相对路径配置
- `configs/mcp_settings_go_run.json.example` - Go run开发模式配置

## 许可证

[LICENSE](LICENSE)

## 贡献

欢迎提交Issue和Pull Request！

## 相关链接

- [MCP协议规范](https://modelcontextprotocol.io/)
- [OpenWeatherMap API](https://openweathermap.org/api)
- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) 
