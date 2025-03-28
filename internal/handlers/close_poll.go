package handlers

import (
	"mattermost-bot/internal/service"
	"mattermost-bot/internal/models"

	"github.com/mattermost/mattermost-server/v6/model"
)

func HandleClosePoll(s service.PollService, bot *models.Bot, post *model.Post, args []string) {
	if len(args) != 1 {
		replyToPost(bot, post, "Использование: /close <ID опроса>")
		return
	}

	err := s.ClosePoll(args[0], post.UserId)
	if err != nil {
		replyToPost(bot, post, "❌ Ошибка: "+err.Error())
		return
	}

	replyToPost(bot, post, "🔒 Опрос закрыт. Используйте /results для просмотра итогов")
}
