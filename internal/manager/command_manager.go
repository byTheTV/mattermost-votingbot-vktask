package manager

import (
	"strings"

	"mattermost-bot/internal/handlers"
	"mattermost-bot/internal/models"
	"mattermost-bot/internal/service"

	"github.com/mattermost/mattermost-server/v6/model"
	"go.uber.org/zap"
)

type CommandManager struct {
    service service.PollService
    logger *zap.Logger
}

func NewCommandManager(service service.PollService, logger *zap.Logger) *CommandManager {
    return &CommandManager{
        service: service,
        logger:      logger.With(zap.String("component", "command_manager")),
    }
}

func (cm *CommandManager) ProcessCommand(bot *models.Bot, post *model.Post) {
    
    cm.logger.Info("Процессим команду...")

    parts := strings.Fields(post.Message)
    if len(parts) == 0 {
        
        return
    }   
    
    switch parts[0] {
    case "/poll":
        handlers.HandleCreatePoll(cm.service, bot, post, parts[1:])
    case "/vote":
        handlers.HandleVote(cm.service, bot, post, parts[1:])
    case "/results":
        handlers.HandleResults(cm.service, bot, post, parts[1:])
    case "/close":
        handlers.HandleClosePoll(cm.service, bot, post, parts[1:])
    case "/polls":
        handlers.HandleListPolls(cm.service, bot, post)
    }
}