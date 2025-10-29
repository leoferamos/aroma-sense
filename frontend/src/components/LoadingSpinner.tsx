import React from 'react';

interface LoadingSpinnerProps {
  className?: string;
  message?: string;
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({ className = '', message = 'Loading...' }) => (
  <div className={`flex flex-col items-center justify-center py-20 ${className}`} role="status" aria-live="polite">
    <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
    <p className="text-gray-600">{message}</p>
  </div>
);

export default LoadingSpinner;
