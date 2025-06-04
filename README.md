### Ollama 健康检测

符合 RESTful 规范的 API

项目环境: Go 1.24.0+

#### 使用方法

使用 `git clone https://github.com/AnitsuriW/ollama-healthcheck.git` 克隆本项目

将本项目和 Ollama 部署在同一台主机上

（提前安装 Go 语言环境后）

在项目根目录下使用 `go run .` 来暂时启动应用

也可以运行 `go build -o ollama-healthcheck` 来构建项目 `ollama-healthcheck.exe`，随后双击启动服务

通过使用 HTTP 检测 `http://localhost:8080/health` 来检测 Ollama 是否运行

```bash
# 健康的情况
$ curl http://localhost:8080/health
{
	"healthy": true,
	"message": "Ollama is healthy",
	"timestamp": "2025-06-04T10:10:01+08:00"
}

# 不健康的情况
$ curl http://localhost:8080/health
{
	"healthy": false,
	"message": "Cannot connect to Ollama: Get \"http://localhost:11434/api/tags\": dial tcp 127.0.0.1:11434: connect: connection refused",
	"timestamp": "2025-06-04T09:42:18+08:00"
}
```



### ~~原理~~

~~向 `http://localhost:11434/api/tags` 发送 HTTP 请求，如果正常返回 json 内容，则说明健康，反之，则不健康~~