import React from 'react';
import type { OrderResponse } from '../types/order';
import { formatCurrency } from '../utils/format';
import { Link } from 'react-router-dom';

interface Props {
  orders: OrderResponse[];
  onSelect: (id: number) => void;
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

const OrderList: React.FC<Props> = ({ orders, onSelect }) => {
  return (
    <div className="space-y-4">
      {orders.map((o) => {
        const thumbs = o.items.slice(0, 4);
        return (
          <div
            key={o.id}
            tabIndex={0}
            role="button"
            onClick={() => onSelect(o.id)}
            onKeyDown={(e) => { if (e.key === 'Enter') onSelect(o.id); }}
            className="p-6 bg-white rounded-lg shadow-sm hover:shadow-md cursor-pointer focus:outline-none focus:ring transition-all"
          >
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              <div className="flex items-center gap-4">
                <div className="flex -space-x-3 items-center">
                  {thumbs.map((it, index) => (
                    <Link key={index} to={`/products/${it.product_slug}`} onClick={(e) => e.stopPropagation()} className="w-20 h-20 rounded-lg overflow-hidden bg-white ring-1 ring-gray-100 block">
                      <img src={it.product_image_url || '/placeholder.png'} alt={it.product_name || ''} className="w-full h-full object-contain" />
                    </Link>
                  ))}
                </div>

                <div>
                  <div className="text-sm text-gray-500">Order #{o.id}</div>
                  <div className="text-base font-semibold">{new Date(o.created_at).toLocaleString()}</div>
                  <div className="text-sm text-gray-600 mt-1">{o.item_count} {o.item_count === 1 ? 'item' : 'items'}</div>
                  <div className="text-sm text-gray-700 mt-2 truncate w-40">
                    {o.items.length > 0 ? (
                      <Link to={`/products/${o.items[0].product_slug}`} onClick={(e) => e.stopPropagation()} className="hover:underline">
                        {o.items[0].product_name}
                      </Link>
                    ) : ''}
                    {o.items.length > 1 ? ` +${o.items.length - 1} more` : ''}
                  </div>
                </div>
              </div>

              <div className="flex items-center justify-between sm:justify-end gap-4">
                <div className="text-right">
                  <div className="text-lg font-semibold">{formatCurrency(o.total_amount)}</div>
                </div>
                <div>
                  <span className={`inline-flex items-center px-3 py-1 text-sm font-medium rounded ${statusColor(o.status)}`}>
                    {o.status}
                  </span>
                </div>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default OrderList;
