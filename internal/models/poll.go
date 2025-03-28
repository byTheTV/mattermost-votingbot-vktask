package models

import (
    "github.com/google/uuid"
)

type Poll struct {
    ID        string
    Question  string
    Options   []string
    CreatedBy string
    ChannelID string
    Active    bool
}

// NewId генерирует уникальный ID для опроса
func NewId() string {
    return uuid.New().String() // Использует UUID v4
}
