package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

type Bot struct {
	Client *model.Client4         // Клиент для REST API Mattermost
	Ws     *model.WebSocketClient // WebSocket-клиент для реальных событий
	User   *model.User            // Информация о самом боте
	URL    string                 // URL сервера Mattermost
	Token  string                 // Токен бота
	Logger *zap.Logger
}

func NewBot(serverURL, token string, logger *zap.Logger) (*Bot, error) {
	logger = logger.With(zap.String("component", "bot"))

	logger.Info("Initializing Mattermost bot",
		zap.String("url", serverURL),
		zap.String("token_prefix", token[:4]+"***"),
	)

	client := model.NewAPIv4Client(serverURL)
	client.SetToken(token)

	user, _, err := client.GetMe("")
	if err != nil {
		return nil, fmt.Errorf("ошибка авторизации: %v", err)
	}

	return &Bot{
		Client: client,
		URL:    serverURL,
		Token:  token,
		User:   user,
		Logger: logger,
	}, nil
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

func (b *Bot) Listen(eventChan chan<- *model.WebSocketEvent) {
	for {
		wsURL := convertToWebsocketURL(b.URL)
		b.Logger.Info("Connecting to WebSocket", zap.String("url", wsURL))

		ws, err := model.NewWebSocketClient4(wsURL, b.Token)
		if err != nil {
			b.Logger.Error("WebSocket connection failed",
				zap.Error(err),
				zap.Duration("retry_in", 5*time.Second),
			)
			time.Sleep(5 * time.Second)
			continue
		}
		b.Ws = ws

		b.Logger.Info("WebSocket connected")
		ws.Listen()

		// Обработчик событий
		go func() {
			defer ws.Close()
			for {
				select {
				case event, ok := <-ws.EventChannel:
					if !ok {
						b.Logger.Info("Канал событий закрыт")
						return
					}
					eventChan <- event
				case <-ws.PingTimeoutChannel:
					b.Logger.Warn("Таймаут пинга, переподключение...")
					return
				case <-ws.ResponseChannel:
					// Игнорируем ответы
				}
			}
		}()

		// Ждем разрыва соединения
		<-ws.PingTimeoutChannel
		b.Logger.Info("Переподключение через 2 секунды...")
		time.Sleep(2 * time.Second)
	}
}
