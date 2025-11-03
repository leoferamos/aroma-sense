import React, { useEffect, useState } from 'react';
import Navbar from '../components/Navbar';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorState from '../components/ErrorState';
import { getUserOrders } from '../services/order';
import type { OrderResponse } from '../types/order';
import OrderDetail from '../components/OrderDetail';

const OrdersPage: React.FC = () => {
  const [orders, setOrders] = useState<OrderResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

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

  // render all orders as full cards (single-column)

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <main className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h1 className="text-2xl font-semibold mb-6">My Orders</h1>

        {loading && <LoadingSpinner message="Loading orders..." />}
        {error && <ErrorState message={error} onRetry={() => window.location.reload()} />}

        {!loading && !error && (
          <div className="space-y-6">
            {orders.length === 0 ? (
              <div className="p-6 bg-white rounded shadow-sm text-center">
                <h2 className="text-lg font-semibold">You have no orders yet</h2>
                <p className="text-sm text-gray-500 mt-2">Browse products and place your first order.</p>
              </div>
            ) : (
              orders.map((o) => (
                <OrderDetail key={o.id} order={o} />
              ))
            )}
          </div>
        )}
      </main>
    </div>
  );
};

export default OrdersPage;
