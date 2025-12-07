import { useEffect, useState } from 'react';
import { isAxiosError } from 'axios';
import { getProductBySlug } from '../services/product';
import type { Product } from '../types/product';

export const useProductDetail = (param: string) => {
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProduct = async () => {
      try {
        setLoading(true);
        setError(null);
        
        // For /product/:slug route, param is always a slug
        const data = await getProductBySlug(param);
        setProduct(data);
      } catch (err: unknown) {
        if (isAxiosError(err)) {
          if (err.response?.status === 404) {
            setError('Product not found');
          } else {
            setError(err.response?.data?.error || 'Failed to load product');
          }
        } else if (err instanceof Error) {
          setError(err.message);
        } else {
          setError('Failed to load product');
        }
      } finally {
        setLoading(false);
      }
    };

    if (param) {
      fetchProduct();
    }
  }, [param]);

  return { product, loading, error };
};
