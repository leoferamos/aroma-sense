import React, { useState, useEffect } from 'react';
import { isAxiosError } from 'axios';
import { useAuth } from '../hooks/useAuth';
import { cancelAccountDeletion, exportMyData, requestContestation } from '../services/profile';
import { useCart } from '../hooks/useCart';

interface AccountBlockOverlayProps {
  deactivationData?: {
    deactivated_at: string;
    deactivated_by: string;
    deactivation_reason: string;
    deactivation_notes?: string;
    suspension_until?: string;
    contestation_deadline?: string;
  } | null;
}

const AccountBlockOverlay: React.FC<AccountBlockOverlayProps> = ({ deactivationData }) => {
  const { user, refreshUser, logout } = useAuth();
  const { refresh: refreshCart } = useCart();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [contestationReason, setContestationReason] = useState('');
  const [showContestationForm, setShowContestationForm] = useState(false);
  const [contestationSubmitted, setContestationSubmitted] = useState(false);

  const isDeletionBlocked = Boolean(user?.deletion_requested_at || user?.deletion_confirmed_at);
  const isDeactivated = Boolean(user?.deactivated_at || deactivationData);
  const isBlocked = isDeletionBlocked || isDeactivated;
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
      } catch {
        // ignore refresh errors
      }
      window.location.replace('/products');
    } catch {
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
    } catch {
      setError('Failed to export data');
    } finally {
      setLoading(false);
    }
  };

  const onContestDeactivation = async () => {
    if (contestationReason.trim().length < 10) {
      setError('Please provide a reason for contestation (minimum 10 characters)');
      return;
    }
    setError(null);
    setLoading(true);
    try {
      await requestContestation({ reason: contestationReason });
      setShowContestationForm(false);
      setContestationReason('');
      setContestationSubmitted(true);
      setError('Contestation submitted successfully. Our team will review it within 5 business days.');
    } catch (err) {
      // Handle validation / already-submitted responses from the API
      if (isAxiosError(err) && err.response && err.response.data) {
        const body = err.response.data as { error?: string; message?: string };
        const serverMsg: string = body?.error || body?.message || '';
        const lower = serverMsg.toLowerCase();
        if (serverMsg.includes('ContestationRequest.Reason') || serverMsg.includes('min')) {
          setContestationSubmitted(true);
          setError('You have already submitted a contestation or the reason is too short (minimum 10 characters).');
          // Close modal so the main overlay shows the error
          setShowContestationForm(false);
        } else if (lower.includes('already') || lower.includes('submitted')) {
          setContestationSubmitted(true);
          setError('You have already submitted a contestation.');
          setShowContestationForm(false);
        } else {
          setError(serverMsg || 'Failed to submit contestation');
        }
      } else {
        setError('Failed to submit contestation');
      }
    } finally {
      setLoading(false);
    }
  };

  const onLogout = async () => {
    setError(null);
    setLoading(true);
    try {
      await logout();
    } catch {
      setError('Failed to logout');
    } finally {
      setLoading(false);
    }
  };

  const requestedAt = user?.deletion_requested_at ? new Date(user.deletion_requested_at).toLocaleString() : null;
  const confirmedAt = user?.deletion_confirmed_at ? new Date(user.deletion_confirmed_at).toLocaleString() : null;
  const deactivatedAt = deactivationData?.deactivated_at ? new Date(deactivationData.deactivated_at).toLocaleString() : user?.deactivated_at ? new Date(user.deactivated_at).toLocaleString() : null;
  const suspensionUntil = deactivationData?.suspension_until ? new Date(deactivationData.suspension_until).toLocaleString() : user?.suspension_until ? new Date(user.suspension_until).toLocaleString() : null;

  let title = '';
  let message = '';

  if (isDeactivated) {
    title = 'Account deactivated';
    const reason = deactivationData?.deactivation_reason || user?.deactivation_reason || 'Not specified';
    const notes = deactivationData?.deactivation_notes || user?.deactivation_notes;
    message = `Your account was deactivated on ${deactivatedAt}. Reason: ${reason}. ${notes ? `Notes: ${notes}.` : ''} ${suspensionUntil ? `Suspended until ${suspensionUntil}.` : ''} You can export your data or contest this deactivation.`;
  } else if (user?.deletion_confirmed_at) {
    title = 'Account deletion confirmed';
    message = `Your account deletion was confirmed on ${confirmedAt}. Your personal data will be retained according to our retention policy.`;
  } else {
    title = 'Account deletion requested';
    message = `We received a request to delete your account on ${requestedAt}. You are in a cooling-off period and can cancel the deletion or export your data.`;
  }

  if (!isBlocked) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-white pointer-events-auto">
      <div className="relative max-w-xl w-full bg-white rounded-lg shadow-xl p-6 border border-gray-100">
        <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
        <button
          type="button"
          onClick={onLogout}
          disabled={loading}
          className="absolute top-3 right-3 inline-flex items-center px-3 py-1.5 bg-red-600 text-white rounded-md text-sm hover:bg-red-700 disabled:opacity-50"
        >
          Logout
        </button>
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
          {isDeactivated ? (
            <button
              type="button"
              onClick={() => setShowContestationForm(true)}
              disabled={loading}
              className="inline-flex items-center px-4 py-2 bg-white border border-gray-300 text-sm rounded-md hover:bg-gray-50 disabled:opacity-50"
            >
              Contest deactivation
            </button>
          ) : !user?.deletion_confirmed_at && (
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

      {showContestationForm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
          <div className="max-w-md w-full bg-white rounded-lg shadow-xl p-6 border border-gray-100">
            <h4 className="text-md font-semibold text-gray-900">Contest Account Deactivation</h4>
            <p className="mt-2 text-sm text-gray-700">Please provide a reason for contesting this deactivation. Our team will review it within 5 business days.</p>
            <textarea
              value={contestationReason}
              onChange={(e) => setContestationReason(e.target.value)}
              placeholder="Explain why you believe this deactivation was incorrect..."
              className="mt-3 w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
              rows={4}
              maxLength={500}
            />
              {contestationSubmitted && (
                <p className="mt-2 text-sm text-gray-600">You have already submitted a contestation.</p>
              )}
            <div className="mt-4 flex gap-3">
              <button
                type="button"
                onClick={onContestDeactivation}
<<<<<<< HEAD
                disabled={loading || contestationReason.trim().length < 10 || contestationSubmitted}
                className="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-md text-sm hover:bg-indigo-700 disabled:opacity-50"
              >
                {contestationSubmitted ? 'Submitted' : 'Submit Contest'}
              </button>
              <button
                type="button"
                onClick={() => setShowContestationForm(false)}
                disabled={loading}
                className="inline-flex items-center px-4 py-2 bg-white border border-gray-300 text-sm rounded-md hover:bg-gray-50 disabled:opacity-50"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AccountBlockOverlay;
