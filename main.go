package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mattermost/mattermost-server/v6/model"
)

func main() {
	// Загружаем переменные из .env файла
	loaderr := godotenv.Load()
	if loaderr != nil {
		fmt.Println("Ошибка загрузки .env файла:", loaderr)
		return
	}

	// Получаем токен из переменной окружения
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		fmt.Println("Ошибка: BOT_TOKEN не задан в .env файле")
		return
	}

	// Настройки бота
	serverURL := "http://localhost:8065" // Замените на ваш URL Mattermost сервера

	// Создаём клиент Mattermost
	client := model.NewAPIv4Client(serverURL)
	client.SetToken(botToken)

	// Создаём сообщение
	post := &model.Post{
		ChannelId: os.Getenv("CHANNEL_ID"), // Замените на ваш ID канала
		Message:   "help me!",
	}

	// Отправляем сообщение
	post, _, err := client.CreatePost(post)
	if err != nil {
		fmt.Println("Ошибка при отправке сообщения:", err)
		return
	}

	fmt.Println("Сообщение успешно отправлено!", post.Id)
}
