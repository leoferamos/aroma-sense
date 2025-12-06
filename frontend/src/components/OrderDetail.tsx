import React from 'react';
import type { OrderResponse } from '../types/order';
import { formatCurrency } from '../utils/format';
import { Link } from 'react-router-dom';

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
  const handleOrderAgain = async () => {
    // This functionality is disabled because product IDs are not exposed to the frontend
    alert('Order Again functionality is currently not available. Please add items manually to your cart.');
  };

  return (
    <div className="mt-3 bg-white p-6 rounded-lg shadow-sm">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6 gap-4">
        <div className="space-y-1">
          <div className="text-sm text-gray-500">Order #{order.id}</div>
          <div className="text-lg font-semibold">{new Date(order.created_at).toLocaleString()}</div>
        </div>
        <div className={`inline-flex items-center px-3 py-1 rounded text-sm font-medium ${statusColor(order.status)}`}>
          {order.status}
        </div>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        {order.items.map((it) => (
          <div key={it.id} className="flex items-center gap-4 p-3 rounded-lg border border-gray-100">
            <Link to={`/products/${it.product_slug}`} className="w-28 h-28 flex-shrink-0 rounded-md overflow-hidden bg-white flex items-center justify-center" onClick={(e) => e.stopPropagation()}>
              <img src={it.product_image_url || '/placeholder.png'} alt={it.product_name || `Product ${it.product_slug}`} className="max-w-full max-h-full object-contain" />
            </Link>
            <div className="flex-1">
              <Link to={`/products/${it.product_slug}`} onClick={(e) => e.stopPropagation()} className="font-medium text-sm text-gray-900 hover:underline">
                {it.product_name || `Product #${it.product_slug}`}
              </Link>
              <div className="text-sm text-gray-500 mt-1">Qty: {it.quantity}</div>
              <div className="text-sm text-gray-600 mt-2">Price: {formatCurrency(it.price_at_purchase)}</div>
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
          className="inline-flex items-center px-4 py-2 bg-gray-400 text-white rounded cursor-not-allowed"
          disabled
        >
          Order Again (Disabled)
        </button>
      </div>
    </div>
  );
};

export default OrderDetail;
