package repository

import (
	model "mattermost-bot/internal/models"
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
)

type VoteRepository interface {
    CreateVote(pollID, userID string, optionIdx int) error
    UpdateVote(pollID, userID string, optionIdx int) error
    GetVote(pollID, userID string) (bool, *model.Vote, error)
    GetVotes(pollID string) ([]*model.Vote, error) 
}

type TarantoolVoteRepo struct {
	conn *tarantool.Connection
}

func NewTarantoolVoteRepo(conn *tarantool.Connection) *TarantoolVoteRepo {
	return &TarantoolVoteRepo{conn: conn}
}

func (r *TarantoolVoteRepo) CreateVote(pollID, userID string, optionIdx int) error {
    _, err := r.conn.Do(
        tarantool.NewInsertRequest("votes").Tuple([]interface{}{
            pollID, 
            userID, 
            optionIdx,
        }),
    ).Get()
    return err
}

func (r *TarantoolVoteRepo) UpdateVote(pollID, userID string, optionIdx int) error {
    _, err := r.conn.Do(
        tarantool.NewUpdateRequest("votes").
            Key([]interface{}{pollID, userID}).
            Operations(tarantool.NewOperations().Assign(2, optionIdx)), // 2 — индекс поля option_idx
    ).Get()
    return err
}

func (r *TarantoolVoteRepo) GetVote(pollID, userID string) (bool, *model.Vote, error) {
    data, err := r.conn.Do(
        tarantool.NewSelectRequest("votes").
            Index("primary").
            Key([]interface{}{pollID, userID}),
    ).Get()

    if err != nil {
        return false, nil, fmt.Errorf("ошибка запроса: %w", err)
    }

    if len(data) == 0 {
        return false, nil, nil
    }

    // Type assertion for the first tuple in response
    tuple := data[0].([]interface{})
    vote := &model.Vote{
        PollID:    tuple[0].(string),
        UserID:    tuple[1].(string),
        OptionIdx: int(tuple[2].(int8)),
    }

    return true, vote, nil
}



func (r *TarantoolVoteRepo) GetVotes(pollID string) ([]*model.Vote, error) {
    resp, err := r.conn.Do(
        tarantool.NewSelectRequest("votes").
            Index("poll_id"). // Используем индекс poll_id
            Key([]interface{}{pollID}),
    ).Get()

    if err != nil {
        return nil, fmt.Errorf("ошибка получения голосов: %w", err)
    }

    var votes []*model.Vote
    for _, item := range resp {
        data := item.([]interface{})
        votes = append(votes, &model.Vote{
            PollID:    data[0].(string),
            UserID:    data[1].(string),
            OptionIdx: int(data[2].(int8)),
        })
    }
    return votes, nil
}