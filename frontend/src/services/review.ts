import api from "./api";

export interface ReviewRequest {
  rating: number;
  comment: string;
}

export interface Review {
  id: string;
  rating: number;
  comment: string;
  author_id: string;
  author_display: string;
  created_at: string;
}

export interface ReviewListResponse {
  items: Review[];
  total: number;
  page: number;
  limit: number;
}

export interface ReviewSummary {
  average: number;
  count: number;
  distribution: Record<number, number>;
}

export async function listReviews(
  productSlug: string,
  params?: { page?: number; limit?: number; signal?: AbortSignal }
): Promise<ReviewListResponse> {
  const { page = 1, limit = 10, signal } = params || {};
  const { data } = await api.get<ReviewListResponse>(
    `/products/${productSlug}/reviews`,
    { params: { page, limit }, signal }
  );
  return data;
}

export async function getSummary(
  productSlug: string,
  opts?: { signal?: AbortSignal }
): Promise<ReviewSummary> {
  const { signal } = opts || {};
  const { data } = await api.get<ReviewSummary>(
    `/products/${productSlug}/reviews/summary`,
    { signal }
  );
  return data;
}

export async function createReview(
  productSlug: string,
  payload: ReviewRequest
): Promise<Review> {
  const { data } = await api.post<Review>(
    `/products/${productSlug}/reviews`,
    payload
  );
  return data;
}

export async function deleteReview(
  productSlug: string,
  reviewId: string
): Promise<void> {
  await api.delete(`/reviews/${reviewId}`);
}
