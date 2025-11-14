import React, { useEffect, useState } from 'react';
import Navbar from '../components/Navbar';
import ProductCard from '../components/ProductCard';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorState from '../components/ErrorState';
import type { Product } from '../types/product';
import { useCart } from '../contexts/CartContext';
import Pagination from '../components/Pagination';
import { useProductSearch } from '../hooks/useProductSearch';
import { listProducts } from '../services/product';
import ProductCardSkeleton from '../components/ProductCardSkeleton';
import { useSearchParams } from 'react-router-dom';
import ChatBubble from '../components/chat/ChatBubble';

const Products: React.FC = () => {
  const [suggestions, setSuggestions] = useState<Product[] | null>(null);
  const [loadingSuggestions, setLoadingSuggestions] = useState(false);
  const { query, setQuery, page, setPage, limit, results, total, isLoading, error, submitNow, isSearching } = useProductSearch({ limit: 12, debounceMs: 600 });
  const [searchParams] = useSearchParams();
  const { addItem } = useCart();

  const handleAddToCart = async (product: Product) => {
    try {
      await addItem(product.id, 1);
    } catch (err) {
      console.error('Failed to add to cart', err);
    }
  };

  // Sync state from URL
  useEffect(() => {
    const q = searchParams.get('q') ?? '';
    const p = Math.max(1, parseInt(searchParams.get('page') ?? '1', 10) || 1);
    if (q !== query) { setQuery(q); }
    if (p !== page) { setPage(p); }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams]);

  // Load suggestions when user searched (>=2 chars) and got 0 results
  useEffect(() => {
    const resultsLen = Array.isArray(results) ? results.length : 0;
    const needSuggestions = query.trim().length >= 2 && !isLoading && resultsLen === 0;
    if (!needSuggestions) {
      setSuggestions(null);
      return;
    }
    let active = true;
    setLoadingSuggestions(true);
    listProducts({ limit: 8 })
      .then((data) => {
        if (!active) return;
        setSuggestions(data);
      })
      .catch(() => {
        if (!active) return;
        setSuggestions([]);
      })
      .finally(() => {
        if (!active) return;
        setLoadingSuggestions(false);
      });
    return () => { active = false; };
  }, [query, isLoading, Array.isArray(results) ? results.length : 0]);

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-16 pb-10">
        {(() => {
          return null;
        })()}
        {/* Header */}
        <div className="mb-10">
          <div className="bg-gradient-to-b from-white to-gray-50 rounded-xl p-6 shadow-sm border border-gray-50">
            <h1 className="text-4xl font-bold tracking-tight text-gray-900 mb-2">Products</h1>
            <p className="text-lg text-gray-500">Find your favorite fragrance</p>
          </div>
        </div>

        {/* Loading State */}
        {isLoading && (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8" aria-live="polite">
            {Array.from({ length: 8 }).map((_, i) => (
              <ProductCardSkeleton key={i} />
            ))}
          </div>
        )}

        {/* Error State */}
        {error && <ErrorState message={error} onRetry={submitNow} />}

        {/* Products Grid */}
        {!isLoading && !error && (
          <>
            {(() => {
              const safeResults = Array.isArray(results)
                ? results.filter((p) => p && typeof (p as any).id === 'number')
                : [];
              const showEmpty = safeResults.length === 0;
              return showEmpty;
            })() ? (
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
                {query.trim().length >= 2 ? (
                  <>
                    <h3 className="text-xl font-semibold text-gray-700 mb-2">We couldn’t find “{query.trim()}”.</h3>
                    <p className="text-gray-500">Here are some suggestions you might like:</p>
                    <div className="mt-8">
                      {loadingSuggestions ? (
                        <LoadingSpinner message="Loading suggestions..." />
                      ) : (
                        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
                          {((Array.isArray(suggestions) ? suggestions : [])
                            .filter((p) => p && typeof (p as any).id === 'number')
                          ).map((p) => (

                            <ProductCard key={(p as any).id} product={p as any} onAddToCart={handleAddToCart} showAddToCart={false} />
                          ))}
                        </div>
                      )}
                    </div>
                  </>
                ) : (
                  <>
                    <h3 className="text-xl font-semibold text-gray-700 mb-2">No products available</h3>
                    <p className="text-gray-500">Check back soon for new fragrances!</p>
                  </>
                )}
              </div>
            ) : (
              <>
                <div className="mb-4" aria-live="polite">
                  <p className="text-gray-600">
                    {isSearching && query.trim().length >= 2 ? (
                      <>
                        Results for <span className="font-semibold">“{query.trim()}”</span>: <span className="font-semibold">{total}</span>
                      </>
                    ) : (
                      (() => {
                        const safeResults = Array.isArray(results)
                          ? results.filter((p) => p && typeof (p as any).id === 'number')
                          : [];
                        return <>
                          Showing <span className="font-semibold">{safeResults.length}</span> {safeResults.length === 1 ? 'product' : 'products'}
                        </>;
                      })()
                    )}
                  </p>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
                  {(Array.isArray(results)
                    ? results.filter((p) => p && typeof (p as any).id === 'number')
                    : []
                  ).map((product: any) => (
                    <ProductCard
                      key={product.id}
                      product={product}
                      onAddToCart={handleAddToCart}
                      showAddToCart={false}
                    />
                  ))}
                </div>

                {isSearching && Number(total) > Number(limit) && (
                  <Pagination
                    page={Number(page) || 1}
                    pageSize={Number(limit) || 12}
                    total={Number(total) || 0}
                    onPageChange={setPage}
                  />
                )}
              </>
            )}
          </>
        )}
      </main>
      <ChatBubble />
    </div>
  );
};

export default Products;
