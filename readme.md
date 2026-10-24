# LinguistAI: AI-Adaptive Language Learning Platform

## 🌟 Vision
**LinguistAI** — это платформа для изучения языков нового поколения. Система использует **Google Gemini** для динамического создания учебного плана. ИИ оценивает уровень пользователя через живое общение и адаптирует сложность контента в реальном времени.

## 🚀 User Flow (The MVP Engine)

### 1. Account & JWT Auth
- Регистрация и логин. Получение JWT токена.
- Хранение токена в **localStorage**; API: `Authorization: Bearer <token>`
- WebSocket: `Sec-WebSocket-Protocol: bearer,<token>`

### Deploy
See [DEPLOY.md](./DEPLOY.md) — Neon (Postgres) + Railway (Go API) + Vercel (Next.js).

### 2. Language Selection & Hub
- Пользователь выбирает язык (например, German).
- Попадает в **Hub**, где ИИ предлагает пройти диагностику уровня.

### 3. AI Diagnostic (Gemini Interview)
- Чат-интерфейс, где Gemini задает 5-8 вопросов разной сложности.
- **Логика:** После диалога ИИ возвращает уровень (A1-C2) и слабые места. Уровень сохраняется в профиль.

### 4. Plan Generation
- На основе уровня Gemini генерирует JSON-структуру курса: "Модуль 1: Грамматика", "Модуль 2: Путешествия" и т.д.
- Веб-интерфейс отрисовывает этот план в виде карточек (Shadcn Card).

### 5. AI Tutor Chat
- Пользователь выбирает тему и начинает общение.
- ИИ выступает в роли тьютора: исправляет ошибки, объясняет правила и поддерживает диалог.

