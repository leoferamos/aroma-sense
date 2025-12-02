import React, { useEffect, useState } from 'react';
import { getAdminUsers } from '../../services/admin';
import type { AdminUser } from '../../services/admin';

const Roles = ['admin', 'client'];
const Statuses = ['active', 'deactivated', 'deleted'];

const UsersPage: React.FC = () => {
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [limit, setLimit] = useState<number>(10);
  const [offset, setOffset] = useState<number>(0);
  const [total, setTotal] = useState<number>(0);
  const [role, setRole] = useState<string>('');
  const [status, setStatus] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchUsers = async () => {
    setLoading(true);
    setError(null);
    try {
      const resp = await getAdminUsers({ limit, offset, role: role || undefined, status: status || undefined });
      setUsers(resp.users || []);
      setTotal(resp.total || 0);
    } catch (err: unknown) {
      console.debug('getAdminUsers error', err);
      setError('Failed to load users');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [limit, offset, role, status]);

  const onPrev = () => setOffset(Math.max(0, offset - limit));
  const onNext = () => setOffset(offset + limit >= total ? offset : offset + limit);

  return (
    <div className="p-6">
      <h2 className="text-2xl font-semibold mb-4">Admin — Users</h2>

      <div className="mb-4 flex gap-3 items-center">
        <label className="text-sm">Role:</label>
        <select value={role} onChange={(e) => { setOffset(0); setRole(e.target.value); }} className="border px-2 py-1 rounded">
          <option value="">All</option>
          {Roles.map(r => <option key={r} value={r}>{r}</option>)}
        </select>

        <label className="text-sm">Status:</label>
        <select value={status} onChange={(e) => { setOffset(0); setStatus(e.target.value); }} className="border px-2 py-1 rounded">
          <option value="">All</option>
          {Statuses.map(s => <option key={s} value={s}>{s}</option>)}
        </select>

        <label className="text-sm">Per page:</label>
        <select value={limit} onChange={(e) => { setOffset(0); setLimit(Number(e.target.value)); }} className="border px-2 py-1 rounded">
          {[5,10,20,50].map(n => <option key={n} value={n}>{n}</option>)}
        </select>
      </div>

      {error && <div className="text-red-600 mb-3">{error}</div>}

      <div className="overflow-x-auto bg-white border rounded">
        <table className="min-w-full text-left text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-2">ID</th>
              <th className="px-4 py-2">Public ID</th>
              <th className="px-4 py-2">Display</th>
              <th className="px-4 py-2">Email</th>
              <th className="px-4 py-2">Role</th>
              <th className="px-4 py-2">Status</th>
              <th className="px-4 py-2">Created</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr><td colSpan={7} className="px-4 py-4">Loading...</td></tr>
            ) : users.length === 0 ? (
              <tr><td colSpan={7} className="px-4 py-4">No users found</td></tr>
            ) : users.map(u => (
              <tr key={u.id} className="border-t">
                <td className="px-4 py-2">{u.id}</td>
                <td className="px-4 py-2">{u.public_id}</td>
                <td className="px-4 py-2">{u.display_name || '-'}</td>
                <td className="px-4 py-2">{u.masked_email || u.email}</td>
                <td className="px-4 py-2">{u.role}</td>
                <td className="px-4 py-2">{u.deactivated_at ? 'deactivated' : 'active'}</td>
                <td className="px-4 py-2">{u.created_at ? new Date(u.created_at).toLocaleString() : '-'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="mt-4 flex items-center justify-between">
        <div className="text-sm text-gray-600">Showing {Math.min(total, offset + 1)} - {Math.min(total, offset + limit)} of {total}</div>
        <div className="flex gap-2">
          <button onClick={onPrev} disabled={offset === 0} className="px-3 py-1 border rounded disabled:opacity-50">Previous</button>
          <button onClick={onNext} disabled={offset + limit >= total} className="px-3 py-1 border rounded disabled:opacity-50">Next</button>
        </div>
      </div>
    </div>
  );
};

export default UsersPage;
