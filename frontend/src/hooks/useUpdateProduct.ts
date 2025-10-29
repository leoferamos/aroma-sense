import { updateProduct } from "../services/product";
import { messages } from "../constants/messages";
import { useProductMutation } from "./useProductMutation";
import type { CreateProductFormData } from "../types/product";

interface UpdateProductData {
  id: number;
  data: CreateProductFormData;
}

export function useUpdateProduct() {
  const { mutate, loading, error, success, product, reset } = useProductMutation<UpdateProductData>(
    ({ id, data }) => updateProduct(id, data),
    {
      onError: (errorMsg) => {
        if (errorMsg.toLowerCase().includes("unauthorized")) {
          return messages.productCreateUnauthorized;
        } else if (errorMsg.toLowerCase().includes("not found")) {
          return "Product not found.";
        }
        return errorMsg || "Failed to update product.";
      },
    }
  );

  const submitUpdate = async (id: number, data: CreateProductFormData) => {
    return mutate({ id, data });
  };

  return {
    submitUpdate,
    loading,
    error,
    success,
    updatedProduct: product,
    reset,
  };
}
