package handlers

import (
	"mattermost-bot/internal/models"
	"mattermost-bot/internal/service"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

func HandleDeletePoll(s service.PollService, bot *models.Bot, post *model.Post, args []string) {
	if len(args) != 1 {
		replyToPost(bot, post, "Использование: /delete_poll <ID опроса>")
		bot.Logger.Info("Неправильно использована команда /delete")
		return
	}

	err := s.DeletePoll(args[0], post.UserId)

	if err != nil {
		replyToPost(bot, post, "❌ Ошибка: "+err.Error())
		bot.Logger.Error("Ошибка удаления опроса",
			zap.String("poll_id", args[0]),
			zap.String("user_id", post.UserId),
			zap.Error(err),
		)
		return
	}

	SendMessageToChannel(bot, post.ChannelId, "🗑 Опрос успешно удален")
	bot.Logger.Info("Опрос удален",
		zap.String("poll_id", args[0]),
		zap.String("user_id", post.UserId),
		zap.String("channel_id", post.ChannelId),
	)
}
