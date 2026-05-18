# LinguistAI Backend

Backend API для платформы изучения языков с использованием AI.

## Технологии

- **Go 1.23**
- **Fiber v2** - HTTP framework
- **GORM** - ORM для PostgreSQL
- **Google Gemini API** - AI движок
- **JWT** - аутентификация
- **Swagger** - API документация

## Архитектура

Проект следует принципам чистой архитектуры согласно `skills/backend/SKILL.md`:

```
internal/
├── app/             # Application bootstrap, маршруты
├── config/          # Конфигурация сервиса
├── dto/             # Request/Response типы
├── handlers/        # HTTP handlers
├── services/        # Бизнес-логика
├── repository/      # Доступ к данным
└── models/          # GORM модели

pkg/                 # Переиспользуемые пакеты
├── apierror/        # Обработка ошибок
├── config/          # Загрузчик конфигурации
├── httputil/        # HTTP утилиты
├── logger/          # Структурированное логирование
├── middleware/      # JWT, CORS, Logger, Recovery
└── postgres/        # Подключение к БД
```

## Быстрый старт

### Установка зависимостей

```bash
# Установить swag CLI для генерации OpenAPI документации
go install github.com/swaggo/swag/cmd/swag@latest

# Загрузить Go модули
go mod download
```

### Генерация Swagger документации

```bash
make swagger
# или
swag init -g cmd/api/docs.go -o docs --parseDependency --parseInternal
```

### Настройка окружения

```bash
cp .env.example .env
# Отредактировать .env и добавить GEMINI_API_KEY
```

### Запуск

```bash
# Development
make dev

# Production build
make build
./bin/api
```

## API Документация

После запуска сервера, Swagger UI доступен по адресу:
- http://localhost:8080/swagger/index.html

### Основные эндпоинты

#### Authentication
- `POST /api/auth/register` - Регистрация
- `POST /api/auth/login` - Вход

#### Profile (требуется JWT)
- `GET /api/profile/` - Список профилей
- `POST /api/profile/select` - Выбор языка
- `GET /api/profile/{id}` - Детали профиля

#### AI (требуется JWT)
- `POST /api/ai/diagnostic/{profileId}` - Диагностический чат
- `POST /api/ai/generate-plan/{profileId}` - Генерация плана
- `GET /api/ai/plan/{profileId}` - Получить план

#### Chat (требуется JWT)
- `POST /api/chat/tutor/{profileId}` - Чат с AI-тьютором

#### System
- `GET /healthz` - Health check

## Переменные окружения

```env
PORT=8080
DB_DSN=host=localhost user=postgres password=postgres dbname=linguistai port=5432 sslmode=disable
JWT_SECRET=your_jwt_secret
GEMINI_API_KEY=your_gemini_api_key
LOG_LEVEL=info
ALLOWED_ORIGINS=http://localhost:3000
```

## Makefile команды

```bash
make swagger    # Генерация OpenAPI документации
make dev        # Запуск в dev режиме
make build      # Сборка бинарника
make test       # Запуск тестов
make clean      # Очистка артефактов
make deps       # Обновление зависимостей
```

## Архитектурные правила

См. `skills/backend/SKILL.md` и `skills/backend/CHECKLIST.md` для деталей:

1. **Слои**: handlers → services → repository → DB
2. **DTO**: разделение между API и DB типами
3. **DI**: сборка зависимостей в `internal/app`
4. **Middleware**: JWT, CORS, Logger в `pkg/middleware`
5. **OpenAPI**: аннотации на всех публичных эндпоинтах
6. **Errors**: структурированные ошибки через `pkg/apierror`

## База данных

Миграции запускаются автоматически при старте через GORM AutoMigrate.

Модели:
- `User` - пользователи
- `LanguageProfile` - языковые профили
- `StudyPlan` - учебные планы (JSONB)
- `ChatMessage` - история сообщений

## Интеграция с Gemini

См. `internal/services/gemini_service.go`:

- **DiagnosticChat** - определение уровня через диалог
- **GeneratePlan** - генерация JSON структуры плана
- **TutorChat** - режим тьютора

System prompts находятся в `internal/services/prompts.go`.
