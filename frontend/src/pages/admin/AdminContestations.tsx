import React, { useEffect, useState } from 'react';
import { getPendingContestations, approveContestation, rejectContestation } from '../../services/adminContestations';
import type { AdminContestation } from '../../services/adminContestations';
import AdminLayout from '../../components/admin/AdminLayout';
import { Link } from 'react-router-dom';

const navItems = [
  { label: 'Dashboard', to: '/admin/dashboard' },
  { label: 'Products', to: '/admin/products' },
  { label: 'Orders', to: '/admin/orders' },
  { label: 'Users', to: '/admin/users' },
  { label: 'Audit Logs', to: '/admin/audit-logs' },
  { label: 'Contestations', to: '/admin/contestations' },
];

const AdminContestations: React.FC = () => {
  const [contestations, setContestations] = useState<AdminContestation[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState<number | null>(null);
  const [reviewNotes, setReviewNotes] = useState<{ [id: number]: string }>({});

  const fetchContestations = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await getPendingContestations();
      setContestations(res.data);
    } catch {
      setError('Failed to load contestations');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchContestations();
  }, []);

  const handleApprove = async (id: number) => {
    setActionLoading(id);
    setError(null);
    try {
      await approveContestation(id, reviewNotes[id]);
      setContestations(contestations.filter(c => c.id !== id));
    } catch {
      setError('Failed to approve contestation');
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
      setError('Failed to reject contestation');
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <AdminLayout
      title="Contestations"
      navItems={navItems}
      actions={<div className="flex items-center gap-2"><Link to="/admin/dashboard" className="px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-200 rounded-lg hover:bg-gray-50">← Dashboard</Link></div>}
    >
      <div className="p-6">
        <h2 className="text-2xl font-bold mb-4">Pending User Contestations</h2>
        {error && <div className="text-red-500 mb-2">{error}</div>}
        {loading ? (
          <div>Loading...</div>
        ) : contestations.length === 0 ? (
          <div>No pending contestations.</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full border bg-white rounded shadow">
              <thead className="bg-gray-100">
                <tr>
                  <th className="border px-3 py-2 text-left">ID</th>
                  <th className="border px-3 py-2 text-left">User ID</th>
                  <th className="border px-3 py-2 text-left">Reason</th>
                  <th className="border px-3 py-2 text-left">Requested At</th>
                  <th className="border px-3 py-2 text-left">Review Notes</th>
                  <th className="border px-3 py-2 text-left">Actions</th>
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
                        placeholder="Review notes"
                        value={reviewNotes[c.id] || ''}
                        onChange={e => setReviewNotes({ ...reviewNotes, [c.id]: e.target.value })}
                      />
                    </td>
                    <td className="border px-3 py-2">
                      <button
                        className="bg-green-600 text-white px-3 py-1 mr-2 rounded hover:bg-green-700 disabled:opacity-50"
                        disabled={actionLoading === c.id}
                        onClick={() => handleApprove(c.id)}
                      >Approve</button>
                      <button
                        className="bg-red-600 text-white px-3 py-1 rounded hover:bg-red-700 disabled:opacity-50"
                        disabled={actionLoading === c.id}
                        onClick={() => handleReject(c.id)}
                      >Reject</button>
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
