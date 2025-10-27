/**
 * Currency and number formatting utilities
 */

const CURRENCY_LOCALE = 'pt-BR';
const CURRENCY_CODE = 'BRL';

/**
 * Formats a number as Brazilian Real (BRL) currency
 * @param value - The numeric value to format
 * @returns Formatted currency string
 */
export function formatCurrency(value: number): string {
  return new Intl.NumberFormat(CURRENCY_LOCALE, {
    style: 'currency',
    currency: CURRENCY_CODE,
  }).format(value);
}

/**
 * Formats a number with locale-specific thousand separators
 * @param value - The numeric value to format
 * @returns Formatted number string
 */
export function formatNumber(value: number): string {
  return new Intl.NumberFormat(CURRENCY_LOCALE).format(value);
}
