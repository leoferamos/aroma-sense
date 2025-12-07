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
    <div className="fixed inset-0 z-50 flex items-center justify-center backdrop-blur-sm">
      <div className="bg-white rounded-lg w-full max-w-3xl p-6 shadow-lg">
        <div className="flex items-start justify-between gap-4">
          <h3 className="text-lg font-semibold">{t('auditLogDetails')}</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors p-1"
            aria-label="Close modal"
          >
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-5 h-5">
              <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="space-y-4 text-sm text-gray-800">
          <div className="p-3 bg-gray-50 rounded-lg">
            <strong className="text-gray-900">{t('action')}:</strong>
            <span className="ml-2 font-medium">{log.action}</span>
          </div>

          <div className="p-3 bg-gray-50 rounded-lg">
            <strong className="text-gray-900">{t('resource')}:</strong>
            <span className="ml-2">{log.resource} {log.resource_id ? `(${log.resource_id})` : ''}</span>
          </div>

          <div className="p-3 bg-gray-50 rounded-lg">
            <strong className="text-gray-900">{t('actor')}:</strong>
            <span className="ml-2">{log.actor?.display_name || log.actor?.public_id || log.actor?.email || '-'}</span>
            <div className="text-xs text-gray-500 mt-1">{t('actorDescription')}</div>
          </div>

          <div className="p-3 bg-gray-50 rounded-lg">
            <strong className="text-gray-900">{t('timestamp')}:</strong>
            <span className="ml-2">{formatTimestamp(log.timestamp || log.created_at)}</span>
          </div>

          <div>
            <strong className="text-gray-900">{t('details')}</strong>
            <div className="text-xs text-gray-500 mt-1 mb-2">{t('detailsDescription')}</div>
            <pre className="max-h-48 overflow-auto bg-gray-50 p-3 rounded-lg text-xs text-gray-700 border border-gray-200">{JSON.stringify(log.details || {}, null, 2)}</pre>
          </div>

          <div>
            <strong className="text-gray-900">{t('oldValues')}</strong>
            <div className="text-xs text-gray-500 mt-1 mb-2">{t('oldValuesDescription')}</div>
            <pre className="max-h-48 overflow-auto bg-gray-50 p-3 rounded-lg text-xs text-gray-700 border border-gray-200">{JSON.stringify(log.old_values || {}, null, 2)}</pre>
          </div>

          <div>
            <strong className="text-gray-900">{t('newValues')}</strong>
            <div className="text-xs text-gray-500 mt-1 mb-2">{t('newValuesDescription')}</div>
            <pre className="max-h-48 overflow-auto bg-gray-50 p-3 rounded-lg text-xs text-gray-700 border border-gray-200">{JSON.stringify(log.new_values || {}, null, 2)}</pre>
          </div>
        </div>

        <div className="mt-6 flex justify-end">
          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium rounded-lg transition-colors"
          >
            {t('close')}
          </button>
        </div>
      </div>
    </div>
  );
};

export default AuditLogDetailsModal;
