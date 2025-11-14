import React, { useEffect, useRef, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useCart } from '../contexts/CartContext';
import { useAuth } from '../contexts/AuthContext';
import { formatCurrency } from '../utils/format';
import { MAX_CART_BADGE_COUNT, LOGO_PATH } from '../constants/app';
import CartItem from './CartItem';
import UserMenu from './UserMenu';
import SearchBar from './SearchBar';

const Navbar: React.FC = () => {
  const { itemCount, cart, removeItem, isRemovingItem } = useCart();
  const { role, isAuthenticated, logout } = useAuth();
  const [open, setOpen] = useState(false); // cart dropdown
  const [userMenuOpen, setUserMenuOpen] = useState(false); // hamburger/user menu
  const dropdownRef = useRef<HTMLDivElement | null>(null);
  const navigate = useNavigate();
  const location = useLocation();
  const [navQuery, setNavQuery] = useState<string>("");

  // Debounced URL navigation for live search (2+ chars). Empty clears query.
  const debounceRef = useRef<number | null>(null);
  useEffect(() => {
    if (debounceRef.current) {
      window.clearTimeout(debounceRef.current);
      debounceRef.current = null;
    }

    const q = navQuery.trim();


    // Require min length to avoid noisy requests
    if (q.length < 2) return;

    debounceRef.current = window.setTimeout(() => {
      const target = `/products?q=${encodeURIComponent(q)}&page=1`;
      const current = `${location.pathname}${location.search}`;
      if (current !== target) {
        navigate(target);
      }
    }, 600);

    return () => {
      if (debounceRef.current) {
        window.clearTimeout(debounceRef.current);
        debounceRef.current = null;
      }
    };
  }, [navQuery, location.pathname, location.search, navigate]);

  useEffect(() => {
    const sp = new URLSearchParams(location.search);
    setNavQuery(sp.get('q') ?? '');
  }, [location.search]);

  useEffect(() => {
    const onClickOutside = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setOpen(false);
        setUserMenuOpen(false);
      }
    };
    if (open || userMenuOpen) {
      document.addEventListener('mousedown', onClickOutside);
      return () => document.removeEventListener('mousedown', onClickOutside);
    }
  }, [open, userMenuOpen]);

  const badgeCount = itemCount > MAX_CART_BADGE_COUNT ? `${MAX_CART_BADGE_COUNT}+` : itemCount;

  return (
    <nav className="bg-white shadow-sm sticky top-0 z-50 border-b border-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16 gap-4">
          {/* Logo */}
          <Link to="/products" className="flex items-center flex-shrink-0">
            <img src={LOGO_PATH} alt="Aroma Sense" className="h-10 w-auto" />
          </Link>

          {/* Search (center) */}
          <div className="hidden md:block flex-1 max-w-2xl">
            <SearchBar
              value={navQuery}
              onChange={setNavQuery}
              onSubmit={() => {
                const q = navQuery.trim();
                if (q.length >= 2) {
                  navigate(`/products?q=${encodeURIComponent(q)}&page=1`);
                } else {
                  navigate('/products');
                }
              }}
              onClear={() => {
                setNavQuery('');
                navigate('/products');
              }}
            />
          </div>

          {/* Right side: User Menu + Cart */}
          <div className="flex items-center gap-2 relative" ref={dropdownRef}>
            {/* Cart Icon */}
            <button
              type="button"
              aria-label="Open cart"
              className="relative p-2 text-gray-700 hover:text-blue-600 hover:bg-gray-100 transition-all duration-200 rounded-lg"
              onClick={() => {
                setOpen((v) => !v);
                setUserMenuOpen(false); // ensure menu closes when opening cart
              }}
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
              <div className="absolute right-0 top-12 w-96 bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden animate-in fade-in slide-in-from-top-2 duration-200">
                <div className="p-4 border-b border-gray-100">
                  <h3 className="text-sm font-semibold text-gray-900">Shopping Cart</h3>
                </div>
                <div className="max-h-80 overflow-auto divide-y divide-gray-100">
                  {cart && cart.items.length > 0 ? (
                    cart.items.map((item) => (
                      <CartItem
                        key={item.id}
                        item={item}
                        onRemove={removeItem}
                        isRemoving={isRemovingItem(item.id)}
                        compact
                        showQuantityControls
                      />
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
                    className="w-full py-2.5 px-4 rounded-lg font-medium bg-blue-600 text-white hover:bg-blue-700 transition-all duration-200"
                  >
                    Go to Checkout
                  </button>
                </div>
              </div>
            )}

            {/* Hamburger Menu*/}
            <button
              type="button"
              aria-label="Open menu"
              className="p-2 text-gray-700 hover:text-blue-600 hover:bg-gray-100 transition-all duration-200 rounded-lg"
              onClick={() => {
                setUserMenuOpen((v) => !v);
                setOpen(false); // close cart if open
              }}
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

            <UserMenu
              open={userMenuOpen}
              role={role}
              isAuthenticated={isAuthenticated}
              onClose={() => setUserMenuOpen(false)}
              onLogout={() => logout()}
              onSignIn={() => navigate('/login')}
            />
          </div>
        </div>
        {/* Mobile search */}
        <div className="block md:hidden py-2">
          <SearchBar
            value={navQuery}
            onChange={setNavQuery}
            onSubmit={() => {
              const q = navQuery.trim();
              if (q.length >= 2) {
                navigate(`/products?q=${encodeURIComponent(q)}&page=1`);
              } else {
                navigate('/products');
              }
            }}
            onClear={() => {
              setNavQuery('');
              navigate('/products');
            }}
          />
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
