# Chat RAG (Aroma Sense)

Architecture document for the conversational retrieval (RAG) flow focused on perfumes, optimized for free-tier and short responses in PT-BR.

## Objective

Deliver perfume recommendations via chat with:

- Determinism and explainability (slots + short reasons)
- LLM only for concise writing (≤ 180 tokens) and at most 1 question

## Overall Flow Overview

```
```
Frontend ── POST /ai/chat ─▶ Backend (Handler)
                     │
                     ▼
                Sanitize (PII)  ──▶ Greeting/Off-topic steer (if applicable)
                     │
                     ▼
                Slots.Parse+Merge (occasion, season, climate, intensity, accords, budget, longevity)
                     │
                     ├──► Query Embedding (Ollama/HuggingFace API) ─┐
                     │                                               │
                     ├──► FTS (Postgres)                             │
                     │                                               ▼
                     └──► Direct Slot Matching ──▶ Combine & Dedup ──▶ top‑k (k ≤ 5)
                                                                 │
                                                                 ▼
                                                Select up to 3 candidates + "reason" explainable
                                                                 │
                                                                 ▼
                                   Short prompt for LLM (≤ 180 tokens; 0–1 question if critical slot missing)
                                                                 │
                                                                 ▼
                                                       Short response in PT-BR + suggestions
```
```

## Components

### 1) Sanitization and Topic Guard

- Remove PII (email, phone, URLs) from `message` before any use (`internal/ai/sanitize.go`).
- If message is a simple greeting ("hi", "hello"…), return a short introduction and `follow_up_hint` (no LLM).
- If off-topic and no preferences yet (empty slots), return polite guidance to talk about perfumes (no LLM).

### 2) Slots (Deterministic Preferences)

- `Slots.Parse(msg)`: extracts signals for occasion, seasons, climate, intensity, accords, budget, longevity.
- `Slots.Merge(a, b)`: merges preferences across conversation (dedup, stable order).
- `NextMissing(s)`: indicates which slot to ask first.
- `ProfileHash(s)`: stable key for cache.

### 3) Embeddings (Query and Products)

- Local model via Ollama for dev (`/api/embeddings`).- Online API via Hugging Face Inference for prod (free-tier, ~30k calls/month).
- Products have precomputed embeddings (automatic on creation via API) stored in `product_embeddings` (pgvector).
- Per query, generate 1 embedding of sanitized text via the configured provider.

### 4) Hybrid Retrieval (top-k)

- Parallel execution: FTS (strict query with all slots), embeddings (semantic similarity), direct slot matching (e.g., exact accords search).
- Combine results, dedup by ID, rank top-k (currently first 5 in order: FTS > Embeddings > Slots).
- No progressive relaxation or fallbacks; always use all sources for diversity.

### 5) Short Prompt for LLM

- Max content: 1–3 candidates (name — brand — reason), user message and summary (if any), PT-BR instructions.
- Policy: ≤ 180 tokens; 0–1 clarification question (if `NextMissing` indicates critical empty slot).
- LLM: Gemma-2-2b-it via Hugging Face Inference API.
- If LLM fails/unavailable: return suggestions with deterministic reasons (no 500 error).

### 6) Caches and Limits

- Cache `(profileHash + query)` → top-k for 5 min.
- Rate-limit per IP in handler.
- Short timeouts.

## Contracts (I/O)

### Request

POST `/ai/chat`

```json
{
  "message": "quero algo cítrico para o verão",
  "session_id": "abc123",
  "history": ["previous messages (optional)"]
}
```

### Response

```json
{
  "reply": "Eu sugiro X e Y…",
  "suggestions": [
    {"id": 1, "name": "...", "brand": "...", "slug": "...", "thumbnail_url": "...", "price": 0, "reason": "..."},
    {"id": 2, "name": "...", "brand": "...", "slug": "...", "thumbnail_url": "...", "price": 0, "reason": "..."}
  ],
  "follow_up_hint": "Prefere algo cítrico, floral ou amadeirado?"
}
```

## Prompt (Example)

General guidelines embedded in prompt:

- "You are a perfume assistant. Respond in Portuguese, direct and short. If the message is off-topic perfumes/fragrances, politely redirect explaining you can only help with perfumes and ask for 1 preference (occasion, accord, intensity or budget)."
- If `NextMissing` ≠ "": ask only 1 specific question.

Candidates block (max 3 lines):

```
- Name — Brand | Reason
```

## Fallbacks

No deterministic fallbacks needed; parallel retrieval from FTS, embeddings, and slots always provides diverse candidates. If a source fails (e.g., embeddings unavailable), others compensate.

## Environment Variables

- `AI_PROVIDER` (huggingface)
- `AI_EMB_MODEL` (Qwen/Qwen3-Embedding-8B)
- `AI_LLM_MODEL` (Gemma-2-2b-it)
- `AI_API_KEY` (for Hugging Face)
- `AI_EMB_BASE_URL` (optional, defaults to Hugging Face)
- `AI_LLM_BASE_URL` (optional, defaults to Hugging Face)

## Performance and Limits (Free-Tier)

- 1 embedding per query, k ≤ 5, prompt ≤ 180 tokens.
- Small/medium catalog: cosine in pgvector is sufficient.
- Parallel goroutines for retrieval to reduce latency (~1.5s average).

## Observability

- Log latencies of: sanitize+slots, embeddings, FTS, combine, LLM.
- Count provider failures for future tuning.

## Security

- Sanitize PII always before embeddings/FTS/LLM.
- Rate-limit in handler; timeouts and concurrency limit in Ollama.
