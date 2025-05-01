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
	Input  int `json:"input_tokens" bson:"input_tokens"`
	Output int `json:"output_tokens" bson:"output_tokens"`
	Total  int `json:"total" bson:"total"`
}

type Trace struct {
	TraceID     string     `json:"trace_id" bson:"trace_id"`
	SessionID   string     `json:"session_id" bson:"session_id"`
	AgentName   string     `json:"agent_name" bson:"agent_name"`
	Timestamp   time.Time  `json:"timestamp" bson:"timestamp"`
	Status      string     `json:"status" bson:"status"`
	InputPrompt string     `json:"input_prompt" bson:"input_prompt"`
	Output      string     `json:"output" bson:"output"`
	LatencyMS   int        `json:"latency_ms" bson:"latency_ms"`
	TokenUsage  TokenUsage `json:"token_usage" bson:"token_usage"`
	SubSteps    []SubStep  `json:"substeps" bson:"substeps"`
}
