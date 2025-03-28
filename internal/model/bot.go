package model

import "github.com/mattermost/mattermost-server/v6/model"

type Bot struct {
	Client *model.Client4         // Клиент для REST API Mattermost
	Ws     *model.WebSocketClient // WebSocket-клиент для реальных событий
	User   *model.User            // Информация о самом боте
}
