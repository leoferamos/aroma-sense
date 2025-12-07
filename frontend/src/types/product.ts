export interface Product {
  id?: number; // Optional for admin routes
  name: string;
  brand: string;
  weight: number;
  description: string;
  price: number;
  image_url: string;
  thumbnail_url?: string;
  slug: string;
  accords?: string[];
  occasions?: string[];
  seasons?: string[];
  intensity?: string;
  gender?: string;
  price_range?: string;
  notes_top?: string[];
  notes_heart?: string[];
  notes_base?: string[];
  category: string;
  stock_quantity: number;
  created_at: string;
  updated_at: string;
  can_review?: boolean;
  cannot_review_reason?: string;
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
  accords: string[];
  occasions: string[];
  seasons: string[];
  intensity: string;
  gender: string;
  price_range: string;
  notes_top: string[];
  notes_heart: string[];
  notes_base: string[];
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
  accords?: string;
  occasions?: string;
  seasons?: string;
  intensity?: string;
  gender?: string;
  price_range?: string;
  notes_top?: string;
  notes_heart?: string;
  notes_base?: string;
  stock_quantity?: string;
  image?: string;
}
