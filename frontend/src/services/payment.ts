import api from './api';

export interface PaymentIntentRequest {
  order_public_id?: string;
  shipping_address: string;
  shipping_selection?: {
    carrier: string;
    service_code: string;
    price: number;
    estimated_days: number;
    quote_id?: string | null;
  };
  customer_email?: string;
}

export interface PaymentIntentResponse {
  payment_intent_id: string;
  client_secret: string;
}

export async function createPaymentIntent(payload: PaymentIntentRequest): Promise<PaymentIntentResponse> {
  const res = await api.post<PaymentIntentResponse>('/payments/intent', payload);
  return res.data;
}

export default { createPaymentIntent };
