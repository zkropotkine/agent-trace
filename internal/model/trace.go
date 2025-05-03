package model

import "time"

type SubStep struct {
	Name   string    `json:"name" bson:"name"`
	Input  string    `json:"input" bson:"input"`
	Output string    `json:"output" bson:"output"`
	Status string    `json:"status" bson:"status"`
	Start  time.Time `json:"start" bson:"start"`
	End    time.Time `json:"end" bson:"end"`
}

type TokenUsage struct {
	Input         int     `json:"input_tokens" bson:"inputTokens"`
	Output        int     `json:"output_tokens" bson:"outputTokens"`
	Total         int     `json:"total" bson:"total"`
	EstimatedCost float64 `json:"estimated_cost_usd" bson:"estimatedCostUSD"`
}

type Trace struct {
	TraceID      string     `json:"trace_id" bson:"traceId"`
	SessionID    string     `json:"session_id" bson:"sessionId"`
	AgentName    string     `json:"agent_name" bson:"agentName"`
	Timestamp    time.Time  `json:"timestamp" bson:"timestamp"`
	Status       string     `json:"status" bson:"status"`
	InputPrompt  string     `json:"input_prompt" bson:"inputPrompt"`
	OutputPrompt string     `json:"output_prompt" bson:"outputPrompt"`
	LatencyMS    int        `json:"latency_ms" bson:"latencyMs"`
	TokenUsage   TokenUsage `json:"token_usage" bson:"tokenUsage"`
	SubSteps     []SubStep  `json:"sub_steps" bson:"sub_Steps"`
	CreatedAt    time.Time  `json:"created_at" bson:"createdAt"`
	Model        string     `json:"model" bson:"model" example:"gpt-3.5-turbo"`
}
