import React from 'react';
import { useTranslation } from 'react-i18next';
import type { AuditLog } from '../../types/audit';

type Props = {
  logs: AuditLog[];
  onView: (log: AuditLog) => void;
};

const AuditLogTable: React.FC<Props> = ({ logs, onView }) => {
  const { t } = useTranslation('admin');

  return (
    <>
      {/* Desktop Table View - Hidden on mobile */}
      <div className="hidden lg:block overflow-x-auto bg-white border border-gray-200 rounded-lg shadow-sm">
        <table className="min-w-full text-left text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-4 py-3 font-semibold text-gray-700">{t('timestamp')}</th>
              <th className="px-4 py-3 font-semibold text-gray-700">{t('action')}</th>
              <th className="px-4 py-3 font-semibold text-gray-700">{t('resource')}</th>
              <th className="px-4 py-3 font-semibold text-gray-700">{t('actor')}</th>
              <th className="px-4 py-3 font-semibold text-gray-700">{t('severity')}</th>
              <th className="px-4 py-3 font-semibold text-gray-700">{t('view')}</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {logs.length === 0 ? (
              <tr><td colSpan={6} className="px-4 py-8 text-center text-gray-500">{t('noLogs')}</td></tr>
            ) : logs.map((l) => (
              <tr key={l.id} className="hover:bg-gray-50 transition-colors">
                <td className="px-4 py-3 text-xs text-gray-600 whitespace-nowrap">{new Date(l.timestamp || l.created_at || '').toLocaleString()}</td>
                <td className="px-4 py-3 font-medium text-gray-900">{l.action}</td>
                <td className="px-4 py-3 text-gray-700">{l.resource}{l.resource_id ? ` (${l.resource_id})` : ''}</td>
                <td className="px-4 py-3 text-gray-600">{l.actor?.display_name || l.user?.display_name || l.actor?.public_id || '-'}</td>
                <td className="px-4 py-3">
                  <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${l.severity === 'high' ? 'bg-red-100 text-red-700' :
                      l.severity === 'medium' ? 'bg-yellow-100 text-yellow-700' :
                        'bg-gray-100 text-gray-700'
                    }`}>
                    {l.severity || 'info'}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <button onClick={() => onView(l)} className="text-blue-600 hover:text-blue-800 font-medium text-sm transition-colors">{t('view')}</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Mobile Card View - Visible only on mobile/tablet */}
      <div className="lg:hidden space-y-3">
        {logs.length === 0 ? (
          <div className="bg-white border border-gray-200 rounded-lg p-6 text-center text-gray-500">{t('noLogs')}</div>
        ) : logs.map((l) => (
          <div key={l.id} className="bg-white border border-gray-200 rounded-lg p-4 shadow-sm">
            <div className="flex items-start justify-between mb-3">
              <div className="flex-1">
                <div className="font-semibold text-gray-900 mb-1">{l.action}</div>
                <div className="text-xs text-gray-500">{new Date(l.timestamp || l.created_at || '').toLocaleString()}</div>
              </div>
              <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full whitespace-nowrap ${l.severity === 'high' ? 'bg-red-100 text-red-700' :
                  l.severity === 'medium' ? 'bg-yellow-100 text-yellow-700' :
                    'bg-gray-100 text-gray-700'
                }`}>
                {l.severity || 'info'}
              </span>
            </div>
            <div className="border-t border-gray-100 pt-3 space-y-2">
              <div className="text-sm">
                <span className="text-gray-500 font-medium">{t('resource')}:</span>
                <span className="text-gray-900 ml-2">{l.resource}{l.resource_id ? ` (${l.resource_id})` : ''}</span>
              </div>
              <div className="text-sm">
                <span className="text-gray-500 font-medium">{t('actor')}:</span>
                <span className="text-gray-900 ml-2">{l.actor?.display_name || l.user?.display_name || l.actor?.public_id || '-'}</span>
              </div>
              <div className="mt-3">
                <button onClick={() => onView(l)} className="w-full px-4 py-2 text-sm font-medium text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors">{t('view')}</button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </>
  );
};

export default AuditLogTable;
