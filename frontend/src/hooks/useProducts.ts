import { useState, useEffect } from "react";
import { isAxiosError } from "axios";
import { getProducts } from "../services/product";
import type { Product } from "../types/product";

export function useProducts() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchProducts() {
      try {
        setLoading(true);
        setError(null);
        const data = await getProducts();
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
    }

    fetchProducts();
  }, []);

  return { products, loading, error };
}
