import React, { useEffect, useState } from 'react';
import { cn } from '../utils/cn';
import { useReviews } from '../hooks/useReviews';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../hooks/useAuth';
import Toast from './Toast';
import ConfirmModal from './ConfirmModal';
import { reportReview, type ReportReviewRequest } from '../services/review';
import { isAxiosError } from 'axios';

interface ProductReviewProps {
    productSlug: string;
    canReview?: boolean;
}

const ProductReview: React.FC<ProductReviewProps> = ({ productSlug, canReview }) => {
    const [rating, setRating] = useState<number>(0);
    const [comment, setComment] = useState<string>('');
    const [hoverRating, setHoverRating] = useState<number>(0);
    const [submitting, setSubmitting] = useState<boolean>(false);
    const [toast, setToast] = useState<{ type: 'success' | 'error' | 'warning' | 'info'; message: string } | null>(null);
    const [confirmModal, setConfirmModal] = useState<{ open: boolean; reviewId: string | null }>({ open: false, reviewId: null });
    const [reportModal, setReportModal] = useState<{ open: boolean; reviewId: string | null }>({ open: false, reviewId: null });
    const [reportForm, setReportForm] = useState<ReportReviewRequest>({ category: 'spam', reason: '' });
    const [reportSubmitting, setReportSubmitting] = useState(false);
    const [deletingReview, setDeletingReview] = useState(false);
    const [canReviewVisible, setCanReviewVisible] = useState<boolean>(!!canReview);
    const { reviews, summary, loading, error, createReview, deleteReview, page, limit, total, setPage, setLimit } = useReviews(productSlug);
    const { t } = useTranslation('common');
    const { user } = useAuth();

    useEffect(() => {
        setCanReviewVisible(!!canReview);
    }, [canReview]);

    const handleStarClick = (value: number) => {
        setRating(value);
    };

    const handleStarHover = (value: number) => {
        setHoverRating(value);
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        if (rating === 0) {
            setToast({ type: 'warning', message: t('reviews.selectRating') });
            return;
        }

        setSubmitting(true);
        const created = await createReview({ rating, comment });
        if (created) {
            setToast({ type: 'success', message: t('reviews.reviewSubmitted') });
            setRating(0);
            setComment('');
            setCanReviewVisible(false);
        }
        setSubmitting(false);
    };

    const handleDeleteReview = async (reviewId: string) => {
        setConfirmModal({ open: true, reviewId });
    };

    const handleConfirmDelete = async () => {
        if (!confirmModal.reviewId) return;

        setDeletingReview(true);
        const deleted = await deleteReview(confirmModal.reviewId);
        setDeletingReview(false);

        if (deleted) {
            setToast({ type: 'success', message: t('reviews.reviewDeleted') });
            setConfirmModal({ open: false, reviewId: null });
            setRating(0);
            setComment('');
            setCanReviewVisible(true);
        } else if (error) {
            setToast({ type: 'error', message: error });
        }
    };

    const handleCancelDelete = () => {
        setConfirmModal({ open: false, reviewId: null });
    };

    const handleOpenReport = (reviewId: string) => {
        setReportForm({ category: 'spam', reason: '' });
        setReportModal({ open: true, reviewId });
    };

    const handleSubmitReport = async () => {
        if (!reportModal.reviewId) return;
        if (!reportForm.reason.trim()) {
            setToast({ type: 'warning', message: t('reviews.provideReportReason') });
            return;
        }
        setReportSubmitting(true);
        try {
            await reportReview(reportModal.reviewId, reportForm);
            setToast({ type: 'success', message: t('reviews.reportSubmitted') });
            setReportModal({ open: false, reviewId: null });
        } catch (err: unknown) {
            let message = t('reviews.reportFailed');
            if (isAxiosError(err)) {
                const code = err.response?.data?.code;
                const status = err.response?.status;
                if (code === 'already_reported') {
                    message = t('reviews.reportAlreadySubmitted');
                } else if (code === 'rate_limited' || status === 429) {
                    message = t('reviews.reportRateLimited');
                } else if (err.response?.data?.error) {
                    message = err.response.data.error;
                }
            }
            setToast({ type: 'error', message: message || t('reviews.reportFailed') });
        } finally {
            setReportSubmitting(false);
        }
    };

    const displayRating = hoverRating || rating;

    return (
        <section className="bg-white shadow rounded-lg p-8 mb-12">
            {canReviewVisible && (
                <h2 className="text-2xl font-semibold text-gray-900 mb-6">{t('reviews.rateThisProduct')}</h2>
            )}

            {/* Summary */}
            {summary && (
                <div className="mb-6">
                    <div className="flex items-center gap-3 mb-3">
                        <div className="text-2xl font-bold text-gray-900">{summary.average.toFixed(1)}</div>
                        <div className="text-sm text-gray-600">{t('reviews.basedOnReviews', { count: summary.count })}</div>
                    </div>
                    {/* Distribution */}
                    <div className="space-y-2">
                        {[5,4,3,2,1].map((stars) => {
                            const count = summary.distribution[stars] || 0;
                            const pct = summary.count > 0 ? Math.round((count / summary.count) * 100) : 0;
                            return (
                                <div key={stars} className="flex items-center gap-3 text-sm">
                                    <div className="w-10 text-right text-gray-700">{stars}★</div>
                                    <div className="flex-1 h-3 bg-gray-100 rounded">
                                        <div className="h-3 bg-yellow-400 rounded" style={{ width: `${pct}%` }} />
                                    </div>
                                    <div className="w-16 text-gray-600 tabular-nums text-right">{count}</div>
                                </div>
                            );
                        })}
                    </div>
                </div>
            )}

            {/* No gating banner: when cannot review, simply hide the form (silent) */}

            {canReviewVisible && (
            <form onSubmit={handleSubmit} noValidate className="space-y-6">
                {/* Star Rating */}
                <div className="space-y-3">
                    <label className="block text-sm font-medium text-gray-900">
                        {t('reviews.rating')}
                        <span className="text-red-500 ml-1">*</span>
                    </label>
                    <div className="flex gap-2">
                        {[1, 2, 3, 4, 5].map((star) => (
                            <button
                                key={star}
                                type="button"
                                onClick={() => handleStarClick(star)}
                                onMouseEnter={() => handleStarHover(star)}
                                onMouseLeave={() => setHoverRating(0)}
                                className="focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md transition-transform hover:scale-110"
                                aria-label={`Rate with ${star} star${star > 1 ? 's' : ''}`}
                            >
                                <svg
                                    className={cn(
                                        'w-8 h-8 transition-colors',
                                        star <= displayRating
                                            ? 'fill-yellow-400 text-yellow-400'
                                            : 'fill-gray-300 text-gray-300'
                                    )}
                                    viewBox="0 0 24 24"
                                    xmlns="http://www.w3.org/2000/svg"
                                >
                                    <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                                </svg>
                            </button>
                        ))}
                    </div>
                    {rating > 0 && (
                        <p className="text-sm text-gray-600">
                            {t('reviews.youSelected', { rating, ratingPlural: rating > 1 ? 's' : '' })}
                        </p>
                    )}
                </div>

                {/* Comment Textarea */}
                <div className="space-y-2">
                    <label htmlFor="comment" className="block text-sm font-medium text-gray-900">
                        {t('reviews.comment')}
                        <span className="text-gray-400 ml-1">({t('reviews.optional')})</span>
                    </label>
                    <textarea
                        id="comment"
                        name="comment"
                        value={comment}
                        onChange={(e) => setComment(e.target.value)}
                        placeholder={t('reviews.placeholder')}
                        rows={4}
                        maxLength={500}
                        className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none text-gray-900 placeholder-gray-500"
                        required={false}
                    />
                    <p className="text-xs text-gray-500">
                        {t('reviews.characters', { length: comment.length })}
                    </p>
                </div>

                {/* Submit Button */}
                <div className="flex gap-3">
                    <button
                        type="submit"
                        disabled={submitting || rating === 0}
                        aria-disabled={submitting || rating === 0}
                        className={cn(
                            'px-6 py-3 rounded-md font-medium transition-colors',
                            (submitting || rating === 0)
                                ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                                : 'bg-blue-600 text-white hover:bg-blue-700'
                        )}
                        title={rating === 0 ? t('reviews.selectRating') : undefined}
                    >
                        {submitting ? t('reviews.submitting') : t('reviews.submitReview')}
                    </button>
                    <button
                        type="reset"
                        onClick={() => {
                            setRating(0);
                            setComment('');
                        }}
                        className="px-6 py-3 rounded-md font-medium border-2 border-gray-300 text-gray-900 hover:bg-gray-50 transition-colors"
                    >
                        {t('common.cancel')}
                    </button>
                </div>
            </form>
            )}

            {/* Toast */}
            {toast && (
                <Toast
                    type={toast.type}
                    message={toast.message}
                    onClose={() => setToast(null)}
                    duration={toast.type === 'error' ? 5000 : 3000}
                />
            )}

            {/* Reviews List */}
            <div className="mt-8">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">{t('reviews.recentReviews')}</h3>
                {loading && <div className="text-sm text-gray-500">{t('reviews.loadingReviews')}</div>}
                {!loading && reviews.length === 0 && (
                    <div className="text-sm text-gray-500">{t('reviews.noReviewsYet')}</div>
                )}
                <ul className="space-y-4">
                    {reviews.map((r) => {
                        // Show delete button only if user is the author of the review
                        const canDelete = user && r.author_id === user.public_id;
                        const canReport = !!user && (!canDelete); // avoid reporting own review
                        
                        return (
                        <li key={r.id} className="border border-gray-200 rounded-md p-4">
                            <div className="flex items-center justify-between">
                                <div className="font-semibold text-gray-900">{r.author_display || 'Anonymous'}</div>
                                <div className="flex items-center gap-2">
                                    <div className="text-xs text-gray-500">{new Date(r.created_at).toLocaleDateString()}</div>
                                    {canDelete && (
                                        <button
                                            onClick={() => handleDeleteReview(r.id)}
                                            className="text-red-600 hover:text-red-800 p-1 rounded transition-colors"
                                            title={t('reviews.deleteReview')}
                                        >
                                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                            </svg>
                                        </button>
                                    )}
                                    {canReport && (
                                        <button
                                            onClick={() => handleOpenReport(r.id)}
                                            className="text-amber-600 hover:text-amber-800 p-1 rounded transition-colors"
                                            aria-label={t('reviews.report')}
                                            title={t('reviews.report')}
                                        >
                                            <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                <path d="M4 21h2" />
                                                <path d="M5 21V5" />
                                                <path d="M5 5l10-2v10l-10 2z" />
                                                <path d="M15 5l4-1v10l-4 1" />
                                            </svg>
                                            <span className="sr-only">{t('reviews.report')}</span>
                                        </button>
                                    )}
                                </div>
                            </div>
                            <div className="flex items-center gap-2 mt-1">
                                {[1,2,3,4,5].map((s) => (
                                    <svg key={s} className={cn('w-4 h-4', s <= r.rating ? 'fill-yellow-400 text-yellow-400' : 'fill-gray-300 text-gray-300')} viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                                    </svg>
                                ))}
                            </div>
                            <div className="mt-2 text-sm text-gray-900 whitespace-pre-line">{r.comment}</div>
                        </li>
                        );
                    })}
                </ul>
                {/* Pagination */}
                {total > limit && (
                    <div className="mt-4 flex items-center justify-between text-sm">
                        <div className="text-gray-600">{t('reviews.pageOf', { page, total: Math.max(1, Math.ceil(total / limit)) })}</div>
                        <div className="flex items-center gap-2">
                            <button
                                type="button"
                                disabled={page <= 1}
                                onClick={() => setPage(Math.max(1, page - 1))}
                                className={cn('px-3 py-1 rounded border', page <= 1 ? 'bg-gray-100 text-gray-400 cursor-not-allowed border-gray-200' : 'bg-white text-gray-700 hover:bg-gray-50 border-gray-300')}
                            >{t('reviews.prev')}</button>
                            <button
                                type="button"
                                disabled={page >= Math.ceil(total / limit)}
                                onClick={() => setPage(Math.min(Math.ceil(total / limit), page + 1))}
                                className={cn('px-3 py-1 rounded border', page >= Math.ceil(total / limit) ? 'bg-gray-100 text-gray-400 cursor-not-allowed border-gray-200' : 'bg-white text-gray-700 hover:bg-gray-50 border-gray-300')}
                            >{t('reviews.next')}</button>
                            <select
                                value={limit}
                                onChange={(e) => setLimit(Number(e.target.value))}
                                className="ml-2 border-gray-300 rounded px-2 py-1 text-gray-700"
                            >
                                {[5,10,20,50].map((n) => (
                                    <option key={n} value={n}>{t('reviews.perPage', { count: n })}</option>
                                ))}
                            </select>
                        </div>
                    </div>
                )}
            </div>

            {/* Confirm Delete Modal */}
            <ConfirmModal
                open={confirmModal.open}
                title={t('reviews.confirmDelete')}
                description={t('reviews.confirmDeleteDescription', 'Esta ação não pode ser desfeita. A avaliação será removida permanentemente.')}
                confirmText={t('common.delete')}
                cancelText={t('common.cancel')}
                onConfirm={handleConfirmDelete}
                onCancel={handleCancelDelete}
                loading={deletingReview}
            />

            {/* Report Review Modal */}
            {reportModal.open && (
                <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-40">
                    <div className="bg-white rounded-lg shadow-lg p-6 w-full max-w-md">
                        <div className="flex items-center justify-between mb-4">
                            <h4 className="text-lg font-semibold text-gray-900">{t('reviews.reportReview')}</h4>
                            <button
                                onClick={() => setReportModal({ open: false, reviewId: null })}
                                className="text-gray-500 hover:text-gray-700"
                                aria-label={t('common.close')}
                            >
                                X
                            </button>
                        </div>
                        <div className="space-y-4">
                            <label className="block text-sm text-gray-800">
                                {t('reviews.reportCategory')}
                                <select
                                    value={reportForm.category}
                                    onChange={(e) => setReportForm((prev) => ({ ...prev, category: e.target.value }))}
                                    className="mt-1 w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                >
                                    <option value="spam">{t('reviews.reportSpam')}</option>
                                    <option value="abuse">{t('reviews.reportAbuse')}</option>
                                    <option value="offensive">{t('reviews.reportOffensive')}</option>
                                    <option value="other">{t('reviews.reportOther')}</option>
                                </select>
                            </label>
                            <label className="block text-sm text-gray-800">
                                {t('reviews.reportReason')}
                                <textarea
                                    value={reportForm.reason}
                                    onChange={(e) => setReportForm((prev) => ({ ...prev, reason: e.target.value }))}
                                    rows={4}
                                    maxLength={500}
                                    className="mt-1 w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    placeholder={t('reviews.reportPlaceholder')}
                                />
                            </label>
                        </div>
                        <div className="mt-6 flex justify-end gap-3">
                            <button
                                onClick={() => setReportModal({ open: false, reviewId: null })}
                                className="px-4 py-2 rounded border border-gray-300 text-gray-700 hover:bg-gray-50"
                                disabled={reportSubmitting}
                            >
                                {t('common.cancel')}
                            </button>
                            <button
                                onClick={handleSubmitReport}
                                disabled={reportSubmitting}
                                className={cn('px-4 py-2 rounded text-white', reportSubmitting ? 'bg-gray-400 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700')}
                            >
                                {reportSubmitting ? t('reviews.submitting') : t('reviews.submitReport')}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </section>
    );
};

export default ProductReview;
