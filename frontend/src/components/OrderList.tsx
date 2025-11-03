import React from 'react';
import type { OrderResponse } from '../types/order';

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
      {orders.map((o) => (
        <div
          key={o.id}
          tabIndex={0}
          role="button"
          onClick={() => onSelect(o.id)}
          onKeyDown={(e) => { if (e.key === 'Enter') onSelect(o.id); }}
          className="p-4 bg-white rounded shadow-sm hover:shadow-md cursor-pointer focus:outline-none focus:ring"
        >
          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm text-gray-500">Order #{o.id}</div>
              <div className="text-lg font-medium">{new Date(o.created_at).toLocaleString()}</div>
            </div>
            <div className="text-right">
              <div className="text-lg font-semibold">${o.total_amount.toFixed(2)}</div>
              <div className={`inline-flex items-center px-2 py-1 mt-2 text-xs font-medium rounded ${statusColor(o.status)}`}>
                {o.status}
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default OrderList;
