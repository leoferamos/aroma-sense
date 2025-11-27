import React, { useCallback } from "react";

type Props = {
  value: string;
  onChange: (value: string) => void;
  onSubmit: () => void;
  onClear?: () => void;
  isLoading?: boolean;
  placeholder?: string;
  className?: string;
};

const SearchBar: React.FC<Props> = ({
  value,
  onChange,
  onSubmit,
  onClear,
  isLoading = false,
  placeholder = "Search products...",
  className = "",
}) => {
  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === "Enter") {
        e.preventDefault();
        onSubmit();
      }
      if (e.key === "Escape" && onClear) {
        onClear();
      }
    },
    [onSubmit, onClear]
  );

  return (
    <div className={`w-full flex items-center gap-2 ${className}`}>
      <div className="relative flex-1">
        <span className="pointer-events-none absolute inset-y-0 left-4 flex items-center text-gray-400">
          <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path fillRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.817-4.817A6 6 0 012 8z" clipRule="evenodd" />
          </svg>
        </span>
        <input
          type="search"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={handleKeyDown}
          aria-label="Search products"
          placeholder={placeholder}
          className="block w-full rounded-2xl border-0 bg-white shadow-sm ring-1 ring-gray-200 py-3.5 pl-12 pr-32 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all duration-300 hover:shadow-md hover:ring-gray-300"
        />
        {value && onClear && (
          <button
            type="button"
            onClick={onClear}
            aria-label="Clear search"
            className="absolute inset-y-0 right-24 flex items-center text-gray-400 hover:text-gray-700 transition-colors duration-200 hover:scale-110"
          >
            <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
              <path fillRule="evenodd" d="M10 8.586l4.95-4.95a1 1 0 011.414 1.414L11.414 10l4.95 4.95a1 1 0 01-1.414 1.414L10 11.414l-4.95 4.95a1 1 0 11-1.414-1.414L8.586 10l-4.95-4.95A1 1 0 115.05 3.636L10 8.586z" clipRule="evenodd" />
            </svg>
          </button>
        )}
        <button
          type="button"
          onClick={onSubmit}
          aria-label="Submit search"
          className="absolute inset-y-0 right-2 my-2 rounded-xl bg-gradient-to-r from-blue-600 to-blue-700 px-5 text-sm font-medium text-white shadow-sm hover:shadow-lg hover:from-blue-700 hover:to-blue-800 transition-all duration-200 active:scale-95"
        >
          Search
        </button>
        {isLoading && (
          <span className="absolute inset-y-0 right-28 my-auto inline-flex items-center text-gray-500" aria-live="polite" aria-busy="true">
            <svg className="animate-spin h-5 w-5 mr-1" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
            </svg>
            <span className="text-xs">Searching...</span>
          </span>
        )}
      </div>
    </div>
  );
};

export default SearchBar;
