import { useContext } from 'react';
import { CartContext } from '../contexts/CartContextData';
import type { CartContextValue } from '../contexts/CartContextData';

export const useCart = (): CartContextValue => {
  const ctx = useContext(CartContext);
  if (!ctx) throw new Error('useCart must be used within a CartProvider');
  return ctx;
};