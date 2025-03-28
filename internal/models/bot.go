package models

import (
    "fmt"
    "log"
    
    "github.com/mattermost/mattermost-server/v6/model"
)

type Bot struct {
	Client *model.Client4         // Клиент для REST API Mattermost
	Ws     *model.WebSocketClient // WebSocket-клиент для реальных событий
	User   *model.User            // Информация о самом боте
}

// NewBot создаёт и настраивает экземпляр бота
func NewBot(serverURL, token string) (*Bot, error) {
    // Инициализация REST-клиента
    client := model.NewAPIv4Client(serverURL)
    client.SetToken(token)

    // Получение информации о боте
    user, _, err := client.GetMe("")
    if err != nil {
        return nil, fmt.Errorf("ошибка аутентификации: %v", err)
    }

    // Подключение WebSocket
    ws, err := model.NewWebSocketClient4(serverURL, token)
    if err != nil {
        return nil, fmt.Errorf("ошибка подключения WebSocket: %v", err)
    }

    return &Bot{
        Client: client,
        Ws:     ws,
        User:   user,
    }, nil
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