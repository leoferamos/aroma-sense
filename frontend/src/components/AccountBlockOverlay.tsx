import React, { useState, useEffect } from 'react';
import { useAuth } from '../hooks/useAuth';
import { cancelAccountDeletion, exportMyData } from '../services/profile';
import { useCart } from '../hooks/useCart';

const AccountBlockOverlay: React.FC = () => {
  const { user, refreshUser } = useAuth();
  const { refresh: refreshCart } = useCart();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const isBlocked = Boolean(user?.deletion_requested_at || user?.deletion_confirmed_at);
  // Prevent background scrolling while overlay is shown.
  useEffect(() => {
    if (!isBlocked) return;
    const prev = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    return () => { document.body.style.overflow = prev; };
  }, [isBlocked]);

  const onCancelDeletion = async () => {
    setError(null);
    setLoading(true);
    try {
      await cancelAccountDeletion();
      await refreshUser();
      try {
        await refreshCart();
      } catch (e) {
      }
      window.location.replace('/products');
    } catch (e) {
      setError('Failed to cancel deletion');
    } finally {
      setLoading(false);
    }
  };

  const onExport = async () => {
    setError(null);
    setLoading(true);
    try {
      const blob = await exportMyData();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `aroma-sense-data-${user?.public_id ?? 'me'}.json`;
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(url);
    } catch (e) {
      setError('Failed to export data');
    } finally {
      setLoading(false);
    }
  };

  const requestedAt = user?.deletion_requested_at ? new Date(user.deletion_requested_at).toLocaleString() : null;
  const confirmedAt = user?.deletion_confirmed_at ? new Date(user.deletion_confirmed_at).toLocaleString() : null;

  const title = user?.deletion_confirmed_at ? 'Account deletion confirmed' : 'Account deletion requested';
  const message = user?.deletion_confirmed_at
    ? `Your account deletion was confirmed on ${confirmedAt}. Your personal data will be retained according to our retention policy.`
    : `We received a request to delete your account on ${requestedAt}. You are in a cooling-off period and can cancel the deletion or export your data.`;

  if (!isBlocked) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-white pointer-events-auto">
      <div className="max-w-xl w-full bg-white rounded-lg shadow-xl p-6 border border-gray-100">
        <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
        <p className="mt-2 text-sm text-gray-700">{message}</p>

        {error && <div className="mt-3 text-sm text-red-600">{error}</div>}

        <div className="mt-4 flex gap-3">
          <button
            type="button"
            onClick={onExport}
            disabled={loading}
            className="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-md text-sm hover:bg-indigo-700 disabled:opacity-50"
          >
            Export my data
          </button>
          {!user?.deletion_confirmed_at && (
            <button
              type="button"
              onClick={onCancelDeletion}
              disabled={loading}
              className="inline-flex items-center px-4 py-2 bg-white border border-gray-300 text-sm rounded-md hover:bg-gray-50 disabled:opacity-50"
            >
              Cancel deletion
            </button>
          )}
        </div>

        <p className="mt-4 text-xs text-gray-500">You can still contact support if you need help.</p>
      </div>
    </div>
  );
};

export default AccountBlockOverlay;
