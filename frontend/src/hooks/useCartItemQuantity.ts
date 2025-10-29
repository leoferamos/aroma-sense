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
 * Hook para gerenciar quantidade de item no carrinho com:
 * - Optimistic UI (atualização imediata na interface)
 * - Debouncing (aguarda X ms antes de enviar requisição)
 * - Rollback automático em caso de erro
 * 
 * Evita múltiplas requisições quando o usuário clica rapidamente
 */
export function useCartItemQuantity({
  itemId,
  initialQuantity,
  debounceMs = 600,
}: UseCartItemQuantityOptions): UseCartItemQuantityReturn {
  const { updateItemQuantity } = useCart();
  
  // Estado local (optimistic)
  const [quantity, setQuantity] = useState(initialQuantity);
  const [isSyncing, setIsSyncing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Refs para controle de debounce
    const debounceTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const previousQuantityRef = useRef(initialQuantity);
  const lastSyncedQuantityRef = useRef(initialQuantity);

  // Atualiza quando a quantidade inicial mudar (ex: refresh do cart)
  useEffect(() => {
    setQuantity(initialQuantity);
    lastSyncedQuantityRef.current = initialQuantity;
    previousQuantityRef.current = initialQuantity;
  }, [initialQuantity]);

  // Cleanup do timer ao desmontar
  useEffect(() => {
    return () => {
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current);
      }
    };
  }, []);

  // Função para sincronizar com backend (com debounce)
  const syncQuantity = useCallback((newQuantity: number) => {
    // Limpa timer anterior
    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current);
    }

    // Cria novo timer
    debounceTimerRef.current = setTimeout(async () => {
      // Se quantidade não mudou em relação ao último sync, ignora
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
        // Rollback para quantidade anterior
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
      // Não permite quantidade menor que 1
      if (prev <= 1) return prev;
      const newQty = prev - 1;
      syncQuantity(newQty);
      return newQty;
    });
  }, [syncQuantity]);

  const setQuantityManual = useCallback((qty: number) => {
    if (qty < 1) return; // Valida quantidade mínima
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
