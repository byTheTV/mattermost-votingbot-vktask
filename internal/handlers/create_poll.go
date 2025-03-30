// Файл: internal/handlers/handlers.go
package handlers

import (
	"fmt"
	"mattermost-bot/internal/models"
	"mattermost-bot/internal/service"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

func HandleCreatePoll(s service.PollService, bot *models.Bot, post *model.Post, args []string) {
	if len(args) < 2 {
		replyToPost(bot, post, `Использование: /poll "Вопрос" "Вариант 1" "Вариант 2" ...`)
		bot.Logger.Info("Ошибка использования /poll")
		return
	}

	bot.Logger.Info("Создание опроса",
		zap.String("опрос", args[0]),            // Логирование названия опроса
		zap.Strings("варианты", args[1:]),       // Логирование вариантов ответа
		zap.String("UserID", post.UserId),       // Логирование ID пользователя
		zap.String("ChannelID", post.ChannelId), // Логирование ID канала
	)
	poll, err := s.CreatePoll(args[0], args[1:], post.UserId, post.ChannelId)
	if err != nil {
		replyToPost(bot, post, "❌ Ошибка: "+err.Error())
		bot.Logger.Error("Ошибка создания опроса", zap.Error(err))
		return
	}

	var sb strings.Builder
	for i, opt := range poll.Options {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, opt))
	}

	SendMessageToChannel(bot, post.ChannelId, fmt.Sprintf(
		"📊 **Новый опрос!**\nВопрос: %s\nВарианты:\n%sID: `%s`",
		poll.Question,
		sb.String(),
		poll.ID,
	))

	bot.Logger.Info("Опрос создан",
		zap.String("poll_id", poll.ID),
		zap.String("question", poll.Question),
		zap.String("channel_id", post.ChannelId),
		zap.Strings("options", poll.Options),
	)
}
