package config

type Config struct {
    MattermostURL  string
    BotToken       string
    TarantoolAddr  string
}

func Load() *Config {
    return &Config{
        MattermostURL: "https://mattermost.example.com",
        BotToken:      "your-bot-token",
        TarantoolAddr: "localhost:3301",
    }
}