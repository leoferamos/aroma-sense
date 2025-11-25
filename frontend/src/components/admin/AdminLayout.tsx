import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

type NavItem = {
  label: string;
  to: string;
};

type Props = {
  title?: string;
  children: React.ReactNode;
  actions?: React.ReactNode;
  navItems?: NavItem[];
};

const defaultNav: NavItem[] = [
  { label: 'Dashboard', to: '/admin/dashboard' },
  { label: 'Products', to: '/admin/products' },
  { label: 'Orders', to: '/admin/orders' },
  { label: 'Users', to: '/admin/users' },
];

const AdminLayout: React.FC<Props> = ({ title, children, actions, navItems }) => {
  const items = navItems ?? defaultNav;
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
                <span className="text-lg font-semibold text-gray-900">{title ?? 'Admin'}</span>
              </Link>

              <nav aria-label="Admin navigation" className="hidden md:flex items-center gap-2">
                {items.map((it) => (
                  <Link key={it.to} to={it.to} className="px-3 py-2 text-sm text-gray-700 hover:bg-gray-50 rounded-md">
                    {it.label}
                  </Link>
                ))}
              </nav>
            </div>

            <div className="flex items-center gap-3">
              <div className="hidden md:block">{actions}</div>

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
                    <Link to="/admin/profile" className="block px-3 py-2 text-sm text-gray-700 hover:bg-gray-50">Profile</Link>
                    <button onClick={logout} className="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-50">Logout</button>
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
              {items.map((it) => (
                <Link key={it.to} to={it.to} className="py-2 text-sm text-gray-700 hover:bg-gray-50 rounded-md" onClick={() => setMobileOpen(false)}>
                  {it.label}
                </Link>
              ))}
              <div className="mt-2 border-t pt-2">
                {actions}
              </div>
              <div className="mt-2 border-t pt-2">
                <Link to="/admin/profile" className="block py-2 text-sm text-gray-700">Profile</Link>
                <button onClick={logout} className="w-full text-left py-2 text-sm text-gray-700">Logout</button>
              </div>
            </div>
          </div>
        )}
      </header>

      <main role="main" className="max-w-7xl mx-auto p-4">{children}</main>
    </div>
  );
};

export default AdminLayout;
