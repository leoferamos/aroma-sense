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
  submitButtonText?: string;
  loadingText?: string;
  isEditMode?: boolean;
  currentImageUrl?: string;
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
  submitButtonText = "Create Product",
  loadingText = "Creating Product...",
  isEditMode = false,
  currentImageUrl,
}) => {
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [noteTopInput, setNoteTopInput] = useState("");
  const [noteHeartInput, setNoteHeartInput] = useState("");
  const [noteBaseInput, setNoteBaseInput] = useState("");
  const [accordInput, setAccordInput] = useState("");
  const [occasionInput, setOccasionInput] = useState("");
  const [seasonInput, setSeasonInput] = useState("");

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: value });
  };

  const handleNumberChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    const normalizedValue = value.replace(',', '.');
    const numValue = normalizedValue === "" ? 0 : parseFloat(normalizedValue);
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

  const handleAddNoteTop = () => {
    if (noteTopInput.trim()) {
      const items = noteTopInput.split(/[;,]/).map(item => item.trim()).filter(item => item.length > 0);
      setForm({ ...form, notes_top: [...form.notes_top, ...items] });
      setNoteTopInput("");
      setTouched({ ...touched, notes_top: true });
    }
  };

  const handleRemoveNoteTop = (index: number) => {
    setForm({ ...form, notes_top: form.notes_top.filter((_, i) => i !== index) });
  };

  const handleAddNoteHeart = () => {
    if (noteHeartInput.trim()) {
      const items = noteHeartInput.split(/[;,]/).map(item => item.trim()).filter(item => item.length > 0);
      setForm({ ...form, notes_heart: [...form.notes_heart, ...items] });
      setNoteHeartInput("");
      setTouched({ ...touched, notes_heart: true });
    }
  };

  const handleRemoveNoteHeart = (index: number) => {
    setForm({ ...form, notes_heart: form.notes_heart.filter((_, i) => i !== index) });
  };

  const handleAddNoteBase = () => {
    if (noteBaseInput.trim()) {
      const items = noteBaseInput.split(/[;,]/).map(item => item.trim()).filter(item => item.length > 0);
      setForm({ ...form, notes_base: [...form.notes_base, ...items] });
      setNoteBaseInput("");
      setTouched({ ...touched, notes_base: true });
    }
  };

  const handleRemoveNoteBase = (index: number) => {
    setForm({ ...form, notes_base: form.notes_base.filter((_, i) => i !== index) });
  };

  const handleAddAccord = () => {
    if (accordInput.trim()) {
      const items = accordInput.split(/[;,]/).map(item => item.trim()).filter(item => item.length > 0);
      setForm({ ...form, accords: [...form.accords, ...items] });
      setAccordInput("");
      setTouched({ ...touched, accords: true });
    }
  };

  const handleRemoveAccord = (index: number) => {
    setForm({ ...form, accords: form.accords.filter((_, i) => i !== index) });
  };

  const handleAddOccasion = () => {
    if (occasionInput.trim()) {
      const items = occasionInput.split(/[;,]/).map(item => item.trim()).filter(item => item.length > 0);
      setForm({ ...form, occasions: [...form.occasions, ...items] });
      setOccasionInput("");
      setTouched({ ...touched, occasions: true });
    }
  };

  const handleRemoveOccasion = (index: number) => {
    setForm({ ...form, occasions: form.occasions.filter((_, i) => i !== index) });
  };

  const handleAddSeason = () => {
    if (seasonInput.trim()) {
      const items = seasonInput.split(/[;,]/).map(item => item.trim()).filter(item => item.length > 0);
      setForm({ ...form, seasons: [...form.seasons, ...items] });
      setSeasonInput("");
      setTouched({ ...touched, seasons: true });
    }
  };

  const handleRemoveSeason = (index: number) => {
    setForm({ ...form, seasons: form.seasons.filter((_, i) => i !== index) });
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

      {/* Product Name & Brand */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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
      </div>

      {/* Weight & Price */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <InputField
            label="Weight (grams)"
            type="number"
            name="weight"
            value={String(form.weight || "")}
            onChange={handleNumberChange}
            onBlur={handleBlur}
            min="0"
            step="0.01"
            required
          />
          <FormError message={errors.weight} />
        </div>
        <div>
          <InputField
            label="Price (R$)"
            type="number"
            name="price"
            value={String(form.price || "")}
            onChange={handleNumberChange}
            onBlur={handleBlur}
            min="0"
            step="0.01"
            required
          />
          <FormError message={errors.price} />
        </div>
      </div>

      {/* Category & Stock Quantity */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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
        <div>
          <InputField
            label="Stock Quantity"
            type="number"
            name="stock_quantity"
            value={String(form.stock_quantity || "")}
            onChange={handleNumberChange}
            onBlur={handleBlur}
            min="0"
            step="1"
            required
          />
          <FormError message={errors.stock_quantity} />
        </div>
      </div>

      {/* Description */}
      <div>
        <label htmlFor="description" className="text-sm font-medium text-gray-700 block mb-2">
          Description
        </label>
        <textarea
          id="description"
          name="description"
          value={form.description}
          onChange={handleChange}
          onBlur={handleBlur}
          className="border border-gray-300 rounded-xl px-4 py-2.5 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 text-sm"
          rows={4}
          required
        />
        <FormError message={errors.description} />
      </div>

      {/* Intensity, Gender & Price Range */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <InputField
            label="Intensity (Optional)"
            type="text"
            name="intensity"
            value={form.intensity}
            onChange={handleChange}
            onBlur={handleBlur}
            placeholder="e.g., Light, Moderate"
          />
          <FormError message={errors.intensity} />
        </div>
        <div>
          <InputField
            label="Gender (Optional)"
            type="text"
            name="gender"
            value={form.gender}
            onChange={handleChange}
            onBlur={handleBlur}
            placeholder="e.g., Unisex, Men"
          />
          <FormError message={errors.gender} />
        </div>
        <div>
          <InputField
            label="Price Range (Optional)"
            type="text"
            name="price_range"
            value={form.price_range}
            onChange={handleChange}
            onBlur={handleBlur}
            placeholder="e.g., Budget, Luxury"
          />
          <FormError message={errors.price_range} />
        </div>
      </div>

      {/* Accords */}
      <div>
        <label className="text-sm font-medium text-gray-700 block mb-2">
          Accords (Optional)
        </label>
        <div className="flex gap-2 mb-2">
          <input
            type="text"
            value={accordInput}
            onChange={(e) => setAccordInput(e.target.value)}
            onKeyPress={(e) => e.key === "Enter" && (e.preventDefault(), handleAddAccord())}
            className="border border-gray-300 rounded-xl px-4 py-2.5 flex-1 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 text-sm"
            placeholder="e.g., Floral; Woody; Citrus or Floral, Woody"
          />
          <button
            type="button"
            onClick={handleAddAccord}
            className="bg-blue-500 text-white px-6 py-2.5 rounded-xl hover:bg-blue-600 transition-colors text-sm font-medium whitespace-nowrap"
          >
            Add
          </button>
        </div>
        {form.accords.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {form.accords.map((accord, index) => (
              <span
                key={index}
                className="bg-gray-100 text-gray-700 px-3 py-1 rounded-full text-sm flex items-center gap-2"
              >
                {accord}
                <button
                  type="button"
                  onClick={() => handleRemoveAccord(index)}
                  className="text-gray-500 hover:text-red-500"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}
        <FormError message={errors.accords} />
      </div>

      {/* Occasions */}
      <div>
        <label className="text-sm font-medium text-gray-700 block mb-2">
          Occasions (Optional)
        </label>
        <div className="flex gap-2 mb-2">
          <input
            type="text"
            value={occasionInput}
            onChange={(e) => setOccasionInput(e.target.value)}
            onKeyPress={(e) => e.key === "Enter" && (e.preventDefault(), handleAddOccasion())}
            className="border border-gray-300 rounded-xl px-4 py-2.5 flex-1 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 text-sm"
            placeholder="e.g., Day; Night; Casual or Day, Night"
          />
          <button
            type="button"
            onClick={handleAddOccasion}
            className="bg-blue-500 text-white px-6 py-2.5 rounded-xl hover:bg-blue-600 transition-colors text-sm font-medium whitespace-nowrap"
          >
            Add
          </button>
        </div>
        {form.occasions.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {form.occasions.map((occasion, index) => (
              <span
                key={index}
                className="bg-gray-100 text-gray-700 px-3 py-1 rounded-full text-sm flex items-center gap-2"
              >
                {occasion}
                <button
                  type="button"
                  onClick={() => handleRemoveOccasion(index)}
                  className="text-gray-500 hover:text-red-500"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}
        <FormError message={errors.occasions} />
      </div>

      {/* Seasons */}
      <div>
        <label className="text-sm font-medium text-gray-700 block mb-2">
          Seasons (Optional)
        </label>
        <div className="flex gap-2 mb-2">
          <input
            type="text"
            value={seasonInput}
            onChange={(e) => setSeasonInput(e.target.value)}
            onKeyPress={(e) => e.key === "Enter" && (e.preventDefault(), handleAddSeason())}
            className="border border-gray-300 rounded-xl px-4 py-2.5 flex-1 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 text-sm"
            placeholder="e.g., Spring; Summer; Fall or Spring, Summer"
          />
          <button
            type="button"
            onClick={handleAddSeason}
            className="bg-blue-500 text-white px-6 py-2.5 rounded-xl hover:bg-blue-600 transition-colors text-sm font-medium whitespace-nowrap"
          >
            Add
          </button>
        </div>
        {form.seasons.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {form.seasons.map((season, index) => (
              <span
                key={index}
                className="bg-gray-100 text-gray-700 px-3 py-1 rounded-full text-sm flex items-center gap-2"
              >
                {season}
                <button
                  type="button"
                  onClick={() => handleRemoveSeason(index)}
                  className="text-gray-500 hover:text-red-500"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}
        <FormError message={errors.seasons} />
      </div>

      {/* Top Notes */}
      <div>
        <label className="text-sm font-medium text-gray-700 block mb-2">
          Top Notes (Optional)
        </label>
        <div className="flex gap-2 mb-2">
          <input
            type="text"
            value={noteTopInput}
            onChange={(e) => setNoteTopInput(e.target.value)}
            onKeyPress={(e) => e.key === "Enter" && (e.preventDefault(), handleAddNoteTop())}
            className="border border-gray-300 rounded-xl px-4 py-2.5 flex-1 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 text-sm"
            placeholder="e.g., Bergamot; Lemon; Orange or Bergamot, Lemon"
          />
          <button
            type="button"
            onClick={handleAddNoteTop}
            className="bg-blue-500 text-white px-6 py-2.5 rounded-xl hover:bg-blue-600 transition-colors text-sm font-medium whitespace-nowrap"
          >
            Add
          </button>
        </div>
        {form.notes_top.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {form.notes_top.map((note, index) => (
              <span
                key={index}
                className="bg-gray-100 text-gray-700 px-3 py-1 rounded-full text-sm flex items-center gap-2"
              >
                {note}
                <button
                  type="button"
                  onClick={() => handleRemoveNoteTop(index)}
                  className="text-gray-500 hover:text-red-500"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}
        <FormError message={errors.notes_top} />
      </div>

      {/* Heart Notes */}
      <div>
        <label className="text-sm font-medium text-gray-700 block mb-2">
          Heart Notes (Optional)
        </label>
        <div className="flex gap-2 mb-2">
          <input
            type="text"
            value={noteHeartInput}
            onChange={(e) => setNoteHeartInput(e.target.value)}
            onKeyPress={(e) => e.key === "Enter" && (e.preventDefault(), handleAddNoteHeart())}
            className="border border-gray-300 rounded-xl px-4 py-2.5 flex-1 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 text-sm"
            placeholder="e.g., Rose; Jasmine; Lavender or Rose, Jasmine"
          />
          <button
            type="button"
            onClick={handleAddNoteHeart}
            className="bg-blue-500 text-white px-6 py-2.5 rounded-xl hover:bg-blue-600 transition-colors text-sm font-medium whitespace-nowrap"
          >
            Add
          </button>
        </div>
        {form.notes_heart.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {form.notes_heart.map((note, index) => (
              <span
                key={index}
                className="bg-gray-100 text-gray-700 px-3 py-1 rounded-full text-sm flex items-center gap-2"
              >
                {note}
                <button
                  type="button"
                  onClick={() => handleRemoveNoteHeart(index)}
                  className="text-gray-500 hover:text-red-500"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}
        <FormError message={errors.notes_heart} />
      </div>

      {/* Base Notes */}
      <div>
        <label className="text-sm font-medium text-gray-700 block mb-2">
          Base Notes (Optional)
        </label>
        <div className="flex gap-2 mb-2">
          <input
            type="text"
            value={noteBaseInput}
            onChange={(e) => setNoteBaseInput(e.target.value)}
            onKeyPress={(e) => e.key === "Enter" && (e.preventDefault(), handleAddNoteBase())}
            className="border border-gray-300 rounded-xl px-4 py-2.5 flex-1 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 text-sm"
            placeholder="e.g., Vanilla; Sandalwood; Musk or Vanilla, Sandalwood"
          />
          <button
            type="button"
            onClick={handleAddNoteBase}
            className="bg-blue-500 text-white px-6 py-2.5 rounded-xl hover:bg-blue-600 transition-colors text-sm font-medium whitespace-nowrap"
          >
            Add
          </button>
        </div>
        {form.notes_base.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {form.notes_base.map((note, index) => (
              <span
                key={index}
                className="bg-gray-100 text-gray-700 px-3 py-1 rounded-full text-sm flex items-center gap-2"
              >
                {note}
                <button
                  type="button"
                  onClick={() => handleRemoveNoteBase(index)}
                  className="text-gray-500 hover:text-red-500"
                >
                  ×
                </button>
              </span>
            ))}
          </div>
        )}
        <FormError message={errors.notes_base} />
      </div>

      {/* Image Upload */}
      {!isEditMode ? (
        <div>
          <label htmlFor="image" className="text-sm font-medium text-gray-700 block mb-2">
            Product Image
          </label>
          <input
            type="file"
            id="image"
            name="image"
            accept="image/jpeg,image/png"
            onChange={handleImageChange}
            onBlur={handleBlur}
            className="border border-gray-300 rounded-xl px-4 py-2.5 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 text-sm"
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
      ) : (
        currentImageUrl && (
          <div>
            <label className="text-sm font-medium text-gray-700 block mb-2">
              Image
            </label>
            <img
              src={currentImageUrl}
              alt="Product"
              className="w-40 h-40 object-cover rounded-lg border border-gray-300"
            />
          </div>
        )
      )}

      {/* Submit Button */}
      <button
        type="submit"
        disabled={loading}
        className="w-full bg-blue-600 text-white py-2.5 rounded-xl font-medium text-sm hover:bg-blue-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed"
      >
        {loading ? loadingText : submitButtonText}
      </button>
    </form>
  );
};

export default ProductForm;
