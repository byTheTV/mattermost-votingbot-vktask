package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	mmModel "github.com/mattermost/mattermost-server/v6/model"

	"mattermost-bot/internal/model"
)

// BotApp wraps the Bot model and provides application-specific functionality
type BotApp struct {
	Bot *model.Bot
}

func NewBot(serverURL, token string) (*BotApp, error) {
	// Преобразуем HTTP URL в WebSocket URL
	wsURL := convertToWebsocketURL(serverURL)

	client := mmModel.NewAPIv4Client(serverURL)
	client.SetToken(token)

	user, _, err := client.GetMe("")
	if err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %v", err)
	}

	ws, err := mmModel.NewWebSocketClient4(wsURL, token)
	if err != nil {
		return nil, fmt.Errorf("ошибка WebSocket: %v", err)
	}

	bot := &model.Bot{
		Client: client,
		Ws:     ws,
		User:   user,
	}

	return &BotApp{Bot: bot}, nil
}

// Преобразование URL для WebSocket
func convertToWebsocketURL(serverURL string) string {
	if strings.HasPrefix(serverURL, "http://") {
		return strings.Replace(serverURL, "http://", "ws://", 1)
	}
	if strings.HasPrefix(serverURL, "https://") {
		return strings.Replace(serverURL, "https://", "wss://", 1)
	}
	return serverURL
}

func (b *BotApp) Listen() {
	b.Bot.Ws.Listen()

	for event := range b.Bot.Ws.EventChannel {
		log.Println(event)
		if event.EventType() == mmModel.WebsocketEventPosted {
			b.handleMessage(event)
		}
	}
}

func (b *BotApp) handleMessage(event *mmModel.WebSocketEvent) {
	var post mmModel.Post
	err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post)
	if err != nil {
		log.Printf("Ошибка парсинга сообщения: %v", err)
		return
	}

	if post.UserId == b.Bot.User.Id {
		return
	}

	if strings.TrimSpace(post.Message) == "/whoisme" {
		user, _, err := b.Bot.Client.GetUser(post.UserId, "")
		if err != nil {
			log.Printf("Ошибка получения пользователя: %v", err)
			return
		}

		reply := &mmModel.Post{
			ChannelId: post.ChannelId,
			Message: fmt.Sprintf(
				" Ваша информация:\n"+
					"- Имя: %s\n"+
					"- Никнейм: @%s\n"+
					"- ID: `%s`",
				user.GetFullName(),
				user.Username,
				user.Id,
			),
		}

		if _, _, err := b.Bot.Client.CreatePost(reply); err != nil {
			log.Printf("Ошибка отправки: %v", err)
		}
	}
}

func main() {
	loaderr := godotenv.Load()
	if loaderr != nil {
		fmt.Println("Ошибка загрузки .env файла:", loaderr)
		return
	}

	// Получаем токен из переменной окружения
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		fmt.Println("Ошибка: BOT_TOKEN не задан в .env файле")
		return
	}

	serverURL := "http://localhost:8065" // Замените на ваш URL Mattermost сервера

	botApp, err := NewBot(serverURL, botToken)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	log.Println("Бот запущен. Для выхода нажмите Ctrl+C")
	go botApp.Listen()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Завершение работы...")
	botApp.Bot.Ws.Close()
}
