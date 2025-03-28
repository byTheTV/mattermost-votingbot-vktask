package handlers

import (
	"fmt"
	"strings"

	"mattermost-bot/internal/service"
	"mattermost-bot/internal/models"

	"github.com/mattermost/mattermost-server/v6/model"
)

func HandleListPolls(s service.PollService, bot *models.Bot, post *model.Post) {
	polls, err := s.ListPolls(post.ChannelId)
	if err != nil {
		replyToPost(bot, post, "❌ Ошибка получения списка опросов")
		return
	}

	if len(polls) == 0 {
		replyToPost(bot, post, "В этом канале нет активных опросов")
		return
	}

	var sb strings.Builder
	sb.WriteString("📋 **Активные опросы:**\n")

	for _, poll := range polls {
		sb.WriteString(fmt.Sprintf(
			"• ID: `%s`\n  Вопрос: %s\n  Создал: <@%s>\n\n",
			poll.ID,
			poll.Question,
			poll.CreatedBy,
		))
	}

	replyToPost(bot, post, sb.String())
}
