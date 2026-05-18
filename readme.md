# LinguistAI: AI-Adaptive Language Learning Platform

## 🌟 Vision
**LinguistAI** — это платформа для изучения языков нового поколения. Система использует **Google Gemini** для динамического создания учебного плана. ИИ оценивает уровень пользователя через живое общение и адаптирует сложность контента в реальном времени.

---

## 🛠 Tech Stack (2026 Edition)
### Backend
- **Language:** Golang 1.26
- **Framework:** Fiber v3
- **Database:** PostgreSQL + GORM
- **AI Engine:** Google Gemini API
- **Auth:** JWT (Standard HS256)

### Frontend
- **Framework:** Next.js (App Router)
- **UI Components:** **Shadcn/UI** + Tailwind CSS
- **Icons:** Lucide React
- **Animations:** Framer Motion

---

## 📊 Database Schema (MVP)
- **User**: ID, Email, PasswordHash, CreatedAt.
- **LanguageProfile**: ID, UserID, TargetLanguage, CurrentLevel (A1-C2), ProgressPercentage.
- **StudyPlan**: ID, ProfileID, Structure (JSONB - модули и темы), IsGenerated.
- **ChatHistory**: ID, UserID, Role (user/assistant), Context (test/lesson), Content.

---

## 🚀 User Flow (The MVP Engine)

### 1. Account & JWT Auth
- Регистрация и логин. Получение JWT токена.
- Хранение токена в Cookie

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

---

## 🛣 API Endpoints

| Method | Route | Description | Auth |
|:-------|:------|:------------|:-----|
| POST | `/api/auth/register` | Регистрация | No |
| POST | `/api/auth/login` | Вход | No |
| POST | `/api/profile/select` | Выбор изучаемого языка | Yes |
| POST | `/api/ai/diagnostic` | Сообщение в режиме теста | Yes |
| POST | `/api/ai/generate-plan` | Создать план на основе уровня | Yes |
| POST | `/api/chat/tutor` | Общение с ИИ-тьютором | Yes |

---

## 🎨 UI/UX (Next.js + Shadcn)
- **Dashboard:** Чистый интерфейс с использованием Tabs и ScrollArea.
- **Chat Window:** Сообщения пользователя и ИИ (Markdown поддержка для правил).
- **Plan View:** Прогресс-бары (Progress) и чекбоксы выполненных уроков.
- **Theme:** Поддержка темной и светлой темы.

---

## ⚙️ Environment Variables (.env)
```env
PORT=8080
DB_DSN=host=localhost user=postgres password=postgres dbname=linguistai port=5432 sslmode=disable
JWT_SECRET=future_key_2026
GEMINI_API_KEY=your_google_ai_key
NEXT_PUBLIC_API_URL=http://localhost:8080
```

# linguist-hub
