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
 * @property {() => void} [onRightIconMouseDown] - Mouse down handler for right icon
 * @property {() => void} [onRightIconMouseUp] - Mouse up handler for right icon
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
  rightIcon,
  onRightIconMouseDown,
  onRightIconMouseUp,
  onRightIconMouseLeave,
}) => {
  return (
    <div className="flex flex-col gap-2">
  <label htmlFor={name} className="text-base font-normal text-gray-800">{label}</label>
      <div className="relative">
        <input
          type={type}
          id={name}
          name={name}
          value={value}
          onChange={onChange}
          {...(onBlur ? { onBlur } : {})}
          className="border border-gray-300 rounded-xl px-4 py-3 w-full focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 pr-10"
          placeholder={placeholder}
          required={required}
          autoComplete={autoComplete}
        />
        {rightIcon && (
          <span
            className="absolute right-3 top-1/2 -translate-y-1/2 cursor-pointer text-gray-400"
            onMouseDown={onRightIconMouseDown}
            onMouseUp={onRightIconMouseUp}
            onMouseLeave={onRightIconMouseLeave}
            tabIndex={0}
            role="button"
          >
            {rightIcon}
          </span>
        )}
      </div>
    </div>
  );
};

export default InputField;
