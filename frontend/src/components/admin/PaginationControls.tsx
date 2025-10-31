import React from 'react';

type Props = {
  page: number;
  totalPages: number;
  onPrev: () => void;
  onNext: () => void;
};

const PaginationControls: React.FC<Props> = ({ page, totalPages, onPrev, onNext }) => {
  return (
    <div className="mt-4 flex items-center gap-2">
      <button className="px-3 py-1 border rounded" onClick={onPrev} disabled={page === 1} aria-label="Previous page">Prev</button>
      <span>Page {page} of {totalPages}</span>
      <button className="px-3 py-1 border rounded" onClick={onNext} disabled={page === totalPages} aria-label="Next page">Next</button>
    </div>
  );
};

export default PaginationControls;
