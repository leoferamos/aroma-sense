import React from 'react';
import { Link } from 'react-router-dom';

export type UserRole = 'admin' | 'client' | null | undefined;

interface UserMenuProps {
  open: boolean;
  role?: UserRole;
  isAuthenticated: boolean;
  onClose: () => void;
  onLogout: () => void;
  onSignIn: () => void;
}

const UserMenu: React.FC<UserMenuProps> = ({ open, role, isAuthenticated, onClose, onLogout, onSignIn }) => {
  if (!open) return null;

  return (
    <div className="absolute right-0 top-12 w-64 bg-white rounded-lg shadow-xl border border-gray-100 overflow-hidden z-10">
      <nav className="py-2" aria-label="User menu">
  <Link to="/profile" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50" onClick={onClose}>Profile</Link>
  <Link to="/orders" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50" onClick={onClose}>My orders</Link>
        {role === 'admin' && (
          <button type="button" className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50">Admin dashboard</button>
        )}
        <div className="my-1 border-t border-gray-100" />
        <Link to="/terms" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50" onClick={onClose}>Terms</Link>
        <Link to="/privacy" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50" onClick={onClose}>Privacy</Link>
        <div className="my-1 border-t border-gray-100" />
        {isAuthenticated ? (
          <button
            type="button"
            onClick={() => { onClose(); onLogout(); }}
            className="w-full text-left px-4 py-2 text-sm font-medium text-red-600 hover:bg-red-50"
          >
            Logout
          </button>
        ) : (
          <button
            type="button"
            onClick={() => { onClose(); onSignIn(); }}
            className="w-full text-left px-4 py-2 text-sm font-medium text-blue-600 hover:bg-blue-50"
          >
            Sign in
          </button>
        )}
      </nav>
    </div>
  );
};

export default UserMenu;
