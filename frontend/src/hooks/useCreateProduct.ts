import { useState } from "react";
import { isAxiosError } from "axios";
import { createProduct } from "../services/product";
import { messages } from "../constants/messages";
import type { CreateProductFormData, Product } from "../types/product";

export function useCreateProduct() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [createdProduct, setCreatedProduct] = useState<Product | null>(null);

  async function submitProduct(data: CreateProductFormData) {
    setLoading(true);
    setError("");
    setSuccess(false);
    setCreatedProduct(null);

    try {
      const product = await createProduct(data);
      setCreatedProduct(product);
      setSuccess(true);
      return product;
    } catch (err) {
      if (isAxiosError<{ error: string }>(err)) {
        const errorMsg = err.response?.data?.error || "";
        
        // Handle specific error cases from backend
        if (errorMsg.toLowerCase().includes("unauthorized")) {
          setError(messages.productCreateUnauthorized);
        } else if (errorMsg.toLowerCase().includes("invalid image")) {
          setError(messages.productCreateInvalidImage);
        } else if (errorMsg.toLowerCase().includes("image too large")) {
          setError(messages.productCreateImageTooLarge);
        } else if (errorMsg.toLowerCase().includes("missing field")) {
          setError(messages.productCreateMissingFields);
        } else {
          setError(errorMsg || messages.productCreateError);
        }
      } else if (err instanceof Error) {
        setError(err.message || messages.productCreateError);
      } else {
        setError(messages.productCreateUnexpectedError);
      }
      return null;
    } finally {
      setLoading(false);
    }
  }

  function reset() {
    setLoading(false);
    setError("");
    setSuccess(false);
    setCreatedProduct(null);
  }

  return {
    submitProduct,
    loading,
    error,
    success,
    createdProduct,
    reset,
  };
}
