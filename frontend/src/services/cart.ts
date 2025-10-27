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
