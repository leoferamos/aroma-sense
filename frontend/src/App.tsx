import React, { lazy, Suspense } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { CartProvider } from './contexts/CartContext';
import ProtectedRoute from './components/ProtectedRoute';
import ErrorBoundary from './components/ErrorBoundary';

// Lazy load pages for better performance
const Register = lazy(() => import('./pages/Register'));
const Login = lazy(() => import('./pages/Login'));
const Terms = lazy(() => import('./pages/Terms'));
const Privacy = lazy(() => import('./pages/Privacy'));
const AdminDashboard = lazy(() => import('./pages/admin/AdminDashboard'));
const AddProduct = lazy(() => import('./pages/admin/AddProduct'));
const Products = lazy(() => import('./pages/Products'));
const ProductDetail = lazy(() => import('./pages/ProductDetail.tsx'));
const Checkout = lazy(() => import('./pages/Checkout.tsx'));
const OrderConfirmation = lazy(() => import('./pages/OrderConfirmation.tsx'));
const NotFound = lazy(() => import('./pages/NotFound'));

// Loading fallback component
const PageLoader: React.FC = () => (
  <div className="min-h-screen flex items-center justify-center bg-gray-50">
    <div className="text-center">
      <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      <p className="mt-4 text-gray-600">Loading...</p>
    </div>
  </div>
);

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <AuthProvider>
        <CartProvider>
          <Suspense fallback={<PageLoader />}>
            <ErrorBoundary>
            <Routes>
              {/* Public routes */}
              <Route path="/register" element={<Register />} />
              <Route path="/login" element={<Login />} />
              <Route path="/terms" element={<Terms />} />
              <Route path="/privacy" element={<Privacy />} />
              
              {/* Protected admin routes */}
              <Route
                path="/admin/dashboard"
                element={
                  <ProtectedRoute allowedRoles={['admin']}>
                    <AdminDashboard />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/admin/products/new"
                element={
                  <ProtectedRoute allowedRoles={['admin']}>
                    <AddProduct />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/products"
                element={
                  <ProtectedRoute>
                    <Products />
                  </ProtectedRoute>
                }
              />
              <Route path="/products/:id" element={<ProductDetail />} />
              <Route
                path="/checkout"
                element={
                  <ProtectedRoute>
                    <Checkout />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/order-confirmation"
                element={
                  <ProtectedRoute>
                    <OrderConfirmation />
                  </ProtectedRoute>
                }
              />
              <Route path="/" element={<Navigate to="/login" replace />} />
              
              {/* 404 - Catch all unmatched routes */}
              <Route path="*" element={<NotFound />} />
            </Routes>
            </ErrorBoundary>
          </Suspense>
        </CartProvider>
      </AuthProvider>
    </BrowserRouter>
  );
};

export default App;
