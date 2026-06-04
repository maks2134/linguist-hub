# Agents (Linguist-hub)

Проектные субагенты для Go backend лежат в [`.cursor/agents/`](.cursor/agents/).

## Доступные агенты

| Имя | Роль |
|-----|------|
| `go-explorer` | Аудит кодовой базы, только чтение |
| `go-architect` | Архитектура, API, БД, план без реализации |
| `go-coder` | Реализация по плану |
| `go-tester` | Тесты, `-race`, покрытие |

Правило делегирования: [`.cursor/rules/go-agents.mdc`](.cursor/rules/go-agents.mdc).

## Примеры вызова

```
Use go-explorer subagent to найти, как устроен AIService и i18n в backend
```

```
Use go-architect subagent to спланировать вынос DI в containers для Linguist-hub
```

```
Use go-coder subagent to реализовать шаг 1 из плана go-architect (вставить план)
```

```
Use go-tester subagent to покрыть тестами internal/services/ai_service.go
```

## Backend

- Модуль: `github.com/makskozlov/linguist-hub`
- Путь: `backend/`
- Команды: `cd backend && go build ./...` и `go test ./...`
