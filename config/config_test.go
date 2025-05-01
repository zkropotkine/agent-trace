package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// save original environment to restore later
	getEnv := func() map[string]string {
		env := make(map[string]string)
		for _, e := range os.Environ() {
			kv := strings.SplitN(e, "=", 2)
			env[kv[0]] = kv[1]
		}
		return env
	}
	initialEnv := getEnv()

	// restore original environment after each case
	resetEnv := func(t *testing.T) {
		for k := range getEnv() {
			if _, ok := initialEnv[k]; !ok {
				assert.NoError(t, os.Unsetenv(k))
			}
		}
		for k, v := range initialEnv {
			assert.NoError(t, os.Setenv(k, v))
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
				assert.Equal(t, "mongodb://localhost:27017", c.MongoURI)
				assert.Equal(t, "dev", c.Env)
			},
		},
		{
			name: "overrides all values from environment",
			envs: func(t *testing.T) map[string]string {
				return map[string]string{
					envPrefix + "_PORT":      ":9090",
					envPrefix + "_MONGO_URI": "mongodb://example:27017",
					envPrefix + "_ENV":       "prod",
				}
			},
			assert: func(t *testing.T, c *Config) {
				assert.Equal(t, ":9090", c.Port)
				assert.Equal(t, "mongodb://example:27017", c.MongoURI)
				assert.Equal(t, "prod", c.Env)
			},
		},
		{
			name: "overrides only port when partially set",
			envs: func(t *testing.T) map[string]string {
				return map[string]string{
					envPrefix + "_PORT": ":3000",
				}
			},
			assert: func(t *testing.T, c *Config) {
				assert.Equal(t, ":3000", c.Port)
				assert.Equal(t, "mongodb://localhost:27017", c.MongoURI)
				assert.Equal(t, "dev", c.Env)
			},
		},
	}

	for _, tc := range tests {
		// set up this test’s environment
		for k, v := range tc.envs(t) {
			assert.NoError(t, os.Setenv(k, v))
		}

		// capture range variable to avoid closure gotcha
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cfg := LoadConfig()
			tc.assert(t, cfg)
		})

		resetEnv(t)
	}
}
