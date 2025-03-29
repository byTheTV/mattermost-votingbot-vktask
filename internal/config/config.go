package config

type Config struct {
	MattermostURL string
	BotToken      string
	TarantoolAddr string
}

func Load() *Config {
	return &Config{
		MattermostURL: "http://localhost:8065",
		BotToken:      "f9zwqh7rsfrb9n5iymeng95rbo",
		TarantoolAddr: "localhost:3301",
	}
}
