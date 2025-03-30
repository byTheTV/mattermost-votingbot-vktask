package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
    "time"
    
	"mattermost-bot/internal/config"
	"mattermost-bot/internal/manager" 
	"mattermost-bot/internal/models"
	"mattermost-bot/internal/repository"
	"mattermost-bot/internal/service"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/tarantool/go-tarantool/v2"
)

func main() {
	cfg := config.Load()

	// Подключение к Tarantool
	dialer := tarantool.NetDialer{
		Address:  cfg.TarantoolAddr,
		User:     "guest",
		Password: "",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := tarantool.Connect(ctx, dialer, tarantool.Opts{})
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer conn.Close()

	log.Println("Successfully connected!")

	

	// Инициализация репозиториев
	pollRepo := repository.NewTarantoolPollRepo(conn)
	voteRepo := repository.NewTarantoolVoteRepo(conn)

	// Сервисный слой
	pollService := service.NewVotingService(pollRepo, voteRepo)

	/*
    schemaManager := repository.NewSchemaManager(conn)
    if err := schemaManager.Init(); err != nil {
        log.Fatal("Ошибка инициализации схемы:", err)
    }
	*/

	// Инициализация CommandManager
	cmdManager := manager.NewCommandManager(pollService) 

	// Настройка бота Mattermost
	bot, err := models.NewBot(cfg.MattermostURL, cfg.BotToken)
	if err != nil {
		log.Fatal(err)
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
					log.Printf("Ошибка парсинга сообщения: %v", err)
					continue
				}
				// Заменяем handleCommand на вызов менеджера
				cmdManager.ProcessCommand(bot, &post)
			}
		}
	}()

	// Ожидание сигнала завершения
	<-sigChan
	log.Println("Завершение работы...")
	close(eventChan)
}
