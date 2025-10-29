import React, { createContext, useContext, useEffect, useMemo, useState, useCallback } from 'react';
import { isAxiosError } from 'axios';
import type { CartResponse } from '../types/cart';
import { addToCart as svcAddToCart, getCart as svcGetCart, removeItem as svcRemoveItem } from '../services/cart';
import { useAuth } from './AuthContext';

interface CartContextValue {
  cart: CartResponse | null;
  itemCount: number;
  loading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
  addItem: (productId: number, quantity?: number) => Promise<void>;
  removeItem: (itemId: number) => Promise<void>;
  isRemovingItem: (itemId: number) => boolean;
}

const CartContext = createContext<CartContextValue | undefined>(undefined);

export const CartProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth();
  const [cart, setCart] = useState<CartResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [removingItemIds, setRemovingItemIds] = useState<Set<number>>(new Set());

  const refresh = useCallback(async () => {
    if (!isAuthenticated) {
      setCart(null);
      return;
    }
    try {
      setLoading(true);
      setError(null);
      const data = await svcGetCart();
      setCart(data);
    } catch (err: unknown) {
      if (isAxiosError(err)) {
        setError(err.response?.data?.error || 'Failed to load cart');
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Failed to load cart');
      }
    } finally {
      setLoading(false);
    }
  }, [isAuthenticated]);

  const addItem = useCallback(async (productId: number, quantity = 1) => {
    try {
      setLoading(true);
      setError(null);
      const data = await svcAddToCart(productId, quantity);
      setCart(data);
    } catch (err: unknown) {
      if (isAxiosError(err)) {
        setError(err.response?.data?.error || 'Failed to add to cart');
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Failed to add to cart');
      }
    } finally {
      setLoading(false);
    }
  }, []);

  const removeItem = useCallback(async (itemId: number) => {
    try {
      setRemovingItemIds(prev => new Set(prev).add(itemId));
      setError(null);
      const data = await svcRemoveItem(itemId);
      setCart(data);
    } catch (err: unknown) {
      if (isAxiosError(err)) {
        setError(err.response?.data?.error || 'Failed to remove item');
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Failed to remove item');
      }
    } finally {
      setRemovingItemIds(prev => {
        const next = new Set(prev);
        next.delete(itemId);
        return next;
      });
    }
  }, []);

  const isRemovingItem = useCallback((itemId: number) => removingItemIds.has(itemId), [removingItemIds]);

  useEffect(() => {
    // Fetch cart on mount or when auth changes
    refresh();
  }, [refresh]);

  const itemCount = useMemo(() => cart?.item_count ?? 0, [cart]);

  const value = useMemo(
    () => ({ cart, itemCount, loading, error, refresh, addItem, removeItem, isRemovingItem }),
    [cart, itemCount, loading, error, refresh, addItem, removeItem, isRemovingItem]
  );

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
};

export const useCart = () => {
  const ctx = useContext(CartContext);
  if (!ctx) throw new Error('useCart must be used within a CartProvider');
  return ctx;
};
