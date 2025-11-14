import React, { useState } from 'react';
import { cn } from '../utils/cn';

interface ProductReviewProps {
    productId: number;
}

const ProductReview: React.FC<ProductReviewProps> = ({ productId }) => {
    const [rating, setRating] = useState<number>(0);
    const [comment, setComment] = useState<string>('');
    const [hoverRating, setHoverRating] = useState<number>(0);
    const [submitting, setSubmitting] = useState<boolean>(false);

    const handleStarClick = (value: number) => {
        setRating(value);
    };

    const handleStarHover = (value: number) => {
        setHoverRating(value);
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        if (rating === 0) {
            alert('Please select a rating');
            return;
        }

        if (!comment.trim()) {
            alert('Please write a comment');
            return;
        }

        setSubmitting(true);
        try {
            // TODO: Make HTTP call here
            console.log({
                productId,
                rating,
                comment,
            });
            alert('Review submitted successfully!');
            setRating(0);
            setComment('');
        } catch (error) {
            console.error('Error submitting review:', error);
        } finally {
            setSubmitting(false);
        }
    };

    const displayRating = hoverRating || rating;

    return (
        <section className="bg-white shadow rounded-lg p-8 mb-12">
            <h2 className="text-2xl font-semibold text-gray-900 mb-6">Rate This Product</h2>

            <form onSubmit={handleSubmit} noValidate className="space-y-6">
                {/* Star Rating */}
                <div className="space-y-3">
                    <label className="block text-sm font-medium text-gray-900">
                        Rating
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
                            You selected {rating} star{rating > 1 ? 's' : ''}
                        </p>
                    )}
                </div>

                {/* Comment Textarea */}
                <div className="space-y-2">
                    <label htmlFor="comment" className="block text-sm font-medium text-gray-900">
                        Comment
                        <span className="text-red-500 ml-1">*</span>
                    </label>
                    <textarea
                        id="comment"
                        name="comment"
                        value={comment}
                        onChange={(e) => setComment(e.target.value)}
                        placeholder="Share your experience with this product..."
                        rows={4}
                        className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none text-gray-900 placeholder-gray-500"
                        required
                    />
                    <p className="text-xs text-gray-500">
                        {comment.length}/500 characters
                    </p>
                </div>

                {/* Submit Button */}
                <div className="flex gap-3">
                    <button
                        type="submit"
                        disabled={submitting}
                        aria-disabled={submitting}
                        className={cn(
                            'px-6 py-3 rounded-md font-medium transition-colors',
                            submitting
                                ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                                : 'bg-blue-600 text-white hover:bg-blue-700'
                        )}
                    >
                        {submitting ? 'Submitting...' : 'Submit Review'}
                    </button>
                    <button
                        type="reset"
                        onClick={() => {
                            setRating(0);
                            setComment('');
                        }}
                        className="px-6 py-3 rounded-md font-medium border-2 border-gray-300 text-gray-900 hover:bg-gray-50 transition-colors"
                    >
                        Cancel
                    </button>
                </div>
            </form>
        </section>
    );
};

export default ProductReview;
