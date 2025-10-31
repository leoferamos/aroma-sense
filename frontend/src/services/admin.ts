import api from './api';
import type { AdminOrdersResponse } from '../types/order';

interface GetAdminOrdersParams {
  page?: number;
  per_page?: number;
  status?: string;
  start_date?: string; // YYYY-MM-DD
  end_date?: string; // YYYY-MM-DD
}

export async function getAdminOrders(params: GetAdminOrdersParams = {}): Promise<AdminOrdersResponse> {
  const res = await api.get<AdminOrdersResponse>('/admin/orders', { params });
  return res.data;
}
