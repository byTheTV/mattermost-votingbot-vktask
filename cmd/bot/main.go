package main

import (
    "log"
    "context"
    "strings"
    "encoding/json" // Добавлен недостающий импорт
 	"os"
    "os/signal"
	"syscall"

    "mattermost-bot/internal/config"
    "mattermost-bot/internal/handlers"
    "mattermost-bot/internal/models"
    "mattermost-bot/internal/repository"
    "mattermost-bot/internal/service"
    "github.com/mattermost/mattermost-server/v6/model"
    "github.com/tarantool/go-tarantool/v2"
)

func main() {
    cfg := config.Load()
    ctx := context.Background()

    // Подключение к Tarantool
    dialer := tarantool.NetDialer{
        Address:  cfg.TarantoolAddr,
        User:     "guest",
        Password: "",
    }

    // Подключение к Tarantool
    conn, err := tarantool.Connect(ctx, dialer, tarantool.Opts{})
    if err != nil {
        log.Fatal("Ошибка подключения к Tarantool:", err)
    }
	defer conn.Close()

    // Инициализация репозиториев
    pollRepo := repository.NewTarantoolPollRepo(conn)
    voteRepo := repository.NewTarantoolVoteRepo(conn)

    // Сервисный слой
    pollService := service.NewVotingService(pollRepo, voteRepo)

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
                handleCommand(&post, bot, pollService)
            }
        }
    }()

    // Ожидание сигнала завершения
    <-sigChan
    log.Println("Завершение работы...")
    close(eventChan) // Корректное закрытие канала

}


func handleCommand(post *model.Post, bot *models.Bot, s service.PollService) {
    parts := strings.Fields(post.Message)
    if len(parts) == 0 {
        return
    }

    switch parts[0] {
    case "/poll":
        handlers.HandleCreatePoll(s, bot, post, parts[1:])
    case "/vote":
        handlers.HandleVote(s, bot, post, parts[1:])
    case "/results":
        handlers.HandleResults(s, bot, post, parts[1:])
    case "/close":
        handlers.HandleClosePoll(s, bot, post, parts[1:])
    case "/polls":
        handlers.HandleListPolls(s, bot, post)
    }
}
