package models

import (
    "fmt"
    "log"
    "time"
    "strings"
    "github.com/mattermost/mattermost-server/v6/model"
)

type Bot struct {
	Client *model.Client4         // Клиент для REST API Mattermost
	Ws     *model.WebSocketClient // WebSocket-клиент для реальных событий
	User   *model.User            // Информация о самом боте
    URL    string                 // URL сервера Mattermost
    Token  string                 // Токен бота
}

func NewBot(serverURL, token string) (*Bot, error) {
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
        ws, err := model.NewWebSocketClient4(wsURL, b.Token)
        if err != nil {
            log.Printf("Ошибка подключения WebSocket: %v. Повтор через 5 секунд...", err)
            time.Sleep(5 * time.Second)
            continue
        }
        b.Ws = ws

        log.Println("WebSocket подключен")
        ws.Listen()

        // Обработчик событий
        go func() {
            defer ws.Close()
            for {
                select {
                case event, ok := <-ws.EventChannel:
                    if !ok {
                        log.Println("Канал событий закрыт")
                        return
                    }
                    eventChan <- event
                case <-ws.PingTimeoutChannel:
                    log.Println("Таймаут пинга, переподключение...")
                    return
                case <-ws.ResponseChannel:
                    // Игнорируем ответы
                }
            }
        }()

        // Ждем разрыва соединения
        <-ws.PingTimeoutChannel
        log.Println("Переподключение через 2 секунды...")
        time.Sleep(2 * time.Second)
    }
}