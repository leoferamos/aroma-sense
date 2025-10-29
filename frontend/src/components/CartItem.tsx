import React, { memo } from 'react';
import TrashIcon from './TrashIcon';
import { formatCurrency } from '../utils/format';
import { PLACEHOLDER_IMAGE } from '../constants/app';
import { cn } from '../utils/cn';
import type { CartItem as CartItemType } from '../types/cart';
import { useCartItemQuantity } from '../hooks/useCartItemQuantity';

interface CartItemProps {
  item: CartItemType;
  onRemove?: (itemId: number) => void;
  isRemoving?: boolean;
  showRemoveButton?: boolean;
  compact?: boolean;
  showQuantityControls?: boolean;
}

const CartItem: React.FC<CartItemProps> = ({ 
  item, 
  onRemove, 
  isRemoving = false,
  showRemoveButton = true,
  compact = false,
  showQuantityControls = false
}) => {
  const { quantity, increment, decrement, error } = useCartItemQuantity({
    itemId: item.id,
    initialQuantity: item.quantity,
  });

  const imageSize = compact ? 'h-12 w-12' : 'h-16 w-16';
  const textSize = compact ? 'text-xs' : 'text-sm';
  const padding = compact ? 'p-3' : 'py-4';
  const isMinQuantity = quantity <= 1;

  // Shared button classes for quantity controls
  const quantityButtonBase = 'h-6 w-6 rounded border flex items-center justify-center font-medium transition-colors';
  const quantityButtonActive = 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50 hover:border-gray-400';
  const quantityButtonDisabled = 'bg-gray-100 text-gray-400 border-gray-200 cursor-not-allowed';

  return (
    <div className={cn('flex flex-col gap-2', padding)}>
      {error && (
        <div className="text-xs text-red-600 bg-red-50 p-2 rounded">
          {error}
        </div>
      )}
      
      <div className="flex gap-3 items-center">
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
        
        {showQuantityControls ? (
          <div className="flex items-center gap-2 mt-1">
            <button
              type="button"
              onClick={decrement}
              disabled={isMinQuantity}
              className={cn(
                quantityButtonBase,
                isMinQuantity ? quantityButtonDisabled : quantityButtonActive
              )}
              aria-label="Decrease quantity"
            >
              −
            </button>
            
            <span className="min-w-[2rem] text-center font-medium text-gray-900">
              {quantity}
            </span>
            
            <button
              type="button"
              onClick={increment}
              className={cn(quantityButtonBase, quantityButtonActive)}
              aria-label="Increase quantity"
            >
              +
            </button>
          </div>
        ) : (
          <p className={cn(textSize, 'text-gray-600')}>
            Qty: {quantity}
          </p>
        )}
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
    </div>
  );
};

// Memoize to prevent unnecessary re-renders when parent updates
export default memo(CartItem);
