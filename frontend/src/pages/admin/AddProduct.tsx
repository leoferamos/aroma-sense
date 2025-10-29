import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import ProductForm from "../../components/ProductForm";
import { useProductFormValidation } from "../../hooks/useProductFormValidation";
import { useCreateProduct } from "../../hooks/useCreateProduct";
import type { CreateProductFormData } from "../../types/product";
import type { ProductFormTouched } from "../../hooks/useProductFormValidation";

const AddProduct: React.FC = () => {
  const navigate = useNavigate();
  const [form, setForm] = useState<CreateProductFormData>({
    name: "",
    brand: "",
    weight: 0,
    description: "",
    price: 0,
    category: "",
    notes: [],
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
    notes: false,
    stock_quantity: false,
    image: false,
  });

  const { errors, isValid } = useProductFormValidation(form, touched);
  const { submitProduct, loading, error, success } = useCreateProduct();

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
      notes: true,
      stock_quantity: true,
      image: true,
    });

    // Check if form is valid
    if (!isValid) {
      return;
    }

    // Submit the form
    const product = await submitProduct(form);

    if (product) {
      // Redirect to products page or show success message
      setTimeout(() => {
        navigate("/products");
      }, 2000);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center gap-3">
              <img src="/logo.png" alt="Aroma Sense" className="h-8" />
              <h1 className="text-lg font-semibold text-gray-900">
                Add New Product
              </h1>
            </div>
            <button
              onClick={() => navigate("/admin/dashboard")}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              ← Back to Dashboard
            </button>
          </div>
        </div>
      </nav>

      <main className="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 md:p-8">
          <div className="mb-6">
            <h2 className="text-2xl font-bold text-gray-900">
              Product Details
            </h2>
            <p className="text-gray-500 text-sm mt-1">Fill in the information below</p>
          </div>

          {success && (
            <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg mb-6 flex items-center gap-2">
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              <span>Product created successfully! Redirecting to products page...</span>
            </div>
          )}

          <ProductForm
            form={form}
            setForm={setForm}
            touched={touched}
            setTouched={setTouched}
            errors={errors}
            onSubmit={handleSubmit}
            loading={loading}
            error={error}
          />
        </div>
      </main>
    </div>
  );
};

export default AddProduct;
