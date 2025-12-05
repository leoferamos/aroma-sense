import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { isAxiosError } from 'axios';
import { getMyProfile, updateMyProfile, type ProfileResponse } from '../services/profile';
import LoadingSpinner from '../components/LoadingSpinner';
import InputField from '../components/InputField';
import Navbar from '../components/Navbar';
import BackButton from '../components/BackButton';
import ConfirmModal from '../components/ConfirmModal';
import { requestAccountDeletion, cancelAccountDeletion, exportMyData } from '../services/profile';

const Profile: React.FC = () => {
  const navigate = useNavigate();
  const [profile, setProfile] = useState<ProfileResponse | null>(null);
  const [name, setName] = useState('');
  const [saving, setSaving] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState(false);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);

  useEffect(() => {
    let mounted = true;
    (async () => {
      try {
        const data = await getMyProfile();
        if (!mounted) return;
        setProfile(data);
        setName(data.display_name ?? '');
      } catch (e: unknown) {
        if (!mounted) return;
        if (isAxiosError(e)) {
          setError(e.response?.data?.error || 'Failed to load profile');
        } else if (e instanceof Error) {
          setError(e.message);
        } else {
          setError('Failed to load profile');
        }
      } finally {
        if (mounted) setLoading(false);
      }
    })();
    return () => { mounted = false; };
  }, []);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    setSaving(true);
    try {
      const updated = await updateMyProfile({ display_name: name.trim() });
      setProfile(updated);
      setSuccess('Profile updated successfully');
    } catch (e: unknown) {
      if (isAxiosError(e) && e.response?.status === 401) {
        setError('Session expired. Please sign in again.');
      } else {
        if (isAxiosError(e)) {
          setError(e.response?.data?.error || 'Failed to update profile');
        } else if (e instanceof Error) {
          setError(e.message);
        } else {
          setError('Failed to update profile');
        }
      }
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <main className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-4">
          <BackButton fallbackPath="/products" />
        </div>
        <h1 className="text-2xl font-semibold text-gray-900 mb-6">My Profile</h1>

        {/* Deletion status banner */}
        {profile?.deletion_requested_at && !profile?.deletion_confirmed_at && (
          <div className="mb-6 rounded-lg border border-yellow-200 bg-yellow-50 p-4">
            <h3 className="text-sm font-semibold text-yellow-800">Account deletion pending</h3>
            <p className="text-sm text-yellow-700 mt-1">
              We received a request to delete your account on{' '}
              <strong>{new Date(profile.deletion_requested_at).toLocaleString()}</strong>.
              You have a 7-day cooling-off period during which you can cancel the deletion or export your data.
            </p>
          </div>
        )}

        {profile?.deletion_confirmed_at && (
          <div className="mb-6 rounded-lg border border-red-200 bg-red-50 p-4">
            <h3 className="text-sm font-semibold text-red-800">Account deletion confirmed</h3>
            <p className="text-sm text-red-700 mt-1">
              Your account deletion was confirmed on <strong>{new Date(profile.deletion_confirmed_at).toLocaleString()}</strong>.
              Your personal data will be retained for the configured retention period and then anonymized.
            </p>
          </div>
        )}

        {error && (
          <div className="mb-4 rounded-md bg-red-50 p-4 text-red-700 border border-red-200">{error}</div>
        )}
        {success && (
          <div className="mb-4 rounded-md bg-green-50 p-4 text-green-700 border border-green-200">{success}</div>
        )}

        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 space-y-6">
          <form onSubmit={onSubmit} className="space-y-5">
            <InputField
              label="Email"
              type="email"
              name="email"
              value={profile?.email ?? ''}
              onChange={() => { /* read-only */ }}
              placeholder="you@example.com"
              disabled
              readOnly
            />
            <InputField
              label="Your name"
              name="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Your name"
              required
            />
            <p className="text-sm text-gray-500">Shown publicly with your reviews.</p>
            <div>
              <button
                type="submit"
                disabled={saving || name.trim().length < 2}
                className="inline-flex items-center px-5 py-2.5 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50"
              >
                {saving ? 'Saving…' : 'Save changes'}
              </button>
            </div>
          </form>
        </div>

        {/* Security */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 space-y-4 mt-6">
          <h2 className="text-lg font-semibold">Security</h2>
          <p className="text-sm text-gray-600">Manage your password and account security.</p>
          <div>
            <button
              type="button"
              onClick={() => navigate('/change-password')}
              className="inline-flex items-center px-5 py-2.5 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50"
            >
              Change your password
            </button>
          </div>
        </div>

        {/* Account & Privacy */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 space-y-4 mt-6">
          <h2 className="text-lg font-semibold">Account & Privacy</h2>
          <p className="text-sm text-gray-600">Export your data, request account deletion, or cancel a pending deletion.</p>
          <div className="flex flex-col sm:flex-row sm:items-center sm:gap-4 mt-4">
            <button
              type="button"
              onClick={async () => {
                setActionLoading(true);
                try {
                  const blob = await exportMyData();
                  const url = window.URL.createObjectURL(blob);
                  const a = document.createElement('a');
                  a.href = url;
                  a.download = `aroma-sense-data-${profile?.public_id ?? 'me'}.json`;
                  document.body.appendChild(a);
                  a.click();
                  a.remove();
                  window.URL.revokeObjectURL(url);
                  setSuccess('Data exported successfully');
                } catch {
                  setError('Failed to export data');
                } finally {
                  setActionLoading(false);
                }
              }}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50"
              disabled={actionLoading}
            >
              {actionLoading ? 'Processing…' : 'Export my data'}
            </button>

            {profile?.deletion_requested_at && !profile?.deletion_confirmed_at ? (
              <button
                type="button"
                onClick={async () => {
                  setActionLoading(true);
                  try {
                    await cancelAccountDeletion();
                    const updated = await getMyProfile();
                    setProfile(updated);
                    setSuccess('Account deletion cancelled');
                  } catch {
                    setError('Failed to cancel deletion');
                  } finally {
                    setActionLoading(false);
                  }
                }}
                className="mt-3 sm:mt-0 inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
                disabled={actionLoading}
              >
                Cancel deletion
              </button>
            ) : (
              <button
                type="button"
                onClick={() => setDeleteModalOpen(true)}
                className="mt-3 sm:mt-0 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 disabled:opacity-50"
                disabled={actionLoading}
              >
                Request account deletion
              </button>
            )}
          </div>
        </div>

        <ConfirmModal
          open={deleteModalOpen}
          title="Request account deletion"
          description="Type DELETE_MY_ACCOUNT to confirm that you want to request account deletion. You will have 7 days to cancel."
          confirmText="Request deletion"
          cancelText="Cancel"
          requirePhrase="DELETE_MY_ACCOUNT"
          onConfirm={async () => {
            setActionLoading(true);
            try {
              await requestAccountDeletion();
              const updated = await getMyProfile();
              setProfile(updated);
              setSuccess('Account deletion requested successfully');
            } catch {
              setError('Failed to request deletion');
            } finally {
              setActionLoading(false);
              setDeleteModalOpen(false);
            }
          }}
          onCancel={() => setDeleteModalOpen(false)}
          loading={actionLoading}
        />
      </main>
    </div>
  );
};

export default Profile;
