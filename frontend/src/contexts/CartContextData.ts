import { createContext } from 'react';
import type { CartResponse } from '../types/cart';

export interface CartContextValue {
  cart: CartResponse | null;
  itemCount: number;
  loading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
  addItem: (productId: number, quantity?: number) => Promise<void>;
  removeItem: (itemId: number) => Promise<void>;
  isRemovingItem: (itemId: number) => boolean;
  updateItemQuantity: (itemId: number, quantity: number) => Promise<void>;
}

export const CartContext = createContext<CartContextValue | undefined>(undefined);
