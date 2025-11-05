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

export interface SearchResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
}

// Form data for creating a new product
export interface CreateProductFormData {
  name: string;
  brand: string;
  weight: number;
  description: string;
  price: number;
  category: string;
  notes: string[];
  stock_quantity: number;
  image: File | null;
}

// Validation errors for the product form
export interface ProductFormErrors {
  name?: string;
  brand?: string;
  weight?: string;
  description?: string;
  price?: string;
  category?: string;
  notes?: string;
  stock_quantity?: string;
  image?: string;
}
