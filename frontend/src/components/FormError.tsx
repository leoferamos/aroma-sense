import React from 'react';

/**
 * FormError component props
 * @property {string} [message] - Error message to display
 */
interface FormErrorProps {
  message?: string;
}

/**
 * Displays a form error message in red text.
 */
const FormError: React.FC<FormErrorProps> = ({ message }) => {
  if (!message) return null;
  return <span className="text-sm text-red-500 mt-1">{message}</span>;
};

export default FormError;
