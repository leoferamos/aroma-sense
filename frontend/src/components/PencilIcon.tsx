import React from 'react';

interface PencilIconProps {
  className?: string;
}

const PencilIcon: React.FC<PencilIconProps> = ({ className }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="currentColor"
    className={className || 'w-5 h-5'}
    aria-hidden="true"
  >
    <path d="M21.731 2.269a2.625 2.625 0 00-3.712 0L8.669 11.62a1.5 1.5 0 00-.372.6l-1.5 4.125a.75.75 0 00.95.95l4.125-1.5a1.5 1.5 0 00.6-.372l9.35-9.35a2.625 2.625 0 000-3.712z" />
    <path d="M5.25 19.5h13.5a.75.75 0 010 1.5H5.25a.75.75 0 010-1.5z" />
  </svg>
);

export default PencilIcon;
