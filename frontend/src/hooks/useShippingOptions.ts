import { useEffect, useMemo, useState } from 'react';
import { getShippingOptions } from '../services/shipping';
import type { ShippingOption } from '../types/shipping';

function normalizeCEP(cep: string): string {
  return (cep || '').replace(/\D/g, '');
}

export function useShippingOptions(postalCode: string) {
  const [options, setOptions] = useState<ShippingOption[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const cep = useMemo(() => normalizeCEP(postalCode), [postalCode]);

  useEffect(() => {
    let active = true;
    async function fetchOptions() {
      if (!cep || cep.length < 8) {
        setOptions([]);
        setError(null);
        return;
      }
      setLoading(true);
      setError(null);
      try {
        const data = await getShippingOptions(cep);
        if (!active) return;
        setOptions(data);
  } catch {
        if (!active) return;
        setError('Failed to load shipping options.');
        setOptions([]);
      } finally {
        if (active) setLoading(false);
      }
    }
    // small debounce: 300ms
    const t = setTimeout(fetchOptions, 300);
    return () => {
      active = false;
      clearTimeout(t);
    };
  }, [cep]);

  return { options, loading, error };
}

export default useShippingOptions;
