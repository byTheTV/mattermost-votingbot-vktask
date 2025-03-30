package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
    "time"


	
    "go.uber.org/zap"
	"mattermost-bot/internal/logger"
	"mattermost-bot/internal/config"
	"mattermost-bot/internal/manager" 
	"mattermost-bot/internal/models"
	"mattermost-bot/internal/repository"
	"mattermost-bot/internal/service"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/tarantool/go-tarantool/v2"
)

func main() {
	logger.Init()
	defer logger.L().Sync()

	log := logger.Component("main")

	cfg := config.Load()

	// Подключение к Tarantool
	dialer := tarantool.NetDialer{
		Address:  cfg.TarantoolAddr,
		User:     "guest",
		Password: "",
	}

	log.Info("Starting application")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := tarantool.Connect(ctx, dialer, tarantool.Opts{})
	if err != nil {
		log.Fatal("Connection error:", zap.Error(err))
	}
	defer conn.Close()

	log.Info("Successfully connected!")

	
	// Инициализация репозиториев
	pollRepo := repository.NewTarantoolPollRepo(conn)
	voteRepo := repository.NewTarantoolVoteRepo(conn)

	// Сервисный слой
	pollService := service.NewVotingService(pollRepo, voteRepo)

	// Инициализация CommandManager
	cmdManager := manager.NewCommandManager(pollService) 

	// Настройка бота Mattermost
	botLogger := logger.Component("bot")
	bot, err := models.NewBot(cfg.MattermostURL, cfg.BotToken, botLogger)
	
	if err != nil {
		log.Fatal("Mattermost bot creation error", zap.Error(err))
	}
	defer bot.Ws.Close()

	// Канал для событий WebSocket
	eventChan := make(chan *model.WebSocketEvent)
	go bot.Listen(eventChan)

	// Обработка сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Главный цикл
	go func() {
		for event := range eventChan {
			if event.EventType() == model.WebsocketEventPosted {
				var post model.Post
				err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post)
				if err != nil {
					log.Error("Ошибка парсинга сообщения", zap.Error(err))
					continue
				}
				cmdManager.ProcessCommand(bot, &post)
			}
		}
	}()

	// Ожидание сигнала завершения
	<-sigChan
	log.Info("Завершение работы...")
	close(eventChan)
}
