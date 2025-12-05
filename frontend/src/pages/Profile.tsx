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
import { useTranslation } from 'react-i18next';

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
  const { t } = useTranslation('common');

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
        <h1 className="text-2xl font-semibold text-gray-900 mb-6">{t('profile.title')}</h1>

        {/* Deletion status banner */}
        {profile?.deletion_requested_at && !profile?.deletion_confirmed_at && (
          <div className="mb-6 rounded-lg border border-yellow-200 bg-yellow-50 p-4">
            <h3 className="text-sm font-semibold text-yellow-800">{t('profile.deletionPendingTitle')}</h3>
            <p className="text-sm text-yellow-700 mt-1">
              {t('profile.deletionPendingText', { date: new Date(profile.deletion_requested_at).toLocaleString() })}
            </p>
          </div>
        )}

        {profile?.deletion_confirmed_at && (
          <div className="mb-6 rounded-lg border border-red-200 bg-red-50 p-4">
            <h3 className="text-sm font-semibold text-red-800">{t('profile.deletionConfirmedTitle')}</h3>
            <p className="text-sm text-red-700 mt-1">
              {t('profile.deletionConfirmedText', { date: new Date(profile.deletion_confirmed_at).toLocaleString() })}
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
              label={t('profile.email')}
              type="email"
              name="email"
              value={profile?.email ?? ''}
              onChange={() => { /* read-only */ }}
              placeholder="you@example.com"
              disabled
              readOnly
            />
            <InputField
              label={t('profile.yourName')}
              name="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t('profile.namePlaceholder')}
              required
            />
            <p className="text-sm text-gray-500">{t('profile.nameDescription')}</p>
            <div>
              <button
                type="submit"
                disabled={saving || name.trim().length < 2}
                className="inline-flex items-center px-5 py-2.5 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50"
              >
                {saving ? t('profile.saving') : t('profile.saveChanges')}
              </button>
            </div>
          </form>
        </div>

        {/* Security */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 space-y-4 mt-6">
          <h2 className="text-lg font-semibold">{t('profile.security')}</h2>
          <p className="text-sm text-gray-600">{t('profile.securityDescription')}</p>
          <div>
            <button
              type="button"
              onClick={() => navigate('/change-password')}
              className="inline-flex items-center px-5 py-2.5 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50"
            >
              {t('profile.changePassword')}
            </button>
          </div>
        </div>

        {/* Account & Privacy */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 space-y-4 mt-6">
          <h2 className="text-lg font-semibold">{t('profile.accountPrivacy')}</h2>
          <p className="text-sm text-gray-600">{t('profile.accountPrivacyDescription')}</p>
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
                  setSuccess(t('profile.dataExported'));
                } catch {
                  setError(t('profile.exportFailed'));
                } finally {
                  setActionLoading(false);
                }
              }}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50"
              disabled={actionLoading}
            >
              {actionLoading ? t('profile.processing') : t('profile.exportData')}
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
                    setSuccess(t('profile.deletionCancelled'));
                  } catch {
                    setError(t('profile.cancelDeletionFailed'));
                  } finally {
                    setActionLoading(false);
                  }
                }}
                className="mt-3 sm:mt-0 inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
                disabled={actionLoading}
              >
                {t('profile.cancelDeletion')}
              </button>
            ) : (
              <button
                type="button"
                onClick={() => setDeleteModalOpen(true)}
                className="mt-3 sm:mt-0 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 disabled:opacity-50"
                disabled={actionLoading}
              >
                {t('profile.requestDeletion')}
              </button>
            )}
          </div>
        </div>

        <ConfirmModal
          open={deleteModalOpen}
          title={t('profile.deletionModalTitle')}
          description={t('profile.deletionModalDescription')}
          confirmText={t('profile.requestDeletionConfirm')}
          cancelText={t('profile.cancel')}
          requirePhrase="DELETE_MY_ACCOUNT"
          onConfirm={async () => {
            setActionLoading(true);
            try {
              await requestAccountDeletion();
              const updated = await getMyProfile();
              setProfile(updated);
              setSuccess(t('profile.deletionRequested'));
            } catch {
              setError(t('profile.requestDeletionFailed'));
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
