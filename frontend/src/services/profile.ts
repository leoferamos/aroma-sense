import api from './api';

export interface ProfileResponse {
  public_id: string;
  email: string;
  role: string;
  display_name?: string | null;
  created_at: string;
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
