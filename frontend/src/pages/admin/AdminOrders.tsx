import React from 'react';
import { Link } from 'react-router-dom';
import FiltersBar from '../../components/admin/FiltersBar';
import OrdersTable from '../../components/admin/OrdersTable';
import PaginationControls from '../../components/admin/PaginationControls';
import { useAdminOrders } from '../../hooks/useAdminOrders';
import { formatCurrency } from '../../utils/format';
import AdminLayout from '../../components/admin/AdminLayout';

const AdminOrders: React.FC = () => {
  const { data, loading, error, params, setPage, setPerPage, setStatus, setDateRange } = useAdminOrders({ page: 1, per_page: 25 });

  const orders = data?.orders ?? [];
  const page = params.page ?? 1;
  const perPage = params.per_page ?? 25;
  const totalPages = data?.meta.pagination.total_pages ?? 1;

  const actions = (
    <div className="flex items-center gap-2">
      <Link to="/admin/dashboard" className="px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-lg hover:bg-gray-50">← Dashboard</Link>
    </div>
  );

  return (
    <AdminLayout title="Orders" actions={actions}>
      <div className="p-6">
        <h1 className="text-2xl font-semibold mb-4">Orders</h1>

        <FiltersBar
          status={params.status ?? ''}
          onStatusChange={(s) => setStatus(s || undefined)}
          startDate={params.start_date}
          endDate={params.end_date}
          onDateChange={(start, end) => setDateRange(start, end)}
          perPage={perPage}
          onPerPageChange={(n) => setPerPage(n)}
        />

        {loading && <div className="py-8">Loading orders…</div>}
        {error && <div className="py-8 text-red-600">{error}</div>}

        {!loading && !error && (
          <>
            <div className="mb-4 text-sm text-gray-700">
              <strong>{data?.meta.pagination.total_count ?? 0}</strong> orders — Total revenue: <strong>{formatCurrency(data?.meta.stats.total_revenue ?? 0)}</strong>
            </div>

            <div className="overflow-x-auto">
              <OrdersTable orders={orders} />
            </div>

            <PaginationControls
              page={page}
              totalPages={totalPages}
              onPrev={() => setPage(Math.max(1, page - 1))}
              onNext={() => setPage(Math.min(totalPages, page + 1))}
            />
          </>
        )}
      </div>
    </AdminLayout>
  );
};

export default AdminOrders;
