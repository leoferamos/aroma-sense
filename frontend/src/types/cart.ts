import type { Product } from "./product";

export interface CartItem {
  id: number;
  cart_id: number;
  product_id: number;
  product?: Product;
  quantity: number;
  price: number;
  total: number;
  created_at: string;
  updated_at: string;
}

export interface CartResponse {
  id: number;
  user_id: string;
  items: CartItem[];
  total: number;
  item_count: number;
  created_at: string;
  updated_at: string;
}
