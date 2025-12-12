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
    <>
      {/* Desktop Table View - Hidden on mobile */}
      <div className="hidden lg:block overflow-x-auto bg-white border border-gray-200 rounded-lg shadow-sm">
        <table className="min-w-full" role="table">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">{t('orderId')}</th>
              <th className="px-4 py-3 text-right font-semibold text-gray-700">{t('total')}</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">{t('status')}</th>
              <th className="px-4 py-3 text-left font-semibold text-gray-700">{t('date')}</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {orders.length === 0 ? (
              <tr><td colSpan={4} className="px-4 py-8 text-center text-gray-500">No orders found</td></tr>
            ) : orders.map((o) => (
              <tr
                key={o.id}
                className="hover:bg-gray-50 cursor-pointer transition-colors"
                role="button"
                tabIndex={0}
                aria-label={`Open order ${o.id}`}
                onClick={() => navigate(`/admin/orders/${o.id}`)}
                onKeyDown={(e) => onRowKeyDown(e, o.id)}
              >
                <td className="px-4 py-3 text-gray-900 font-medium">#{o.id}</td>
                <td className="px-4 py-3 text-right text-gray-900 font-semibold">{formatCurrency(o.total_amount ?? 0)}</td>
                <td className="px-4 py-3">
                  <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${o.status === 'delivered' ? 'bg-green-100 text-green-700' :
                      o.status === 'shipped' ? 'bg-blue-100 text-blue-700' :
                        o.status === 'processing' ? 'bg-yellow-100 text-yellow-700' :
                          o.status === 'cancelled' ? 'bg-red-100 text-red-700' :
                            'bg-gray-100 text-gray-700'
                    }`}>
                    {o.status || t('unknown')}
                  </span>
                </td>
                <td className="px-4 py-3 text-gray-600 text-sm">{o.created_at ? new Date(o.created_at).toLocaleString() : '-'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Mobile Card View - Visible only on mobile/tablet */}
      <div className="lg:hidden space-y-3">
        {orders.length === 0 ? (
          <div className="bg-white border border-gray-200 rounded-lg p-6 text-center text-gray-500">{t('noOrdersFound')}</div>
        ) : orders.map((o) => (
          <div
            key={o.id}
            className="bg-white border border-gray-200 rounded-lg p-4 shadow-sm cursor-pointer hover:shadow-md transition-shadow"
            role="button"
            tabIndex={0}
            aria-label={`Open order ${o.id}`}
            onClick={() => navigate(`/admin/orders/${o.id}`)}
            onKeyDown={(e) => onRowKeyDown(e, o.id)}
          >
            <div className="flex items-start justify-between mb-3">
              <div>
                <div className="text-sm text-gray-500 mb-1">{t('orderId')}</div>
                <div className="font-semibold text-gray-900">#{o.id}</div>
              </div>
              <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full whitespace-nowrap ${o.status === 'delivered' ? 'bg-green-100 text-green-700' :
                  o.status === 'shipped' ? 'bg-blue-100 text-blue-700' :
                    o.status === 'processing' ? 'bg-yellow-100 text-yellow-700' :
                      o.status === 'cancelled' ? 'bg-red-100 text-red-700' :
                        'bg-gray-100 text-gray-700'
                }`}>
                {o.status || t('unknown')}
              </span>
            </div>
            <div className="border-t border-gray-100 pt-3 space-y-2">
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-500">{t('total')}:</span>
                <span className="text-lg font-bold text-gray-900">{formatCurrency(o.total_amount ?? 0)}</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-500">{t('date')}:</span>
                <span className="text-xs text-gray-600">{o.created_at ? new Date(o.created_at).toLocaleString() : '-'}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </>
  );
};

export default OrdersTable;
