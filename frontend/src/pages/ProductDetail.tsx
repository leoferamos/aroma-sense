import React, { useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import Navbar from '../components/Navbar';
import ProductCard from '../components/ProductCard';
import ErrorState from '../components/ErrorState';
import LoadingSpinner from '../components/LoadingSpinner';
import { useProductDetail } from '../hooks/useProductDetail';
import { useProducts } from '../hooks/useProducts';
import { useCart } from '../contexts/CartContext';
import { formatCurrency } from '../utils/format';
import { cn } from '../utils/cn';
import { PLACEHOLDER_IMAGE, LOW_STOCK_THRESHOLD } from '../constants/app';

const ProductDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const productId = parseInt(id || '0', 10);
  const { product, loading, error } = useProductDetail(productId);
  const { products: relatedProducts } = useProducts();
  const { addItem } = useCart();
  const [addingToCart, setAddingToCart] = useState(false);

  const isOutOfStock = useMemo(() => !product || product.stock_quantity === 0, [product]);
  const isLowStock = useMemo(
    () => product && product.stock_quantity > 0 && product.stock_quantity <= LOW_STOCK_THRESHOLD,
    [product]
  );

  const related = useMemo(
    () => relatedProducts.filter((p) => p.id !== productId).slice(0, 4),
    [relatedProducts, productId]
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
          <LoadingSpinner message="Loading product..." />
        </main>
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <ErrorState message={error || 'Product not found.'} />
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Product Detail Section */}
        <div className="bg-white shadow rounded-lg overflow-hidden mb-12">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-0">
            {/* Left: Image - Full height, half width */}
            <div className="flex items-center justify-center bg-white p-12 min-h-[600px]">
              <img
                src={product.image_url || PLACEHOLDER_IMAGE}
                alt={product.name}
                className="w-full h-full max-h-[500px] object-contain"
              />
            </div>

            {/* Right: Product Info */}
            <div className="flex flex-col p-8 lg:p-12">
              <div className="mb-2">
                <p className="text-sm text-blue-600 font-medium uppercase tracking-wide">
                  {product.brand}
                </p>
              </div>
              <h1 className="text-3xl font-bold text-gray-900 mb-4">{product.name}</h1>

              {/* Price */}
              <div className="mb-6">
                <p className="text-3xl font-bold text-gray-900">{formatCurrency(product.price)}</p>
              </div>

              {/* Stock Status */}
              <div className="mb-6">
                {isOutOfStock ? (
                  <p className="text-red-600 font-medium">Out of stock</p>
                ) : isLowStock ? (
                  <p className="text-orange-600 font-medium">Only {product.stock_quantity} left in stock</p>
                ) : (
                  <p className="text-green-600 font-medium">In stock</p>
                )}
              </div>

              {/* Description */}
              {product.description && (
                <div className="mb-6">
                  <h2 className="text-lg font-semibold text-gray-900 mb-2">Description</h2>
                  <p className="text-gray-700 whitespace-pre-line">{product.description}</p>
                </div>
              )}

              {/* Product Details */}
              <div className="mb-6 space-y-2">
                <h2 className="text-lg font-semibold text-gray-900 mb-2">Product Details</h2>
                <div className="grid grid-cols-2 gap-2 text-sm">
                  <div className="text-gray-600">Weight:</div>
                  <div className="text-gray-900 font-medium">{product.weight} ml</div>
                  <div className="text-gray-600">Category:</div>
                  <div className="text-gray-900 font-medium">{product.category}</div>
                  {product.notes && (
                    <>
                      <div className="text-gray-600">Notes:</div>
                      <div className="text-gray-900 font-medium">{product.notes}</div>
                    </>
                  )}
                </div>
              </div>

              {/* Action Buttons */}
              <div className="mt-auto space-y-3">
                <button
                  onClick={handleAddToCart}
                  disabled={isOutOfStock || addingToCart}
                  aria-disabled={isOutOfStock || addingToCart}
                  className={cn(
                    'w-full px-6 py-3 rounded-md font-medium transition-colors',
                    isOutOfStock
                      ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                      : 'bg-blue-600 text-white hover:bg-blue-700'
                  )}
                >
                  {addingToCart ? 'Adding...' : isOutOfStock ? 'Out of Stock' : 'Add to Cart'}
                </button>
                <button
                  disabled={isOutOfStock}
                  aria-disabled={isOutOfStock}
                  className={cn(
                    'w-full px-6 py-3 rounded-md font-medium transition-colors border-2',
                    isOutOfStock
                      ? 'bg-gray-100 text-gray-400 border-gray-300 cursor-not-allowed'
                      : 'bg-white text-blue-600 border-blue-600 hover:bg-blue-50'
                  )}
                >
                  Buy Now
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Related Products */}
        {related.length > 0 && (
          <section>
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Related Products</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
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
