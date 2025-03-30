package handlers

import (
	"fmt"
	botModel "mattermost-bot/internal/models"
	"strconv"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

// ParseOptionIndex преобразует строку в индекс опции (начиная с 0)
func ParseOptionIndex(input string) (int, error) {
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 {
		return -1, fmt.Errorf("некорректный номер варианта")
	}
	return idx - 1, nil
}

// ReplyToPost отправляет ответ в тред
func replyToPost(bot *botModel.Bot, post *model.Post, message string) {
	reply := &model.Post{
		ChannelId: post.ChannelId,
		Message:   message,
		RootId:    post.Id,
	}

	if _, _, err := bot.Client.CreatePost(reply); err != nil {
		bot.Logger.Error("[ERROR] Ошибка отправки:", zap.Error(err))
	}
}

// SendMessageToChannel отправляет сообщение в указанный канал
func SendMessageToChannel(bot *botModel.Bot, channelID string, message string) {
	post := &model.Post{
		ChannelId: channelID,
		Message:   message,
	}

	if _, _, err := bot.Client.CreatePost(post); err != nil {
		bot.Logger.Error("[ERROR] Ошибка отправки:", zap.Error(err))
	}
}

func SplitArgs(input string, logger *zap.Logger) []string {
	var args []string
	var currentArg strings.Builder
	inQuotes := false

	// Логирование входных данных
	logger.Debug("[DEBUG] SplitArgs input:", zap.String("input", input))

	for i, r := range input {
		switch r {
		case '"':
			inQuotes = !inQuotes
			// Логирование обнаружения кавычки
		case ' ':
			if inQuotes {
				currentArg.WriteRune(r)
				// Логирование пробела внутри кавычек
			} else {
				if currentArg.Len() > 0 {
					args = append(args, currentArg.String())
					// Логирование сохраненного аргумента
					logger.Debug("[DEBUG] Saved argument", zap.Int("pos", i), zap.String("argument", currentArg.String()))
					currentArg.Reset()
				}
			}
		default:
			currentArg.WriteRune(r)
			// Логирование добавленного символа
		}
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
		// Логирование финального аргумента
		logger.Debug("[DEBUG] Final argument", zap.String("argument", currentArg.String()))
	}

	// Логирование результата
	logger.Debug("[DEBUG] Result", zap.Strings("args", args))
	return args
}
