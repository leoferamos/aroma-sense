import React from 'react';
import { useNavigate } from 'react-router-dom';

interface BackButtonProps {
    className?: string;
    label?: string;
    fallbackPath?: string;
}

const BackButton: React.FC<BackButtonProps> = ({
    className = '',
    label = 'Back',
    fallbackPath = '/'
}) => {
    const navigate = useNavigate();

    const handleBack = () => {
        if (window.history.length > 1) {
            navigate(-1);
        } else {
            navigate(fallbackPath);
        }
    };

    return (
        <button
            onClick={handleBack}
            className={`inline-flex items-center gap-2 text-sm font-medium text-gray-600 hover:text-gray-900 transition-colors ${className}`}
            aria-label="Go back"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={2}
                stroke="currentColor"
                className="w-5 h-5"
            >
                <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
                />
            </svg>
            {label}
        </button>
    );
};

export default BackButton;
