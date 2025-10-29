import React from 'react';
import TrashIcon from './TrashIcon';
import { formatCurrency } from '../utils/format';
import { PLACEHOLDER_IMAGE } from '../constants/app';
import { cn } from '../utils/cn';
import type { CartItem as CartItemType } from '../types/cart';

interface CartItemProps {
  item: CartItemType;
  onRemove?: (itemId: number) => void;
  isRemoving?: boolean;
  showRemoveButton?: boolean;
  compact?: boolean;
}

const CartItem: React.FC<CartItemProps> = ({ 
  item, 
  onRemove, 
  isRemoving = false,
  showRemoveButton = true,
  compact = false 
}) => {
  const imageSize = compact ? 'h-12 w-12' : 'h-16 w-16';
  const textSize = compact ? 'text-xs' : 'text-sm';
  const padding = compact ? 'p-3' : 'py-4';

  return (
    <div className={cn('flex gap-3 items-center', padding)}>
      <img
        src={item.product?.image_url || PLACEHOLDER_IMAGE}
        alt={item.product?.name || 'Product image'}
        className={cn(imageSize, 'object-contain bg-gray-50 rounded border border-gray-200 p-1')}
        onError={(e) => {
          (e.currentTarget as HTMLImageElement).src = PLACEHOLDER_IMAGE;
        }}
      />
      
      <div className="flex-1 min-w-0">
        <p className={cn(textSize, 'font-medium text-gray-900 truncate')}>
          {item.product?.name || 'Product'}
        </p>
        <p className={cn(textSize, 'text-gray-600')}>
          Qty: {item.quantity}
        </p>
      </div>
      
      <div className="flex items-center gap-2">
        <div className={cn(textSize, 'font-semibold text-gray-900 whitespace-nowrap')}>
          {formatCurrency(item.total)}
        </div>
        
        {showRemoveButton && onRemove && (
          <button
            type="button"
            onClick={() => onRemove(item.id)}
            disabled={isRemoving}
            title="Remove item"
            aria-label={`Remove ${item.product?.name || 'item'}`}
            className={cn(
              'inline-flex items-center justify-center rounded-md border p-2 transition-colors',
              isRemoving 
                ? 'opacity-60 cursor-not-allowed border-gray-200' 
                : 'border-gray-200 hover:bg-red-50 hover:text-red-600 hover:border-red-200'
            )}
          >
            <TrashIcon className={compact ? 'w-4 h-4' : 'w-5 h-5'} />
          </button>
        )}
      </div>
    </div>
  );
};

export default CartItem;
