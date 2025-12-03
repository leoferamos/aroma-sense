export interface AdminOrderItem {
  id: number;
  user_id: string;
  total_amount: number;
  status: string;
  created_at: string;
}

export interface PaginationMeta {
  page: number;
  per_page: number;
  total_pages: number;
  total_count: number;
}

export interface StatsMeta {
  total_revenue: number;
  average_order_value: number;
}


// User-facing order types
export interface OrderItemResponse {
  id: number;
  product_id: number;
  product_name?: string;
  product_image_url?: string;
  quantity: number;
  price_at_purchase: number;
  subtotal: number;
}

export interface OrderResponse {
  id: number;
  user_id: string;
  total_amount: number;
  status: string;
  shipping_address: string;
  payment_method: string;
  items: OrderItemResponse[];
  item_count: number;
  created_at: string;
  updated_at: string;
}
