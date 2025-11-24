# Conversational AI Flow

This document describes the architecture and deterministic flow used in the perfume recommendation chat.

## Slots Order

Slots represent user preferences extracted from messages:
- Occasions (work, date, party, casual, day/night)
- Climate (cold, hot, humid, dry)
- Intensity (soft, moderate, strong)
- Accords (citrus, floral, woody, oriental, etc.)
- Budget (low, medium, high)
- Longevity (short, medium, long)
- Seasons (summer, winter, autumn, spring)

Clarification order (NextMissing):
1) Occasions → 2) Climate → 3) Intensity → 4) Accords → 5) Budget → 6) Longevity → 7) Seasons

## NextMissing

Pure function that returns the next empty slot according to the order above. Used for:
- Guiding the clarification question (at most 1 per turn)
- Showing a "next step hint" in the client

## Profile Hash

- Each slots profile is hashed (short SHA-1) in `ProfileHash(slots)`.
- The hash is used as a cache key to avoid recomputing candidates for the same profile within a short TTL (5 minutes).
- The hash includes all relevant slots, keeping caching deterministic.

## Prompt Format

Minimalist prompt in Portuguese, with clear objectives:
- Context: preferences summary (if any) + user message
- Candidates: up to 3 suggestions (name, brand, reason)
- Task: respond shortly; if a critical slot is missing, ask only 1 focused question on that slot

The model receives instructions to avoid redundancy and be objective.

## Implementation and Maintenance

- Module `internal/ai/slots`:
  - `Parse/Merge`: extraction and evolution of preferences (pure).
  - `NextMissing`, `ProfileHash`, `BuildSearchQuery`: pure utility functions.
- Service `ChatService`:
  - Uses cache by `ProfileHash`.
  - Runs parallel retrieval from FTS, embeddings, and slots; combines and takes top-k.
  - Reduces token cap (180) and keeps prompt short.

This approach prioritizes predictability, low cost, and diversity from multiple sources.
