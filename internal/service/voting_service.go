package service

import (
	"fmt"
	"mattermost-bot/internal/models"
	"mattermost-bot/internal/repository"
)

type VotingService struct {
	pollRepo repository.PollRepository
	voteRepo repository.VoteRepository
}

func NewVotingService(pollRepo repository.PollRepository, voteRepo repository.VoteRepository) *VotingService {
	return &VotingService{pollRepo, voteRepo}
}

func (s *VotingService) CreatePoll(question string, options []string, userID, channelID string) (*models.Poll, error) {
	poll := &models.Poll{
		ID:        models.NewId(),
		Question:  question,
		Options:   options,
		CreatedBy: userID,
		ChannelID: channelID,
		Active:    true,
	}
	return poll, s.pollRepo.CreatePoll(poll)
}

func (s *VotingService) ClosePoll(pollID, userID string) error {
	poll, err := s.pollRepo.GetPoll(pollID)
	if err != nil {
		return fmt.Errorf("опрос не найден: %w", err)
	}

	if poll.CreatedBy != userID {
		return fmt.Errorf("недостаточно прав")
	}

	return s.pollRepo.ClosePoll(pollID)
}

func (s *VotingService) Vote(pollID, userID string, optionIdx int) error {

	exists, _, err := s.voteRepo.GetVote(pollID, userID)

	if err != nil {
		return fmt.Errorf("ошибка проверки голоса: %w", err)
	}

	if exists {
		return s.voteRepo.UpdateVote(pollID, userID, optionIdx)
	}
	return s.voteRepo.CreateVote(pollID, userID, optionIdx)
}

// GetResults возвращает результаты опроса
func (s *VotingService) GetResults(pollID string) (*PollResults, error) {
	poll, err := s.pollRepo.GetPoll(pollID)
	if err != nil {
		return nil, fmt.Errorf("опрос не найден: %w", err)
	}

	votes, err := s.voteRepo.GetVotes(pollID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения голосов: %w", err)
	}

	counts := make(map[int]int)
	for _, vote := range votes {
		counts[vote.OptionIdx]++
	}

	return &PollResults{
		Poll:   poll,
		Votes:  votes,
		Counts: counts,
	}, nil
}

// ListPolls возвращает список активных опросов в канале
func (s *VotingService) ListPolls(channelID string) ([]*models.Poll, error) {
	return s.pollRepo.GetPollsByChannel(channelID)
}

// GetPoll возвращает опрос по ID
func (s *VotingService) GetPoll(pollID string) (*models.Poll, error) {
	return s.pollRepo.GetPoll(pollID)
}

// UpdateVote обновляет голос
func (s *VotingService) UpdateVote(pollID, userID string, optionIdx int) error {
	return s.voteRepo.UpdateVote(pollID, userID, optionIdx)
}

// Добавьте этот метод в структуру VotingService
func (s *VotingService) DeletePoll(pollID, userID string) error {
	poll, err := s.pollRepo.GetPoll(pollID)
	if err != nil {
		return fmt.Errorf("опрос не найден: %w", err)
	}

	if poll.CreatedBy != userID {
		return fmt.Errorf("недостаточно прав для удаления")
	}

	return s.pollRepo.DeletePoll(pollID)
}
