import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import ProductForm from "../../components/ProductForm";
import LoadingSpinner from "../../components/LoadingSpinner";
import ErrorState from "../../components/ErrorState";
import { useProductFormValidation } from "../../hooks/useProductFormValidation";
import { useUpdateProduct } from "../../hooks/useUpdateProduct";
import { getProductById } from "../../services/product";
import type { CreateProductFormData, Product } from "../../types/product";
import type { ProductFormTouched } from "../../hooks/useProductFormValidation";
import AdminLayout from '../../components/admin/AdminLayout';
const EditProduct: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const productId = parseInt(id || "0", 10);

  const [loadingProduct, setLoadingProduct] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [product, setProduct] = useState<Product | null>(null);

  const [form, setForm] = useState<CreateProductFormData>({
    name: "",
    brand: "",
    weight: 0,
    description: "",
    price: 0,
    category: "",
    accords: [],
    occasions: [],
    seasons: [],
    intensity: "",
    gender: "",
    price_range: "",
    notes_top: [],
    notes_heart: [],
    notes_base: [],
    stock_quantity: 0,
    image: null,
  });

  const [touched, setTouched] = useState<ProductFormTouched>({
    name: false,
    brand: false,
    weight: false,
    description: false,
    price: false,
    category: false,
    accords: false,
    occasions: false,
    seasons: false,
    intensity: false,
    gender: false,
    price_range: false,
    notes_top: false,
    notes_heart: false,
    notes_base: false,
    stock_quantity: false,
    image: false,
  });

  const { errors, isValid } = useProductFormValidation({
    form,
    touched,
    isEditMode: true,
  });
  const { submitUpdate, loading, error, success } = useUpdateProduct();

  // Fetch product on mount
  useEffect(() => {
    async function fetchProduct() {
      try {
        setLoadingProduct(true);
        const data = await getProductById(productId);
        setProduct(data);
        // Prefill form with product data
        setForm({
          name: data.name,
          brand: data.brand,
          weight: data.weight,
          description: data.description,
          price: data.price,
          category: data.category,
          accords: data.accords || [],
          occasions: data.occasions || [],
          seasons: data.seasons || [],
          intensity: data.intensity || "",
          gender: data.gender || "",
          price_range: data.price_range || "",
          notes_top: data.notes_top || [],
          notes_heart: data.notes_heart || [],
          notes_base: data.notes_base || [],
          stock_quantity: data.stock_quantity,
          image: null,
        });
      } catch (err: unknown) {
        const e = err as { response?: { data?: { error?: string } }; message?: string };
        setLoadError(e?.response?.data?.error || e?.message || "Failed to load product.");
      } finally {
        setLoadingProduct(false);
      }
    }

    fetchProduct();
  }, [productId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Mark all fields as touched
    setTouched({
      name: true,
      brand: true,
      weight: true,
      description: true,
      price: true,
      category: true,
      accords: true,
      occasions: true,
      seasons: true,
      intensity: true,
      gender: true,
      price_range: true,
      notes_top: true,
      notes_heart: true,
      notes_base: true,
      stock_quantity: true,
      image: true,
    });

    // Check if form is valid
    if (!isValid) {
      return;
    }

    // Submit the form
    const updatedProduct = await submitUpdate(productId, form);

    if (updatedProduct) {
      // Redirect to products page after success
      setTimeout(() => {
        navigate("/admin/products");
      }, 2000);
    }
  };

  if (loadingProduct) {
    return (
      <AdminLayout title="Edit Product">
        <div className="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
          <LoadingSpinner message="Loading product..." />
        </div>
      </AdminLayout>
    );
  }

  if (loadError || !product) {
    return (
      <AdminLayout title="Edit Product">
        <div className="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
          <ErrorState
            message={loadError || "Product not found."}
            onRetry={() => navigate("/admin/products")}
          />
        </div>
      </AdminLayout>
    );
  }

  const actions = (
    <button
      onClick={() => navigate('/admin/products')}
      className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
    >
      ← Back to Products
    </button>
  );

  return (
    <AdminLayout actions={actions}>
      <div className="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 md:p-8">
          <div className="mb-6">
            <h2 className="text-2xl font-bold text-gray-900">Edit Product Details</h2>
            <p className="text-gray-500 text-sm mt-1">Update the information below</p>
          </div>

          {success && (
            <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg mb-6 flex items-center gap-2">
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              <span>Product updated successfully! Redirecting to products page...</span>
            </div>
          )}

          {error && <ErrorState message={error} />}

          <ProductForm
            form={form}
            setForm={setForm}
            touched={touched}
            setTouched={setTouched}
            errors={errors}
            onSubmit={handleSubmit}
            loading={loading}
            error={undefined}
            submitButtonText="Update Product"
            loadingText="Updating Product..."
            isEditMode={true}
            currentImageUrl={product?.image_url}
          />
        </div>
      </div>
    </AdminLayout>
  );
};

export default EditProduct;
