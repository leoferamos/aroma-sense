import api from './api';

export interface ProfileResponse {
  public_id: string;
  email: string;
  role: string;
  display_name?: string | null;
  created_at: string;
  deletion_requested_at?: string | null;
  deletion_confirmed_at?: string | null;
  contestation_deadline?: string | null;
}

export interface UpdateProfileRequest {
  display_name: string;
}

export async function getMyProfile(): Promise<ProfileResponse> {
  const { data } = await api.get('/users/me');
  return data as ProfileResponse;
}

export async function updateMyProfile(payload: UpdateProfileRequest): Promise<ProfileResponse> {
  const { data } = await api.patch('/users/me/profile', payload);
  return data as ProfileResponse;
}

export interface DeleteAccountRequest {
  confirmation: string;
}

export async function requestAccountDeletion(): Promise<{ message: string }> {
  const payload: DeleteAccountRequest = { confirmation: 'DELETE_MY_ACCOUNT' };
  const { data } = await api.post('/users/me/deletion', payload);
  return data as { message: string };
}

export async function cancelAccountDeletion(): Promise<{ message: string }> {
  const { data } = await api.post('/users/me/deletion/cancel');
  return data as { message: string };
}

export interface ContestRequest {
  reason: string;
}

export async function requestContestation(payload: ContestRequest): Promise<{ message: string }> {
  const { data } = await api.post('/users/me/contest', payload);
  return data as { message: string };
}

export async function exportMyData(): Promise<Blob> {
  const res = await api.get('/users/me/export', { responseType: 'blob' });
  return res.data as Blob;
}
