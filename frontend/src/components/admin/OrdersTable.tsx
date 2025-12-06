import React from 'react';
import { useNavigate } from 'react-router-dom';
import type { AdminOrder } from '../../services/admin';
import { formatCurrency } from '../../utils/format';
import { useTranslation } from 'react-i18next';

type Props = {
  orders: AdminOrder[];
};

const OrdersTable: React.FC<Props> = ({ orders }) => {
  const navigate = useNavigate();
  const { t } = useTranslation('admin');

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
          <th className="px-4 py-2 border text-left">{t('orderId')}</th>
          <th className="px-4 py-2 border text-right">{t('total')}</th>
          <th className="px-4 py-2 border text-left">{t('status')}</th>
          <th className="px-4 py-2 border text-left">{t('date')}</th>
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
            <td className="px-4 py-2 border">{o.status || 'Unknown'}</td>
            <td className="px-4 py-2 border">{o.created_at ? new Date(o.created_at).toLocaleString() : '-'}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default OrdersTable;
