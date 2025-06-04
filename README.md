### Ollama 健康检测

符合 RESTful 规范的 API

项目环境: Go 1.24.0+

#### 使用方法

使用 `git clone https://github.com/AnitsuriW/ollama-healthcheck.git` 克隆本项目

将本项目和 Ollama 部署在同一台主机上

（提前安装 Go 语言环境后）

在项目根目录下使用 `go run .` 来暂时启动应用

也可以运行 `go build -o ollama-healthcheck` 来构建项目 `ollama-healthcheck.exe`，随后双击启动服务

1. 通过使用 HTTP 检测 `http://localhost:8080/health` 来检测 Ollama 是否运行

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

   

2. 通过使用 HTTP 检测 `http://localhost:8080/predict-failure` 来检测运行某一模型时，系统是否可能发生故障

   - POST 请求：接收用户提交的 `json` 监控数据， CPU 使用率、内存使用率、响应延迟和最近一分钟的错误数，根据这些数据进行动态预测，返回故障可能性和置信度

     ```bash
     $ curl -X POST http://localhost:8080/predict-failure -H "Content-Type: application/json" -d '{
       "cpu_usage": 45.3,
       "memory_usage": 60.7,
       "response_latency_ms": 800,
       "errors_last_minute": 0
     }'
     ```

     ```json
     {
     	"prediction": "failure_unlikely",
     	"confidence": 0.3
     }
     ```

   - GET 请求：直接读取系统的实时资源使用情况，根据实时数据进行预测，返回故障可能性和置信度

     ```bash
     $ curl http://localhost:8080/predict-failure
     ```

     ```json
     {
     	"prediction": "failure_likely",
     	"confidence": 0.9
     }
     ```

     



### ~~原理~~

~~向 `http://localhost:11434/api/tags` 发送 HTTP 请求，如果正常返回 json 内容，则说明健康，反之，则不健康~~