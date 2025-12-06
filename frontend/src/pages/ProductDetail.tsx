import React, { useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import Navbar from '../components/Navbar';
import BackButton from '../components/BackButton';
import ProductCard from '../components/ProductCard';
import ProductReview from '../components/ProductReview';
import ErrorState from '../components/ErrorState';
import LoadingSpinner from '../components/LoadingSpinner';
import { useProductDetail } from '../hooks/useProductDetail';
import { useProducts } from '../hooks/useProducts';
import { useCart } from '../hooks/useCart';
import { formatCurrency } from '../utils/format';
import { cn } from '../utils/cn';
import { PLACEHOLDER_IMAGE, LOW_STOCK_THRESHOLD } from '../constants/app';
import { useTranslation } from 'react-i18next';

const ProductDetail: React.FC = () => {
  const { t } = useTranslation('common');
  const { slug } = useParams<{ slug: string }>();
  const { product, loading, error } = useProductDetail(slug || '');
  const { products: relatedProducts } = useProducts();
  const { addItem } = useCart();
  const [addingToCart, setAddingToCart] = useState(false);

  const isOutOfStock = useMemo(() => !product || product.stock_quantity === 0, [product]);
  const isLowStock = useMemo(
    () => product && product.stock_quantity > 0 && product.stock_quantity <= LOW_STOCK_THRESHOLD,
    [product]
  );

  const related = useMemo(
    () => relatedProducts.filter((p) => p.id !== product?.id).slice(0, 4),
    [relatedProducts, product?.id]
  );

  const handleAddToCart = async () => {
    if (!product || isOutOfStock) return;
    setAddingToCart(true);
    try {
      await addItem(product.id, 1);
    } finally {
      setAddingToCart(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <LoadingSpinner message={t('common.loading')} />
        </main>
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <ErrorState message={error || t('error.somethingWrong')} />
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="mb-4">
          <BackButton fallbackPath="/products" />
        </div>
        {/* Product Detail Section */}
        <div className="bg-white shadow-sm rounded-xl overflow-hidden mb-12 border border-gray-100">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Left: Image - Full height, half width */}
            <div className="flex items-center justify-center bg-white p-8 min-h-[520px]">
              <img
                src={product.image_url || PLACEHOLDER_IMAGE}
                alt={product.name}
                className="w-full h-full max-h-[480px] object-contain transition-transform duration-300 hover:scale-105"
              />
            </div>

            {/* Right: Product Info */}
            <div className="flex flex-col p-8 lg:p-12 gap-4">
              <div>
                <p className="text-sm text-blue-600 font-semibold uppercase tracking-wide">{product.brand}</p>
              </div>
              <h1 className="text-3xl font-extrabold text-gray-900 mb-2">{product.name}</h1>

              {/* Price */}
              <div className="mb-4">
                <p className="text-3xl font-extrabold text-gray-900">{formatCurrency(product.price)}</p>
              </div>

              {/* Stock Status */}
              <div className="mb-4">
                {isOutOfStock ? (
                  <p className="text-red-600 font-medium">{t('products.outOfStock')}</p>
                ) : isLowStock ? (
                  <p className="text-orange-600 font-medium">{t('products.lowStock', { quantity: product.stock_quantity })}</p>
                ) : (
                  <p className="text-green-600 font-medium">{t('products.inStock')}</p>
                )}
              </div>

              {/* Description */}
              {product.description && (
                <div className="mb-4">
                  <h2 className="text-lg font-semibold text-gray-900 mb-2">{t('products.description')}</h2>
                  <p className="text-gray-700 whitespace-pre-line">{product.description}</p>
                </div>
              )}

              {/* Product Details */}
              <div className="mb-6 space-y-2">
                <h2 className="text-lg font-semibold text-gray-900 mb-2">{t('products.productDetails')}</h2>
                <div className="grid grid-cols-2 gap-3 text-sm">
                  <div className="text-gray-600">{t('products.weight')}</div>
                  <div className="text-gray-900 font-medium">{product.weight} ml</div>
                  <div className="text-gray-600">{t('products.category')}</div>
                  <div className="text-gray-900 font-medium">{product.category}</div>
                </div>
              </div>

              {/* Action Buttons */}
              <div className="mt-auto space-y-3">
                <button
                  onClick={handleAddToCart}
                  disabled={isOutOfStock || addingToCart}
                  aria-disabled={isOutOfStock || addingToCart}
                  className={cn(
                    'w-full px-6 py-3 rounded-lg font-semibold transition-all duration-200 shadow-sm',
                    isOutOfStock
                      ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                      : 'bg-blue-600 text-white hover:bg-blue-700 hover:shadow-md'
                  )}
                >
                  {addingToCart ? t('products.adding') : isOutOfStock ? t('products.outOfStock') : t('products.addToCart')}
                </button>
                <button
                  disabled={isOutOfStock}
                  aria-disabled={isOutOfStock}
                  className={cn(
                    'w-full px-6 py-3 rounded-lg font-semibold transition-all duration-200 border-2',
                    isOutOfStock
                      ? 'bg-white text-gray-400 border-gray-200 cursor-not-allowed'
                      : 'bg-white text-blue-600 border-blue-600 hover:bg-blue-50'
                  )}
                >
                  {t('products.buyNow')}
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Product Review Section */}
        {product && product.slug && (
          <ProductReview
            productSlug={product.slug}
            canReview={product.can_review}
          />
        )}

        {/* Related Products */}
        {related.length > 0 && (
          <section>
            <h2 className="text-2xl font-bold text-gray-900 mb-6">{t('products.relatedProducts')}</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
              {related.map((relatedProduct) => (
                <ProductCard
                  key={relatedProduct.id}
                  product={relatedProduct}
                  onAddToCart={() => addItem(relatedProduct.id, 1)}
                />
              ))}
            </div>
          </section>
        )}
      </main>
    </div>
  );
};

export default ProductDetail;
