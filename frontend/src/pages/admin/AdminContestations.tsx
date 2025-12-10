import React, { useEffect, useState } from 'react';
import { getPendingContestations, approveContestation, rejectContestation } from '../../services/adminContestations';
import type { AdminContestation } from '../../services/adminContestations';
import AdminLayout from '../../components/admin/AdminLayout';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const AdminContestations: React.FC = () => {
  const [contestations, setContestations] = useState<AdminContestation[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { t } = useTranslation('admin');
  const [actionLoading, setActionLoading] = useState<number | null>(null);
  const [reviewNotes, setReviewNotes] = useState<{ [id: number]: string }>({});

  const navItems = [
    { label: t('dashboard'), to: '/admin/dashboard' },
    { label: t('products'), to: '/admin/products' },
    { label: t('orders'), to: '/admin/orders' },
    { label: t('users'), to: '/admin/users' },
    { label: t('auditLogs'), to: '/admin/audit-logs' },
    { label: t('contestations'), to: '/admin/contestations' },
  ];

  const fetchContestations = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await getPendingContestations();
      setContestations(res.data);
    } catch {
      setError(t('failedToLoadContestations'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchContestations();
  }, [fetchContestations]);

  const handleApprove = async (id: number) => {
    setActionLoading(id);
    setError(null);
    try {
      await approveContestation(id, reviewNotes[id]);
      setContestations(contestations.filter(c => c.id !== id));
    } catch {
      setError(t('failedToApproveContestation'));
    } finally {
      setActionLoading(null);
    }
  };

  const handleReject = async (id: number) => {
    setActionLoading(id);
    setError(null);
    try {
      await rejectContestation(id, reviewNotes[id]);
      setContestations(contestations.filter(c => c.id !== id));
    } catch {
      setError(t('failedToRejectContestation'));
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <AdminLayout
      navItems={navItems}
      actions={<div className="flex items-center gap-2"><Link to="/admin/dashboard" className="px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-lg hover:bg-gray-50">← {t('dashboard')}</Link></div>}
    >
      <div className="p-6">
        <h2 className="text-2xl font-bold mb-4">{t('pendingUserContestations')}</h2>
        {error && <div className="text-red-500 mb-2">{error}</div>}
        {loading ? (
          <div>{t('loadingContestations')}</div>
        ) : contestations.length === 0 ? (
          <div>{t('noPendingContestations')}</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full border bg-white rounded shadow">
              <thead className="bg-gray-100">
                <tr>
                  <th className="border px-3 py-2 text-left">{t('contestationId')}</th>
                  <th className="border px-3 py-2 text-left">{t('userId')}</th>
                  <th className="border px-3 py-2 text-left">{t('reason')}</th>
                  <th className="border px-3 py-2 text-left">{t('requestedAt')}</th>
                  <th className="border px-3 py-2 text-left">{t('reviewNotes')}</th>
                  <th className="border px-3 py-2 text-left">{t('actions')}</th>
                </tr>
              </thead>
              <tbody>
                {contestations.map(c => (
                  <tr key={c.id} className="hover:bg-gray-50">
                    <td className="border px-3 py-2">{c.id}</td>
                    <td className="border px-3 py-2">{c.user_id}</td>
                    <td className="border px-3 py-2 max-w-xs truncate" title={c.reason}>{c.reason}</td>
                    <td className="border px-3 py-2">{new Date(c.requested_at).toLocaleString()}</td>
                    <td className="border px-3 py-2">
                      <input
                        type="text"
                        className="border px-2 py-1 w-40 rounded"
                        placeholder={t('reviewNotesPlaceholder')}
                        value={reviewNotes[c.id] || ''}
                        onChange={e => setReviewNotes({ ...reviewNotes, [c.id]: e.target.value })}
                      />
                    </td>
                    <td className="border px-3 py-2">
                      <button
                        className="bg-green-600 text-white px-3 py-1 mr-2 rounded hover:bg-green-700 disabled:opacity-50"
                        disabled={actionLoading === c.id}
                        onClick={() => handleApprove(c.id)}
                      >{t('approve')}</button>
                      <button
                        className="bg-red-600 text-white px-3 py-1 rounded hover:bg-red-700 disabled:opacity-50"
                        disabled={actionLoading === c.id}
                        onClick={() => handleReject(c.id)}
                      >{t('reject')}</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </AdminLayout>
  );
};

export default AdminContestations;
