import api from "./api";

export interface ReviewRequest {
  rating: number;
  comment: string;
}

export interface Review {
  id: string;
  rating: number;
  comment: string;
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
  productId: number,
  params?: { page?: number; limit?: number; signal?: AbortSignal }
): Promise<ReviewListResponse> {
  const { page = 1, limit = 10, signal } = params || {};
  const { data } = await api.get<ReviewListResponse>(
    `/products/${productId}/reviews`,
    { params: { page, limit }, signal }
  );
  return data;
}

export async function getSummary(
  productId: number,
  opts?: { signal?: AbortSignal }
): Promise<ReviewSummary> {
  const { signal } = opts || {};
  const { data } = await api.get<ReviewSummary>(
    `/products/${productId}/reviews/summary`,
    { signal }
  );
  return data;
}

export async function createReview(
  productId: number,
  payload: ReviewRequest
): Promise<Review> {
  const { data } = await api.post<Review>(
    `/products/${productId}/reviews`,
    payload
  );
  return data;
}
