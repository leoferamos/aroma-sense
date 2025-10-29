import { createProduct } from "../services/product";
import { messages } from "../constants/messages";
import { useProductMutation } from "./useProductMutation";
import type { CreateProductFormData } from "../types/product";

export function useCreateProduct() {
  const { mutate, loading, error, success, product, reset } = useProductMutation<CreateProductFormData>(
    createProduct,
    {
      onError: (errorMsg) => {
        if (errorMsg.toLowerCase().includes("unauthorized")) {
          return messages.productCreateUnauthorized;
        } else if (errorMsg.toLowerCase().includes("invalid image")) {
          return messages.productCreateInvalidImage;
        } else if (errorMsg.toLowerCase().includes("image too large")) {
          return messages.productCreateImageTooLarge;
        } else if (errorMsg.toLowerCase().includes("missing field")) {
          return messages.productCreateMissingFields;
        }
        return errorMsg || messages.productCreateError;
      },
    }
  );

  return {
    submitProduct: mutate,
    loading,
    error,
    success,
    createdProduct: product,
    reset,
  };
}
