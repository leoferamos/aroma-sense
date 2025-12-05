import React from 'react';
import type { Product } from '../../types/product';
import { PLACEHOLDER_IMAGE } from '../../constants/app';
import { formatCurrency } from '../../utils/format';
import { cn } from '../../utils/cn';
import TrashIcon from '../TrashIcon';
import PencilIcon from '../PencilIcon';
import { useTranslation } from 'react-i18next';

interface AdminProductCardProps {
  product: Product;
  onEdit?: (product: Product) => void;
  onDelete?: (product: Product) => void;
}

const AdminProductCard: React.FC<AdminProductCardProps> = ({ product, onEdit, onDelete }) => {
  const { t } = useTranslation('common');
  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
  <div className="relative h-48 bg-white flex items-center justify-center p-3">
        <img
          src={product.image_url || PLACEHOLDER_IMAGE}
          alt={product.name}
          className="max-h-full max-w-full object-contain"
          onError={(e) => {
            (e.currentTarget as HTMLImageElement).src = PLACEHOLDER_IMAGE;
          }}
        />

        {/* Overlay actions */}
        <div className="absolute inset-x-0 top-0 p-2 flex items-center justify-between pointer-events-none">
          <button
            type="button"
            onClick={() => onEdit?.(product)}
            className="pointer-events-auto inline-flex items-center justify-center rounded-md border border-gray-200 bg-white/90 hover:bg-blue-50 hover:border-blue-200 text-blue-600 p-2 shadow-sm"
            title="Edit product"
            aria-label="Edit product"
          >
            <PencilIcon className="w-4 h-4" />
          </button>
          <button
            type="button"
            onClick={() => onDelete?.(product)}
            className="pointer-events-auto inline-flex items-center justify-center rounded-md border border-gray-200 bg-white/90 hover:bg-red-50 hover:border-red-200 text-red-600 p-2 shadow-sm"
            title="Delete product"
            aria-label="Delete product"
          >
            <TrashIcon className="w-4 h-4" />
          </button>
        </div>
      </div>

      <div className="p-4">
        <h3 className="font-semibold text-gray-900 line-clamp-2 min-h-[2.5rem]">{product.name}</h3>
        <div className="mt-1 text-sm text-gray-600 flex items-center gap-2">
          <span>{product.brand}</span>
          <span>•</span>
          <span>{product.category}</span>
          <span>•</span>
          <span>{product.weight}ml</span>
        </div>
        <div className="mt-2 flex items-center justify-between">
          <span className="text-blue-600 font-semibold">{formatCurrency(product.price)}</span>
          <span className={cn('text-xs font-medium', product.stock_quantity === 0 ? 'text-red-600' : 'text-gray-600')}>
            {product.stock_quantity === 0 ? t('products.outOfStock') : `${product.stock_quantity} ${t('products.inStock')}`}
          </span>
        </div>
      </div>
    </div>
  );
};

export default AdminProductCard;
