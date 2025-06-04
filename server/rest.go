package server

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type HealthResponse struct {
	Healthy   bool   `json:"healthy"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type PredictRequest struct {
	CPUUsage          float64 `json:"cpu_usage"`
	MemoryUsage       float64 `json:"memory_usage"`
	ResponseLatencyMs int     `json:"response_latency_ms"`
	ErrorsLastMinute  int     `json:"errors_last_minute"`
}

type PredictResponse struct {
	Prediction string  `json:"prediction"`
	Confidence float64 `json:"confidence"`
}

func CheckOllamaHealth() (bool, string) {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false, "Cannot connect to Ollama: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "Ollama returned non-200 status: " + resp.Status
	}
	return true, "Ollama is healthy"
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	healthy, msg := CheckOllamaHealth()
	response := HealthResponse{
		Healthy:   healthy,
		Message:   msg,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func PredictFailureHandler(w http.ResponseWriter, r *http.Request) {
	var req PredictRequest
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
			return
		}
	} else {
		// 读取系统资源的逻辑
		req.CPUUsage = getSystemCPUUsage()
		req.MemoryUsage = getSystemMemoryUsage()
		req.ResponseLatencyMs = getSystemResponseLatency()
		req.ErrorsLastMinute = getSystemErrorsLastMinute()
	}

	// 动态预测逻辑
	var prediction string
	var confidence float64

	if req.CPUUsage > 80 || req.MemoryUsage > 85 {
		prediction = "failure_likely"
		confidence = 0.9
	} else if req.ResponseLatencyMs > 1000 || req.ErrorsLastMinute > 1 {
		prediction = "failure_possible"
		confidence = 0.7
	} else {
		prediction = "failure_unlikely"
		confidence = 0.3
	}

	response := PredictResponse{
		Prediction: prediction,
		Confidence: confidence,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getSystemCPUUsage() float64 {
	// 使用 Windows 的 wmic 命令动态获取 CPU 使用率
	cmd := exec.Command("wmic", "cpu", "get", "loadpercentage")
	output, err := cmd.Output()
	if err != nil {
		return 0.0 // 如果出错，返回默认值
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		value := strings.TrimSpace(lines[1])
		cpuUsage, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return cpuUsage
		}
	}
	return 0.0
}

func getSystemMemoryUsage() float64 {
	// 使用 Windows 的 wmic 命令动态获取内存使用率
	cmd := exec.Command("wmic", "OS", "get", "FreePhysicalMemory,TotalVisibleMemorySize")
	output, err := cmd.Output()
	if err != nil {
		return 0.0 // 如果出错，返回默认值
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		fields := strings.Fields(lines[1])
		if len(fields) == 2 {
			freeMemory, err1 := strconv.ParseFloat(fields[0], 64)
			totalMemory, err2 := strconv.ParseFloat(fields[1], 64)
			if err1 == nil && err2 == nil && totalMemory > 0 {
				usedMemory := totalMemory - freeMemory
				return (usedMemory / totalMemory) * 100
			}
		}
	}
	return 0.0
}

func getSystemResponseLatency() int {
	// 模拟系统响应延迟获取逻辑
	return 500
}

func getSystemErrorsLastMinute() int {
	// 模拟系统最近一分钟错误数获取逻辑
	return 0
}
