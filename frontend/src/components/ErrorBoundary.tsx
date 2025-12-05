import React, { useState, useEffect } from 'react';
import { useAuth } from '../hooks/useAuth';
import AccountBlockOverlay from './AccountBlockOverlay';
import { useTranslation } from 'react-i18next';

interface ErrorBoundaryProps {
  children: React.ReactNode;
}

const ErrorBoundary: React.FC<ErrorBoundaryProps> = ({ children }) => {
  const [hasError, setHasError] = useState(false);
  const { user } = useAuth();
  const { t } = useTranslation('errors');

  useEffect(() => {
    const handleError = (error: unknown, errorInfo: unknown) => {
      console.error('Uncaught error in component tree:', error, errorInfo);
      setHasError(true);
    };

    // Listen for unhandled errors
    window.addEventListener('error', (event) => handleError(event.error, event));
    window.addEventListener('unhandledrejection', (event) => handleError(event.reason, event));

    return () => {
      window.removeEventListener('error', (event) => handleError(event.error, event));
      window.removeEventListener('unhandledrejection', (event) => handleError(event.reason, event));
    };
  }, []);

  if (hasError) {
    if (user && (user.deletion_requested_at || user.deletion_confirmed_at)) {
      return <AccountBlockOverlay />;
    }

    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="bg-white shadow rounded-lg p-8 text-center max-w-md">
          <h1 className="text-xl font-semibold text-gray-900">{t('somethingWrong')}</h1>
          <p className="mt-2 text-gray-600">{t('refreshOrTryLater')}</p>
          <div className="mt-6 flex flex-col gap-3">
            <button
              type="button"
              onClick={() => setHasError(false)}
              className="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors"
            >
              {t('tryAgain')}
            </button>
            <button
              type="button"
              onClick={() => { window.location.assign('/products'); }}
              className="px-4 py-2 rounded-lg border border-gray-300 text-gray-800 hover:bg-gray-50 transition-colors"
            >
              {t('goToProducts')}
            </button>
          </div>
        </div>
      </div>
    );
  }

  return <>{children}</>;
};

export default ErrorBoundary;
