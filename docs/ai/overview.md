# AI Perfume Recommendations: Overview

This document summarizes the approach to build a free-tier friendly AI chat for perfume recommendations using Retrieval-Augmented Generation (RAG).

## Architecture Overview

- **Data Model**: Products have tags (accords, occasions, seasons, intensity) and presentation fields (slug, thumbnail) filled during product registration. Embeddings are precomputed automatically on creation.
- **Embeddings**: Vector representations of each product's description/tags, precomputed automatically on product creation via Hugging Face API (Qwen/Qwen3-Embedding-8B) and stored in Postgres (pgvector). At query time, embed the user message via the configured provider and compute cosine similarity to retrieve top-k candidates.
- **Hybrid Retrieval**: Combines Full-Text Search (FTS), semantic embeddings, and direct slot matching (e.g., accords) in parallel, then ranks top-k. No deterministic fallbacks; always uses all available sources.
- **LLM Layer**: Small instruct model (Gemma-2-2b-it via Hugging Face Inference API) to generate short, friendly responses in Portuguese.
- **Slots Parsing**: Deterministic extraction of preferences (occasions, climate, intensity, accords, budget, longevity, seasons) from user messages.
- **Caching**: Profile-based caching (hash of slots + query) for 5 minutes to reduce computations.
- **Safety**: Sanitize user input to remove PII; rate limiting; timeouts.

## Key Components

- **Providers**: Abstracted AI providers (Ollama for dev, Hugging Face for prod) for embeddings and LLM.
- **Chat Service**: Orchestrates sanitization, slots parsing/merging, parallel retrieval, prompt building, and LLM generation.
- **Repository**: Handles product search (FTS + embeddings), embedding storage.
- **Slots Module**: Pure functions for parsing, merging, and query building.

## Free-Tier Optimization

- **Dev**: Ollama local models (tinyllama for LLM, nomic-embed-text for embeddings).
- **Prod**: Hugging Face free Inference API (~30k calls/month for embeddings/LLM).
- **Low Token Usage**: Prompts ≤ 180 tokens; responses short and focused.