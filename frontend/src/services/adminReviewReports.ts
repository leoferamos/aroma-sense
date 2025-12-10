import api from './api';

export type ReviewReportStatus = 'pending' | 'accepted' | 'rejected';

export interface ReviewReportAdminReview {
  id: string;
  comment: string;
  rating: number;
  user_id: string;
}

export interface ReviewReportAdminReporter {
  public_id: string;
  display_name?: string | null;
}

export interface ReviewReportAdminItem {
  id: string;
  review_id: string;
  reported_by: string;
  reason_category: string;
  reason_text: string;
  status: ReviewReportStatus;
  created_at: string;
  review?: ReviewReportAdminReview;
  reporter?: ReviewReportAdminReporter;
}

export interface ReviewReportAdminResponse {
  items: ReviewReportAdminItem[];
  total: number;
  limit: number;
  offset: number;
}

export interface ListReviewReportsParams {
  status?: ReviewReportStatus;
  limit?: number;
  offset?: number;
}

export interface ResolveReviewReportRequest {
  action: 'accept' | 'reject';
  deactivate_user?: boolean;
  suspension_until?: string | null;
  notes?: string | null;
}

export async function listReviewReports(params: ListReviewReportsParams = {}): Promise<ReviewReportAdminResponse> {
  const { status = 'pending', limit = 20, offset = 0 } = params;
  const { data } = await api.get<ReviewReportAdminResponse>('/admin/review-reports', {
    params: { status, limit, offset },
  });
  return data;
}

export async function resolveReviewReport(id: string, payload: ResolveReviewReportRequest): Promise<void> {
  await api.post(`/admin/review-reports/${id}/resolve`, payload);
}
