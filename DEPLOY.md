# Deploy: Neon + Fly.io + Vercel

```
Browser → Vercel (Next.js) → Fly.io (Go Fiber) → Neon (Postgres)
                          ↘ OpenRouter LLM API (from Fly only)
```

> Railway trial expired; backend targets **Fly.io** (`backend/fly.toml`).  
> Fly requires a payment method on the account (free allowance after card).

## 1. Neon

1. Create project + database (done: `linguist-hub` / `dry-union-33518524`).
2. Connection string with `sslmode=require` → `DATABASE_URL` or `DB_DSN`.

## 2. Fly.io (backend)

```bash
cd backend
fly auth login
fly apps create linguist-hub-api   # needs billing on account
fly secrets set \
  APP_ENV=production \
  ENABLE_SWAGGER=false \
  DATABASE_URL='postgresql://...' \
  JWT_SECRET='...' \
  LLM_API_KEY='...' \
  LLM_BASE_URL='https://openrouter.ai/api/v1' \
  LLM_MODEL='openai/gpt-4o-mini' \
  ALLOWED_ORIGINS='https://your-app.vercel.app'
fly deploy
```

Healthcheck: `GET /healthz`  
Public URL: `https://linguist-hub-api.fly.dev`

## 3. Vercel (frontend)

- Root Directory: `frontend`
- Framework: Next.js
- Env (**Production + Preview**, at **build**):

| Key | Value |
|-----|--------|
| `NEXT_PUBLIC_API_URL` | `https://linguist-hub-api.fly.dev` |
| `NEXT_PUBLIC_WS_URL` | `wss://linguist-hub-api.fly.dev` |

Redeploy frontend after Fly URL is live.

## 4. Smoke test

1. Open Vercel URL → register / login  
2. Hub → diagnostic chat  
3. Voice/tutor WS (if used)  
4. Confirm CORS (no browser console block)  

## Auth note

JWT is stored in **localStorage** and sent as `Authorization: Bearer`.  
WebSocket uses `Sec-WebSocket-Protocol: bearer,<jwt>` (query `?token=` still accepted as fallback).
