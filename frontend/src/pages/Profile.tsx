import React, { useEffect, useState } from 'react';
import { getMyProfile, updateMyProfile, type ProfileResponse } from '../services/profile';
import LoadingSpinner from '../components/LoadingSpinner';
import InputField from '../components/InputField';
import Navbar from '../components/Navbar';

const Profile: React.FC = () => {
  const [profile, setProfile] = useState<ProfileResponse | null>(null);
  const [name, setName] = useState('');
  const [saving, setSaving] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;
    (async () => {
      try {
    const data = await getMyProfile();
    if (!mounted) return;
		setProfile(data);
		setName(data.display_name ?? '');
      } catch (e: any) {
        if (!mounted) return;
        setError(e?.response?.data?.error || 'Failed to load profile');
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
    } catch (e: any) {
      if (e?.response?.status === 401) {
        setError('Session expired. Please sign in again.');
      } else {
        setError(e?.response?.data?.error || 'Failed to update profile');
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
        <h1 className="text-2xl font-semibold text-gray-900 mb-6">My Profile</h1>

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
      </main>
    </div>
  );
};

export default Profile;
