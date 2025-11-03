import React, { useEffect, useState } from 'react';
import { getUserOrders } from '../services/order';
import type { OrderResponse } from '../types/order';
import OrderList from '../components/OrderList';
import OrderDetail from '../components/OrderDetail';

const OrdersPage: React.FC = () => {
  const [orders, setOrders] = useState<OrderResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedOrderId, setSelectedOrderId] = useState<number | null>(null);

  useEffect(() => {
    let mounted = true;
    (async () => {
      try {
        setLoading(true);
        const data = await getUserOrders();
        if (!mounted) return;
        // sort newest first
        data.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
        setOrders(data);
      } catch {
        setError('Failed to load orders');
      } finally {
        setLoading(false);
      }
    })();
    return () => { mounted = false; };
  }, []);

  const selectedOrder = orders.find((o) => o.id === selectedOrderId) ?? null;

  return (
    <div className="max-w-5xl mx-auto py-8">
      <h1 className="text-2xl font-semibold mb-6">My Orders</h1>
      {loading && <div className="text-gray-600">Loading...</div>}
      {error && <div className="text-red-500">{error}</div>}

      {!loading && !error && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-1">
            <OrderList orders={orders} onSelect={(id) => setSelectedOrderId((prev) => (prev === id ? null : id))} />
          </div>
          <div className="lg:col-span-2">
            {selectedOrder ? (
              <OrderDetail order={selectedOrder} />
            ) : (
              <div className="p-6 bg-white rounded shadow-sm">Select an order to view details</div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default OrdersPage;
