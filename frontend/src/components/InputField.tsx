import React from 'react';

/**
 * InputField component props
 * @property {string} label - Label for the input field
 * @property {string} [type] - Input type (default: 'text')
 * @property {string} name - Name of the input field
 * @property {string} value - Value of the input field
 * @property {(e: React.ChangeEvent<HTMLInputElement>) => void} onChange - Change handler
 * @property {(e: React.FocusEvent<HTMLInputElement>) => void} [onBlur] - Blur handler
 * @property {string} [placeholder] - Placeholder text
 * @property {boolean} [required] - Whether the field is required
 * @property {string} [autoComplete] - Autocomplete attribute
 * @property {React.ReactNode} [rightIcon] - Optional icon rendered on the right
 * @property {() => void} [onRightIconMouseDown] - Mouse/Touch down handler for right icon
 * @property {() => void} [onRightIconMouseUp] - Mouse/Touch up handler for right icon
 * @property {() => void} [onRightIconMouseLeave] - Mouse leave handler for right icon
 */
interface InputFieldProps {
  label: string;
  type?: string;
  name: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onBlur?: (e: React.FocusEvent<HTMLInputElement>) => void;
  placeholder?: string;
  required?: boolean;
  autoComplete?: string;
  disabled?: boolean;
  readOnly?: boolean;
  rightIcon?: React.ReactNode;
  onRightIconMouseDown?: () => void;
  onRightIconMouseUp?: () => void;
  onRightIconMouseLeave?: () => void;
}

/**
 * Reusable input field component with optional right icon and event handlers.
 */
const InputField: React.FC<InputFieldProps> = ({
  label,
  type = 'text',
  name,
  value,
  onChange,
  onBlur,
  placeholder,
  required = false,
  autoComplete,
  disabled = false,
  readOnly = false,
  rightIcon,
  onRightIconMouseDown,
  onRightIconMouseUp,
  onRightIconMouseLeave,
}) => {
  const baseClasses = "w-full px-4 py-2.5 border rounded-lg transition-all duration-200 placeholder-gray-400";
  const enabledClasses = "border-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-gray-50 hover:bg-white";
  const disabledClasses = "border-gray-200 bg-gray-100 text-gray-500 cursor-not-allowed";
  return (
    <div className="flex flex-col gap-2">
      <label htmlFor={name} className="text-sm font-medium text-gray-700">{label}</label>
      <div className="relative">
        <input
          type={type}
          id={name}
          name={name}
          value={value}
          onChange={onChange}
          {...(onBlur ? { onBlur } : {})}
          className={`${baseClasses} ${disabled ? disabledClasses : enabledClasses}`}
          placeholder={placeholder}
          required={required}
          autoComplete={autoComplete}
          disabled={disabled}
          readOnly={readOnly}
        />
        {rightIcon && (
          <span
            className="absolute right-3 top-1/2 -translate-y-1/2 cursor-pointer text-gray-400 hover:text-gray-600 select-none transition-colors duration-200"
            onMouseDown={onRightIconMouseDown}
            onMouseUp={onRightIconMouseUp}
            onMouseLeave={onRightIconMouseLeave}
            onTouchStart={(e) => {
              e.preventDefault();
              onRightIconMouseDown?.();
            }}
            onTouchEnd={(e) => {
              e.preventDefault();
              onRightIconMouseUp?.();
            }}
            tabIndex={0}
            role="button"
            aria-label="Toggle password visibility"
          >
            {rightIcon}
          </span>
        )}
      </div>
    </div>
  );
};

export default InputField;
