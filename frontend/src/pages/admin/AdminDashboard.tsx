import React from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { Link } from 'react-router-dom';

const AdminDashboard: React.FC = () => {
  const { role, logout } = useAuth();

  const handleLogout = async () => {
    await logout();
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <img src="/logo.png" alt="Aroma Sense" className="h-10" />
              <h1 className="ml-4 text-xl font-semibold text-gray-900">Admin Dashboard</h1>
            </div>
            <div className="flex items-center gap-4">
              <span className="text-sm text-gray-600">
                Role: <span className="font-semibold text-blue-600">{role}</span>
              </span>
              <button
                onClick={handleLogout}
                className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 transition-colors"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-2xl font-bold text-gray-900 mb-4">
              Welcome, Admin
            </h2>
            <p className="text-gray-600 mb-4">
              Role: <span className="font-semibold">{role}</span>
            </p>
            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
              <div className="bg-blue-50 p-6 rounded-lg border border-blue-200">
                <h3 className="text-lg font-semibold text-blue-900 mb-2">Products</h3>
                <p className="text-blue-700">Manage your product catalog</p>
                <Link
                  to="/admin/products"
                  className="mt-4 inline-block text-blue-600 hover:text-blue-800 font-medium"
                >
                  View Products →
                </Link>
              </div>

              <div className="bg-green-50 p-6 rounded-lg border border-green-200">
                <h3 className="text-lg font-semibold text-green-900 mb-2">Orders</h3>
                <p className="text-green-700">Track and manage orders</p>
                <Link
                  to="/admin/orders"
                  className="mt-4 inline-block text-green-600 hover:text-green-800 font-medium"
                >
                  View Orders →
                </Link>
              </div>

              <div className="bg-purple-50 p-6 rounded-lg border border-purple-200">
                <h3 className="text-lg font-semibold text-purple-900 mb-2">Users</h3>
                <p className="text-purple-700">Manage user accounts</p>
                <Link
                  to="/admin/users"
                  className="mt-4 inline-block text-purple-600 hover:text-purple-800 font-medium"
                >
                  View Users →
                </Link>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default AdminDashboard;
