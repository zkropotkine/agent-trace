package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "AGENT_TRACE"

type Config struct {
	Port     string `envconfig:"PORT" default:":8080"`
	MongoURI string `envconfig:"MONGO_URI" default:"mongodb://localhost:27017" split_words:"true"`
	Env      string `envconfig:"ENV" default:"dev"` // e.g. dev, prod
}

func LoadConfig() *Config {
	var cfg Config
	if err := envconfig.Process(envPrefix, &cfg); err != nil {
		log.Fatalf("failed to load environment config: %v", err)
	}

	return &cfg
}
