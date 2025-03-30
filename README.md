## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/byTheTV/https://github.com/byTheTV/mattermost-votingbot-vktask.git
```
2. 
Создайте .env файл как в примере .env.sample
```bash
cp .env.sample .env
```

3. Заполните .env
MATTERMOST_URL=http://your-mattermost-server:8065
BOT_TOKEN=your_bot_token_here
TARANTOOL_ADDR=tarantool:3301


## Запуск: 
```bash
docker-compose up --build -d
```

### Логи бота
```bash
docker-compose logs -f app
```
