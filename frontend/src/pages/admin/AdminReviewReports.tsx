import React, { useCallback, useEffect, useMemo, useState } from 'react';
import AdminLayout from '../../components/admin/AdminLayout';
import { listReviewReports, resolveReviewReport, type ReviewReportAdminItem, type ReviewReportStatus } from '../../services/adminReviewReports';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

const AdminReviewReports: React.FC = () => {
  const { t } = useTranslation('admin');
  const [reports, setReports] = useState<ReviewReportAdminItem[]>([]);
  const [total, setTotal] = useState(0);
  const [statusFilter, setStatusFilter] = useState<ReviewReportStatus>('pending');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [notes, setNotes] = useState<Record<string, string>>({});
  const [deactivate, setDeactivate] = useState<Record<string, boolean>>({});
  const [suspensionUntil, setSuspensionUntil] = useState<Record<string, string>>({});
  const limit = 20;
  const [offset, setOffset] = useState(0);
  const totalPages = useMemo(() => Math.max(1, Math.ceil(total / limit)), [total, limit]);
  const currentPage = useMemo(() => Math.floor(offset / limit) + 1, [offset, limit]);

  const navItems = [
    { label: t('dashboard'), to: '/admin/dashboard' },
    { label: t('products'), to: '/admin/products' },
    { label: t('orders'), to: '/admin/orders' },
    { label: t('users'), to: '/admin/users' },
    { label: t('auditLogs'), to: '/admin/audit-logs' },
    { label: t('contestations'), to: '/admin/contestations' },
    { label: t('reviewReports'), to: '/admin/review-reports' },
  ];

  const fetchReports = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await listReviewReports({ status: statusFilter, limit, offset });
      setReports(res.items);
      setTotal(res.total);
    } catch (err) {
      setError(t('failedToLoadReviewReports'));
    } finally {
      setLoading(false);
    }
  }, [statusFilter, limit, offset, t]);

  useEffect(() => {
    setOffset(0);
  }, [statusFilter]);

  useEffect(() => {
    fetchReports();
  }, [fetchReports]);

  const handleResolve = async (id: string, action: 'accept' | 'reject') => {
    setActionLoading(id);
    setError(null);
    try {
      const suspendValue = suspensionUntil[id];
      await resolveReviewReport(id, {
        action,
        deactivate_user: !!deactivate[id],
        suspension_until: suspendValue ? new Date(suspendValue).toISOString() : null,
        notes: notes[id] ? notes[id] : null,
      });
      await fetchReports();
    } catch (err) {
      setError(t('failedToResolveReviewReport'));
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
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between mb-4">
          <h2 className="text-2xl font-bold">{t('reviewReportsTitle')}</h2>
          <div className="flex items-center gap-3">
            <label className="text-sm text-gray-700">
              {t('status')}:
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value as ReviewReportStatus)}
                className="ml-2 border rounded px-2 py-1 text-sm"
              >
                <option value="pending">{t('pending')}</option>
                <option value="accepted">{t('accepted')}</option>
                <option value="rejected">{t('rejected')}</option>
              </select>
            </label>
          </div>
        </div>

          {error && <div className="text-red-600 mb-3 text-sm">{error}</div>}

        {loading ? (
          <div>{t('loadingDots')}</div>
        ) : reports.length === 0 ? (
          <div className="text-gray-700">{t('noReviewReports')}</div>
        ) : (
          <div className="space-y-3">
            <div className="overflow-x-auto border rounded bg-white shadow">
              <table className="min-w-full text-sm">
                <thead className="bg-gray-100 text-left">
                  <tr>
                    <th className="px-3 py-2 border">{t('id')}</th>
                    <th className="px-3 py-2 border">{t('review')}</th>
                    <th className="px-3 py-2 border">{t('reporter')}</th>
                    <th className="px-3 py-2 border">{t('reason')}</th>
                    <th className="px-3 py-2 border">{t('status')}</th>
                    <th className="px-3 py-2 border">{t('requestedAt')}</th>
                    <th className="px-3 py-2 border">{t('actions')}</th>
                  </tr>
                </thead>
                <tbody>
                  {reports.map((r) => {
                    const pending = r.status === 'pending';
                    return (
                      <tr key={r.id} className="hover:bg-gray-50 align-top">
                        <td className="px-3 py-2 border text-gray-800">{r.id}</td>
                        <td className="px-3 py-2 border text-gray-800">
                          <div className="font-semibold">{t('rating')}: {r.review?.rating ?? '-'}</div>
                          <div className="text-gray-700 whitespace-pre-line">{r.review?.comment || t('noComment')}</div>
                          <div className="text-xs text-gray-500">ID: {r.review?.id}</div>
                          <div className="text-xs text-gray-500">User: {r.review?.user_id}</div>
                        </td>
                        <td className="px-3 py-2 border text-gray-800">
                          <div>{r.reporter?.display_name || t('unknownUser')}</div>
                          <div className="text-xs text-gray-500">{r.reporter?.public_id}</div>
                        </td>
                        <td className="px-3 py-2 border text-gray-800">
                          <div className="font-medium">{t(`reviewCategory.${r.reason_category}`, { defaultValue: r.reason_category })}</div>
                          <div className="text-gray-700 whitespace-pre-line mt-1">{r.reason_text || '-'}</div>
                        </td>
                        <td className="px-3 py-2 border text-gray-800 capitalize">{t(r.status)}</td>
                        <td className="px-3 py-2 border text-gray-800">{new Date(r.created_at).toLocaleString()}</td>
                        <td className="px-3 py-2 border w-64">
                          {pending ? (
                            <div className="space-y-2">
                              <textarea
                                className="w-full border rounded px-2 py-1 text-sm"
                                placeholder={t('reviewNotesPlaceholder')}
                                value={notes[r.id] || ''}
                                onChange={(e) => setNotes({ ...notes, [r.id]: e.target.value })}
                              />
                              <label className="flex items-center gap-2 text-sm text-gray-800">
                                <input
                                  type="checkbox"
                                  checked={!!deactivate[r.id]}
                                  onChange={(e) => setDeactivate({ ...deactivate, [r.id]: e.target.checked })}
                                />
                                {t('deactivateUser')}
                              </label>
                              <label className="block text-sm text-gray-800">
                                {t('suspensionUntil')}
                                <input
                                  type="datetime-local"
                                  className="mt-1 w-full border rounded px-2 py-1"
                                  value={suspensionUntil[r.id] || ''}
                                  onChange={(e) => setSuspensionUntil({ ...suspensionUntil, [r.id]: e.target.value })}
                                />
                              </label>
                              <div className="flex gap-2">
                                <button
                                  className="bg-green-600 text-white px-3 py-1 rounded hover:bg-green-700 disabled:opacity-50"
                                  disabled={actionLoading === r.id}
                                  onClick={() => handleResolve(r.id, 'accept')}
                                >{t('accept')}</button>
                                <button
                                  className="bg-red-600 text-white px-3 py-1 rounded hover:bg-red-700 disabled:opacity-50"
                                  disabled={actionLoading === r.id}
                                  onClick={() => handleResolve(r.id, 'reject')}
                                >{t('reject')}</button>
                              </div>
                            </div>
                          ) : (
                            <div className="text-sm text-gray-700">{t('alreadyResolved')}</div>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
            {total > limit && (
              <div className="flex justify-between items-center mt-3 text-sm text-gray-700">
                <div>{t('pageOf', { page: currentPage, total: totalPages })}</div>
                <div className="flex gap-2">
                  <button
                    className="px-3 py-1 border rounded disabled:opacity-50"
                    disabled={offset <= 0}
                    onClick={() => setOffset(Math.max(0, offset - limit))}
                  >{t('prev')}</button>
                  <button
                    className="px-3 py-1 border rounded disabled:opacity-50"
                    disabled={(offset + limit) >= total}
                    onClick={() => setOffset(Math.min(total - 1, offset + limit))}
                  >{t('next')}</button>
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </AdminLayout>
  );
};

export default AdminReviewReports;
