import React from 'react';
import { Link } from 'react-router-dom';
import FiltersBar from '../../components/admin/FiltersBar';
import OrdersTable from '../../components/admin/OrdersTable';
import PaginationControls from '../../components/admin/PaginationControls';
import { useAdminOrders } from '../../hooks/useAdminOrders';
import AdminLayout from '../../components/admin/AdminLayout';
import { useTranslation } from 'react-i18next';

const AdminOrders: React.FC = () => {
  const { data, loading, error, params, setPage, setPerPage, setStatus, setDateRange } = useAdminOrders({ page: 1, per_page: 25 });
  const { t } = useTranslation('admin');

  const orders = data?.orders ?? [];
  const page = data?.page ?? 1;
  const perPage = data?.per_page ?? 25;
  const totalPages = Math.ceil((data?.total ?? 0) / perPage) || 1;

  const actions = (
    <div className="flex items-center gap-2">
      <Link to="/admin/dashboard" className="px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-lg hover:bg-gray-50">← {t('nav.dashboard')}</Link>
    </div>
  );

  return (
    <AdminLayout actions={actions}>
      <div className="p-4 sm:p-6">
        <h1 className="text-xl sm:text-2xl font-semibold mb-4">{t('orders')}</h1>

        <FiltersBar
          status={params.status ?? ''}
          onStatusChange={(s) => setStatus(s || undefined)}
          startDate={params.start_date}
          endDate={params.end_date}
          onDateChange={(start, end) => setDateRange(start, end)}
          perPage={perPage}
          onPerPageChange={(n) => setPerPage(n)}
        />

        {loading && <div className="py-8 text-center text-gray-500">{t('loadingOrders')}</div>}
        {error && <div className="py-8 text-red-600 p-3 bg-red-50 border border-red-200 rounded-lg text-sm">{error}</div>}

        {!loading && !error && (
          <>
            <div className="mb-4 text-sm text-gray-700">
              <strong>{data?.total ?? 0}</strong> {t('totalOrders', { count: data?.total ?? 0 })}
            </div>

            <OrdersTable orders={orders} />

            <div className="mt-4">
              <PaginationControls
                page={page}
                totalPages={totalPages}
                onPrev={() => setPage(Math.max(1, page - 1))}
                onNext={() => setPage(Math.min(totalPages, page + 1))}
              />
            </div>
          </>
        )}
      </div>
    </AdminLayout>
  );
};

export default AdminOrders;
