package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	getEnv := func() map[string]string {
		env := make(map[string]string)
		for _, e := range os.Environ() {
			kv := strings.SplitN(e, "=", 2)
			env[kv[0]] = kv[1]
		}
		return env
	}
	initialEnv := getEnv()

	resetEnv := func(t *testing.T) {
		for k := range getEnv() {
			if _, ok := initialEnv[k]; !ok {
				_ = os.Unsetenv(k)
			}
		}
		for k, v := range initialEnv {
			_ = os.Setenv(k, v)
		}
	}

	tests := []struct {
		name   string
		envs   func(t *testing.T) map[string]string
		assert func(t *testing.T, c *Config)
	}{
		{
			name: "loads defaults when no env vars set",
			envs: func(t *testing.T) map[string]string { return nil },
			assert: func(t *testing.T, c *Config) {
				assert.Equal(t, ":8080", c.Port)
				assert.Equal(t, "mongodb://localhost:27017", c.Mongo.URI)
				assert.Equal(t, "agentTrace", c.Mongo.DB)
				assert.Equal(t, "traces", c.Mongo.Collection)
				assert.Equal(t, "info", c.Log.Level)
				assert.Equal(t, "text", c.Log.Format)
				assert.Equal(t, "dev", c.Env)
			},
		},
		{
			name: "overrides all values from environment",
			envs: func(t *testing.T) map[string]string {
				return map[string]string{
					"AGENT_TRACE_PORT":             ":9090",
					"AGENT_TRACE_ENV":              "prod",
					"AGENT_TRACE_MONGO_URI":        "mongodb://override:27017",
					"AGENT_TRACE_MONGO_DB":         "overrideDB",
					"AGENT_TRACE_MONGO_COLLECTION": "logs",
					"AGENT_TRACE_LOG_LEVEL":        "debug",
					"AGENT_TRACE_LOG_FORMAT":       "json",
				}
			},
			assert: func(t *testing.T, c *Config) {
				assert.Equal(t, ":9090", c.Port)
				assert.Equal(t, "prod", c.Env)
				assert.Equal(t, "mongodb://override:27017", c.Mongo.URI)
				assert.Equal(t, "overrideDB", c.Mongo.DB)
				assert.Equal(t, "logs", c.Mongo.Collection)
				assert.Equal(t, "debug", c.Log.Level)
				assert.Equal(t, "json", c.Log.Format)
			},
		},
		{
			name: "overrides only port when partially set",
			envs: func(t *testing.T) map[string]string {
				return map[string]string{
					"AGENT_TRACE_PORT": ":3000",
				}
			},
			assert: func(t *testing.T, c *Config) {
				assert.Equal(t, ":3000", c.Port)
				assert.Equal(t, "mongodb://localhost:27017", c.Mongo.URI)
				assert.Equal(t, "agentTrace", c.Mongo.DB)
				assert.Equal(t, "traces", c.Mongo.Collection)
				assert.Equal(t, "info", c.Log.Level)
				assert.Equal(t, "text", c.Log.Format)
				assert.Equal(t, "dev", c.Env)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs(t) {
				_ = os.Setenv(k, v)
			}
			cfg := LoadConfig()
			tc.assert(t, cfg)
			resetEnv(t)
		})
	}
}
