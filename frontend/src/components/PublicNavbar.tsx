import React from 'react';
import { Link } from 'react-router-dom';
import { LOGO_PATH } from '../constants/app';

const PublicNavbar: React.FC = () => {
  return (
    <nav className="bg-white shadow-md sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link to="/products" className="flex items-center">
            <img src={LOGO_PATH} alt="Aroma Sense" className="h-10 w-auto" />
          </Link>

          {/* Simple Menu Icon*/}
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
    </nav>
  );
};

export default PublicNavbar;
