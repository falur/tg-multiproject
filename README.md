# tg-multiproject

Telegram-бот для удалённой разработки через Claude Code CLI. Позволяет управлять несколькими проектами, отправлять задачи Claude, стримить прогресс в чат, переключаться между plan/edit режимами и деплоить через Ansible.

## Требования

- Go 1.22+
- Claude Code CLI (`npm install -g @anthropic-ai/claude-code`)
- Telegram Bot Token (получить у [@BotFather](https://t.me/BotFather))
- GitHub CLI `gh` (опционально, для создания PR)

## Быстрый старт

### 1. Клонирование и настройка

```bash
git clone <repo-url> tg-multiproject
cd tg-multiproject
cp .env.example .env
```

Отредактируйте `.env`:

```
TELEGRAM_TOKEN=ваш-токен-от-botfather
ALLOWED_USER_ID=ваш-telegram-id
CLAUDE_BINARY=claude
PROJECTS_DIR=./projects
DATABASE_PATH=./data/bot.db
```

Узнать свой Telegram ID можно у бота [@userinfobot](https://t.me/userinfobot).

### 2. Сборка и запуск

```bash
make build
make run
```

Или напрямую:

```bash
go mod tidy
go run ./cmd/bot
```

## Использование бота

### Команды

| Команда | Описание |
|---------|----------|
| `/start` | Главное меню |
| `/cancel` | Отмена текущей задачи Claude |

### Работа с проектами

1. `/start` — открывает главное меню
2. **Create Project** — создание нового проекта:
   - Введите имя проекта
   - Введите GitHub URL для клонирования или нажмите **Skip** для пустого проекта
3. **My Projects** — список проектов, выбор активного

### Выполнение задач

1. Выберите проект из списка
2. Отправьте текст задачи в чат
3. Бот запустит Claude Code CLI и будет стримить прогресс
4. По завершении — результат в чате (или файлом, если > 4096 символов)

### Режимы работы

- **Plan** (по умолчанию) — Claude работает в режиме `--permission-mode plan`, только читает и планирует
- **Edit** — Claude может читать, редактировать, создавать файлы и выполнять команды

Переключение через кнопку в контексте проекта.

### Сессии

Каждый запуск Claude создаёт сессию. Можно возобновить предыдущую сессию через кнопку **Sessions** — Claude продолжит с сохранённым контекстом (`--resume`).

## Структура проекта

```
tg-multiproject/
├── cmd/bot/main.go                 # Точка входа
├── internal/
│   ├── config/config.go            # Конфигурация из переменных окружения
│   ├── bot/
│   │   ├── bot.go                  # Инициализация бота, роутинг
│   │   ├── middleware.go           # Авторизация по Telegram ID
│   │   ├── keyboards.go           # Inline-клавиатуры
│   │   ├── handler_start.go       # /start
│   │   ├── handler_projects.go    # Создание и выбор проектов
│   │   ├── handler_task.go        # Отправка задач, стриминг результатов
│   │   ├── handler_mode.go        # Переключение plan/edit
│   │   ├── handler_cancel.go      # Отмена задачи
│   │   └── handler_session.go     # Список и возобновление сессий
│   ├── claude/
│   │   ├── runner.go              # Запуск Claude CLI, парсинг NDJSON-стрима
│   │   └── events.go              # Структуры событий стрима
│   ├── storage/
│   │   ├── sqlite.go              # Инициализация SQLite, миграции
│   │   ├── project.go             # CRUD проектов
│   │   └── session.go             # CRUD сессий
│   ├── github/github.go           # git clone, git pull, gh pr create
│   └── state/state.go             # In-memory FSM (состояния пользователя)
├── deploy/ansible/                 # Ansible-плейбуки для деплоя
├── Makefile
├── .env.example
└── .gitignore
```

## Деплой на сервер

### Подготовка

1. Укажите IP сервера в `deploy/ansible/inventory.ini`
2. Задайте переменные в `deploy/ansible/deploy.yml` или через `-e`:
   - `telegram_token`
   - `allowed_user_id`
   - `repo_url` (для деплоя через GitHub)

### Настройка сервера (первый раз)

```bash
make setup-server
```

Установит Go, Node.js, Claude CLI, gh CLI, создаст пользователя `bot` и директории.

### Деплой

Через GitHub (по умолчанию):

```bash
make deploy
```

Через rsync (с локальной машины):

```bash
make deploy-local
```

### Управление сервисом

```bash
make status    # Статус systemd-сервиса
make logs      # Логи в реальном времени
make restart   # Перезапуск бота
```

## Makefile-команды

| Команда | Описание |
|---------|----------|
| `make build` | Сборка бинарника в `bin/` |
| `make run` | Запуск локально |
| `make test` | Запуск тестов |
| `make lint` | Линтер (golangci-lint) |
| `make setup-server` | Настройка сервера через Ansible |
| `make deploy` | Деплой через GitHub |
| `make deploy-local` | Деплой через rsync |
| `make status` | Статус сервиса |
| `make logs` | Логи сервиса |
| `make restart` | Перезапуск сервиса |

## Безопасность

- Бот доступен только одному пользователю (проверка по `ALLOWED_USER_ID`)
- Файл `.env` не попадает в git
- На сервере `.env` имеет права `0600`
- SQLite работает в режиме WAL с одним подключением
