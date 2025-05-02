package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "AGENT_TRACE"

type Config struct {
	Env   string `envconfig:"ENV" default:"dev"`
	Log   Log    `envconfig:"LOG"`
	Mongo Mongo  `envconfig:"MONGO"`
	Port  string `envconfig:"PORT" default:":8080"`
}

type Mongo struct {
	URI        string `envconfig:"URI" default:"mongodb://localhost:27017"`
	DB         string `envconfig:"DB" default:"agentTrace"`
	Collection string `envconfig:"COLLECTION" default:"traces"`
}

type Log struct {
	Format string `envconfig:"FORMAT" default:"text"`
	Level  string `envconfig:"LEVEL" default:"info"`
}

func LoadConfig() *Config {
	var cfg Config
	if err := envconfig.Process(envPrefix, &cfg); err != nil {
		log.Fatalf("failed to load environment config: %v", err)
	}

	return &cfg
}
