import React, { useEffect, useState } from 'react';
import AdminLayout from '../../components/admin/AdminLayout';
import { getAuditLogs, getAuditLog } from '../../services/audit';
import type { GetAuditLogsParams, AuditLog } from '../../types/audit';
import { useSearchParams } from 'react-router-dom';
import AuditLogTable from '../../components/admin/AuditLogTable';
import AuditLogDetailsModal from '../../components/admin/AuditLogDetailsModal';
import PaginationControls from '../../components/admin/PaginationControls';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const AdminAuditLogs: React.FC = () => {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [limit, setLimit] = useState<number>(25);
  const [offset, setOffset] = useState<number>(0);
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // filters
  const [userId, setUserId] = useState<string>('');
  const [actorId, setActorId] = useState<string>('');
  const [action, setAction] = useState<string>('');
  const [resource, setResource] = useState<string>('');
  const [resourceId, setResourceId] = useState<string>('');
  const [startDate, setStartDate] = useState<string>('');
  const [endDate, setEndDate] = useState<string>('');

  const { t } = useTranslation('admin');

  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  const fetch = async () => {
    setLoading(true);
    setError(null);
    const params: GetAuditLogsParams = {
      limit,
      offset,
      user_id: userId ? Number(userId) : undefined,
      actor_id: actorId ? Number(actorId) : undefined,
      action: action || undefined,
      resource: resource || undefined,
      resource_id: resourceId || undefined,
      start_date: startDate || undefined,
      end_date: endDate || undefined,
    };

    try {
      const resp = await getAuditLogs(params);
      setLogs(resp.audit_logs || []);
      setTotal(resp.total || 0);
    } catch (err) {
      console.debug('getAuditLogs error', err);
      setError(t('failedToLoadAuditLogs'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetch();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [limit, offset]);

  // auto-open modal when ?id= is present
  const [searchParams] = useSearchParams();
  useEffect(() => {
    const id = searchParams.get('id');
    if (!id) return;
    let mounted = true;
    (async () => {
      try {
        const single = await getAuditLog(id);
        if (!mounted) return;
        setSelectedLog(single);
        setModalOpen(true);
      } catch (err) {
        console.debug('getAuditLog error', err);
      }
    })();
    return () => { mounted = false; };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const openLog = (l: AuditLog) => { setSelectedLog(l); setModalOpen(true); };

  const totalPages = Math.max(1, Math.ceil(total / limit));
  const page = Math.floor(offset / limit) + 1;

  const actions = (
    <div className="flex items-center gap-2">
      <Link to="/admin/dashboard" className="px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-lg hover:bg-gray-50">← {t('nav.dashboard')}</Link>
    </div>
  );

  return (
    <AdminLayout actions={actions}>
      <div className="p-6">
        <h1 className="text-2xl font-semibold mb-4">{t('auditLogsTitle')}</h1>

        <div className="mb-4 grid grid-cols-1 md:grid-cols-3 gap-3">
          <input value={userId} onChange={(e) => setUserId(e.target.value)} placeholder={t('userId')} className="border rounded px-2 py-1" />
          <input value={actorId} onChange={(e) => setActorId(e.target.value)} placeholder={t('actorId')} className="border rounded px-2 py-1" />
          <input value={action} onChange={(e) => setAction(e.target.value)} placeholder={t('action')} className="border rounded px-2 py-1" />
          <input value={resource} onChange={(e) => setResource(e.target.value)} placeholder={t('resource')} className="border rounded px-2 py-1" />
          <input value={resourceId} onChange={(e) => setResourceId(e.target.value)} placeholder={t('resourceId')} className="border rounded px-2 py-1" />
          <div className="flex gap-2">
            <input type="datetime-local" value={startDate} onChange={(e) => setStartDate(e.target.value)} className="border rounded px-2 py-1 w-full" />
            <input type="datetime-local" value={endDate} onChange={(e) => setEndDate(e.target.value)} className="border rounded px-2 py-1 w-full" />
          </div>
        </div>

        <div className="mb-4 flex items-center gap-3">
          <button onClick={() => { setOffset(0); fetch(); }} className="px-3 py-2 bg-blue-600 text-white rounded">{t('apply')}</button>
          <select value={limit} onChange={(e) => { setLimit(Number(e.target.value)); setOffset(0); }} className="border px-2 py-1 rounded">
            {[10,25,50,100].map(n => <option key={n} value={n}>{n} {t('perPage')}</option>)}
          </select>
        </div>

        {loading && <div className="py-8">{t('loading')}</div>}
        {error && <div className="py-8 text-red-600">{error}</div>}

        {!loading && !error && (
          <>
            <AuditLogTable logs={logs} onView={openLog} />

            <div className="mt-4 flex items-center justify-between">
              <div className="text-sm text-gray-600">{t('showingRange', { start: Math.min(total, offset + 1), end: Math.min(total, offset + limit), total })}</div>
              <PaginationControls
                page={page}
                totalPages={totalPages}
                onPrev={() => setOffset(Math.max(0, offset - limit))}
                onNext={() => setOffset(Math.min(totalPages * limit, offset + limit))}
              />
            </div>
          </>
        )}

        <AuditLogDetailsModal open={modalOpen} onClose={() => setModalOpen(false)} log={selectedLog} />
      </div>
    </AdminLayout>
  );
};

export default AdminAuditLogs;
