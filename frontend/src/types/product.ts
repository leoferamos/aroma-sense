export interface Product {
  id: number;
  name: string;
  brand: string;
  weight: number;
  description: string;
  price: number;
  image_url: string;
  category: string;
  notes: string;
  stock_quantity: number;
  created_at: string;
  updated_at: string;
}

export interface ProductsResponse {
  products: Product[];
}
