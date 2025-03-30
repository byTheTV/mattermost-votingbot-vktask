package manager

import (
    "strings"
    
    "mattermost-bot/internal/handlers"
    "mattermost-bot/internal/models"
    "mattermost-bot/internal/service"

	"github.com/mattermost/mattermost-server/v6/model"

)

type CommandManager struct {
    service service.PollService
}

func NewCommandManager(service service.PollService) *CommandManager {
    return &CommandManager{service: service}
}

func (cm *CommandManager) ProcessCommand(bot *models.Bot, post *model.Post) {
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