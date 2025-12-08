import { useCallback, useEffect, useMemo, useState } from 'react';
import { isAxiosError } from 'axios';
import { useTranslation } from 'react-i18next';
import { listReviews, getSummary, createReview, deleteReview, type Review, type ReviewSummary, type ReviewRequest } from '../services/review';

export interface UseReviewsResult {
  reviews: Review[];
  summary: ReviewSummary | null;
  page: number;
  limit: number;
  total: number;
  loading: boolean;
  error: string | null;
  setPage: (p: number) => void;
  setLimit: (l: number) => void;
  refresh: () => Promise<void>;
  createReview: (payload: ReviewRequest) => Promise<Review | null>;
  deleteReview: (reviewId: string) => Promise<boolean>;
}

export function useReviews(productSlug: string, opts?: { initialPage?: number; initialLimit?: number }) : UseReviewsResult {
  const { t } = useTranslation('common');
  const [reviews, setReviews] = useState<Review[]>([]);
  const [summaryState, setSummaryState] = useState<ReviewSummary | null>(null);
  const [page, setPage] = useState(opts?.initialPage ?? 1);
  const [limit, setLimit] = useState(opts?.initialLimit ?? 10);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async (signal?: AbortSignal) => {
    setLoading(true);
    setError(null);
    try {
      const [listResp, summaryResp] = await Promise.all([
        listReviews(productSlug, { page, limit, signal }),
        getSummary(productSlug, { signal }),
      ]);
      setReviews(listResp.items);
      setTotal(listResp.total);
      setSummaryState(summaryResp);
    } catch (err: unknown) {
      if (isAxiosError(err) && (err.code === 'ERR_CANCELED' || err.message?.toLowerCase().includes('canceled'))) {
        return;
      }
      if (isAxiosError(err)) {
        setError(err.response?.data?.error || 'Failed to load reviews');
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Failed to load reviews');
      }
    } finally {
      setLoading(false);
    }
  }, [productSlug, page, limit]);

  useEffect(() => {
    const controller = new AbortController();
    load(controller.signal);
    return () => controller.abort();
  }, [load]);

  const refresh = useCallback(async () => {
    await load();
  }, [load]);

  const doCreate = useCallback(async (payload: ReviewRequest) => {
    try {
      const created = await createReview(productSlug, payload);
      // Optimistic refresh: prepend or reload depending on page
      if (page === 1) {
        setReviews((prev) => [created, ...prev]);
        setTotal((t) => t + 1);
        // Update summary locally if available
        setSummaryState((s) => {
          if (!s) return s;
          const newCount = s.count + 1;
          const newAvg = (s.average * s.count + created.rating) / newCount;
          const newDist = { ...s.distribution };
          newDist[created.rating] = (newDist[created.rating] || 0) + 1;
          return { average: newAvg, count: newCount, distribution: newDist };
        });
      } else {
        await refresh();
      }
      return created;
    } catch (err: unknown) {
      if (isAxiosError(err)) {
        const msg = err.response?.data?.error || 'Failed to create review';
        setError(msg);
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Failed to create review');
      }
      return null;
    }
  }, [productSlug, page, refresh]);

  const doDelete = useCallback(async (reviewId: string): Promise<boolean> => {
    // Validate input
    if (!reviewId || typeof reviewId !== 'string') {
      return false;
    }

    try {
      // Make API call first
      await deleteReview(reviewId);

      // Refresh data from server to ensure consistency
      await refresh();

      return true;

    } catch (err: unknown) {
      // Handle different error types with i18n
      if (isAxiosError(err)) {
        const status = err.response?.status;

        if (status === 403) {
          setError(t('reviews.deleteForbidden', 'You do not have permission to delete this review'));
        } else if (status === 404) {
          setError(t('errors.notFound', 'Review not found or already deleted'));
        } else if (status === 401) {
          setError(t('errors.unauthorized', 'Please log in to delete reviews'));
        } else {
          setError(t('reviews.deleteFailed', 'Failed to delete review. Please try again.'));
        }
      } else {
        setError(t('errors.unexpected', 'An unexpected error occurred'));
      }

      return false;
    }
  }, [refresh, t]);

  const summary = useMemo(() => summaryState, [summaryState]);

  return { reviews, summary, page, limit, total, loading, error, setPage, setLimit, refresh, createReview: doCreate, deleteReview: doDelete };
}
