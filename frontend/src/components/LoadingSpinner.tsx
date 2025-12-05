import React from 'react';
import { useTranslation } from 'react-i18next';

interface LoadingSpinnerProps {
  className?: string;
  message?: string;
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({ className = '', message }) => {
  const { t } = useTranslation('common');
  const displayMessage = message || t('common.loading');

  return (
    <div className={`flex flex-col items-center justify-center py-20 ${className}`} role="status" aria-live="polite">
      <div className="relative w-12 h-12 mb-4">
        <div className="absolute inset-0 rounded-full border-2 border-gray-200"></div>
        <div className="absolute inset-0 rounded-full border-2 border-transparent border-t-blue-600 border-r-blue-600 animate-spin"></div>
      </div>
      <p className="text-gray-600 text-sm font-medium">{displayMessage}</p>
    </div>
  );
};

export default LoadingSpinner;
