package handlers

import (
	"fmt"
	"log"
	"mattermost-bot/internal/service"
	"strings"
	"mattermost-bot/internal/models"

	"github.com/mattermost/mattermost-server/v6/model"
)

func HandleCreatePoll(s service.PollService, bot *models.Bot, post *model.Post, args []string) {
	if len(args) < 2 {
		replyToPost(bot, post, "Использование: /poll 'Вопрос' 'Вариант 1' 'Вариант 2' ...")
		return
	}
	
	log.Println(args[0], args[1:], post.UserId, post.ChannelId)

	poll, err := s.CreatePoll(args[0], args[1:], post.UserId, post.ChannelId)
	if err != nil {
		replyToPost(bot, post, "❌ Ошибка: "+err.Error())
		return
	}

	var sb strings.Builder
	for i, opt := range poll.Options {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, opt))
	}

	replyToPost(bot, post, fmt.Sprintf(
		"📊 **Новый опрос!**\nВопрос: %s\nВарианты:\n%sID: `%s`",
		poll.Question,
		sb.String(),
		poll.ID,
	))
}
