package config

import(
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	MattermostURL string
	BotToken      string
	TarantoolAddr string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		MattermostURL: getEnv("MATTERMOST_URL", "http://localhost:8065"),
		BotToken:      getEnv("BOT_TOKEN", "cjrfc5yy8bbj9ff4j3wnb34the"),
		TarantoolAddr: getEnv("TARANTOOL_ADDR", "localhost:3301"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
