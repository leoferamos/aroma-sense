import api from './api';
import type { ShippingOption } from '../types/shipping';

export async function getShippingOptions(postal_code: string): Promise<ShippingOption[]> {
  const res = await api.get<ShippingOption[]>('/shipping/options', {
    params: { postal_code },
  });
  return res.data;
}

export default { getShippingOptions };
