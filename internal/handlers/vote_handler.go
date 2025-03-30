package handlers

import (
	"mattermost-bot/internal/models"
	"mattermost-bot/internal/service"
	"strconv"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

func HandleVote(s service.PollService, bot *models.Bot, post *model.Post, args []string) {
	if len(args) != 2 {
		replyToPost(bot, post, "Использование: /vote <ID опроса> <номер варианта>")
		bot.Logger.Info("Неправильно использована команда /vote")
		return
	}

	pollID := args[0]
	optionIdx, err := strconv.Atoi(args[1])
	if err != nil || optionIdx < 1 {
		replyToPost(bot, post, "❌ Некорректный номер варианта")
		bot.Logger.Info("Некорректный номер варианта у пользователя")

		return
	}
	bot.Logger.Info("Выбранный номер варианта корректен",
		zap.Int("option_index", optionIdx),
	)
	// Проверка существования опроса
	poll, err := s.GetPoll(pollID)
	if err != nil {
		replyToPost(bot, post, "❌ Опрос не найден")
		bot.Logger.Info("Опрос пользователя не найден")

		return
	}
	bot.Logger.Info("Выбранный опрос существует", zap.Any("poll", poll))

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
	bot.Logger.Info("Голос учтён",
		zap.String("user_id", post.UserId),
		zap.String("pollID", pollID),
	)

}
