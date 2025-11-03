import api from './api';
import type { OrderResponse } from '../types/order';

export async function getUserOrders(): Promise<OrderResponse[]> {
  const res = await api.get<OrderResponse[]>('/orders');
  return res.data;
}

export default { getUserOrders };
