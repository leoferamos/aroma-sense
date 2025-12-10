import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getAdminUsers, createAdminUser } from '../../services/admin';
import type { AdminUser } from '../../services/admin';
import AdminLayout from '../../components/admin/AdminLayout';
import PaginationControls from '../../components/admin/PaginationControls';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../../hooks/useAuth';

const Roles = ['admin', 'client'];
const Statuses = ['active', 'deactivated', 'deleted'];

const UsersPage: React.FC = () => {
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [limit, setLimit] = useState<number>(10);
  const [offset, setOffset] = useState<number>(0);
  const [total, setTotal] = useState<number>(0);
  const currentPage = Math.floor(offset / limit) + 1;
  const totalPages = Math.max(1, Math.ceil(total / limit));
  const [filterRole, setFilterRole] = useState<string>('');
  const [status, setStatus] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { t } = useTranslation('admin');
  const { role: authRole } = useAuth();

  const [newAdmin, setNewAdmin] = useState({ name: '', email: '', password: '' });
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState<string | null>(null);
  const [createSuccess, setCreateSuccess] = useState<string | null>(null);

  const fetchUsers = async () => {
    setLoading(true);
    setError(null);
    try {
      const resp = await getAdminUsers({ limit, offset, role: filterRole || undefined, status: status || undefined });
      setUsers(resp.users || []);
      setTotal(resp.total || 0);
    } catch (err: unknown) {
      console.debug('getAdminUsers error', err);
      setError(t('failedToLoadUsers'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [limit, offset, filterRole, status]);

  const handleCreateAdmin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newAdmin.name || !newAdmin.email || !newAdmin.password) {
      setCreateError(t('adminCreateRequiredFields'));
      return;
    }
    setCreating(true);
    setCreateError(null);
    setCreateSuccess(null);
    try {
      await createAdminUser({
        name: newAdmin.name,
        email: newAdmin.email,
        password: newAdmin.password,
      });
      setCreateSuccess(t('adminCreateSuccess'));
      setNewAdmin({ name: '', email: '', password: '' });
      fetchUsers();
    } catch (err) {
      console.error('createAdminUser error', err);
      setCreateError(t('adminCreateFailure'));
    } finally {
      setCreating(false);
    }
  };

  return (
    <AdminLayout actions={<div className="flex items-center gap-2"><Link to="/admin/dashboard" className="px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-lg hover:bg-gray-50">← {t('nav.dashboard')}</Link></div>}>
      <div className="p-4 sm:p-6 space-y-6">
        <div className="flex flex-col gap-2">
          <h1 className="text-xl sm:text-2xl font-semibold">{t('users')}</h1>
          <p className="text-sm text-gray-600">{t('usersDescription')}</p>
        </div>

        {authRole === 'super_admin' && (
          <form onSubmit={handleCreateAdmin} className="bg-white border border-gray-200 rounded-lg p-4 shadow-sm space-y-3">
            <div className="flex flex-col sm:flex-row sm:items-end gap-3">
              <div className="flex-1">
                <label className="text-sm font-medium text-gray-700">{t('adminNameLabel')}</label>
                <input
                  type="text"
                  className="mt-1 w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  value={newAdmin.name}
                  onChange={(e) => setNewAdmin({ ...newAdmin, name: e.target.value })}
                  placeholder={t('adminNamePlaceholder')}
                />
              </div>
              <div className="flex-1">
                <label className="text-sm font-medium text-gray-700">{t('adminEmailLabel')}</label>
                <input
                  type="email"
                  className="mt-1 w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  value={newAdmin.email}
                  onChange={(e) => setNewAdmin({ ...newAdmin, email: e.target.value })}
                  placeholder={t('adminEmailPlaceholder')}
                />
              </div>
              <div className="flex-1">
                <label className="text-sm font-medium text-gray-700">{t('adminPasswordLabel')}</label>
                <input
                  type="password"
                  className="mt-1 w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  value={newAdmin.password}
                  onChange={(e) => setNewAdmin({ ...newAdmin, password: e.target.value })}
                  placeholder={t('adminPasswordPlaceholder')}
                />
              </div>
              <div className="sm:self-end">
                <button
                  type="submit"
                  disabled={creating}
                  className="inline-flex items-center justify-center px-4 py-2 rounded-lg bg-blue-600 text-white text-sm font-semibold shadow-sm hover:bg-blue-700 disabled:opacity-60"
                >
                  {creating ? t('creatingAdmin') : t('createAdminTitle')}
                </button>
              </div>
            </div>
            {createError && <div className="text-sm text-red-600">{createError}</div>}
            {createSuccess && <div className="text-sm text-green-600">{createSuccess}</div>}
          </form>
        )}

        {/* Filters - Responsive Grid */}
        <div className="mb-4 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
          <div className="flex flex-col">
            <label className="text-sm font-medium text-gray-700 mb-1">{t('role')}:</label>
            <select value={filterRole} onChange={(e) => { setOffset(0); setFilterRole(e.target.value); }} className="border border-gray-300 px-3 py-2 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
              <option value="">{t('all')}</option>
              {Roles.map(r => <option key={r} value={r}>{r}</option>)}
            </select>
          </div>

          <div className="flex flex-col">
            <label className="text-sm font-medium text-gray-700 mb-1">{t('status')}:</label>
            <select value={status} onChange={(e) => { setOffset(0); setStatus(e.target.value); }} className="border border-gray-300 px-3 py-2 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
              <option value="">{t('all')}</option>
              {Statuses.map(s => <option key={s} value={s}>{s}</option>)}
            </select>
          </div>

          <div className="flex flex-col">
            <label className="text-sm font-medium text-gray-700 mb-1">{t('perPage')}:</label>
            <select value={limit} onChange={(e) => { setOffset(0); setLimit(Number(e.target.value)); }} className="border border-gray-300 px-3 py-2 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
              {[5, 10, 20, 50].map(n => <option key={n} value={n}>{n}</option>)}
            </select>
          </div>
        </div>

        {error && <div className="text-red-600 mb-3 p-3 bg-red-50 border border-red-200 rounded-lg text-sm">{error}</div>}

        {/* Desktop Table View - Hidden on mobile */}
        <div className="hidden lg:block overflow-x-auto bg-white border border-gray-200 rounded-lg shadow-sm">
          <table className="min-w-full text-left text-sm">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-4 py-3 font-semibold text-gray-700">{t('id')}</th>
                <th className="px-4 py-3 font-semibold text-gray-700">{t('publicId')}</th>
                <th className="px-4 py-3 font-semibold text-gray-700">{t('display')}</th>
                <th className="px-4 py-3 font-semibold text-gray-700">{t('email')}</th>
                <th className="px-4 py-3 font-semibold text-gray-700">{t('role')}</th>
                <th className="px-4 py-3 font-semibold text-gray-700">{t('status')}</th>
                <th className="px-4 py-3 font-semibold text-gray-700">{t('created')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {loading ? (
                <tr><td colSpan={7} className="px-4 py-8 text-center text-gray-500">{t('loadingUsers')}</td></tr>
              ) : users.length === 0 ? (
                <tr><td colSpan={7} className="px-4 py-8 text-center text-gray-500">{t('noUsersFound')}</td></tr>
              ) : users.map(u => (
                <tr key={u.id} className="hover:bg-gray-50 transition-colors">
                  <td className="px-4 py-3 text-gray-900">{u.id}</td>
                  <td className="px-4 py-3 text-gray-900 font-mono text-xs">{u.public_id}</td>
                  <td className="px-4 py-3 text-gray-900">{u.display_name || '-'}</td>
                  <td className="px-4 py-3 text-gray-600">{u.masked_email}</td>
                  <td className="px-4 py-3">
                    <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${u.role === 'admin' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'}`}>
                      {u.role}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${u.deactivated_at ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
                      {u.deactivated_at ? 'deactivated' : 'active'}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-gray-600 text-xs">{u.created_at ? new Date(u.created_at).toLocaleString() : '-'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Mobile Card View - Visible only on mobile/tablet */}
        <div className="lg:hidden space-y-3">
          {loading ? (
            <div className="bg-white border border-gray-200 rounded-lg p-6 text-center text-gray-500">{t('loadingUsers')}</div>
          ) : users.length === 0 ? (
            <div className="bg-white border border-gray-200 rounded-lg p-6 text-center text-gray-500">{t('noUsersFound')}</div>
          ) : users.map(u => (
            <div key={u.id} className="bg-white border border-gray-200 rounded-lg p-4 shadow-sm">
              <div className="flex items-start justify-between mb-3">
                <div>
                  <div className="text-xs text-gray-500 mb-1">{t('id')}: {u.id}</div>
                  <div className="font-medium text-gray-900">{u.display_name || 'No name'}</div>
                  <div className="text-sm text-gray-600 mt-1">{u.masked_email}</div>
                </div>
                <div className="flex flex-col gap-1 items-end">
                  <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${u.role === 'admin' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'}`}>
                    {u.role}
                  </span>
                  <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${u.deactivated_at ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
                    {u.deactivated_at ? 'deactivated' : 'active'}
                  </span>
                </div>
              </div>
              <div className="border-t border-gray-100 pt-3 space-y-1">
                <div className="text-xs text-gray-500">
                  <span className="font-medium">{t('publicId')}:</span> <span className="font-mono">{u.public_id}</span>
                </div>
                <div className="text-xs text-gray-500">
                  <span className="font-medium">{t('created')}:</span> {u.created_at ? new Date(u.created_at).toLocaleString() : '-'}
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="mt-4 flex flex-col sm:flex-row items-center justify-between gap-3">
          <div className="text-xs sm:text-sm text-gray-600">{t('showingRange', { from: Math.min(total, offset + 1), to: Math.min(total, offset + limit), total })}</div>
          <PaginationControls
            page={currentPage}
            totalPages={totalPages}
            onPrev={() => { if (currentPage > 1) { setOffset(Math.max(0, offset - limit)); } }}
            onNext={() => { if (currentPage < totalPages) { setOffset(offset + limit); } }}
          />
        </div>
      </div>
    </AdminLayout>
  );
};

export default UsersPage;
