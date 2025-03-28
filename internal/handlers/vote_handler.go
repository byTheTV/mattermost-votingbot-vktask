package handlers

import (
	"strconv"

	"mattermost-bot/internal/service"
	"mattermost-bot/internal/models"
	
	"github.com/mattermost/mattermost-server/v6/model"
)

func HandleVote(s service.PollService, bot *models.Bot, post *model.Post, args []string) {
	if len(args) != 2 {
		replyToPost(bot, post, "Использование: /vote <ID опроса> <номер варианта>")
		return
	}

	pollID := args[0]
	optionIdx, err := strconv.Atoi(args[1])
	if err != nil || optionIdx < 1 {
		replyToPost(bot, post, "❌ Некорректный номер варианта")
		return
	}

	// Проверка существования опроса
	poll, err := s.GetPoll(pollID)
	if err != nil {
		replyToPost(bot, post, "❌ Опрос не найден")
		return
	}

	if !poll.Active {
		replyToPost(bot, post, "❌ Опрос закрыт")
		return
	}

	// Сохранение голоса
	err = s.Vote(pollID, post.UserId, optionIdx-1)
	if err != nil {
		replyToPost(bot, post, "❌ Ошибка голосования: "+err.Error())
		return
	}

	replyToPost(bot, post, "✅ Ваш голос учтён!")
}
