package models

type Vote struct {
    PollID    string
    UserID    string
    OptionIdx int
}