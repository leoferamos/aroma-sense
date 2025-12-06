import api from "./api";
import type { Product, CreateProductFormData, SearchResponse } from "../types/product";

export async function getProducts(): Promise<Product[]> {
  const response = await api.get<Product[]>("/products");
  return response.data;
}

export async function getProductById(id: number): Promise<Product> {
  const response = await api.get<Product>(`/products/${id}`);
  return response.data;
}

export async function getProductBySlug(slug: string): Promise<Product> {
  const response = await api.get<Product>(`/products/${slug}`);
  return response.data;
}

/**
 * Creates a new product with image upload
 * Sends multipart/form-data to POST /admin/products
 */
export async function createProduct(
  formData: CreateProductFormData
): Promise<Product> {
  const data = new FormData();

  // Add text fields
  data.append("name", formData.name);
  data.append("brand", formData.brand);
  data.append("description", formData.description);
  data.append("category", formData.category);
  
  // Add optional AI fields
  if (formData.intensity) data.append("intensity", formData.intensity);
  if (formData.gender) data.append("gender", formData.gender);
  if (formData.price_range) data.append("price_range", formData.price_range);
  
  // Add numeric fields as strings for multipart/form-data
  data.append("weight", formData.weight.toString());
  data.append("price", formData.price.toString());
  data.append("stock_quantity", Math.floor(formData.stock_quantity).toString());
  
  // Add array fields as separate form fields
  formData.accords.forEach((accord) => {
    data.append("accords", accord);
  });
  formData.occasions.forEach((occasion) => {
    data.append("occasions", occasion);
  });
  formData.seasons.forEach((season) => {
    data.append("seasons", season);
  });
  formData.notes_top.forEach((note) => {
    data.append("notes_top", note);
  });
  formData.notes_heart.forEach((note) => {
    data.append("notes_heart", note);
  });
  formData.notes_base.forEach((note) => {
    data.append("notes_base", note);
  });

  // Add image file if present
  if (formData.image) {
    data.append("image", formData.image);
  }

  const response = await api.post<Product>("/admin/products", data, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
  
  return response.data;
}

/**
 * Updates an existing product
 * Sends JSON payload to PATCH /admin/products/:id
 * Note: Image updates are not supported yet and will be handled in the future
 */
export async function updateProduct(
  id: number,
  formData: CreateProductFormData
): Promise<Product> {
  const payload = {
    name: formData.name,
    brand: formData.brand,
    description: formData.description,
    category: formData.category,
    weight: formData.weight,
    price: formData.price,
    stock_quantity: Math.floor(formData.stock_quantity),
    accords: formData.accords,
    occasions: formData.occasions,
    seasons: formData.seasons,
    intensity: formData.intensity,
    gender: formData.gender,
    price_range: formData.price_range,
    notes_top: formData.notes_top,
    notes_heart: formData.notes_heart,
    notes_base: formData.notes_base,
  };

  const response = await api.patch<Product>(`/admin/products/${id}`, payload);
  
  return response.data;
}

/**
 * Deletes a product by id
 */
export async function deleteProduct(id: number): Promise<void> {
  await api.delete(`/admin/products/${id}`);
}

// Searches products when `query` is provided. Returns a paginated envelope
export async function searchProducts(params: {
  query: string;
  page?: number;
  limit?: number;
  sort?: "relevance" | "latest";
  signal?: AbortSignal;
}): Promise<SearchResponse<Product>> {
  const { query, page = 1, limit = 10, sort = "relevance", signal } = params;
  const response = await api.get<SearchResponse<Product>>("/products", {
    params: { query, page, limit, sort },
    signal,
  });
  return response.data;
}

// Lists latest products when `query` is absent.
export async function listProducts(params?: {
  limit?: number;
  page?: number;
  signal?: AbortSignal;
}): Promise<Product[]> {
  const { limit = 10, page = 1, signal } = params || {};
  const response = await api.get("/products", {
    params: { limit, page },
    signal,
  });
  
  // Handle both array format (page=1) and paginated format (page>1)
  const data = response.data;
  if (Array.isArray(data)) {
    return data;
  } else if (data && typeof data === 'object' && Array.isArray(data.items)) {
    return data.items;
  } else {
    return [];
  }
}

export async function adminListProducts(params?: {
  limit?: number;
  page?: number;
}): Promise<Product[]> {
  const { limit = 50, page = 1 } = params || {};
  const response = await api.get("/admin/products", {
    params: { limit, page },
  });
  
  // Handle paginated format
  const data = response.data;
  if (data && typeof data === 'object' && Array.isArray(data.items)) {
    return data.items;
  } else {
    return [];
  }
}
