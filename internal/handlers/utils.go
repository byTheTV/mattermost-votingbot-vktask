package handlers

import (
	"fmt"
	"strconv"

	botModel "mattermost-bot/internal/models"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

// ParseOptionIndex преобразует строку в индекс опции (начиная с 0)
func ParseOptionIndex(input string) (int, error) {
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 {
		return -1, fmt.Errorf("некорректный номер варианта")
	}
	return idx - 1, nil
}

// ReplyToPost отправляет ответ в тред
func replyToPost(bot *botModel.Bot, post *model.Post, message string) {
	reply := &model.Post{
		ChannelId: post.ChannelId,
		Message:   message,
		RootId:    post.Id,
	}

	if _, _, err := bot.Client.CreatePost(reply); err != nil {
		bot.Logger.Error("[ERROR] Ошибка отправки:", zap.Error(err))
	}
}

// SendMessageToChannel отправляет сообщение в указанный канал
func SendMessageToChannel(bot *botModel.Bot, channelID string, message string) {
	post := &model.Post{
		ChannelId: channelID,
		Message:   message,
	}

	if _, _, err := bot.Client.CreatePost(post); err != nil {
		bot.Logger.Error("[ERROR] Ошибка отправки:", zap.Error(err))
	}
}
