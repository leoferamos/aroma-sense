import React from 'react';

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
    <div className="mb-4 flex flex-col md:flex-row md:items-center md:gap-4">
      <div className="flex items-center gap-2">
        <label className="text-sm">Status</label>
        <select value={status} onChange={(e) => onStatusChange(e.target.value)} className="border rounded px-2 py-1">
          <option value="">All</option>
          <option value="pending">pending</option>
          <option value="processing">processing</option>
          <option value="shipped">shipped</option>
          <option value="delivered">delivered</option>
          <option value="cancelled">cancelled</option>
        </select>
      </div>

      <div className="flex items-center gap-2 mt-2 md:mt-0">
        <label className="text-sm">From</label>
        <input type="date" value={startDate || ''} onChange={(e) => onDateChange(e.target.value || undefined, endDate)} className="border rounded px-2 py-1" />
        <label className="text-sm">To</label>
        <input type="date" value={endDate || ''} onChange={(e) => onDateChange(startDate, e.target.value || undefined)} className="border rounded px-2 py-1" />
      </div>

      <div className="flex items-center gap-2 mt-2 md:mt-0 md:ml-auto">
        <label className="text-sm">Per page</label>
        <select value={perPage} onChange={(e) => onPerPageChange(Number(e.target.value))} className="border rounded px-2 py-1">
          {perPageOptions.map((o) => (
            <option key={o} value={o}>{o}</option>
          ))}
        </select>
      
        <div className="flex items-center gap-2 ml-4">
          <button type="button" onClick={clearFilters} className="text-sm px-3 py-1 border rounded hover:bg-gray-50">Clear</button>

          <div className="flex items-center gap-2">
            <button type="button" onClick={() => applyPreset(1)} className="text-sm px-2 py-1 border rounded hover:bg-gray-50">Today</button>
            <button type="button" onClick={() => applyPreset(7)} className="text-sm px-2 py-1 border rounded hover:bg-gray-50">7d</button>
            <button type="button" onClick={() => applyPreset(30)} className="text-sm px-2 py-1 border rounded hover:bg-gray-50">30d</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default FiltersBar;
