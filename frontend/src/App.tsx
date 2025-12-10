import React, { lazy, Suspense } from 'react';
import { BrowserRouter, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import GuestRoute from './components/GuestRoute';
import { AuthProvider } from './contexts/AuthContext';
import { CartProvider } from './contexts/CartContext';
import ProtectedRoute from './components/ProtectedRoute';
import ErrorBoundary from './components/ErrorBoundary';
import ChatBubble from './components/chat/ChatBubble';
import AccountBlockOverlay from './components/AccountBlockOverlay';
import { useTranslation } from 'react-i18next';
// Lazy load pages for better performance
const AdminContestations = lazy(() => import('./pages/admin/AdminContestations'));
const AdminReviewReports = lazy(() => import('./pages/admin/AdminReviewReports'));
const Register = lazy(() => import('./pages/Register'));
const Login = lazy(() => import('./pages/Login'));
const Terms = lazy(() => import('./pages/Terms'));
const Privacy = lazy(() => import('./pages/Privacy'));
const ForgotPassword = lazy(() => import('./pages/ForgotPassword'));
const AdminDashboard = lazy(() => import('./pages/admin/AdminDashboard'));
const AddProduct = lazy(() => import('./pages/admin/AddProduct'));
const EditProduct = lazy(() => import('./pages/admin/EditProduct'));
const AdminProducts = lazy(() => import('./pages/admin/AdminProducts'));
const AdminOrders = lazy(() => import('./pages/admin/AdminOrders'));
const AdminUsers = lazy(() => import('./pages/admin/Users'));
const AdminAuditLogs = lazy(() => import('./pages/admin/AdminAuditLogs'));
const Orders = lazy(() => import('./pages/Orders'));
const Products = lazy(() => import('./pages/Products'));
const ProductDetail = lazy(() => import('./pages/ProductDetail.tsx'));
const Checkout = lazy(() => import('./pages/Checkout.tsx'));
const OrderConfirmation = lazy(() => import('./pages/OrderConfirmation.tsx'));
const NotFound = lazy(() => import('./pages/NotFound'));
const Profile = lazy(() => import('./pages/Profile'));
const ChangePassword = lazy(() => import('./pages/ChangePassword'));

// Loading fallback component
const PageLoader: React.FC = () => {
  const { t } = useTranslation('common');
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <p className="mt-4 text-gray-600">{t('errors.loading')}</p>
      </div>
    </div>
  );
};

const ChatMount: React.FC = () => {
  const location = useLocation();
  const path = location.pathname;
  const isAuthPage = path.startsWith('/login') || path.startsWith('/register') || path.startsWith('/forgot-password');
  const isAdmin = path.startsWith('/admin');
  const showChat = !isAuthPage && !isAdmin;
  return (
    <>
      <Routes>
        {/* Public Routes */}
        <Route
          path="/register"
          element={<GuestRoute><Register /></GuestRoute>}
        />
        <Route path="/login" element={<GuestRoute><Login /></GuestRoute>} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route
          path="/admin/users"
          element={
            <ProtectedRoute allowedRoles={["admin", "super_admin"]}>
              <AdminUsers />
            </ProtectedRoute>
          }
        />
        <Route path="/terms" element={<Terms />} />
        <Route path="/privacy" element={<Privacy />} />

        {/* Protected admin routes */}
        <Route
          path="/admin/dashboard"
          element={
            <ProtectedRoute allowedRoles={['admin', 'super_admin']}>
              <AdminDashboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/products/new"
          element={
            <ProtectedRoute allowedRoles={['admin', 'super_admin']}>
              <AddProduct />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/products/:id/edit"
          element={
            <ProtectedRoute allowedRoles={['admin', 'super_admin']}>
              <EditProduct />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/products"
          element={
            <ProtectedRoute allowedRoles={['admin', 'super_admin']}>
              <AdminProducts />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/orders"
          element={
            <ProtectedRoute allowedRoles={['admin', 'super_admin']}>
              <AdminOrders />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/audit-logs"
          element={
            <ProtectedRoute allowedRoles={['admin', 'super_admin']}>
              <AdminAuditLogs />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/contestations"
          element={
            <ProtectedRoute allowedRoles={["admin", "super_admin"]}>
              <AdminContestations />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin/review-reports"
          element={
            <ProtectedRoute allowedRoles={['admin', 'super_admin']}>
              <AdminReviewReports />
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
        <Route
          path="/orders"
          element={
            <ProtectedRoute>
              <Orders />
            </ProtectedRoute>
          }
        />
        <Route
          path="/profile"
          element={
            <ProtectedRoute>
              <Profile />
            </ProtectedRoute>
          }
        />
        <Route
          path="/change-password"
          element={
            <ProtectedRoute>
              <ChangePassword />
            </ProtectedRoute>
          }
        />
        <Route path="/products/:slug" element={<ProductDetail />} />
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
      {showChat && <ChatBubble />}
    </>
  );
}

const App: React.FC = () => {
  return (
    <BrowserRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
      <AuthProvider>
        <CartProvider>
          <Suspense fallback={<PageLoader />}>
            <ErrorBoundary>
              <ChatMount />
              {/* Global overlay shown when account is blocked */}
              <AccountBlockOverlay />
            </ErrorBoundary>
          </Suspense>
        </CartProvider>
      </AuthProvider>
    </BrowserRouter>
  );
};

export default App;
