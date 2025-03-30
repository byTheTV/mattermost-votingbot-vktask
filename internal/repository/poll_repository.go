package repository

import (
	"fmt"
	"strings"
	
	"mattermost-bot/internal/models"
	"github.com/tarantool/go-tarantool/v2"
)

const (
	PollSpace     = "polls"
	VotesSpace    = "votes"
	OptionSep     = "|"
	ActiveStatus  = 5
)

type PollRepository interface {
	CreatePoll(poll *models.Poll) error
	GetPoll(pollID string) (*models.Poll, error)
	ClosePoll(pollID string) error
	GetPollsByChannel(channelID string) ([]*models.Poll, error)
}

type TarantoolPollRepo struct {
	conn *tarantool.Connection
}

func NewTarantoolPollRepo(conn *tarantool.Connection) *TarantoolPollRepo {
	return &TarantoolPollRepo{conn: conn}
}

func (r *TarantoolPollRepo) CreatePoll(poll *models.Poll) error {
	optionsStr := strings.Join(poll.Options, OptionSep)
	_, err := r.conn.Do(
		tarantool.NewInsertRequest(PollSpace).Tuple([]interface{}{
			poll.ID,
			poll.Question,
			optionsStr,
			poll.CreatedBy,
			poll.ChannelID,
			true,
		}),
	).Get()
	return err
}

func (r *TarantoolPollRepo) GetPoll(pollID string) (*models.Poll, error) {
	resp, err := r.conn.Do(
		tarantool.NewCallRequest("box.space.polls:get").Args([]interface{}{pollID}),
	).Get()

	if err != nil || len(resp) == 0 {
		return nil, fmt.Errorf("poll not found")
	}

	data := resp[0].([]interface{})
	return &models.Poll{
		ID:        data[0].(string),
		Question:  data[1].(string),
		Options:   strings.Split(data[2].(string), OptionSep),
		CreatedBy: data[3].(string),
		ChannelID: data[4].(string),
		Active:    data[5].(bool),
	}, nil
}

func (r *TarantoolPollRepo) ClosePoll(pollID string) error {
    _, err := r.conn.Do(
        tarantool.NewUpdateRequest("polls").
            Key([]interface{}{pollID}).
            Operations(tarantool.NewOperations().Assign(5, false)), // 5 — индекс поля "active"
    ).Get()
    return err
}

func (r *TarantoolPollRepo) GetPollsByChannel(channelID string) ([]*models.Poll, error) {
    resp, err := r.conn.Do(
        tarantool.NewSelectRequest("polls").
            Index("channel"). // Указываем индекс
            Key([]interface{}{channelID}),
    ).Get()

    if err != nil {
        return nil, fmt.Errorf("ошибка получения опросов: %w", err)
    }

    var polls []*models.Poll
    for _, item := range resp {
        data := item.([]interface{})
        options := strings.Split(data[2].(string), OptionSep)
        polls = append(polls, &models.Poll{
            ID:        data[0].(string),
            Question:  data[1].(string),
            Options:   options,
            CreatedBy: data[3].(string),
            ChannelID: data[4].(string),
            Active:    data[5].(bool),
        })
    }
    return polls, nil
}