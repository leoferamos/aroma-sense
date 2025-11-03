import React, { useState } from 'react';
import type { OrderResponse } from '../types/order';
import { useCart } from '../contexts/CartContext';
import { useNavigate } from 'react-router-dom';

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
  const navigate = useNavigate();
  const [adding, setAdding] = useState(false);

  const handleOrderAgain = async () => {
    setAdding(true);
    try {
      // Add each item to the cart
      for (const it of order.items) {
        await addItem(it.product_id, it.quantity);
      }
      // Redirect to checkout
      navigate('/checkout');
    } catch {
      // ignore
    } finally {
      setAdding(false);
    }
  };

  return (
    <div className="mt-3 bg-white p-4 rounded shadow-sm">
      <div className="flex items-center justify-between mb-4">
        <div className="space-y-1">
          <div className="text-sm text-gray-500">Order #{order.id}</div>
          <div className="text-sm text-gray-600">{new Date(order.created_at).toLocaleString()}</div>
        </div>
        <div className={`inline-flex items-center px-3 py-1 rounded text-sm font-medium ${statusColor(order.status)}`}>
          {order.status}
        </div>
      </div>

      <div className="space-y-4">
        {order.items.map((it) => (
          <div key={it.id} className="flex items-center gap-4">
            <img src={it.product_image_url || '/placeholder.png'} alt={it.product_name} className="w-20 h-20 object-cover rounded" />
            <div className="flex-1">
              <div className="font-medium">{it.product_name || `Product #${it.product_id}`}</div>
              <div className="text-sm text-gray-500">Qty: {it.quantity}</div>
            </div>
            <div className="text-right w-32">
              <div className="font-medium">${it.price_at_purchase.toFixed(2)}</div>
              <div className="text-sm text-gray-500">Subtotal: ${it.subtotal.toFixed(2)}</div>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-6 flex justify-end">
        <button
          onClick={handleOrderAgain}
          disabled={adding}
          className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-60"
        >
          {adding ? 'Adding...' : 'Order Again'}
        </button>
      </div>
    </div>
  );
};

export default OrderDetail;
