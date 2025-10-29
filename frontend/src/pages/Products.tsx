import React from 'react';
import Navbar from '../components/Navbar';
import ProductCard from '../components/ProductCard';
import { useProducts } from '../hooks/useProducts';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorState from '../components/ErrorState';
import type { Product } from '../types/product';
import { useCart } from '../contexts/CartContext';

const Products: React.FC = () => {
  const { products, loading, error, refetch } = useProducts();
  const { addItem } = useCart();

  const handleAddToCart = async (product: Product) => {
    try {
      await addItem(product.id, 1);
    } catch (err) {
      console.error('Failed to add to cart', err);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl md:text-4xl font-bold text-gray-900 mb-2">
            Discover Our Fragrances
          </h1>
          <p className="text-gray-600">
            Explore our curated collection of premium perfumes
          </p>
        </div>

        {/* Loading State */}
        {loading && <LoadingSpinner message="Loading products..." />}

        {/* Error State */}
        {error && (
          <ErrorState message={error} onRetry={refetch} />
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
                  No products available
                </h3>
                <p className="text-gray-500">
                  Check back soon for new fragrances!
                </p>
              </div>
            ) : (
              <>
                <div className="mb-4">
                  <p className="text-gray-600">
                    Showing <span className="font-semibold">{products.length}</span>{' '}
                    {products.length === 1 ? 'product' : 'products'}
                  </p>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                  {products.map((product) => (
                    <ProductCard
                      key={product.id}
                      product={product}
                      onAddToCart={handleAddToCart}
                    />
                  ))}
                </div>
              </>
            )}
          </>
        )}
      </main>
    </div>
  );
};

export default Products;
