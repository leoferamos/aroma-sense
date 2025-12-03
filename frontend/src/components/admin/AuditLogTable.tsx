import React from 'react';
import type { AuditLog } from '../../types/audit';

type Props = {
  logs: AuditLog[];
  onView: (log: AuditLog) => void;
};

const AuditLogTable: React.FC<Props> = ({ logs, onView }) => {
  return (
    <div className="overflow-x-auto bg-white border rounded">
      <table className="min-w-full text-left text-sm">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-4 py-2">Timestamp</th>
            <th className="px-4 py-2">Action</th>
            <th className="px-4 py-2">Resource</th>
            <th className="px-4 py-2">Actor</th>
            <th className="px-4 py-2">Severity</th>
            <th className="px-4 py-2">Details</th>
          </tr>
        </thead>
        <tbody>
          {logs.length === 0 ? (
            <tr><td colSpan={6} className="px-4 py-4">No logs</td></tr>
          ) : logs.map((l) => (
            <tr key={l.id} className="border-t">
              <td className="px-4 py-2 align-top text-xs text-gray-600">{new Date(l.timestamp || l.created_at || '').toLocaleString()}</td>
              <td className="px-4 py-2 font-medium">{l.action}</td>
              <td className="px-4 py-2">{l.resource}{l.resource_id ? ` (${l.resource_id})` : ''}</td>
              <td className="px-4 py-2">{l.actor?.display_name || l.user?.display_name || l.actor?.public_id || '-'}</td>
              <td className="px-4 py-2"><span className={`px-2 py-0.5 text-xs rounded ${l.severity === 'high' ? 'bg-red-100 text-red-700' : 'bg-gray-100 text-gray-700'}`}>{l.severity || 'info'}</span></td>
              <td className="px-4 py-2">
                <button onClick={() => onView(l)} className="text-blue-600 hover:underline text-sm">View</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AuditLogTable;
