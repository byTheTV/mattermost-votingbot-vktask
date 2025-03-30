package config

type Config struct {
	MattermostURL string
	BotToken      string
	TarantoolAddr string
}

func Load() *Config {
	return &Config{
		MattermostURL: "http://localhost:8065",
		BotToken:      "cjrfc5yy8bbj9ff4j3wnb34the",
		TarantoolAddr: "localhost:3301",
	}
}
