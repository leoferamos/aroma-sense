import React from 'react';
import { Link } from 'react-router-dom';
import type { Product } from '../types/product';
import { formatCurrency } from '../utils/format';
import { cn } from '../utils/cn';
import { PLACEHOLDER_IMAGE, LOW_STOCK_THRESHOLD } from '../constants/app';

interface ProductCardProps {
  product: Product;
  onAddToCart?: (product: Product) => void;
}

const ProductCard: React.FC<ProductCardProps> = ({ product, onAddToCart }) => {

  const handleAddToCart = (e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent navigation when clicking "Add to Cart"
    if (onAddToCart) {
      onAddToCart(product);
    }
  };

  const isOutOfStock = product.stock_quantity === 0;
  const isLowStock = product.stock_quantity > 0 && product.stock_quantity < LOW_STOCK_THRESHOLD;

  return (
    <div className="bg-white rounded-lg shadow-md hover:shadow-xl transition-shadow duration-300 overflow-hidden group">
      {/* Clickable area*/}
      <Link
        to={`/products/${product.id}`}
        aria-label={`View details for ${product.name}`}
        className="block cursor-pointer"
      >
        {/* Image Container */}
        <div className="relative h-64 bg-white overflow-hidden flex items-center justify-center p-4">
          <img
            src={product.image_url || PLACEHOLDER_IMAGE}
            alt={product.name}
            className="max-h-full max-w-full object-contain transition-transform duration-300 group-hover:scale-105 bg-white"
            onError={(e) => {
              e.currentTarget.src = PLACEHOLDER_IMAGE;
            }}
          />
          {isOutOfStock && (
            <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center">
              <span className="text-white font-bold text-lg">Out of Stock</span>
            </div>
          )}
        </div>

        {/* Content */}
        <div className="px-4 pb-2">
          {/* Brand */}
          <p className="text-sm text-gray-500 font-medium mb-1">{product.brand}</p>

          {/* Product Name */}
          <h3 className="text-lg font-semibold text-gray-900 mb-2 line-clamp-2 min-h-[3.5rem]">
            {product.name}
          </h3>

          {/* Category & Weight */}
          <div className="flex items-center gap-2 text-sm text-gray-600 mb-3">
            <span>{product.category}</span>
            <span>•</span>
            <span>{product.weight}ml</span>
          </div>

          {/* Price */}
          <div className="flex items-center justify-between mb-3">
            <span className="text-2xl font-bold text-blue-600">
              {formatCurrency(product.price)}
            </span>
            {isLowStock && (
              <span className="text-xs text-orange-600 font-medium">
                Only {product.stock_quantity} left
              </span>
            )}
          </div>
        </div>
      </Link>

      {/* Actions */}
      <div className="px-4 pb-4">
        <button
          onClick={handleAddToCart}
          disabled={isOutOfStock}
          className={cn(
            'w-full py-2.5 px-4 rounded-lg font-medium transition-colors duration-200',
            isOutOfStock
              ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
              : 'bg-blue-600 text-white hover:bg-blue-700 active:bg-blue-800'
          )}
        >
          {isOutOfStock ? 'Out of Stock' : 'Add to Cart'}
        </button>
      </div>
    </div>
  );
};

export default React.memo(ProductCard);
