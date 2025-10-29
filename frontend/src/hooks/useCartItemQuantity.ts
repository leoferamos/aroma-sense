import { useState, useEffect, useCallback, useRef } from 'react';
import { useCart } from '../contexts/CartContext';

interface UseCartItemQuantityOptions {
  itemId: number;
  initialQuantity: number;
  debounceMs?: number;
}

interface UseCartItemQuantityReturn {
  quantity: number;
  increment: () => void;
  decrement: () => void;
  setQuantity: (qty: number) => void;
  isSyncing: boolean;
  error: string | null;
}

/**
 * Hook to manage cart item quantity with:
 * - Optimistic UI (immediate interface update)
 * - Debouncing (waits X ms before sending request)
 * - Automatic rollback on error
 * 
 * Prevents multiple requests when user clicks rapidly
 */
export function useCartItemQuantity({
  itemId,
  initialQuantity,
  debounceMs = 600,
}: UseCartItemQuantityOptions): UseCartItemQuantityReturn {
  const { updateItemQuantity } = useCart();
  
  // Local state (optimistic)
  const [quantity, setQuantity] = useState(initialQuantity);
  const [isSyncing, setIsSyncing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Refs for debounce control
    const debounceTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const previousQuantityRef = useRef(initialQuantity);
  const lastSyncedQuantityRef = useRef(initialQuantity);

  // Update when initial quantity changes (e.g., cart refresh)
  useEffect(() => {
    setQuantity(initialQuantity);
    lastSyncedQuantityRef.current = initialQuantity;
    previousQuantityRef.current = initialQuantity;
  }, [initialQuantity]);

  // Cleanup timer on unmount
  useEffect(() => {
    return () => {
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current);
      }
    };
  }, []);

  // Function to sync with backend (with debounce)
  const syncQuantity = useCallback((newQuantity: number) => {
    // Clear previous timer
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current);
    }

    // Create new timer
    debounceTimerRef.current = setTimeout(async () => {
      // If quantity hasn't changed from last sync, ignore
      if (newQuantity === lastSyncedQuantityRef.current) {
        return;
      }

      setIsSyncing(true);
      setError(null);

      try {
        await updateItemQuantity(itemId, newQuantity);
        lastSyncedQuantityRef.current = newQuantity;
        previousQuantityRef.current = newQuantity;
        setError(null); 
      } catch (err) {
        // Rollback to previous quantity
        setQuantity(previousQuantityRef.current);
        setError(err instanceof Error ? err.message : 'Failed to update quantity');
      } finally {
        setIsSyncing(false);
      }
    }, debounceMs);
  }, [itemId, updateItemQuantity, debounceMs]);

  const increment = useCallback(() => {
    setError(null); // Clear error on new interaction
    setQuantity(prev => {
      const newQty = prev + 1;
      syncQuantity(newQty);
      return newQty;
    });
  }, [syncQuantity]);

  const decrement = useCallback(() => {
    setError(null); // Clear error on new interaction
    setQuantity(prev => {
      // Don't allow quantity less than 1
      if (prev <= 1) return prev;
      const newQty = prev - 1;
      syncQuantity(newQty);
      return newQty;
    });
  }, [syncQuantity]);

  const setQuantityManual = useCallback((qty: number) => {
    if (qty < 1) return; // Validate minimum quantity
    setError(null); // Clear error on new interaction
    setQuantity(qty);
    syncQuantity(qty);
  }, [syncQuantity]);

  return {
    quantity,
    increment,
    decrement,
    setQuantity: setQuantityManual,
    isSyncing,
    error,
  };
}
