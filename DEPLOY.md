# Deploy: Neon + Railway + Vercel

```
Browser → Vercel (Next.js) → Railway (Go Fiber) → Neon (Postgres)
                          ↘ Google Gemini (from Railway only)
```

## 1. Neon

1. Create project + database.
2. Copy connection string (pooled or direct) with `sslmode=require`.
3. You can set it as `DATABASE_URL` **or** `DB_DSN` on Railway.

## 2. Railway (backend)

Root directory / Dockerfile: `backend/`

Env:

| Key | Value |
|-----|--------|
| `APP_ENV` | `production` |
| `ENABLE_SWAGGER` | `false` |
| `DATABASE_URL` or `DB_DSN` | Neon URL |
| `JWT_SECRET` | long random |
| `GEMINI_API_KEY` | Google AI key |
| `ALLOWED_ORIGINS` | `https://your-app.vercel.app` (comma-separated OK) |
| `PORT` | Railway sets this |

Healthcheck: `GET /healthz`

## 3. Vercel (frontend)

- Root Directory: `frontend`
- Framework: Next.js
- Env (**Production + Preview**, available at **build**):

| Key | Value |
|-----|--------|
| `NEXT_PUBLIC_API_URL` | `https://<railway-public-host>` |
| `NEXT_PUBLIC_WS_URL` | `wss://<railway-public-host>` |

Redeploy frontend after Railway URL is known.

## 4. Smoke test

1. Open Vercel URL → register / login  
2. Hub → diagnostic chat  
3. Voice/tutor WS (if used)  
4. Confirm CORS (no browser console block)  

## Auth note

JWT is stored in **localStorage** and sent as `Authorization: Bearer`.  
WebSocket uses `Sec-WebSocket-Protocol: bearer,<jwt>` (query `?token=` still accepted as fallback).
