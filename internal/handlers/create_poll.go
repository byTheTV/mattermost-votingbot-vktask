package handlers

import (
	"fmt"
	"log"
	"mattermost-bot/internal/service"
	"strings"
	"mattermost-bot/internal/models"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

func HandleCreatePoll(s service.PollService, bot *models.Bot, post *model.Post, args []string) {
	if len(args) < 2 {
		replyToPost(bot, post, "Использование: /poll 'Вопрос' 'Вариант 1' 'Вариант 2' ...")
		bot.Logger.Info("Ошибка использования /poll")
		return
	}
	
	log.Println(args[0], args[1:], post.UserId, post.ChannelId)

	poll, err := s.CreatePoll(args[0], args[1:], post.UserId, post.ChannelId)
	if err != nil {
		replyToPost(bot, post, "❌ Ошибка: "+err.Error())
		bot.Logger.Error("Возникла ошибка у пользователя:", zap.Error(err))

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

	bot.Logger.Info("Новый опрос создан",
		zap.String("poll_id", poll.ID),
		zap.String("question", poll.Question),
		zap.String("channel_id", post.ChannelId),
		zap.String("options", sb.String()),
	)
}
