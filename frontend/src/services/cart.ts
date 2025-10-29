import api from "./api";
import type { CartResponse } from "../types/cart";

export async function getCart(): Promise<CartResponse> {
  const res = await api.get<CartResponse>("/cart");
  return res.data;
}

export async function addToCart(productId: number, quantity = 1): Promise<CartResponse> {
  const response = await api.post<CartResponse>("/cart", {
    product_id: productId,
    quantity,
  });
  return response.data;
}

export async function removeItem(itemId: number): Promise<CartResponse> {
  const response = await api.delete<CartResponse>(`/cart/items/${itemId}`);
  return response.data;
}

export async function updateItemQuantity(itemId: number, quantity: number): Promise<CartResponse> {
  const response = await api.patch<CartResponse>(`/cart/items/${itemId}`, {
    quantity,
  });
  return response.data;
}
