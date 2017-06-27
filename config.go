package letsrest

import "os"

type Config struct {
	MongoURL string
}

func ReadConfigFromEnv() *Config {
	return &Config{
		MongoURL: os.Getenv("LETSREST_MONGO_URL"),
	}
}
