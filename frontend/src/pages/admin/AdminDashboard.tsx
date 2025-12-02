import React from 'react';
import { useAuth } from '../../hooks/useAuth';
import { Link } from 'react-router-dom';
import AdminLayout from '../../components/admin/AdminLayout';

const AdminDashboard: React.FC = () => {
  const { role, logout } = useAuth();

  const handleLogout = async () => {
    await logout();
  };

  const actions = (
    <div className="flex items-center gap-3">
      <span className="text-sm text-gray-500">
        <span className="font-medium text-blue-600">{role}</span>
      </span>
      <button
        onClick={handleLogout}
        className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"
      >
        Logout
      </button>
    </div>
  );

  return (
    <AdminLayout title="Dashboard" actions={actions}>
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-900">Welcome, Admin</h2>
        <p className="text-gray-500 text-sm mt-1">Manage your store</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {/* Products Card */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5 hover:shadow-md transition-shadow">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-lg font-semibold text-gray-900">Products</h3>
            <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
              <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
          </div>
          <p className="text-gray-600 text-sm mb-4">Manage your product catalog</p>
          <div className="flex flex-col gap-2">
            <Link
              to="/admin/products"
              className="text-blue-600 hover:text-blue-700 font-medium text-center text-sm"
            >
              View All Products →
            </Link>
          </div>
        </div>

        {/* Orders Card */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5 hover:shadow-md transition-shadow">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-lg font-semibold text-gray-900">Orders</h3>
            <div className="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
              <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
          </div>
          <p className="text-gray-600 text-sm mb-4">Track and manage orders</p>
          <Link
            to="/admin/orders"
            className="text-green-600 hover:text-green-700 font-medium text-sm inline-block"
          >
            View Orders →
          </Link>
        </div>

        {/* Users Card */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5 hover:shadow-md transition-shadow">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-lg font-semibold text-gray-900">Users</h3>
            <div className="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
              <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </div>
          </div>
          <p className="text-gray-600 text-sm mb-4">Manage user accounts</p>
          <Link
            to="/admin/users"
            className="text-purple-600 hover:text-purple-700 font-medium text-sm inline-block"
          >
            View Users →
          </Link>
        </div>
      </div>
    </AdminLayout>
  );
};

export default AdminDashboard;
