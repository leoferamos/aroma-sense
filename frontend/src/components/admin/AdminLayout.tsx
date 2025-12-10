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
      <header role="banner" className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center gap-4">
              <button className="md:hidden p-2" onClick={toggleMobile} aria-label="Toggle menu">
                <svg className="w-6 h-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>

              <Link to="/admin/dashboard" className="flex items-center gap-3">
                <img src="/logo.png" alt="Aroma Sense" className="h-8" />
                <span className="text-lg font-semibold text-gray-900">{getPageTitle()}</span>
              </Link>

              <nav aria-label="Admin navigation" className="hidden md:flex items-center gap-2">
                {items.map((it) => {
                  const active = pathname === it.to || pathname.startsWith(it.to + '/');
                  return (
                    <Link
                      key={it.to}
                      to={it.to}
                      className={`px-3 py-2 text-sm rounded-md ${active ? 'text-blue-700 font-semibold bg-blue-50' : 'text-gray-700 hover:bg-gray-50'}`}
                    >
                      <span className="inline-flex items-center gap-2">
                        {active && <span className="w-2 h-2 bg-blue-600 rounded-full" aria-hidden />}
                        <span>{it.label}</span>
                      </span>
                    </Link>
                  );
                })}
              </nav>
            </div>

            <div className="flex items-center gap-3">
              <div className="hidden md:block">{actions}</div>

              <LanguageSelector />

              <div className="relative">
                <button
                  onClick={toggleUserMenu}
                  className="flex items-center gap-2 px-3 py-1 rounded-md hover:bg-gray-50"
                  aria-haspopup="true"
                  aria-expanded={userMenuOpen}
                >
                  <span className="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center text-sm text-gray-700">
                    {user?.email ? user.email.charAt(0).toUpperCase() : 'U'}
                  </span>
                </button>

                {userMenuOpen && (
                  <div className="absolute right-0 mt-2 w-40 bg-white border rounded-md shadow-lg z-20">
                    <Link to="/admin/profile" className="block px-3 py-2 text-sm text-gray-700 hover:bg-gray-50">{t('nav.profile')}</Link>
                    <button onClick={logout} className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50">{t('nav.logout')}</button>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Mobile menu */}
        {mobileOpen && (
          <div className="md:hidden border-t bg-white">
            <div className="px-4 py-3 flex flex-col">
              {items.map((it) => {
                const active = pathname === it.to || pathname.startsWith(it.to + '/');
                return (
                  <Link
                    key={it.to}
                    to={it.to}
                    className={`py-2 text-sm rounded-md ${active ? 'text-blue-700 font-semibold bg-blue-50' : 'text-gray-700 hover:bg-gray-50'}`}
                    onClick={() => setMobileOpen(false)}
                  >
                    <span className="inline-flex items-center gap-2">
                      {active && <span className="w-2 h-2 bg-blue-600 rounded-full" aria-hidden />}
                      <span>{it.label}</span>
                    </span>
                  </Link>
                );
              })}
              <div className="mt-2 border-t pt-2">
                {actions}
              </div>
              <div className="mt-2 border-t pt-2">
                <Link to="/admin/profile" className="block py-2 text-sm text-gray-700">{t('nav.profile')}</Link>
                <button onClick={logout} className="w-full text-left py-2 text-sm text-gray-700">{t('nav.logout')}</button>
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
