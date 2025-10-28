import React, { useEffect, useRef, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useCart } from '../contexts/CartContext';
import { formatCurrency } from '../utils/format';
import { PLACEHOLDER_IMAGE, MAX_CART_BADGE_COUNT, LOGO_PATH } from '../constants/app';

const Navbar: React.FC = () => {
  const { itemCount, cart } = useCart();
  const [open, setOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const onClickOutside = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    if (open) document.addEventListener('mousedown', onClickOutside);
    return () => document.removeEventListener('mousedown', onClickOutside);
  }, [open]);

  const badgeCount = itemCount > MAX_CART_BADGE_COUNT ? `${MAX_CART_BADGE_COUNT}+` : itemCount;

  return (
    <nav className="bg-white shadow-md sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link to="/products" className="flex items-center">
            <img src={LOGO_PATH} alt="Aroma Sense" className="h-10 w-auto" />
          </Link>

          {/* Right side: Hamburger Menu + Cart */}
          <div className="flex items-center gap-4 relative" ref={dropdownRef}>
            {/* Cart Icon */}
            <button
              type="button"
              aria-label="Open cart"
              className="relative p-2 text-gray-700 hover:text-blue-600 transition-colors"
              onClick={() => setOpen((v) => !v)}
            >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={1.5}
                  stroke="currentColor"
                  className="w-6 h-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M2.25 3h1.386c.51 0 .955.343 1.087.835l.383 1.437M7.5 14.25a3 3 0 00-3 3h15.75m-12.75-3h11.218c1.121-2.3 2.1-4.684 2.924-7.138a60.114 60.114 0 00-16.536-1.84M7.5 14.25L5.106 5.272M6 20.25a.75.75 0 11-1.5 0 .75.75 0 011.5 0zm12.75 0a.75.75 0 11-1.5 0 .75.75 0 011.5 0z"
                  />
                </svg>
                {itemCount > 0 && (
                  <span className="absolute -top-1 -right-1 bg-blue-600 text-white text-xs font-bold rounded-full h-5 w-5 flex items-center justify-center">
                  {badgeCount}
                </span>
              )}
            </button>

            {/* Dropdown */}
            {open && (
              <div className="absolute right-0 top-12 w-80 bg-white rounded-lg shadow-xl border border-gray-100 overflow-hidden">
                <div className="p-4 border-b border-gray-100">
                  <h3 className="text-sm font-semibold text-gray-900">Cart</h3>
                </div>
                <div className="max-h-80 overflow-auto divide-y">
                  {cart && cart.items.length > 0 ? (
                    cart.items.map((item) => (
                      <div key={item.id} className="p-4 flex gap-3 items-center">
                        <img
                          src={item.product?.image_url || PLACEHOLDER_IMAGE}
                          alt={item.product?.name || 'Product'}
                          className="h-12 w-12 object-contain bg-gray-50 rounded"
                          onError={(e) => { 
                            (e.currentTarget as HTMLImageElement).src = PLACEHOLDER_IMAGE;
                          }}
                        />
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium text-gray-900 truncate">
                            {item.product?.name || 'Product'}
                          </p>
                          <p className="text-xs text-gray-500">Qty: {item.quantity}</p>
                        </div>
                        <div className="text-sm font-semibold text-gray-900 whitespace-nowrap">
                          {formatCurrency(item.total)}
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="p-6 text-center text-gray-500 text-sm">Your cart is empty</div>
                  )}
                </div>
                <div className="p-4 border-t border-gray-100">
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-sm text-gray-600">Subtotal</span>
                    <span className="text-base font-semibold text-gray-900">
                      {formatCurrency(cart?.total || 0)}
                    </span>
                  </div>
                  <button
                    type="button"
                    onClick={() => { setOpen(false); navigate('/checkout'); }}
                    className="w-full py-2.5 px-4 rounded-lg font-medium bg-blue-600 text-white hover:bg-blue-700 transition-colors"
                  >
                    Go to Checkout
                  </button>
                </div>
              </div>
            )}

            {/* Hamburger Menu */}
            <button
              type="button"
              aria-label="Open menu"
              className="p-2 text-gray-700 hover:text-blue-600 transition-colors"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-6 h-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
