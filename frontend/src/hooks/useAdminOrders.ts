import { useEffect, useState, useCallback, useRef } from 'react';
import type { AdminOrdersResponse } from '@/services/admin';
import { getAdminOrders } from '../services/admin';

export type AdminOrdersParams = {
  page?: number;
  per_page?: number;
  status?: string;
  start_date?: string; // YYYY-MM-DD
  end_date?: string; // YYYY-MM-DD
};

export function useAdminOrders(initial = { page: 1, per_page: 25 }) {
  const [data, setData] = useState<AdminOrdersResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [params, setParams] = useState<AdminOrdersParams>({
    page: initial.page,
    per_page: initial.per_page,
  });

  const abortRef = useRef<AbortController | null>(null);
  const debounceRef = useRef<number | null>(null);

  const fetch = useCallback(async (p: AdminOrdersParams) => {
    setLoading(true);
    setError(null);
    // Cancel previous
    if (abortRef.current) abortRef.current.abort();
    const controller = new AbortController();
    abortRef.current = controller;
    try {
      const res = await getAdminOrders(p);
      setData(res);
    } catch (err) {
  const e = err as Error | { name?: string; message?: string };
  if (e && 'name' in e && (e.name === 'CanceledError' || e.name === 'AbortError')) return;
  setError(e?.message || 'Failed to fetch orders');
      setData(null);
    } finally {
      setLoading(false);
    }
  }, []);

  // Trigger fetch when params change (debounced for filters)
  useEffect(() => {
    if (debounceRef.current) window.clearTimeout(debounceRef.current);
    // debounce 250ms
    debounceRef.current = window.setTimeout(() => {
      fetch(params as AdminOrdersParams);
    }, 250);
    return () => {
      if (debounceRef.current) window.clearTimeout(debounceRef.current);
    };
  }, [params, fetch]);

  useEffect(() => {
    return () => {
      if (abortRef.current) abortRef.current.abort();
    };
  }, []);

  const setPage = (page: number) => setParams((s) => ({ ...(s || {}), page }));
  const setPerPage = (per_page: number) => setParams((s) => ({ ...(s || {}), per_page, page: 1 }));
  const setStatus = (status?: string) => setParams((s) => ({ ...(s || {}), status, page: 1 }));
  const setDateRange = (start_date?: string, end_date?: string) => setParams((s) => ({ ...(s || {}), start_date, end_date, page: 1 }));

  return {
    data,
    loading,
    error,
    params,
    setPage,
    setPerPage,
    setStatus,
    setDateRange,
  };
}
