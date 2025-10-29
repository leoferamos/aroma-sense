import { useState } from "react";
import { isAxiosError } from "axios";
import type { Product } from "../types/product";

interface UseProductMutationOptions {
  onError?: (errorMsg: string) => string;
}

/**
 * Generic hook for product mutations
 */
export function useProductMutation<TData>(
  mutationFn: (data: TData) => Promise<Product>,
  options?: UseProductMutationOptions
) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [product, setProduct] = useState<Product | null>(null);

  async function mutate(data: TData) {
    setLoading(true);
    setError("");
    setSuccess(false);
    setProduct(null);

    try {
      const result = await mutationFn(data);
      setProduct(result);
      setSuccess(true);
      return result;
    } catch (err) {
      let errorMessage = "An unexpected error occurred.";

      if (isAxiosError<{ error: string }>(err)) {
        const errorMsg = err.response?.data?.error || "";
        errorMessage = options?.onError?.(errorMsg) || errorMsg || errorMessage;
      } else if (err instanceof Error) {
        errorMessage = err.message;
      }

      setError(errorMessage);
      return null;
    } finally {
      setLoading(false);
    }
  }

  function reset() {
    setLoading(false);
    setError("");
    setSuccess(false);
    setProduct(null);
  }

  return {
    mutate,
    loading,
    error,
    success,
    product,
    reset,
  };
}
