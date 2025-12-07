import api from "./api";
import type { CartResponse } from "../types/cart";

export async function getCart(): Promise<CartResponse> {
  const res = await api.get<CartResponse>("/cart");
  return res.data;
}

export async function addToCart(productSlug: string, quantity = 1): Promise<CartResponse> {
  const response = await api.post<CartResponse>("/cart", {
    product_slug: productSlug,
    quantity,
  });
  return response.data;
}

export async function removeItem(productSlug: string): Promise<CartResponse> {
  const response = await api.delete<CartResponse>(`/cart/items/${productSlug}`);
  return response.data;
}

export async function updateItemQuantity(productSlug: string, quantity: number): Promise<CartResponse> {
  const response = await api.patch<CartResponse>(`/cart/items/${productSlug}`, {
    quantity,
  });
  return response.data;
}
