package errors

import "fmt"

type PollNotFound struct {
    PollID string
}

func (e PollNotFound) Error() string {
    return fmt.Sprintf("опрос с ID '%s' не найден", e.PollID)
}

type VoteConflict struct {
    UserID string
}

func (e VoteConflict) Error() string {
    return fmt.Sprintf("пользователь '%s' уже проголосовал", e.UserID)
}