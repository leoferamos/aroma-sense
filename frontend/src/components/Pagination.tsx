import React from "react";

type Props = {
  page: number;
  pageSize: number; // items per page
  total: number; // total items
  onPageChange: (page: number) => void;
  className?: string;
};

const Pagination: React.FC<Props> = ({ page, pageSize, total, onPageChange, className = "" }) => {
  const totalPages = Math.max(1, Math.ceil(total / Math.max(1, pageSize)));
  const canPrev = page > 1;
  const canNext = page < totalPages;

  const goTo = (p: number) => () => {
    if (p >= 1 && p <= totalPages) onPageChange(p);
  };
  const windowSize = 5;
  const half = Math.floor(windowSize / 2);
  const start = Math.max(1, Math.min(page - half, totalPages - windowSize + 1));
  const end = Math.min(totalPages, start + windowSize - 1);

  return (
    <nav className={`mt-6 flex items-center justify-center gap-2 ${className}`} aria-label="Pagination">
      <button
        className="px-3 py-1 border rounded disabled:opacity-50"
        onClick={goTo(page - 1)}
        disabled={!canPrev}
        aria-label="Previous page"
      >
        Previous
      </button>
      {Array.from({ length: end - start + 1 }, (_, i) => start + i).map((p) => (
        <button
          key={p}
          className={`px-3 py-1 border rounded ${p === page ? "bg-gray-900 text-white" : ""}`}
          onClick={goTo(p)}
          aria-current={p === page ? "page" : undefined}
        >
          {p}
        </button>
      ))}
      <button
        className="px-3 py-1 border rounded disabled:opacity-50"
        onClick={goTo(page + 1)}
        disabled={!canNext}
        aria-label="Next page"
      >
        Next
      </button>
    </nav>
  );
};

export default Pagination;
