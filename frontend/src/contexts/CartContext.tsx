import React, { useEffect, useMemo, useState, useCallback } from 'react';
import { CartContext } from './CartContextData';
import { isAxiosError } from 'axios';
import type { CartResponse } from '../types/cart';
import { addToCart as svcAddToCart, getCart as svcGetCart, removeItem as svcRemoveItem, updateItemQuantity as svcUpdateItemQuantity } from '../services/cart';
import { getAccessToken } from '../services/api';
import { useAuth } from '../hooks/useAuth';
import { useTranslation } from 'react-i18next';

export const CartProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isReady } = useAuth();
  const { t } = useTranslation('common');
  const [cart, setCart] = useState<CartResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [removingItemIds, setRemovingItemIds] = useState<Set<number>>(new Set());

  const refresh = useCallback(async () => {
    const token = getAccessToken();
    if (!token) {
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
        setError(err.response?.data?.error || t('errors.failedToLoadCart'));
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError(t('errors.failedToLoadCart'));
      }
    } finally {
      setLoading(false);
    }
  }, []);

  const addItem = useCallback(async (productId: number, quantity = 1) => {
    try {
      setLoading(true);
      setError(null);
      const data = await svcAddToCart(productId, quantity);
      setCart(data);
    } catch (err: unknown) {
      if (isAxiosError(err)) {
        setError(err.response?.data?.error || t('errors.failedToAddToCart'));
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError(t('errors.failedToAddToCart'));
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
        setError(err.response?.data?.error || t('errors.failedToRemoveItem'));
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError(t('errors.failedToRemoveItem'));
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

  const updateItemQuantity = useCallback(async (itemId: number, quantity: number) => {
    try {
      const data = await svcUpdateItemQuantity(itemId, quantity);
      setCart(data);
      setError(null); // Clear error on success
    } catch (err: unknown) {
      if (isAxiosError(err)) {
        const errorMsg = err.response?.data?.error || t('errors.failedToUpdateQuantity');
        setError(errorMsg);
        throw new Error(errorMsg);
      } else if (err instanceof Error) {
        setError(err.message);
        throw err;
      } else {
        const errorMsg = t('errors.failedToUpdateQuantity');
        setError(errorMsg);
        throw new Error(errorMsg);
      }
    }
  }, []);

  useEffect(() => {
    // Only fetch cart when auth is ready
    if (isReady) {
      refresh();
    }
  }, [isReady, refresh]);

  const value = useMemo(
    () => ({
      cart,
      itemCount: cart?.item_count ?? 0,
      loading,
      error,
      refresh,
      addItem,
      removeItem,
      isRemovingItem,
      updateItemQuantity
    }),
    [cart, loading, error, refresh, addItem, removeItem, isRemovingItem, updateItemQuantity]
  );

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
};
