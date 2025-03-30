package handlers

import (
	"mattermost-bot/internal/models"
	"mattermost-bot/internal/service"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

func HandleClosePoll(s service.PollService, bot *models.Bot, post *model.Post, args []string) {
	if len(args) != 1 {
		replyToPost(bot, post, "Использование: /close <ID опроса>")
		bot.Logger.Info("Неправильно использована команда /close")
		return
	}

	err := s.ClosePoll(args[0], post.UserId)
	if err != nil {
		replyToPost(bot, post, "❌ Ошибка: "+err.Error())
		bot.Logger.Info("Возникла ошибка у пользователя:", zap.Error(err))
		return
	}

	SendMessageToChannel(bot, post.ChannelId, "🔒 Опрос закрыт. Используйте /results для просмотра итогов")
	bot.Logger.Info("Опрос закрыт",
		zap.String("poll_id", args[0]),
		zap.String("user_id", post.UserId),
		zap.String("channel_id", post.ChannelId),
	)
}
