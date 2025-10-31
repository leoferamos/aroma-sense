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

export interface AdminOrdersResponse {
  orders: AdminOrderItem[];
  meta: {
    pagination: PaginationMeta;
    stats: StatsMeta;
  };
}
