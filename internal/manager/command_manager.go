// Файл: internal/manager/command_manager.go
package manager

import (
	"mattermost-bot/internal/handlers"
	"mattermost-bot/internal/models"
	"mattermost-bot/internal/service"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

type CommandManager struct {
	service service.PollService
	logger  *zap.Logger
}

func NewCommandManager(service service.PollService, logger *zap.Logger) *CommandManager {
	return &CommandManager{
		service: service,
		logger:  logger.With(zap.String("component", "command_manager")),
	}
}

func (cm *CommandManager) ProcessCommand(bot *models.Bot, post *model.Post) {
	cm.logger.Info("Обработка команды", zap.String("message", post.Message))

	// Извлекаем команду
	parts := strings.Fields(post.Message)
	if len(parts) == 0 {
		return
	}
	command := parts[0]

	// Получаем оставшуюся часть сообщения
	idx := strings.Index(post.Message, command)
	if idx == -1 {
		return
	}
	remaining := post.Message[idx+len(command):]
	messageWithoutCommand := strings.TrimSpace(remaining)

	// Парсим аргументы
	args := handlers.SplitArgs(messageWithoutCommand, cm.logger)

	switch command {
	case "/delete_poll":
		handlers.HandleDeletePoll(cm.service, bot, post, args)
	case "/poll":
		handlers.HandleCreatePoll(cm.service, bot, post, args)
	case "/vote":
		handlers.HandleVote(cm.service, bot, post, args)
	case "/results":
		handlers.HandleResults(cm.service, bot, post, args)
	case "/close":
		handlers.HandleClosePoll(cm.service, bot, post, args)
	case "/polls":
		handlers.HandleListPolls(cm.service, bot, post)
	default:
		cm.logger.Warn("Неизвестная команда", zap.String("command", command))
	}
}
