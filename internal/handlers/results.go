package handlers

import (
	"fmt"
	"strings"

	"mattermost-bot/internal/service"
	"mattermost-bot/internal/models"

	"github.com/mattermost/mattermost-server/v6/model"
)

func HandleResults(s service.PollService, bot *models.Bot, post *model.Post, args []string) {
	if len(args) != 1 {
		replyToPost(bot, post, "Использование: /results <ID опроса>")
		return
	}

	results, err := s.GetResults(args[0])
	if err != nil {
		replyToPost(bot, post, "❌ Ошибка: "+err.Error())
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 **Результаты опроса:** %s\n", results.Poll.Question))

	for i, opt := range results.Poll.Options {
		sb.WriteString(fmt.Sprintf("%d. %s — %d голосов\n", i+1, opt, results.Counts[i]))
	}

	sb.WriteString(fmt.Sprintf("\nВсего участников: %d", len(results.Votes)))
	SendMessageToChannel(bot, post.ChannelId, sb.String())
}
