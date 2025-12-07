import React from 'react';
import { useTranslation } from 'react-i18next';
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
  const { t } = useTranslation('admin');

  if (!open || !log) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center backdrop-blur-sm p-4 overflow-y-auto">
      <div className="bg-white rounded-lg w-full max-w-3xl p-4 sm:p-6 shadow-lg my-8 max-h-[90vh] overflow-y-auto">
        <div className="flex items-start justify-between gap-4 mb-4">
          <h3 className="text-lg font-semibold">{t('auditLogDetails')}</h3>
          <button onClick={onClose} className="text-sm text-gray-500">{t('close')}</button>
        </div>

        <div className="mt-4 space-y-3 text-sm text-gray-800">
          <div>
            <strong>{t('action')}:</strong> <span className="ml-2">{log.action}</span>
          </div>
          <div>
            <strong>{t('resource')}:</strong> <span className="ml-2">{log.resource} {log.resource_id ? `(${log.resource_id})` : ''}</span>
          </div>
          <div>
            <strong>{t('actor')}:</strong>
            <span className="ml-2">{log.actor?.display_name || log.actor?.public_id || log.actor?.email || '-'}</span>
            <div className="text-xs text-gray-500 mt-1">{t('actorDescription')}</div>
          </div>
          <div>
            <strong>{t('timestamp')}:</strong>
            <span className="ml-2">{formatTimestamp(log.timestamp || log.created_at)}</span>
          </div>

          <div>
            <strong>{t('details')}</strong>
            <div className="text-xs text-gray-500 mt-1">{t('detailsDescription')}</div>
            <pre className="mt-2 max-h-64 overflow-auto bg-gray-50 p-3 rounded text-xs text-gray-700">{JSON.stringify(log.details || {}, null, 2)}</pre>
          </div>

          <div>
            <strong>{t('oldValues')}</strong>
            <div className="text-xs text-gray-500 mt-1">{t('oldValuesDescription')}</div>
            <pre className="mt-2 max-h-64 overflow-auto bg-gray-50 p-3 rounded text-xs text-gray-700">{JSON.stringify(log.old_values || {}, null, 2)}</pre>
          </div>

          <div>
            <strong>{t('newValues')}</strong>
            <div className="text-xs text-gray-500 mt-1">{t('newValuesDescription')}</div>
            <pre className="mt-2 max-h-64 overflow-auto bg-gray-50 p-3 rounded text-xs text-gray-700">{JSON.stringify(log.new_values || {}, null, 2)}</pre>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AuditLogDetailsModal;
