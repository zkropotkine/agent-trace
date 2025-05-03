package tokens

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		model       string
		expectedMin int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Basic count",
			text:        "Hello world",
			model:       "gpt-3.5-turbo",
			expectedMin: 2,
			expectError: false,
		},
		{
			name:        "Empty text",
			text:        "",
			model:       "gpt-3.5-turbo",
			expectedMin: 0,
			expectError: false,
		},
		{
			name:        "Long text",
			text:        "This is a much longer text that should contain significantly more tokens than the simple hello world example. It includes various words and phrases to ensure proper token counting.",
			model:       "gpt-4",
			expectedMin: 20,
			expectError: false,
		},
		{
			name:        "Non-English text",
			text:        "こんにちは世界",
			model:       "gpt-4",
			expectedMin: 2,
			expectError: false,
		},
		{
			name:        "Unknown model",
			text:        "Hello world",
			model:       "unknown-model",
			expectedMin: 0,
			expectError: true,
			errorMsg:    "no encoding for model unknown-model",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			count, err := Count(tc.text, tc.model)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, count, tc.expectedMin)
			}
		})
	}
}

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		output         string
		model          string
		expectedMinIn  int
		expectedMinOut int
		expectError    bool
		errorMsg       string
	}{
		{
			name:           "Basic analysis",
			input:          "Hello",
			output:         "world",
			model:          "gpt-4-turbo",
			expectedMinIn:  1,
			expectedMinOut: 1,
			expectError:    false,
		},
		{
			name:           "Longer texts",
			input:          "This is a longer input text that should be properly analyzed by our function.",
			output:         "And this is a longer output that would typically come from an AI model in response.",
			model:          "gpt-3.5-turbo-0125",
			expectedMinIn:  10,
			expectedMinOut: 10,
			expectError:    false,
		},
		{
			name:           "Different OpenAI model",
			input:          "Testing with different model",
			output:         "Model response",
			model:          "gpt-4",
			expectedMinIn:  3,
			expectedMinOut: 2,
			expectError:    false,
		},
		{
			name:           "Empty input",
			input:          "",
			output:         "Response to empty input",
			model:          "gpt-4-turbo",
			expectedMinIn:  0,
			expectedMinOut: 3,
			expectError:    false,
		},
		{
			name:           "Unknown model",
			input:          "Hello",
			output:         "world",
			model:          "unknown-model",
			expectedMinIn:  0,
			expectedMinOut: 0,
			expectError:    true,
			errorMsg:       "pricing info not available for model: unknown-model",
		},
		{
			name:           "Missing pricing info",
			input:          "Hello",
			output:         "world",
			model:          "custom-model-without-pricing",
			expectedMinIn:  0,
			expectedMinOut: 0,
			expectError:    true,
			errorMsg:       "pricing info not available for model: custom-model-without-pricing",
		},
	}

	// Temporarily add a custom model to the tiktoken-go's known models for testing
	// This would typically be done with a mock, but for simplicity we'll just test the error case

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Analyze(tc.input, tc.output, tc.model)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.GreaterOrEqual(t, result.InputTokens, tc.expectedMinIn)
				assert.GreaterOrEqual(t, result.OutputTokens, tc.expectedMinOut)
				assert.Equal(t, result.TotalTokens, result.InputTokens+result.OutputTokens)

				// Calculate expected cost manually to verify logic
				pricing := modelPricing[tc.model]
				expectedCost := (float64(result.InputTokens)/1000)*pricing.InputCostPer1K +
					(float64(result.OutputTokens)/1000)*pricing.OutputCostPer1K
				assert.InDelta(t, expectedCost, result.EstimatedCost, 0.0001)
			}
		})
	}
}

func TestModelPricing(t *testing.T) {
	tests := []struct {
		name                    string
		model                   string
		expectedInputCostPer1K  float64
		expectedOutputCostPer1K float64
	}{
		{
			name:                    "OpenAI GPT-4o",
			model:                   "gpt-4o",
			expectedInputCostPer1K:  0.005,
			expectedOutputCostPer1K: 0.015,
		},
		{
			name:                    "Anthropic Claude",
			model:                   "claude-3-opus",
			expectedInputCostPer1K:  0.015,
			expectedOutputCostPer1K: 0.075,
		},
		{
			name:                    "Gemini",
			model:                   "gemini-1.5-pro",
			expectedInputCostPer1K:  0.00125,
			expectedOutputCostPer1K: 0.005,
		},
		{
			name:                    "Mistral",
			model:                   "mistral-large",
			expectedInputCostPer1K:  0.002,
			expectedOutputCostPer1K: 0.006,
		},
		{
			name:                    "Cohere",
			model:                   "command-r",
			expectedInputCostPer1K:  0.00015,
			expectedOutputCostPer1K: 0.0006,
		},
		{
			name:                    "Meta Llama",
			model:                   "llama-3.1-70b",
			expectedInputCostPer1K:  0.00059,
			expectedOutputCostPer1K: 0.00079,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pricing, exists := modelPricing[tc.model]
			assert.True(t, exists, fmt.Sprintf("Model %s should exist in pricing map", tc.model))
			assert.Equal(t, tc.expectedInputCostPer1K, pricing.InputCostPer1K)
			assert.Equal(t, tc.expectedOutputCostPer1K, pricing.OutputCostPer1K)
		})
	}
}
