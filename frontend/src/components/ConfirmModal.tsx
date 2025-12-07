import React from 'react';
import { useTranslation } from 'react-i18next';

interface ConfirmModalProps {
  open: boolean;
  title?: string;
  description?: string;
  confirmText?: string;
  cancelText?: string;
  requirePhrase?: string;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
  error?: string | null;
}

const ConfirmModal: React.FC<ConfirmModalProps> = ({
  open,
  title = 'Are you sure?',
  description = 'This action cannot be undone.',
  confirmText = 'Delete',
  cancelText = 'Cancel',
  requirePhrase,
  onConfirm,
  onCancel,
  loading = false,
  error = null,
}) => {
  const { t } = useTranslation('common');
  const [phrase, setPhrase] = React.useState('');

  React.useEffect(() => {
    if (!open) setPhrase('');
  }, [open]);
  if (!open) return null;
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center backdrop-blur-sm" onClick={onCancel}>
      <div className="bg-white rounded-lg shadow-lg max-w-sm w-full p-6 relative animate-fade-in" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-lg font-semibold text-gray-900 mb-2">{title}</h2>
        <p className="text-gray-700 mb-6 text-sm">{description}</p>
        {error && (
          <div className="mb-4 rounded-md bg-red-50 p-3 text-red-700 border border-red-200 text-sm">
            {error}
          </div>
        )}
        {requirePhrase ? (
          <div className="mb-4">
            <p className="text-sm text-gray-600 mb-2" dangerouslySetInnerHTML={{ __html: t('confirmTypePhrase', { phrase: requirePhrase }) }} />
            <input
              value={phrase}
              onChange={(e) => setPhrase(e.target.value)}
              className="w-full px-3 py-2 border rounded-md"
              placeholder={requirePhrase}
            />
          </div>
        ) : null}

        <div className="flex justify-end gap-2">
          <button
            type="button"
            className="px-4 py-2 rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200 font-medium disabled:opacity-60 disabled:cursor-not-allowed"
            onClick={onCancel}
            disabled={loading}
          >
            {cancelText}
          </button>
          <button
            type="button"
            className="px-4 py-2 rounded-md bg-red-600 text-white hover:bg-red-700 font-medium disabled:opacity-60 disabled:cursor-not-allowed"
            onClick={onConfirm}
            disabled={loading || (requirePhrase ? phrase !== requirePhrase : false)}
          >
            {loading ? t('deleting') : confirmText}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ConfirmModal;
