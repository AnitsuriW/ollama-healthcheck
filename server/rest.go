package server

import (
	"encoding/json"
	"net/http"
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// 动态计算预测结果（基于实际指标）
	prediction, confidence := calculateFailureProbability(req)

	response := PredictResponse{
		Prediction: prediction,
		Confidence: confidence,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 根据输入指标计算故障概率
func calculateFailureProbability(req PredictRequest) (string, float64) {
	// 定义阈值（可根据实际需求调整）
	const (
		highCPUTreshold     = 80.0 // CPU使用率 >80% 视为高风险
		highMemoryTreshold  = 85.0 // 内存使用率 >85% 视为高风险
		highLatencyTreshold = 1000 // 延迟 >1000ms 视为高风险
		maxErrorsPerMinute  = 1    // 每分钟错误数 >1 视为高风险
	)

	// 检查各项指标是否超阈值
	cpuRisk := req.CPUUsage > highCPUTreshold
	memoryRisk := req.MemoryUsage > highMemoryTreshold
	latencyRisk := req.ResponseLatencyMs > highLatencyTreshold
	errorRisk := req.ErrorsLastMinute > maxErrorsPerMinute

	// 计算风险等级
	riskFactors := 0
	if cpuRisk {
		riskFactors++
	}
	if memoryRisk {
		riskFactors++
	}
	if latencyRisk {
		riskFactors++
	}
	if errorRisk {
		riskFactors++
	}

	// 根据风险因素数量决定预测结果
	switch {
	case riskFactors >= 3:
		return "failure_likely", 0.9 // 高风险，置信度90%
	case riskFactors >= 2:
		return "failure_possible", 0.7 // 中等风险，置信度70%
	default:
		return "failure_unlikely", 0.2 // 低风险，置信度20%
	}
}
