import { useEffect, useRef, useState, useCallback } from "react";
import { isAxiosError } from "axios";
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
 * - Infinite scroll: for latest products, supports loading more pages
 */
export function useProductSearch(options?: {
  initialQuery?: string;
  initialPage?: number;
  limit?: number;
  sort?: "relevance" | "latest";
  debounceMs?: number;
  enableInfiniteScroll?: boolean;
}) {
  const {
    initialQuery = "",
    initialPage = 1,
    limit = 12,
    sort = "relevance",
    debounceMs = 350,
    enableInfiniteScroll = false,
  } = options || {};

  const [query, setQuery] = useState<string>(initialQuery);
  const [page, setPage] = useState<number>(initialPage);
  const [results, setResults] = useState<Product[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isSearching, setIsSearching] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState<boolean>(true);
  const [isLoadingMore, setIsLoadingMore] = useState<boolean>(false);

  const debounceRef = useRef<number | null>(null);
  const abortRef = useRef<AbortController | null>(null);
  const requestSeq = useRef(0);
  const lastKeyRef = useRef<string>("");
  const cacheRef = useRef<Map<string, { items: Product[]; total: number; expiresAt: number }>>(new Map());

  const cancelInFlight = useCallback(() => {
    if (abortRef.current) {
      abortRef.current.abort();
    }
    abortRef.current = null;
  }, []);

  const runFetch = useCallback(
    async (opts: { q: string; p: number; append?: boolean }) => {
      const currentSeq = ++requestSeq.current;
      const { q, p, append = false } = opts;
      const trimmed = q.trim();
      const tooShort = trimmed.length < 2;
      const key = `${trimmed}|${p}|${limit}|${sort}`;

      if (!append && lastKeyRef.current === key) {
        return;
      }
      lastKeyRef.current = key;

      // Set appropriate loading state
      if (append) {
        setIsLoadingMore(true);
      } else {
        setIsLoading(true);
      }
      setIsSearching(!tooShort);
      setIsSearching(!tooShort);

      // Try cache first
      if (!append) {
        const cached = cacheRef.current.get(key);
        const now = Date.now();
        if (cached && cached.expiresAt > now) {
          setResults(cached.items);
          setTotal(cached.total);
          setHasMore(cached.items.length < cached.total);
          setIsLoading(false);
          return;
        }
      }

      cancelInFlight();
      const controller = new AbortController();
      abortRef.current = controller;

  try {
        if (tooShort) {
          // Default: latest products
          const data = await listProducts({ limit, page: p, signal: controller.signal });
          if (requestSeq.current !== currentSeq) return; // stale
          
          const products = Array.isArray(data) ? data : [];
          
          if (append && enableInfiniteScroll) {
            // For infinite scroll, append new products
            setResults(prev => [...prev, ...products]);
            setHasMore(products.length === limit);
          } else {
            // Regular load or reset
            const payload = { items: products, total: products.length };
            // Cache
            cacheRef.current.set(key, { ...payload, expiresAt: Date.now() + 5 * 60_000 });
            if (cacheRef.current.size > 50) {
              const firstKey = cacheRef.current.keys().next().value as string | undefined;
              if (firstKey) cacheRef.current.delete(firstKey);
            }
            setResults(payload.items);
            setTotal(payload.total);
            setHasMore(products.length === limit);
          }
        } else {
          const data = await searchProducts({ query: trimmed, page: p, limit, sort, signal: controller.signal });
          if (requestSeq.current !== currentSeq) return; // stale
          const items = Array.isArray(data.items) ? data.items : [];
          const totalVal = typeof data.total === 'number' ? data.total : items.length;
          
          if (append && enableInfiniteScroll) {
            // For infinite scroll search results
            setResults(prev => [...prev, ...items]);
            setHasMore(items.length === limit);
          } else {
            // Cache
            cacheRef.current.set(key, { items, total: totalVal, expiresAt: Date.now() + 5 * 60_000 });
            if (cacheRef.current.size > 50) {
              const firstKey = cacheRef.current.keys().next().value as string | undefined;
              if (firstKey) cacheRef.current.delete(firstKey);
            }
            setResults(items);
            setTotal(totalVal);
            setHasMore(items.length === limit);
          }
        }
      } catch (e: unknown) {
        // Ignore aborts/cancels
        if (e && typeof e === 'object' && 'name' in e && (e as { name?: string }).name && ((e as { name: string }).name === "CanceledError" || (e as { name: string }).name === "AbortError")) {
          return;
        }
        // Allow retries for same key after an error
        lastKeyRef.current = "";
        const status = isAxiosError(e) ? e.response?.status : undefined;

        // If we were running a search (>=2 chars), degrade gracefully to empty
        if (!tooShort) {
          if (append) {
            setHasMore(false);
          } else {
            const payload = { items: [] as Product[], total: 0 };
            cacheRef.current.set(key, { ...payload, expiresAt: Date.now() + 60_000 });
            setResults(payload.items);
            setTotal(payload.total);
            setHasMore(false);
          }
          setError(null);
          return;
        }

        // For listing latest (tooShort), keep a friendly error
        if (status === 404) {
          if (append) {
            setHasMore(false);
          } else {
            const payload = { items: [] as Product[], total: 0 };
            cacheRef.current.set(key, { ...payload, expiresAt: Date.now() + 60_000 });
            setResults(payload.items);
            setTotal(payload.total);
            setHasMore(false);
          }
          setError(null);
        } else {
          setError("Failed to load products. Please try again.");
          setHasMore(false);
        }
      } finally {
        if (requestSeq.current === currentSeq) {
          if (append) {
            setIsLoadingMore(false);
          } else {
            setIsLoading(false);
          }
        }
      }
    },
    [cancelInFlight, limit, sort, enableInfiniteScroll]
  );

  // Load more products for infinite scroll
  const loadMore = useCallback(() => {
    if (isLoading || !hasMore) return;
    const nextPage = Math.floor(results.length / limit) + 1;
    runFetch({ q: query, p: nextPage, append: true });
  }, [isLoading, hasMore, results.length, limit, query, runFetch]);

  // Debounced effect
  useEffect(() => {
    // reset to page 1 when query changes
    setPage(1);
    setResults([]);
    setTotal(0);
    setHasMore(true);
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

  // Clear utilities
  const clear = useCallback(() => {
    setQuery("");
    setPage(1);
    setResults([]);
    setTotal(0);
    setHasMore(true);
  }, []);

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
    hasMore,
    isLoadingMore,
    // actions
    clear,
    submitNow,
    loadMore,
  } as const;
}
