import React, { useState } from 'react';
import type { OrderResponse } from '../types/order';
import { formatCurrency } from '../utils/format';
import { Link } from 'react-router-dom';
import { useCart } from '../hooks/useCart';
import { useTranslation } from 'react-i18next';
import { cn } from '../utils/cn';

interface Props {
  order: OrderResponse;
}

const statusColor = (status: string) => {
  switch (status) {
    case 'pending':
      return 'bg-yellow-100 text-yellow-800';
    case 'shipped':
      return 'bg-blue-100 text-blue-800';
    case 'delivered':
      return 'bg-green-100 text-green-800';
    case 'cancelled':
      return 'bg-red-100 text-red-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

const OrderDetail: React.FC<Props> = ({ order }) => {
  const { addItem } = useCart();
  const { t } = useTranslation('common');
  const [reordering, setReordering] = useState(false);
  const [toast, setToast] = useState<{ type: 'success' | 'error'; message: string } | null>(null);

  const handleOrderAgain = async () => {
    if (reordering) return;

    setReordering(true);
    try {
      // Add each item from the order to cart sequentially
      for (const item of order.items) {
        await addItem(item.product_slug, item.quantity);
      }
      setToast({ type: 'success', message: t('order.orderAgainSuccess', 'Items added to cart successfully!') });
      setTimeout(() => setToast(null), 2500);
    } catch (error) {
      console.error('Failed to reorder items:', error);
      setToast({ type: 'error', message: t('order.orderAgainError', 'Failed to add items to cart. Please try again.') });
      setTimeout(() => setToast(null), 2500);
    } finally {
      setReordering(false);
    }
  };

  return (
    <div className="mt-3 bg-white p-6 rounded-lg shadow-sm">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6 gap-4">
        <div className="space-y-1">
          <div className="text-sm text-gray-500">{t('order.orderNumber', { id: order.id })}</div>
          <div className="text-lg font-semibold">{new Date(order.created_at).toLocaleString()}</div>
        </div>
        <div className={`inline-flex items-center px-3 py-1 rounded text-sm font-medium ${statusColor(order.status)}`}>
          {t(`order.status.${order.status}`, order.status)}
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {order.items.map((it, index) => (
          <div key={`${it.product_slug}-${index}`} className="flex items-center gap-4 p-3 rounded-lg border border-gray-100">
            <Link to={`/products/${it.product_slug}`} className="w-28 h-28 flex-shrink-0 rounded-md overflow-hidden bg-white flex items-center justify-center" onClick={(e) => e.stopPropagation()}>
              <img src={it.product_image_url || '/placeholder.png'} alt={it.product_name || `Product ${it.product_slug}`} className="max-w-full max-h-full object-contain" />
            </Link>
            <div className="flex-1">
              <Link to={`/products/${it.product_slug}`} onClick={(e) => e.stopPropagation()} className="font-medium text-sm text-gray-900 hover:underline">
                {it.product_name || `Product #${it.product_slug}`}
              </Link>
              <div className="text-sm text-gray-500 mt-1">{t('order.quantity', { quantity: it.quantity })}</div>
              <div className="text-sm text-gray-600 mt-2">{t('order.price', { price: formatCurrency(it.price_at_purchase) })}</div>
            </div>
            <div className="text-right w-36">
              <div className="font-medium">{formatCurrency(it.subtotal)}</div>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-6 flex justify-end">
        <button
          onClick={handleOrderAgain}
          disabled={reordering}
          className={`inline-flex items-center px-4 py-2 rounded ${
            reordering
              ? 'bg-gray-400 cursor-not-allowed'
              : 'bg-blue-600 hover:bg-blue-700 text-white'
          }`}
        >
          {reordering ? t('order.reordering', 'Reordering...') : t('order.orderAgain', 'Order Again')}
        </button>
      </div>

      {/* Toast */}
      {toast && (
        <div
          role="status"
          className={cn(
            'fixed top-4 right-4 z-50 px-4 py-3 rounded shadow-lg text-sm',
            toast.type === 'success' ? 'bg-green-600 text-white' : 'bg-red-600 text-white'
          )}
        >
          {toast.message}
        </div>
      )}
    </div>
  );
};

export default OrderDetail;
