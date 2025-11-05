import { useEffect, useRef, useState, useCallback } from "react";
import type { AxiosError } from "axios";
import type { Product } from "../types/product";
import { searchProducts, listProducts } from "../services/product";

/**
 * Debounced product search hook with cancelation and stale-response protection.
 * Strategy:
 * - Debounce: setTimeout (~350ms) managed via ref; cleared on unmount/changes
 * - Min length: require >= 2 chars to search; else show latest list
 * - Cancelation: AbortController aborts prior in-flight requests on changes
 * - Stale guard: sequence counter ensures only the last response updates state
 * - Reset page: when query changes, page resets to 1
 * - Dedupe: ignore repeated requests for the same (query|page|limit|sort)
 */
export function useProductSearch(options?: {
  initialQuery?: string;
  initialPage?: number;
  limit?: number;
  sort?: "relevance" | "latest";
  debounceMs?: number;
}) {
  const {
    initialQuery = "",
    initialPage = 1,
    limit = 12,
    sort = "relevance",
    debounceMs = 350,
  } = options || {};

  const [query, setQuery] = useState<string>(initialQuery);
  const [page, setPage] = useState<number>(initialPage);
  const [results, setResults] = useState<Product[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isSearching, setIsSearching] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const debounceRef = useRef<number | null>(null);
  const abortRef = useRef<AbortController | null>(null);
  const requestSeq = useRef(0);
  const lastKeyRef = useRef<string>("");
  const cacheRef = useRef<Map<string, { items: Product[]; total: number; expiresAt: number }>>(new Map());

  // Clear utilities
  const clear = useCallback(() => {
    setQuery("");
    setPage(1);
  }, []);

  const cancelInFlight = useCallback(() => {
    if (abortRef.current) {
      abortRef.current.abort();
    }
    abortRef.current = null;
  }, []);

  const runFetch = useCallback(
    async (opts: { q: string; p: number }) => {
      const currentSeq = ++requestSeq.current;
      const { q, p } = opts;
      const trimmed = q.trim();
      const tooShort = trimmed.length < 2;
      const key = `${trimmed}|${p}|${limit}|${sort}`;

      if (lastKeyRef.current === key) {
        return;
      }
      lastKeyRef.current = key;

      setIsLoading(true);
      setError(null);
      setIsSearching(!tooShort);

      // Try cache first
      const cached = cacheRef.current.get(key);
      const now = Date.now();
      if (cached && cached.expiresAt > now) {
        setResults(cached.items);
        setTotal(cached.total);
        setIsLoading(false);
        return;
      }

      cancelInFlight();
      const controller = new AbortController();
      abortRef.current = controller;

  try {
        if (tooShort) {
          // Default: latest products
          const data = await listProducts({ limit, signal: controller.signal });
          if (requestSeq.current !== currentSeq) return; // stale
          const payload = { items: data, total: data.length };
          // Cache
          cacheRef.current.set(key, { ...payload, expiresAt: now + 5 * 60_000 });
          if (cacheRef.current.size > 50) {
            const firstKey = cacheRef.current.keys().next().value as string | undefined;
            if (firstKey) cacheRef.current.delete(firstKey);
          }
          setResults(payload.items);
          setTotal(payload.total);
        } else {
          const data = await searchProducts({ query: trimmed, page: p, limit, sort, signal: controller.signal });
          if (requestSeq.current !== currentSeq) return; // stale
          const items = Array.isArray((data as any)?.items) ? (data as any).items as Product[] : [];
          const totalVal = typeof (data as any)?.total === 'number' ? (data as any).total as number : items.length;
          // Cache
          cacheRef.current.set(key, { items, total: totalVal, expiresAt: now + 5 * 60_000 });
          if (cacheRef.current.size > 50) {
            const firstKey = cacheRef.current.keys().next().value as string | undefined;
            if (firstKey) cacheRef.current.delete(firstKey);
          }
          setResults(items);
          setTotal(totalVal);
        }
      } catch (e: any) {
        // Ignore aborts/cancels
        if (e?.name === "CanceledError" || e?.name === "AbortError") {
          return;
        }
        // Allow retries for same key after an error
        lastKeyRef.current = "";

        const err = e as AxiosError;
        const status = (err.response?.status as number | undefined) ?? undefined;

        // If we were running a search (>=2 chars), degrade gracefully to empty
        if (!tooShort) {
          const payload = { items: [] as Product[], total: 0 };
          cacheRef.current.set(key, { ...payload, expiresAt: now + 60_000 });
          setResults(payload.items);
          setTotal(payload.total);
          setError(null);
          return;
        }

        // For listing latest (tooShort), keep a friendly error
        if (status === 404) {
          const payload = { items: [] as Product[], total: 0 };
          cacheRef.current.set(key, { ...payload, expiresAt: now + 60_000 });
          setResults(payload.items);
          setTotal(payload.total);
          setError(null);
        } else {
          setError("Failed to load products. Please try again.");
        }
      } finally {
        if (requestSeq.current === currentSeq) {
          setIsLoading(false);
        }
      }
    },
    [cancelInFlight, limit, sort]
  );

  // Debounced effect
  useEffect(() => {
    // reset to page 1 when query changes
    setPage(1);
  }, [query]);

  useEffect(() => {
    if (debounceRef.current) {
      window.clearTimeout(debounceRef.current);
    }
    debounceRef.current = window.setTimeout(() => {
      runFetch({ q: query, p: page });
    }, debounceMs);

    return () => {
      if (debounceRef.current) {
        window.clearTimeout(debounceRef.current);
      }
    };
  }, [query, page, debounceMs, runFetch]);

  // Immediate submit
  const submitNow = useCallback(() => {
    if (debounceRef.current) {
      window.clearTimeout(debounceRef.current);
    }
    runFetch({ q: query, p: 1 });
    setPage(1);
  }, [query, runFetch]);

  // Clean up on unmount
  useEffect(() => {
    return () => {
      if (debounceRef.current) {
        window.clearTimeout(debounceRef.current);
      }
      cancelInFlight();
    };
  }, [cancelInFlight]);

  return {
    // state
    query,
    setQuery,
    page,
    setPage,
    limit,
    results,
    total,
    isLoading,
    isSearching,
    error,
    // actions
    clear,
    submitNow,
  } as const;
}
