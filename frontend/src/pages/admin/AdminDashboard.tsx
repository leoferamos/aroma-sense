import React from 'react';
import { useAuth } from '../../hooks/useAuth';
import { Link } from 'react-router-dom';
import AdminLayout from '../../components/admin/AdminLayout';
import { getAuditLogs, getAuditLogsSummary } from '../../services/audit';
import type { AuditLog, AuditLogSummary } from '../../types/audit';
import { useTranslation } from 'react-i18next';

const AdminDashboard: React.FC = () => {
  const { role, logout } = useAuth();
  const { t } = useTranslation('admin');
  const { t: tCommon } = useTranslation('common');

  const handleLogout = async () => {
    await logout();
  };

  const actions = (
    <div className="flex items-center gap-3">
      <span className="text-sm text-gray-500">
        <span className="font-medium text-blue-600">{role}</span>
      </span>
      <button
        onClick={handleLogout}
        className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors uppercase"
      >
        {t('nav.logout')}
      </button>
    </div>
  );

  // RecentAuditLogs component (inline) — shows a few latest logs
  const RecentAuditLogs: React.FC = () => {
    const [logs, setLogs] = React.useState<AuditLog[]>([]);
    const [summary, setSummary] = React.useState<AuditLogSummary | null>(null);
    const [loading, setLoading] = React.useState(false);
    const [error, setError] = React.useState<string | null>(null);

    React.useEffect(() => {
      let mounted = true;
      const fetch = async () => {
        setLoading(true);
        setError(null);
        try {
          const resp = await getAuditLogs({ limit: 5, offset: 0 });
          if (!mounted) return;
          setLogs(resp.audit_logs || []);
          // fetch summary in background
          try {
            const s = await getAuditLogsSummary();
            if (mounted) setSummary(s || null);
          } catch (e) {
            // ignore summary errors
            console.debug('getAuditLogsSummary error', e);
          }
        } catch (err) {
          console.debug('getAuditLogs recent error', err);
          if (!mounted) return;
          setError(tCommon('errors.failedToLoadAuditLogs'));
        } finally {
          if (mounted) setLoading(false);
        }
      };
      fetch();
      return () => { mounted = false; };
    }, []);

    function fmt(ts?: string | undefined | null) {
      if (!ts) return '-';
      const d = new Date(ts);
      if (Number.isNaN(d.getTime())) return ts;
      return d.toLocaleString();
    }

    if (loading) return <div className="py-4">{t('loading')}</div>;
    if (error) return <div className="py-4 text-red-600">{error}</div>;

    return (
      <div className="space-y-2">
        {summary && (
          <div className="flex items-center gap-3 mb-2">
            <div className="text-sm text-gray-700">{t('totalActions')}: <span className="font-medium">{summary.total_actions}</span></div>
            {Object.entries(summary.actions_by_type || {}).slice(0,3).map(([k,v]) => (
              <div key={k} className="text-xs bg-gray-100 px-2 py-1 rounded">{k}: {v}</div>
            ))}
          </div>
        )}
        {logs.length === 0 ? (
          <div className="text-sm text-gray-600">{t('noRecentAuditLogs')}</div>
        ) : (
          logs.map((l) => (
            <div key={l.id} className="flex items-center justify-between">
              <div className="flex items-center gap-4">
                <div className="text-xs text-gray-500 w-44">{fmt(l.timestamp || l.created_at)}</div>
                <div>
                  <div className="font-medium text-sm">{l.action}</div>
                  <div className="text-xs text-gray-500">{l.resource}{l.resource_id ? ` · ${l.resource_id}` : ''}</div>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div className="text-sm text-gray-600">{l.actor?.display_name || l.user?.display_name || '-'}</div>
                <Link to={`/admin/audit-logs?id=${l.id}`} className="text-blue-600 hover:underline text-sm">{t('details')}</Link>
              </div>
            </div>
          ))
        )}
      </div>
    );
  };

  return (
    <AdminLayout actions={actions}>
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-900">{t('welcomeAdmin')}</h2>
        <p className="text-gray-500 text-sm mt-1">{t('manageStore')}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {/* Products Card */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5 hover:shadow-md transition-shadow">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-lg font-semibold text-gray-900">{t('products')}</h3>
            <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
              <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
          </div>
          <p className="text-gray-600 text-sm mb-4">{t('manageProducts')}</p>
          <div className="flex flex-col gap-2">
              <Link
                to="/admin/products"
                className="text-blue-600 hover:text-blue-700 font-medium text-sm inline-block"
              >
                {t('viewAllProducts')}
              </Link>
          </div>
        </div>


        {/* Orders Card */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5 hover:shadow-md transition-shadow">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-lg font-semibold text-gray-900">{t('orders')}</h3>
            <div className="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
              <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
          </div>
          <p className="text-gray-600 text-sm mb-4">{t('trackOrders')}</p>
          <Link
            to="/admin/orders"
            className="text-green-600 hover:text-green-700 font-medium text-sm inline-block"
          >
            {t('viewOrders')}
          </Link>
        </div>

        {/* Users Card */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5 hover:shadow-md transition-shadow">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-lg font-semibold text-gray-900">{t('users')}</h3>
            <div className="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
              <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </div>
          </div>
          <p className="text-gray-600 text-sm mb-4">{t('manageUsers')}</p>
          <Link
            to="/admin/users"
            className="text-purple-600 hover:text-purple-700 font-medium text-sm inline-block"
          >
            {t('viewUsers')}
          </Link>
        </div>
      </div>
      
      {/* Recent Audit Logs - full width under the cards */}
      <div className="mt-6">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5 hover:shadow-md transition-shadow">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-lg font-semibold text-gray-900">{t('recentAuditLogs')}</h3>
            <div className="w-10 h-10 bg-yellow-100 rounded-lg flex items-center justify-center">
              <svg className="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
          <p className="text-gray-600 text-sm mb-4">{t('latestAuditEvents')}</p>
          <RecentAuditLogs />
          <div className="mt-4">
            <Link to="/admin/audit-logs" className="text-yellow-600 hover:text-yellow-700 font-medium text-sm inline-block">{t('viewAllAuditLogs')}</Link>
          </div>
        </div>
      </div>
    </AdminLayout>
  );
};

export default AdminDashboard;
