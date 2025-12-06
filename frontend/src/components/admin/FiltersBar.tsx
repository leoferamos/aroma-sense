import React from 'react';
import { useTranslation } from 'react-i18next';

type Props = {
  status: string;
  onStatusChange: (s: string) => void;
  startDate?: string;
  endDate?: string;
  onDateChange: (start?: string, end?: string) => void;
  perPage: number;
  onPerPageChange: (n: number) => void;
};

const perPageOptions = [10, 25, 50, 100];

const FiltersBar: React.FC<Props> = ({ status, onStatusChange, startDate, endDate, onDateChange, perPage, onPerPageChange }) => {
  const { t } = useTranslation('admin');

  const formatDate = (d: Date) => {
    const yyyy = d.getFullYear();
    const mm = String(d.getMonth() + 1).padStart(2, '0');
    const dd = String(d.getDate()).padStart(2, '0');
    return `${yyyy}-${mm}-${dd}`;
  };

  const applyPreset = (days?: number) => {
    if (!days) {
      // today
      const today = new Date();
      onDateChange(formatDate(today), formatDate(today));
      return;
    }
    const end = new Date();
    const start = new Date();
    start.setDate(end.getDate() - (days - 1));
    onDateChange(formatDate(start), formatDate(end));
  };

  const clearFilters = () => {
    onStatusChange('');
    onDateChange(undefined, undefined);
  };

  return (
    <div className="mb-4 space-y-3">
      {/* First Row - Status and Date Range */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        <div className="flex flex-col">
          <label className="text-xs font-medium text-gray-700 mb-1">{t('status')}</label>
          <select value={status} onChange={(e) => onStatusChange(e.target.value)} className="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
            <option value="">{t('all')}</option>
            <option value="pending">{t('pending')}</option>
            <option value="processing">{t('processing')}</option>
            <option value="shipped">{t('shipped')}</option>
            <option value="delivered">{t('delivered')}</option>
            <option value="cancelled">{t('cancelled')}</option>
          </select>
        </div>

        <div className="flex flex-col">
          <label className="text-xs font-medium text-gray-700 mb-1">{t('from')}</label>
          <input type="date" value={startDate || ''} onChange={(e) => onDateChange(e.target.value || undefined, endDate)} className="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500" />
        </div>

        <div className="flex flex-col">
          <label className="text-xs font-medium text-gray-700 mb-1">{t('to')}</label>
          <input type="date" value={endDate || ''} onChange={(e) => onDateChange(startDate, e.target.value || undefined)} className="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500" />
        </div>
      </div>

      {/* Second Row - Per Page and Actions */}
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
        <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-2">
          <label className="text-xs sm:text-sm font-medium text-gray-700">{t('perPage')}</label>
          <select value={perPage} onChange={(e) => onPerPageChange(Number(e.target.value))} className="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
            {perPageOptions.map((o) => (
              <option key={o} value={o}>{o}</option>
            ))}
          </select>
        </div>

        <div className="flex flex-col sm:flex-row items-stretch gap-2 sm:ml-auto">
          <button type="button" onClick={clearFilters} className="text-sm px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 font-medium transition-colors">{t('clear')}</button>

          <div className="grid grid-cols-3 gap-2">
            <button type="button" onClick={() => applyPreset(1)} className="text-xs sm:text-sm px-3 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 font-medium transition-colors">{t('today')}</button>
            <button type="button" onClick={() => applyPreset(7)} className="text-xs sm:text-sm px-3 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 font-medium transition-colors">{t('7days')}</button>
            <button type="button" onClick={() => applyPreset(30)} className="text-xs sm:text-sm px-3 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 font-medium transition-colors">{t('30days')}</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default FiltersBar;
