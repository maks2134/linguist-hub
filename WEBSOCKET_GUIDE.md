# WebSocket Real-Time AI Integration Guide

## Overview
Backend и Frontend теперь поддерживают **real-time streaming** ответов от AI через WebSocket.

## Backend Changes

### WebSocket Endpoint
- **URL**: `ws://localhost:8080/api/ws`
- **Auth**: JWT token через `Authorization: Bearer <token>` header

### Message Format

**Client → Server:**
```json
{
  "type": "diagnostic",  // или "tutor", "generate_plan"
  "message": "Hello",
  "topic": "grammar",    // только для tutor
  "profile_id": 1
}
```

**Server → Client (streaming):**
```json
// Chunks (streaming)
{"type": "diagnostic", "chunk": "Hello", "done": false}
{"type": "diagnostic", "chunk": " there!", "done": false}

// Final response
{
  "type": "diagnostic",
  "content": "Full response text",
  "level": "B1",
  "weaknesses": ["grammar", "vocabulary"],
  "is_assessment": true,
  "done": true
}
```

## Frontend Changes

### New Files
1. **`src/lib/websocket.ts`** - WebSocket client class
2. **`src/hooks/useWebSocket.ts`** - React hook для WebSocket
3. Updated **`diagnostic/[profileId]/page.tsx`** - streaming UI
4. Updated **`tutor/[profileId]/page.tsx`** - streaming UI

### Environment Variables
Добавь в `.env.local`:
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

## Features

### 1. Real-time Streaming
- Ответы от AI появляются **по мере генерации** (word by word)
- Анимированный cursor показывает процесс печати
- Auto-scroll при получении новых chunks

### 2. Automatic Reconnection
- Автоматический реконнект при потере соединения
- Exponential backoff (1s → 2s → 4s → 8s → 16s)
- Максимум 5 попыток

### 3. Error Handling
- Отображение ошибок через toast notifications
- Graceful fallback при проблемах с WebSocket

## Usage Example

```typescript
import { useWebSocket } from "@/hooks/useWebSocket";

function ChatPage() {
  const { connected, sendMessage, onMessage } = useWebSocket();
  const [streamingContent, setStreamingContent] = useState("");

  useEffect(() => {
    if (!connected) return;

    const unsubscribe = onMessage((response) => {
      if (response.chunk) {
        setStreamingContent(prev => prev + response.chunk);
      }
      if (response.done) {
        // Final message received
        console.log("Complete:", response.content);
      }
    });

    return () => unsubscribe();
  }, [connected, onMessage]);

  const handleSend = () => {
    sendMessage({
      type: "diagnostic",
      message: "Test message",
      profile_id: 1,
    });
  };
}
```

## Testing

### 1. Start Backend
```bash
cd backend
go run cmd/api/main.go
```

### 2. Start Frontend
```bash
cd frontend
npm run dev
```

### 3. Open Browser
- Navigate to `/diagnostic/1` or `/tutor/1`
- Start typing - you'll see AI responses stream in real-time!

## Production Deployment

### Environment Variables
```bash
# Production
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
NEXT_PUBLIC_WS_URL=wss://api.yourdomain.com  # Note: wss:// for HTTPS

# Backend
GEMINI_API_KEY=your_production_key
```

### Nginx Configuration
```nginx
location /api/ws {
    proxy_pass http://backend:8080;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

## Troubleshooting

### WebSocket не подключается
1. Проверь что backend запущен
2. Проверь `.env.local` - правильный ли WS_URL
3. Проверь JWT token в localStorage
4. Open DevTools → Network → WS tab

### Streaming не работает
1. Check Backend logs: `[Gemini] operation failed`
2. Проверь Gemini API key
3. Проверь quota на Gemini API

### Connection refused
- Backend должен быть на том же хосте или CORS настроен
- Для production используй `wss://` вместо `ws://`
