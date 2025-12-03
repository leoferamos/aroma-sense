import React from 'react';
import { useNavigate } from 'react-router-dom';
import type { AdminOrder } from '../../services/admin';
import { formatCurrency } from '../../utils/format';

type Props = {
  orders: AdminOrder[];
};

const OrdersTable: React.FC<Props> = ({ orders }) => {
  const navigate = useNavigate();

  const onRowKeyDown = (e: React.KeyboardEvent, id: string | number) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      navigate(`/admin/orders/${id}`);
    }
  };

  return (
    <table className="min-w-full bg-white border" role="table">
      <thead>
        <tr>
          <th className="px-4 py-2 border text-left">Order ID</th>
          <th className="px-4 py-2 border text-right">Total</th>
          <th className="px-4 py-2 border text-left">Status</th>
          <th className="px-4 py-2 border text-left">Date</th>
        </tr>
      </thead>
      <tbody>
        {orders.map((o) => (
          <tr
            key={o.id}
            className="hover:bg-gray-50 cursor-pointer"
            role="button"
            tabIndex={0}
            aria-label={`Open order ${o.id}`}
            onClick={() => navigate(`/admin/orders/${o.id}`)}
            onKeyDown={(e) => onRowKeyDown(e, o.id)}
          >
            <td className="px-4 py-2 border">{o.id}</td>
            <td className="px-4 py-2 border text-right">{formatCurrency(o.total_amount ?? 0)}</td>
            <td className="px-4 py-2 border">{o.status}</td>
            <td className="px-4 py-2 border">{o.created_at ? new Date(o.created_at).toLocaleString() : '-'}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default OrdersTable;
