import { useState, useEffect } from "react";
import { isAxiosError } from "axios";
import { listProducts } from "../services/product";
import type { Product } from "../types/product";

export function useProducts() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchProducts = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await listProducts({ limit: 1000 });
      setProducts(data);
    } catch (err: unknown) {
      if (isAxiosError(err)) {
        setError(err.response?.data?.error || "Failed to load products");
      } else if (err instanceof Error) {
        setError(err.message || "Failed to load products");
      } else {
        setError("Failed to load products");
      }
    } finally {
      setLoading(false);
    }
  };

  // Fetch on mount
  useEffect(() => {
    fetchProducts();
  }, []);

  return { products, loading, error, refetch: fetchProducts };
}
