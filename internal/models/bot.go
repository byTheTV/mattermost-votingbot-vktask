package models

import (
    "fmt"
    "log"
    "strings"
    "github.com/mattermost/mattermost-server/v6/model"
)

type Bot struct {
	Client *model.Client4         // Клиент для REST API Mattermost
	Ws     *model.WebSocketClient // WebSocket-клиент для реальных событий
	User   *model.User            // Информация о самом боте
}

func NewBot(serverURL, token string) (*Bot, error) {
    // Преобразуем HTTP URL в WebSocket URL
    wsURL := convertToWebsocketURL(serverURL)

    client := model.NewAPIv4Client(serverURL)
    client.SetToken(token)

    user, _, err := client.GetMe("")
    if err != nil {
        return nil, fmt.Errorf("ошибка авторизации: %v", err)
    }

    ws, err := model.NewWebSocketClient4(wsURL, token)
    if err != nil {
        return nil, fmt.Errorf("ошибка WebSocket: %v", err)
    }

    return &Bot{
        Client: client,
        Ws:     ws,
        User:   user,
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

// Listen запускает прослушивание событий через WebSocket
func (b *Bot) Listen(eventChan chan *model.WebSocketEvent) {
    b.Ws.Listen()
    
    go func() {
        for {
            select {
            case event := <-b.Ws.EventChannel:
                eventChan <- event
            case <-b.Ws.PingTimeoutChannel:
                log.Println("Соединение с WebSocket разорвано")
                return
            }
        }
    }()
}