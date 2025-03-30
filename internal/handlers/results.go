package handlers

import (
	"fmt"
	"strings"

	"mattermost-bot/internal/service"
	"mattermost-bot/internal/models"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"

)

func HandleResults(s service.PollService, bot *models.Bot, post *model.Post, args []string) {
	if len(args) != 1 {
		replyToPost(bot, post, "Использование: /results <ID опроса>")
		bot.Logger.Info("Неправильно использована команда /result")

		return
	}

	results, err := s.GetResults(args[0])
	if err != nil {
		replyToPost(bot, post, "❌ Ошибка: "+err.Error())
		bot.Logger.Error("Возникла ошибка у пользователя:", zap.Error(err))
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 **Результаты опроса:** %s\n", results.Poll.Question))

	for i, opt := range results.Poll.Options {
		sb.WriteString(fmt.Sprintf("%d. %s — %d голосов\n", i+1, opt, results.Counts[i]))
	}

	sb.WriteString(fmt.Sprintf("\nВсего участников: %d", len(results.Votes)))
	
	SendMessageToChannel(bot, post.ChannelId, sb.String())
		bot.Logger.Info("Результаты опроса отправлены",
		zap.String("poll_question", results.Poll.Question),
		zap.Int("total_participants", len(results.Votes)),
		zap.String("channel_id", post.ChannelId),
	)


}
