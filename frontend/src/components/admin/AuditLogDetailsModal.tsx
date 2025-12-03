import React from 'react';
import type { AuditLog } from '../../types/audit';

function formatTimestamp(ts?: string | undefined | null) {
  if (!ts) return '-';
  const d = new Date(ts);
  if (Number.isNaN(d.getTime())) return ts;
  const formatted = d.toLocaleString(undefined, {
    year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
  });

  try {
    const diff = Date.now() - d.getTime();
    const seconds = Math.round(diff / 1000);
    const rtf = new Intl.RelativeTimeFormat(undefined, { numeric: 'auto' });
    if (seconds >= 86400) return `${formatted} · ${rtf.format(-Math.round(seconds / 86400), 'day')}`;
    if (seconds >= 3600) return `${formatted} · ${rtf.format(-Math.round(seconds / 3600), 'hour')}`;
    if (seconds >= 60) return `${formatted} · ${rtf.format(-Math.round(seconds / 60), 'minute')}`;
    return `${formatted} · ${rtf.format(-seconds, 'second')}`;
  } catch {
    return formatted;
  }
}

type Props = {
  open: boolean;
  onClose: () => void;
  log?: AuditLog | null;
};

const AuditLogDetailsModal: React.FC<Props> = ({ open, onClose, log }) => {
  if (!open || !log) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="bg-white rounded-lg w-full max-w-3xl p-6 shadow-lg">
        <div className="flex items-start justify-between gap-4">
          <h3 className="text-lg font-semibold">Audit Log Details</h3>
          <button onClick={onClose} className="text-sm text-gray-500">Close</button>
        </div>

        <div className="mt-4 space-y-3 text-sm text-gray-800">
          <div>
            <strong>Action:</strong> <span className="ml-2">{log.action}</span>
          </div>
          <div>
            <strong>Resource:</strong> <span className="ml-2">{log.resource} {log.resource_id ? `(${log.resource_id})` : ''}</span>
          </div>
          <div>
            <strong>Actor:</strong>
            <span className="ml-2">{log.actor?.display_name || log.actor?.public_id || log.actor?.email || '-'}</span>
            <div className="text-xs text-gray-500 mt-1">Actor: the entity that performed the action (user, service or system). This shows the actor's display name, public id or email when available.</div>
          </div>
          <div>
            <strong>Timestamp:</strong>
            <span className="ml-2">{formatTimestamp(log.timestamp || log.created_at)}</span>
          </div>

          <div>
            <strong>Details</strong>
            <div className="text-xs text-gray-500 mt-1">Additional contextual data captured for the event.</div>
            <pre className="mt-2 max-h-64 overflow-auto bg-gray-50 p-3 rounded text-xs text-gray-700">{JSON.stringify(log.details || {}, null, 2)}</pre>
          </div>

          <div>
            <strong>Old values</strong>
            <div className="text-xs text-gray-500 mt-1">Old values: the previous state of the resource before this action (may be empty for create events).</div>
            <pre className="mt-2 max-h-64 overflow-auto bg-gray-50 p-3 rounded text-xs text-gray-700">{JSON.stringify(log.old_values || {}, null, 2)}</pre>
          </div>

          <div>
            <strong>New values</strong>
            <div className="text-xs text-gray-500 mt-1">New values: the state of the resource after the action (useful for update events to see what changed).</div>
            <pre className="mt-2 max-h-64 overflow-auto bg-gray-50 p-3 rounded text-xs text-gray-700">{JSON.stringify(log.new_values || {}, null, 2)}</pre>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AuditLogDetailsModal;
