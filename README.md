# 漫展查询 API (Go 版本)

> 基于 Go + Gin 重构的漫展查询 API，数据来源于"无差别同人站（CPP）"

## ✨ 特性

- 🔍 支持关键词搜索漫展信息
- 📅 自动标注展会状态（进行中/倒计时/已取消）
- 📊 多页数据自动聚合
- 🎯 智能时间排序（从近到远）
- 🖼️ 完整的图片 URL 处理
- 🐛 调试模式支持
- ⚡ 高性能，低内存占用

## 📋 要求

- Go 1.20 或更高版本

## 🚀 快速开始

### 安装依赖

```bash
cd cpp_search_go
go mod tidy
```

### 运行

```bash
go run main.go
```

### 编译

```bash
# Windows
go build -o cpp_search.exe main.go

# Linux/macOS
go build -o cpp_search main.go
```

## 📖 API 文档

### 搜索漫展

```
GET /search?msg=关键词
```

#### 参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| msg | string | ✅ | 搜索关键词 |
| debug | string | ❌ | 调试模式：`raw` 返回原始数据 |

#### 响应示例

```json
{
  "code": 200,
  "msg": "上海",
  "data": [
    {
      "id": 12345,
      "name": "某某漫展(还有3天开始)",
      "tag": "综合展",
      "location": "上海 浦东新区",
      "address": "上海新国际博览中心",
      "url": "https://www.allcpp.cn/allcpp/event/event.do?event=12345",
      "type": "综合展",
      "wannaGoCount": 1234,
      "circleCount": 100,
      "doujinshiCount": 50,
      "time": "2025-01-01 09:00:00",
      "appLogoPicUrl": "https://imagecdn3.allcpp.cn/upload/xxx.jpg",
      "logoPicUrl": "https://imagecdn3.allcpp.cn/upload/xxx.jpg",
      "ended": "未结束",
      "isOnline": "线下"
    }
  ]
}
```

### 首页

```
GET /
```

返回 API 基本信息和使用说明。

## ⚙️ 配置

### 配置文件 (推荐)

创建 `config.yml` 文件：

```yaml
# CPP Search API 配置文件

# 服务器配置
host: "0.0.0.0"
port: 51225

# 调试模式
debug: false
```

程序会按以下顺序查找配置文件：
1. 当前工作目录的 `config.yml` / `config.yaml`
2. 可执行文件所在目录的 `config.yml` / `config.yaml`

### 环境变量

环境变量可以覆盖配置文件中的设置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| HOST | 0.0.0.0 | 监听地址 |
| PORT | 51225 | 监听端口 |
| DEBUG | false | 调试模式 |

示例：

```bash
# Windows PowerShell
$env:PORT=3000; go run .

# Linux/macOS
PORT=3000 go run .
```

## 📁 项目结构

```
cpp_search_go/
├── main.go                 # 入口文件
├── go.mod                  # Go 模块定义
├── go.sum                  # 依赖校验
├── README.md               # 说明文档
└── internal/
    ├── api/
    │   └── handler.go      # HTTP 处理器
    ├── config/
    │   └── config.go       # 配置管理
    ├── models/
    │   └── event.go        # 数据模型
    └── service/
        └── search.go       # 业务逻辑
```

## 📝 License

MIT License
