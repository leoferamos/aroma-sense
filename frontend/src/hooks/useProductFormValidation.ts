import { useState, useEffect } from "react";
import { messages } from "../constants/messages";
import type { CreateProductFormData, ProductFormErrors } from "../types/product";

export interface ProductFormTouched {
  name: boolean;
  brand: boolean;
  weight: boolean;
  description: boolean;
  price: boolean;
  category: boolean;
  accords: boolean;
  occasions: boolean;
  seasons: boolean;
  intensity: boolean;
  gender: boolean;
  price_range: boolean;
  notes_top: boolean;
  notes_heart: boolean;
  notes_base: boolean;
  stock_quantity: boolean;
  image: boolean;
}

interface UseProductFormValidationProps {
  form: CreateProductFormData;
  touched: ProductFormTouched;
  isEditMode?: boolean;
}

export function useProductFormValidation({
  form,
  touched,
  isEditMode = false,
}: UseProductFormValidationProps) {
  const [errors, setErrors] = useState<ProductFormErrors>({});

  useEffect(() => {
    const newErrors: ProductFormErrors = {};

    // Name validation
    if (touched.name) {
      if (!form.name || form.name.trim() === "") {
        newErrors.name = messages.productNameRequired;
      } else if (form.name.length < 3) {
        newErrors.name = messages.productNameMinLength;
      } else if (form.name.length > 100) {
        newErrors.name = messages.productNameMaxLength;
      }
    }

    // Brand validation
    if (touched.brand) {
      if (!form.brand || form.brand.trim() === "") {
        newErrors.brand = messages.productBrandRequired;
      } else if (form.brand.length < 2) {
        newErrors.brand = messages.productBrandMinLength;
      }
    }

    // Weight validation
    if (touched.weight) {
      if (form.weight <= 0) {
        newErrors.weight = messages.productWeightInvalid;
      } else if (form.weight > 10000) {
        newErrors.weight = messages.productWeightMax;
      }
    }

    // Description validation
    if (touched.description) {
      if (!form.description || form.description.trim() === "") {
        newErrors.description = messages.productDescriptionRequired;
      } else if (form.description.length < 10) {
        newErrors.description = messages.productDescriptionMinLength;
      } else if (form.description.length > 1000) {
        newErrors.description = messages.productDescriptionMaxLength;
      }
    }

    // Price validation
    if (touched.price) {
      if (form.price <= 0) {
        newErrors.price = messages.productPriceInvalid;
      } else if (form.price > 1000000) {
        newErrors.price = messages.productPriceMax;
      }
    }

    // Category validation
    if (touched.category) {
      if (!form.category || form.category.trim() === "") {
        newErrors.category = messages.productCategoryRequired;
      }
    }

    // Notes validation (separate for top, heart, base)
    if (touched.notes_top) {
      if (form.notes_top.length > 0) {
        const hasEmptyNote = form.notes_top.some((note) => note.trim() === "");
        if (hasEmptyNote) {
          newErrors.notes_top = messages.productNotesInvalid;
        }
      }
    }

    if (touched.notes_heart) {
      if (form.notes_heart.length > 0) {
        const hasEmptyNote = form.notes_heart.some((note) => note.trim() === "");
        if (hasEmptyNote) {
          newErrors.notes_heart = messages.productNotesInvalid;
        }
      }
    }

    if (touched.notes_base) {
      if (form.notes_base.length > 0) {
        const hasEmptyNote = form.notes_base.some((note) => note.trim() === "");
        if (hasEmptyNote) {
          newErrors.notes_base = messages.productNotesInvalid;
        }
      }
    }

    // Stock quantity validation
    if (touched.stock_quantity) {
      if (form.stock_quantity < 0) {
        newErrors.stock_quantity = messages.productStockNegative;
      } else if (!Number.isInteger(form.stock_quantity)) {
        newErrors.stock_quantity = messages.productStockInteger;
      }
    }

    // Image validation
    if (touched.image) {
      if (!isEditMode && !form.image) {
        newErrors.image = messages.productImageRequired;
      } else if (form.image) {
        const validTypes = ["image/jpeg", "image/png"];
        if (!validTypes.includes(form.image.type)) {
          newErrors.image = messages.productImageTypeInvalid;
        } else if (form.image.size > 5 * 1024 * 1024) {
          newErrors.image = messages.productImageSizeExceeded;
        }
      }
    }

    setErrors(newErrors);
  }, [form, touched, isEditMode]);

  const isValid = Object.keys(errors).length === 0;

  return { errors, isValid };
}
