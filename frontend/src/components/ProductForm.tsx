import React, { useState } from "react";
import InputField from "./InputField";
import FormError from "./FormError";
import type { CreateProductFormData, ProductFormErrors } from "../types/product";
import type { ProductFormTouched } from "../hooks/useProductFormValidation";

interface ProductFormProps {
  form: CreateProductFormData;
  setForm: React.Dispatch<React.SetStateAction<CreateProductFormData>>;
  touched: ProductFormTouched;
  setTouched: React.Dispatch<React.SetStateAction<ProductFormTouched>>;
  errors: ProductFormErrors;
  onSubmit: (e: React.FormEvent) => void;
  loading: boolean;
  error?: string;
}

const ProductForm: React.FC<ProductFormProps> = ({
  form,
  setForm,
  touched,
  setTouched,
  errors,
  onSubmit,
  loading,
  error,
}) => {
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [noteInput, setNoteInput] = useState("");

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: value });
  };

  const handleNumberChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    const numValue = value === "" ? 0 : parseFloat(value);
    setForm({ ...form, [name]: isNaN(numValue) ? 0 : numValue });
  };

  const handleBlur = (
    e: React.FocusEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => {
    setTouched({ ...touched, [e.target.name]: true });
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setForm({ ...form, image: file });
      setTouched({ ...touched, image: true });

      // Create preview
      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleAddNote = () => {
    if (noteInput.trim()) {
      setForm({ ...form, notes: [...form.notes, noteInput.trim()] });
      setNoteInput("");
      setTouched({ ...touched, notes: true });
    }
  };

  const handleRemoveNote = (index: number) => {
    setForm({ ...form, notes: form.notes.filter((_, i) => i !== index) });
  };

  return (
    <form onSubmit={onSubmit} className="space-y-6" noValidate>
      {error && (
        <div 
          role="alert"
          className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg"
        >
          {error}
        </div>
      )}

      {/* Product Name */}
      <div>
        <InputField
          label="Product Name"
          type="text"
          name="name"
          value={form.name}
          onChange={handleChange}
          onBlur={handleBlur}
          required
        />
        <FormError message={errors.name} />
      </div>

      {/* Brand */}
      <div>
        <InputField
          label="Brand"
          type="text"
          name="brand"
          value={form.brand}
          onChange={handleChange}
          onBlur={handleBlur}
          required
        />
        <FormError message={errors.brand} />
      </div>

      {/* Weight */}
      <div>
        <label htmlFor="weight" className="text-base font-normal text-gray-800 block mb-2">
          Weight (grams)
        </label>
        <input
          type="number"
          id="weight"
          name="weight"
          value={form.weight || ""}
          onChange={handleNumberChange}
          onBlur={handleBlur}
          className="border border-gray-300 rounded-xl px-4 py-3 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900"
          min="0"
          step="0.01"
          required
        />
        <FormError message={errors.weight} />
      </div>

      {/* Price */}
      <div>
        <label htmlFor="price" className="text-base font-normal text-gray-800 block mb-2">
          Price (R$)
        </label>
        <input
          type="number"
          id="price"
          name="price"
          value={form.price || ""}
          onChange={handleNumberChange}
          onBlur={handleBlur}
          className="border border-gray-300 rounded-xl px-4 py-3 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900"
          min="0"
          step="0.01"
          required
        />
        <FormError message={errors.price} />
      </div>

      {/* Category */}
      <div>
        <InputField
          label="Category"
          type="text"
          name="category"
          value={form.category}
          onChange={handleChange}
          onBlur={handleBlur}
          required
        />
        <FormError message={errors.category} />
      </div>

      {/* Description */}
      <div>
        <label htmlFor="description" className="text-base font-normal text-gray-800 block mb-2">
          Description
        </label>
        <textarea
          id="description"
          name="description"
          value={form.description}
          onChange={handleChange}
          onBlur={handleBlur}
          className="border border-gray-300 rounded-xl px-4 py-3 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900"
          rows={4}
          required
        />
        <FormError message={errors.description} />
      </div>

      {/* Stock Quantity */}
      <div>
        <label htmlFor="stock_quantity" className="text-base font-normal text-gray-800 block mb-2">
          Stock Quantity
        </label>
        <input
          type="number"
          id="stock_quantity"
          name="stock_quantity"
          value={form.stock_quantity || ""}
          onChange={handleNumberChange}
          onBlur={handleBlur}
          className="border border-gray-300 rounded-xl px-4 py-3 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900"
          min="0"
          step="1"
          required
        />
        <FormError message={errors.stock_quantity} />
      </div>

      {/* Notes */}
      <div>
        <label className="text-base font-normal text-gray-800 block mb-2">
          Notes (Optional)
        </label>
        <div className="flex gap-2 mb-2">
          <input
            type="text"
            value={noteInput}
            onChange={(e) => setNoteInput(e.target.value)}
            onKeyPress={(e) => e.key === "Enter" && (e.preventDefault(), handleAddNote())}
            className="border border-gray-300 rounded-xl px-4 py-3 flex-1 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900"
          />
          <button
            type="button"
            onClick={handleAddNote}
            className="bg-blue-500 text-white px-6 py-3 rounded-xl hover:bg-blue-600 transition-colors"
          >
            Add
          </button>
        </div>
        {form.notes.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {form.notes.map((note, index) => (
              <span
                key={index}
                className="bg-gray-100 text-gray-700 px-3 py-1 rounded-full text-sm flex items-center gap-2"
              >
                {note}
                <button
                  type="button"
                  onClick={() => handleRemoveNote(index)}
                  className="text-gray-500 hover:text-red-500"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}
        <FormError message={errors.notes} />
      </div>

      {/* Image Upload */}
      <div>
        <label htmlFor="image" className="text-base font-normal text-gray-800 block mb-2">
          Product Image
        </label>
        <input
          type="file"
          id="image"
          name="image"
          accept="image/jpeg,image/png"
          onChange={handleImageChange}
          onBlur={handleBlur}
          className="border border-gray-300 rounded-xl px-4 py-3 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900"
          required
        />
        {imagePreview && (
          <div className="mt-4">
            <img
              src={imagePreview}
              alt="Preview"
              className="w-40 h-40 object-cover rounded-lg border border-gray-300"
            />
          </div>
        )}
        <FormError message={errors.image} />
      </div>

      {/* Submit Button */}
      <button
        type="submit"
        disabled={loading}
        className="w-full bg-blue-600 text-white py-3 rounded-xl font-medium hover:bg-blue-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed"
      >
        {loading ? "Creating Product..." : "Create Product"}
      </button>
    </form>
  );
};

export default ProductForm;
