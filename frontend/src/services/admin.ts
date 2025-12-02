import api from './api';

export interface AdminUser {
  contestation_deadline?: string | null;
  created_at?: string;
  deactivated_at?: string | null;
  deactivated_by?: string | null;
  deactivation_notes?: string | null;
  deactivation_reason?: string | null;
  display_name?: string | null;
  id: number;
  last_login_at?: string | null;
  masked_email?: string | null;
  public_id?: string;
  reactivation_requested?: boolean;
  role?: string;
  suspension_until?: string | null;
}

export interface AdminUsersResponse {
  limit: number;
  offset: number;
  total: number;
  users: AdminUser[];
}

export interface GetAdminUsersParams {
  limit?: number;
  offset?: number;
  role?: string;
  status?: string;
}

export async function getAdminUsers(params: GetAdminUsersParams = {}): Promise<AdminUsersResponse> {
  const { data } = await api.get('/admin/users', { params });
  return data as AdminUsersResponse;
}

// Orders
interface GetAdminOrdersParams {
  page?: number;
  per_page?: number;
  status?: string;
  start_date?: string; // YYYY-MM-DD
  end_date?: string; // YYYY-MM-DD
}

export async function getAdminOrders(params: GetAdminOrdersParams = {}): Promise<any> {
  const res = await api.get('/admin/orders', { params });
  return res.data;
}

export default {
  getAdminUsers,
  getAdminOrders,
};

