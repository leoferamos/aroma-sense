import { createContext } from 'react';
import type { CartResponse } from '../types/cart';

export interface CartContextValue {
  cart: CartResponse | null;
  itemCount: number;
  loading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
  addItem: (productSlug: string, quantity?: number) => Promise<void>;
  removeItem: (productSlug: string) => Promise<void>;
  isRemovingItem: (productSlug: string) => boolean;
  updateItemQuantity: (productSlug: string, quantity: number) => Promise<void>;
}

export const CartContext = createContext<CartContextValue | undefined>(undefined);
