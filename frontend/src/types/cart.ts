import type { Product } from "./product";

export interface CartItem {
  product?: Product;
  quantity: number;
  price: number;
  total: number;
}

export interface CartResponse {
  items: CartItem[];
  total: number;
  item_count: number;
}
