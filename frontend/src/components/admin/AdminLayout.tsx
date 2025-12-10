import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useTranslation } from 'react-i18next';
import LanguageSelector from '../LanguageSelector';

type NavItem = {
  label: string;
  to: string;
};

type Props = {
  children: React.ReactNode;
  actions?: React.ReactNode;
  navItems?: NavItem[];
};

const AdminLayout: React.FC<Props> = ({ children, actions, navItems }) => {
  const { t } = useTranslation('admin');
  const location = useLocation();
  const items = navItems ?? [
    { label: t('dashboard'), to: '/admin/dashboard' },
    { label: t('products'), to: '/admin/products' },
    { label: t('orders'), to: '/admin/orders' },
    { label: t('users'), to: '/admin/users' },
    { label: t('auditLogs'), to: '/admin/audit-logs' },
    { label: t('contestations'), to: '/admin/contestations' },
    { label: t('reviewReports'), to: '/admin/review-reports' },
  ];
  const pathname = location.pathname;
  const [isAnimating, setIsAnimating] = React.useState(false);

  // Determine page title based on current route
  const getPageTitle = () => {
    if (pathname.startsWith('/admin/products')) return t('products');
    if (pathname.startsWith('/admin/orders')) return t('orders');
    if (pathname.startsWith('/admin/users')) return t('users');
    if (pathname.startsWith('/admin/audit-logs')) return t('auditLogs');
    if (pathname.startsWith('/admin/contestations')) return t('contestations');
    if (pathname.startsWith('/admin/review-reports')) return t('reviewReports');
    if (pathname.startsWith('/admin/dashboard')) return t('dashboard');
    return t('title'); // fallback to "Admin"
  };

  React.useEffect(() => {
    // Respect users who prefer reduced motion
    const prefersReduced = typeof window !== 'undefined' && window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    if (prefersReduced) return;

    // Trigger a short fade/slide animation whenever pathname changes
    setIsAnimating(true);
    const t = window.setTimeout(() => setIsAnimating(false), 220);
    return () => window.clearTimeout(t);
  }, [pathname]);
  const { user, logout } = useAuth();
  const [mobileOpen, setMobileOpen] = React.useState(false);
  const [userMenuOpen, setUserMenuOpen] = React.useState(false);

  const toggleMobile = () => setMobileOpen((v) => !v);
  const toggleUserMenu = () => setUserMenuOpen((v) => !v);

  return (
    <div className="min-h-screen bg-gray-50">
      <header role="banner" className="bg-white shadow-sm border-b sticky top-0 z-30">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-2">
          <div className="flex justify-between items-center gap-4">
            <div className="flex items-center gap-3 min-w-0">
              <button className="md:hidden p-2 rounded-md border border-gray-200 hover:bg-gray-50" onClick={toggleMobile} aria-label="Toggle menu">
                <svg className="w-6 h-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>

              <Link to="/admin/dashboard" className="flex items-center gap-3 min-w-0">
                <img src="/logo.png" alt="Aroma Sense" className="h-8" />
                <span className="text-lg font-semibold text-gray-900 truncate">{getPageTitle()}</span>
              </Link>
            </div>

            <div className="hidden md:flex items-center gap-3">
              {actions}
              <LanguageSelector />
              <div className="relative">
                <button
                  onClick={toggleUserMenu}
                  className="flex items-center gap-2 px-3 py-1.5 rounded-md border border-gray-200 hover:bg-gray-50"
                  aria-haspopup="true"
                  aria-expanded={userMenuOpen}
                >
                  <span className="w-9 h-9 bg-gray-100 rounded-full flex items-center justify-center text-sm text-gray-700">
                    {user?.email ? user.email.charAt(0).toUpperCase() : 'U'}
                  </span>
                </button>

                {userMenuOpen && (
                  <div className="absolute right-0 mt-2 w-44 bg-white border rounded-md shadow-lg z-20">
                    <Link to="/admin/profile" className="block px-3 py-2 text-sm text-gray-700 hover:bg-gray-50">{t('nav.profile')}</Link>
                    <button onClick={logout} className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50">{t('nav.logout')}</button>
                  </div>
                )}
              </div>
            </div>
          </div>

          <div className="mt-2 hidden md:block">
            <div className="flex items-center gap-2 overflow-x-auto pb-1 no-scrollbar">
              {items.map((it) => {
                const active = pathname === it.to || pathname.startsWith(it.to + '/');
                return (
                  <Link
                    key={it.to}
                    to={it.to}
                    className={`flex-none h-9 px-3 text-sm font-medium rounded-full border transition-all ${active ? 'text-blue-800 bg-blue-50 border-blue-200 shadow-sm' : 'text-gray-700 bg-white border-gray-200 hover:bg-gray-50'}`}
                  >
                    <span className="inline-flex items-center gap-2 justify-center h-full">
                      {active && <span className="w-2 h-2 bg-blue-600 rounded-full" aria-hidden />}
                      <span className="truncate">{it.label}</span>
                    </span>
                  </Link>
                );
              })}
            </div>
          </div>
        </div>

        {/* Mobile menu drawer */}
        {mobileOpen && (
          <div className="md:hidden fixed inset-0 z-40 bg-black bg-opacity-40" onClick={() => setMobileOpen(false)}>
            <div className="absolute left-0 top-0 h-full w-72 max-w-full bg-white shadow-xl p-4 flex flex-col" onClick={(e) => e.stopPropagation()}>
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-3">
                  <img src="/logo.png" alt="Aroma Sense" className="h-8" />
                  <span className="text-lg font-semibold text-gray-900">Admin</span>
                </div>
                <button className="p-2 rounded-md border border-gray-200 hover:bg-gray-50" onClick={() => setMobileOpen(false)} aria-label="Close menu">
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <div className="flex flex-col gap-1">
                {items.map((it) => {
                  const active = pathname === it.to || pathname.startsWith(it.to + '/');
                  return (
                    <Link
                      key={it.to}
                      to={it.to}
                      className={`w-full px-3 py-2 text-sm font-medium rounded-md border ${active ? 'text-blue-800 bg-blue-50 border-blue-200 shadow-sm' : 'text-gray-800 bg-white border-gray-200 hover:bg-gray-50'}`}
                      onClick={() => setMobileOpen(false)}
                    >
                      <span className="inline-flex items-center gap-2">
                        {active && <span className="w-2 h-2 bg-blue-600 rounded-full" aria-hidden />}
                        <span className="truncate">{it.label}</span>
                      </span>
                    </Link>
                  );
                })}
              </div>

              <div className="mt-4 border-t pt-4 flex flex-col gap-2">
                {actions}
                <LanguageSelector />
                <div className="flex items-center justify-between mt-2">
                  <div className="text-sm text-gray-700 truncate">{user?.email}</div>
                  <button onClick={logout} className="text-sm font-medium text-red-600">{t('nav.logout')}</button>
                </div>
              </div>
            </div>
          </div>
        )}
      </header>

      <main role="main" className="max-w-7xl mx-auto p-4">
        <div
          className={`transition-all duration-200 ease-out ${isAnimating ? 'opacity-0 -translate-y-2' : 'opacity-100 translate-y-0'}`}
        >
          {children}
        </div>
      </main>
    </div>
  );
};

export default AdminLayout;
