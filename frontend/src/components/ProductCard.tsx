import React from 'react';
import { Link } from 'react-router-dom';
import type { Product } from '../types/product';
import { formatCurrency } from '../utils/format';
import { cn } from '../utils/cn';
import { PLACEHOLDER_IMAGE, LOW_STOCK_THRESHOLD } from '../constants/app';
import { useTranslation } from 'react-i18next';

interface ProductCardProps {
  product: Product;
  onAddToCart?: (product: Product) => void;
  showAddToCart?: boolean;
}

const ProductCard: React.FC<ProductCardProps> = ({ product, onAddToCart, showAddToCart = true }) => {
  const { t } = useTranslation('common');

  const handleAddToCart = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (onAddToCart) {
      onAddToCart(product);
    }
  };

  const isOutOfStock = product.stock_quantity === 0;
  const isLowStock = product.stock_quantity > 0 && product.stock_quantity < LOW_STOCK_THRESHOLD;

  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-lg transition-all duration-200 transform overflow-hidden hover:-translate-y-1 hover:shadow-xl">
      {/* Clickable area*/}
      <Link
        to={`/product/${product.slug}`}
        aria-label={`View details for ${product.name}`}
        className="block cursor-pointer"
      >
        {/* Image Container */}
        <div className="relative h-60 bg-transparent overflow-hidden flex items-center justify-center p-6">
          <img
            src={product.image_url || PLACEHOLDER_IMAGE}
            alt={product.name}
            className="max-h-full max-w-full object-contain transition-transform duration-300 hover:scale-105"
            onError={(e) => {
              e.currentTarget.src = PLACEHOLDER_IMAGE;
            }}
          />
          {isOutOfStock && (
            <div className="absolute inset-0 bg-black bg-opacity-40 flex items-center justify-center">
              <span className="text-white font-semibold text-base">Out of Stock</span>
            </div>
          )}
        </div>

        {/* Content */}
        <div className="px-6 pt-4 pb-4">
          {/* Brand */}
          <p className="text-xs tracking-widest text-gray-500 uppercase font-semibold mb-1">{product.brand}</p>

          {/* Product Name */}
          <h3 className="text-lg font-semibold text-gray-900 mb-3 line-clamp-2 leading-tight min-h-[3rem]">
            {product.name}
          </h3>

          {/* Category & Weight (chips) */}
          <div className="flex items-center gap-2 text-xs text-gray-700 mb-3 flex-wrap">
            <span className="bg-gray-100 text-gray-700 border border-gray-200 rounded-full px-3 py-1 text-xs font-medium">{product.category}</span>
            <span className="bg-gray-100 text-gray-700 border border-gray-200 rounded-full px-3 py-1 text-xs font-medium">{product.weight}ml</span>
          </div>

          {/* Price */}
          <div className="mb-3">
            <div className="flex items-center justify-between mb-1">
              <span className="text-xl font-bold text-blue-600">{formatCurrency(product.price)}</span>
              {isLowStock && (
                <span className="text-xs font-medium px-2 py-1 bg-orange-50 text-orange-700 rounded-full">{product.stock_quantity} left</span>
              )}
            </div>
            <p className="text-xs text-gray-600">ou 10x de {formatCurrency(product.price / 10)} sem juros</p>
          </div>
        </div>
      </Link>

      {/* Actions (optional) */}
      {showAddToCart && (
        <div className="px-4 pb-4">
          <button
            onClick={handleAddToCart}
            disabled={isOutOfStock}
            className={cn(
              'w-full py-2.5 px-4 rounded-lg font-medium transition-all duration-200 text-sm',
              isOutOfStock
                ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                : 'bg-blue-600 text-white hover:bg-blue-700 active:scale-95 shadow-sm hover:shadow-md'
            )}
          >
            {isOutOfStock ? t('products.outOfStock') : t('products.addToCart')}
          </button>
        </div>
      )}
    </div>
  );
};

export default React.memo(ProductCard);
