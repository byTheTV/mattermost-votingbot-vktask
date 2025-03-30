package handlers

import (
    "fmt"
    "log"
    "strconv"
    
    "github.com/mattermost/mattermost-server/v6/model"
    botModel "mattermost-bot/internal/models"
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
        log.Printf("[ERROR] Ошибка отправки: %v", err)
    }
}

// SendMessageToChannel отправляет сообщение в указанный канал
func SendMessageToChannel(bot *botModel.Bot, channelID string, message string) {
    post := &model.Post{
        ChannelId: channelID,
        Message:   message,
    }

    if _, _, err := bot.Client.CreatePost(post); err != nil {
        log.Printf("[ERROR] Ошибка отправки сообщения в канал: %v", err)
    }
}



// ValidatePoll проверяет существование и активность опроса
func ValidatePoll(poll *botModel.Poll, pollID string) error {
    if poll == nil {
        return fmt.Errorf("опрос с ID '%s' не найден", pollID)
    }
    if !poll.Active {
        return fmt.Errorf("опрос закрыт")
    }
    return nil
}

// CheckPermissions проверяет права пользователя
func CheckPermissions(userID, creatorID string) error {
    if userID != creatorID {
        return fmt.Errorf("доступ запрещён")
    }
    return nil
}