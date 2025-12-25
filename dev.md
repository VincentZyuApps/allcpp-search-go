# 开发笔记

## 依赖管理

```shell
go mod tidy
```

## 编译

### Python 编译脚本（推荐）

```shell
# 编译 Windows x64 和 Linux x64
python build.py --version v1.0.0

# 指定输出目录
python build.py --version v1.0.0 --output dist
```

### 手动编译

#### PowerShell
```powershell
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o allcpp_search_windows_amd64.exe .
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o allcpp_search_linux_amd64 .
```

#### Bash
```bash
# Windows x64
GOOS=windows GOARCH=amd64 go build -o allcpp_search_windows_amd64.exe .
# Linux x64
GOOS=linux GOARCH=amd64 go build -o allcpp_search_linux_amd64 .
```

## 运行

```shell
# Windows
.\allcpp_search_windows_amd64.exe

# Linux
chmod +x allcpp_search_linux_amd64
./allcpp_search_linux_amd64
```
