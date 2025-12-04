import api from './api';

export interface AdminContestation {
  id: number;
  user_id: number;
  reason: string;
  status: string;
  requested_at: string;
  reviewed_at?: string;
  reviewed_by?: number;
  review_notes?: string;
}

export interface ContestationListResponse {
  data: AdminContestation[];
  total: number;
}

export async function getPendingContestations(limit = 20, offset = 0): Promise<ContestationListResponse> {
  const { data } = await api.get('/admin/contestations', { params: { limit, offset } });
  return data as ContestationListResponse;
}

export async function approveContestation(id: number, reviewNotes?: string): Promise<{ message: string }> {
  const { data } = await api.post(`/admin/contestations/${id}/approve`, reviewNotes ? { review_notes: reviewNotes } : {});
  return data as { message: string };
}

export async function rejectContestation(id: number, reviewNotes?: string): Promise<{ message: string }> {
  const { data } = await api.post(`/admin/contestations/${id}/reject`, reviewNotes ? { review_notes: reviewNotes } : {});
  return data as { message: string };
}
