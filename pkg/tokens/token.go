package tokens

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

type TokenInfo struct {
	InputTokens   int
	OutputTokens  int
	TotalTokens   int
	EstimatedCost float64
}

var modelPricing = map[string]struct {
	InputCostPer1K  float64
	OutputCostPer1K float64
}{
	// OpenAI
	"gpt-4o":                 {InputCostPer1K: 0.005, OutputCostPer1K: 0.015},
	"gpt-4o-mini":            {InputCostPer1K: 0.00015, OutputCostPer1K: 0.0006},
	"gpt-4-turbo":            {InputCostPer1K: 0.01, OutputCostPer1K: 0.03},
	"gpt-4":                  {InputCostPer1K: 0.03, OutputCostPer1K: 0.06},
	"gpt-4-32k":              {InputCostPer1K: 0.06, OutputCostPer1K: 0.12},
	"gpt-3.5-turbo-0125":     {InputCostPer1K: 0.0005, OutputCostPer1K: 0.0015},
	"gpt-3.5-turbo-1106":     {InputCostPer1K: 0.001, OutputCostPer1K: 0.002},
	"gpt-3.5-turbo-0613":     {InputCostPer1K: 0.0015, OutputCostPer1K: 0.002},
	"gpt-3.5-turbo-instruct": {InputCostPer1K: 0.0015, OutputCostPer1K: 0.002},

	// Google Gemini
	"gemini-1.5-pro":            {InputCostPer1K: 0.00125, OutputCostPer1K: 0.005},
	"gemini-1.5-pro-extended":   {InputCostPer1K: 0.0025, OutputCostPer1K: 0.01},
	"gemini-1.5-flash":          {InputCostPer1K: 0.000075, OutputCostPer1K: 0.0003},
	"gemini-1.5-flash-extended": {InputCostPer1K: 0.00015, OutputCostPer1K: 0.0006},
	"gemini-1.0-pro":            {InputCostPer1K: 0.0005, OutputCostPer1K: 0.0015},

	// Anthropic Claude
	"claude-3.7-sonnet":  {InputCostPer1K: 0.003, OutputCostPer1K: 0.015},
	"claude-3.5-sonnet":  {InputCostPer1K: 0.003, OutputCostPer1K: 0.015},
	"claude-3-opus":      {InputCostPer1K: 0.015, OutputCostPer1K: 0.075},
	"claude-3.5-haiku":   {InputCostPer1K: 0.0008, OutputCostPer1K: 0.004},
	"claude-3-haiku":     {InputCostPer1K: 0.00025, OutputCostPer1K: 0.00125},
	"claude-instant-1.2": {InputCostPer1K: 0.00163, OutputCostPer1K: 0.00551},

	// Mistral
	"mistral-large":  {InputCostPer1K: 0.002, OutputCostPer1K: 0.006},
	"mistral-medium": {InputCostPer1K: 0.00275, OutputCostPer1K: 0.0081},
	"mistral-small":  {InputCostPer1K: 0.0006, OutputCostPer1K: 0.0018},
	"mistral-nemo":   {InputCostPer1K: 0.00015, OutputCostPer1K: 0.00015},
	"codestral":      {InputCostPer1K: 0.0002, OutputCostPer1K: 0.0006},
	"mixtral-8x7b":   {InputCostPer1K: 0.00024, OutputCostPer1K: 0.00024},

	// Cohere
	"command-a":            {InputCostPer1K: 0.0025, OutputCostPer1K: 0.01},
	"command-r+":           {InputCostPer1K: 0.0025, OutputCostPer1K: 0.01},
	"command-r":            {InputCostPer1K: 0.00015, OutputCostPer1K: 0.0006},
	"command-r-ft":         {InputCostPer1K: 0.0003, OutputCostPer1K: 0.0012},
	"command-r7b":          {InputCostPer1K: 0.0000375, OutputCostPer1K: 0.00015},
	"command-legacy":       {InputCostPer1K: 0.001, OutputCostPer1K: 0.002},
	"command-light-legacy": {InputCostPer1K: 0.0003, OutputCostPer1K: 0.0006},

	// Meta (Llama via Groq)
	"llama-3.1-70b": {InputCostPer1K: 0.00059, OutputCostPer1K: 0.00079},
	"llama-3-8b":    {InputCostPer1K: 0.00005, OutputCostPer1K: 0.00008},
	"llama-3.1-8b":  {InputCostPer1K: 0.00005, OutputCostPer1K: 0.00008},
}

func Count(text, model string) (int, error) {
	enc, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return 0, err
	}

	tokens := enc.Encode(text, nil, nil)

	return len(tokens), nil
}

func Analyze(input, output, model string) (*TokenInfo, error) {
	pricing, ok := modelPricing[model]
	if !ok {
		return nil, fmt.Errorf("pricing info not available for model: %s", model)
	}

	// We don't care for the error as the only error can be that the model is not supported
	inCount, _ := Count(input, model)
	outCount, _ := Count(output, model)

	total := inCount + outCount
	cost := (float64(inCount)/1000)*pricing.InputCostPer1K + (float64(outCount)/1000)*pricing.OutputCostPer1K

	return &TokenInfo{
		InputTokens:   inCount,
		OutputTokens:  outCount,
		TotalTokens:   total,
		EstimatedCost: cost,
	}, nil
}
