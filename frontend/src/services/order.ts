import api from './api';
import type { OrderResponse } from '../types/order';

export async function getUserOrders(): Promise<OrderResponse[]> {
  const res = await api.get<OrderResponse[]>('/orders');
  return res.data;
}

export interface OrderCreateRequest {
  payment_method: string;
  shipping_address: string;
  shipping_selection: {
    carrier: string;
    service_code: string;
    price: number;
    estimated_days: number;
    quote_id?: string | null;
  };
}

export async function createOrder(payload: OrderCreateRequest): Promise<OrderResponse> {
  const res = await api.post<OrderResponse>('/orders', payload);
  return res.data;
}

export default { getUserOrders, createOrder };
