import React from 'react';
import { useTranslation } from 'react-i18next';

type Props = {
  page: number;
  totalPages: number;
  onPrev: () => void;
  onNext: () => void;
};

const PaginationControls: React.FC<Props> = ({ page, totalPages, onPrev, onNext }) => {
  const { t } = useTranslation('admin');
  return (
    <div className="mt-4 flex items-center gap-2">
      <button className="px-3 py-1 border rounded" onClick={onPrev} disabled={page === 1} aria-label="Previous page">{t('prev')}</button>
      <span>{t('pageOf', { page, total: totalPages })}</span>
      <button className="px-3 py-1 border rounded" onClick={onNext} disabled={page === totalPages} aria-label="Next page">{t('next')}</button>
    </div>
  );
};

export default PaginationControls;
