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

3. Настройка .env
   - Запустить mattermost + tarantool
   - Зарегистрируйтесь в Mattermost
   - Разрешите создание ботов: System Console -> Bot Account (Integration)
   - Создайте бота: Integrations -> Bot Accounts
   - Получите BOT_TOKEN
   - Вставьте токен в файл .env

4. Запустить бота

## Запуск: 

1. Запустить mattermost + tarantool
```bash
docker-compose up --build mattermost tarantool -d
```

2. Запустить бота
```bash
docker-compose up --build app -d
```

### Логи бота
```bash
docker-compose logs -f app
```

## Команды бота

- /poll Вопрос Вариант1 Вариант2 Вариант3 - Создать новое голосование
- /vote [poll_id] [option_index] - Проголосовать (нумерация вариантов начинается с 1)
- /results [poll_id] - Посмотреть текущие результаты
- /close [poll_id] - Завершить голосование
- /delete_poll [poll_id] - Удалить голосование
- /polls - Показать список всех голосований

P.S. Если что по оформлению команд бот всё подскажет

 <3 thetv

