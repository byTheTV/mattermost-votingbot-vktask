package service

import (
    "mattermost-bot/internal/models"
)

// PollService определяет интерфейс для работы с опросами
type PollService interface {
    // Создание нового опроса
    CreatePoll(question string, options []string, userID, channelID string) (*models.Poll, error)
    
    // Закрытие опроса по ID
    ClosePoll(pollID, userID string) error
    
    // Голосование в опросе
    Vote(pollID, userID string, optionIdx int) error
    
    // Получение результатов опроса
    GetResults(pollID string) (*PollResults, error)
    
    // Список активных опросов в канале
    ListPolls(channelID string) ([]*models.Poll, error)
    
    // Проверка существования опроса
    GetPoll(pollID string) (*models.Poll, error)
    
    // Обновление голоса пользователя
    UpdateVote(pollID, userID string, optionIdx int) error
}

// PollResults содержит результаты голосования
type PollResults struct {
    Poll    *models.Poll       // Данные опроса
    Votes   []*models.Vote     // Все голоса
    Counts  map[int]int       // Количество голосов за каждый вариант
}