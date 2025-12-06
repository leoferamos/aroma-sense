import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useProducts } from '../../hooks/useProducts';
import LoadingSpinner from '../../components/LoadingSpinner';
import ErrorState from '../../components/ErrorState';
import AdminProductCard from '../../components/admin/AdminProductCard';
import type { Product } from '../../types/product';
import ConfirmModal from '../../components/ConfirmModal';
import { deleteProduct } from '../../services/product';
import AdminLayout from '../../components/admin/AdminLayout';
import { useTranslation } from 'react-i18next';


const AdminProducts: React.FC = () => {
  const { products, loading, error, refetch } = useProducts();
  const navigate = useNavigate();
  const { t } = useTranslation('admin');
  const { t } = useTranslation('admin');
  const [modalOpen, setModalOpen] = React.useState(false);
  const [selectedProduct, setSelectedProduct] = React.useState<Product | null>(null);
  const [deleting, setDeleting] = React.useState(false);
  const [successMsg, setSuccessMsg] = React.useState<string | null>(null);
  const [fadeOut, setFadeOut] = React.useState(false);
  const [deleteError, setDeleteError] = React.useState<string | null>(null);

  const handleEdit = (product: Product) => {
    navigate(`/admin/products/${product.id}/edit`);
  };

  const handleDelete = (product: Product) => {
    setSelectedProduct(product);
    setModalOpen(true);
    setDeleteError(null);
  };

  const handleConfirmDelete = async () => {
    if (!selectedProduct) return;
    setDeleting(true);
    setDeleteError(null);
    try {
      await deleteProduct(selectedProduct.id);
      setSuccessMsg('Product deleted successfully.');
      setFadeOut(false);
      setModalOpen(false);
      setSelectedProduct(null);
      await refetch();
    } catch (err: unknown) {
      const e = err as { response?: { data?: { error?: string } }; message?: string };
      setDeleteError(e?.response?.data?.error || e?.message || 'Failed to delete product.');
    } finally {
      setDeleting(false);
    }
  };

  const handleCloseModal = () => {
    setModalOpen(false);
    setSelectedProduct(null);
    setDeleteError(null);
  };

  React.useEffect(() => {
    if (!successMsg) return;
    setFadeOut(false);
    const fadeTimer = setTimeout(() => setFadeOut(true), 2300);
    const removeTimer = setTimeout(() => setSuccessMsg(null), 3000);
    return () => {
      clearTimeout(fadeTimer);
      clearTimeout(removeTimer);
    };
  }, [successMsg]);

  const actions = (
    <div className="flex items-center gap-2">
      <Link
        to="/admin/dashboard"
        className="px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-lg hover:bg-gray-50"
      >
        ← {t('nav.dashboard')}
      </Link>
      <Link
        to="/admin/products/new"
        className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700"
      >
        + {t('addProduct')}
      </Link>
    </div>
  );

  return (
    <AdminLayout actions={actions}>
      {/* Loading State */}
      {loading && <LoadingSpinner message="Loading products..." />}

      {/* Error State */}
      {error && (
        <ErrorState message={error} onRetry={refetch} />
      )}

      {/* Success message */}
      {successMsg && (
        <div
          className={`mb-4 rounded-md border border-green-200 bg-green-50 px-3 py-2 text-green-700 text-center font-medium transition-opacity duration-700 ${fadeOut ? 'opacity-0' : 'opacity-100'}`}
          aria-live="polite"
        >
          {successMsg}
        </div>
      )}
      {/* Delete error */}
      {deleteError && (
        <div className="mb-4 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-red-700 text-center font-medium">
          {deleteError}
          <button className="ml-4 text-red-800 underline" onClick={() => setDeleteError(null)}>
            Dismiss
          </button>
        </div>
      )}
      {/* Products Grid */}
      {!loading && !error && Array.isArray(products) && (
        <>
          {products.length === 0 ? (
            <div className="text-center py-20">
              <svg
                className="w-16 h-16 text-gray-400 mx-auto mb-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
                />
              </svg>
              <h3 className="text-xl font-semibold text-gray-700 mb-2">
                {t('noProductsYet')}
              </h3>
              <p className="text-gray-500">
                {t('clickAddNewProduct')}
              </p>
            </div>
          ) : (
            <>
              <div className="mb-4">
                <p className="text-gray-600">
                  {t('showingProducts', { count: products.length })}
                </p>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                {products.map((product) => (
                  <AdminProductCard
                    key={product.id}
                    product={product}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                  />
                ))}
              </div>
            </>
          )}
        </>
      )}

      {/* Confirm Delete Modal */}
      <ConfirmModal
        open={modalOpen}
        title="Delete Product"
        description={selectedProduct ? `Are you sure you want to delete "${selectedProduct.name}"? This action cannot be undone.` : ''}
        confirmText="Delete"
        cancelText="Cancel"
        onConfirm={handleConfirmDelete}
        onCancel={handleCloseModal}
        loading={deleting}
      />
    </AdminLayout>
  );
};

export default AdminProducts;
