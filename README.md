# Aroma Sense

A full-stack ecommerce platform for fragrances with a Go/Gin API and a React + Vite frontend. It integrates payments, shipping quotes, AI-powered search/chat, email notifications, and S3-compatible media storage.

## Tech Stack
- Backend: Go 1.25, Gin, GORM, PostgreSQL, Swagger, Air (dev hot reload)
- Frontend: React 18, Vite, TypeScript, Tailwind
- Integrations: Stripe payments, SuperFrete shipping, Supabase S3 storage, SMTP, AI providers (Ollama/Hugging Face-compatible)

## Project Layout
- `backend/` – API (`cmd/api`), services, integrations, migrations
- `frontend/` – React app (Vite), Tailwind, routing
- `docs/` – Swagger assets and release/process docs

## Quick Start (Docker)
1. Create `backend/.env` with the environment variables below.
2. Start everything: `docker compose -f docker-compose.dev.yml up --build`
3. Frontend: http://localhost:5173
4. API health: http://localhost:8080/healthz

## Backend Development (without Docker)
1. Install Go 1.25 and PostgreSQL; create a database.
2. Export environment variables (or use a `backend/.env` file).
3. From `backend/`: `go run ./cmd/api`
4. Optional: `go install github.com/air-verse/air@latest && air -c .air.toml` for hot reload.

## Frontend Development (without Docker)
1. From `frontend/`: `npm install`
2. Create `.env.local` with `VITE_API_URL=http://localhost:8080` and `VITE_APP_ENV=development`.
3. Run dev server: `npm run dev` (defaults to port 5173).

## Environment Variables (backend)
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=aroma_sense

# Auth & CORS
JWT_SECRET=change-me-32-chars-min
ALLOWED_ORIGINS=http://localhost:5173
FRONTEND_URL=http://localhost:5173
ENABLE_SWAGGER=true

# Storage (Supabase S3)
SUPABASE_S3_ENDPOINT=https://xxx.supabase.co/storage/v1/s3
SUPABASE_S3_REGION=us-east-1
SUPABASE_S3_ACCESS_KEY=...
SUPABASE_S3_SECRET_KEY=...
SUPABASE_BUCKET=aroma-sense
SUPABASE_PUBLIC_URL=https://xxx.supabase.co/storage/v1/object/public/aroma-sense

# Email (SMTP)
SMTP_HOST=smtp.yourprovider.com
SMTP_PORT=587
SMTP_USERNAME=...
SMTP_PASSWORD=...
SMTP_FROM=no-reply@yourdomain.com

# Payments (Stripe)
STRIPE_SECRET_KEY=sk_test_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx

# Shipping (SuperFrete)
SHIPPING_BASE_URL=https://sandbox.superfrete.com
SHIPPING_QUOTES_PATH=/api/v0/calculator
SHIPPING_USER_AGENT=aroma-sense-dev
SHIPPING_SERVICES=1,2,17
SHIPPING_ORIGIN_CEP=00000000
SHIPPING_TOKEN=replace-with-token
SHIPPING_TIMEOUT=30

# AI Providers
AI_PROVIDER=ollama
AI_LLM_BASE_URL=http://localhost:11434
AI_LLM_MODEL=tinyllama:latest
AI_EMB_BASE_URL=http://localhost:11434
AI_EMB_MODEL=nomic-embed-text:latest
AI_API_KEY=
AI_TIMEOUT=30
```

## Database Migrations
SQL migrations live in `backend/migrations`. With the `migrate` CLI available, run (example):
```bash
migrate -path backend/migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=require" up
```

## API Docs
Set `ENABLE_SWAGGER=true` to expose Swagger at `/swagger/index.html`. Generated assets are in `backend/docs`.

## Testing
- Backend: from `backend/`, run `go test ./...`
- Frontend: `npm run lint` (tests not defined yet)

## Health Checks
- API: `GET /healthz` returns `{ "status": "ok" }`

## License
Not specified yet. Add one before distributing.
