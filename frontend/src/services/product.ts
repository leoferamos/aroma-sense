import api from "./api";
import type { Product, CreateProductFormData } from "../types/product";

export async function getProducts(): Promise<Product[]> {
  const response = await api.get<Product[]>("/products");
  return response.data;
}

export async function getProductById(id: number): Promise<Product> {
  const response = await api.get<Product>(`/products/${id}`);
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
  
  // Add numeric fields as strings for multipart/form-data
  data.append("weight", formData.weight.toString());
  data.append("price", formData.price.toString());
  data.append("stock_quantity", Math.floor(formData.stock_quantity).toString());
  
  // Add notes as separate form fields
  formData.notes.forEach((note) => {
    data.append("notes", note);
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
    notes: formData.notes,
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
