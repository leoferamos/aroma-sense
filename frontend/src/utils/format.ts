/**
 * Currency and number formatting utilities
 */
import { APP_LOCALE, CURRENCY_CODE } from '../constants/app';

/**
 * Formats a number as Brazilian Real (BRL) currency
 * @param value - The numeric value to format
 * @returns Formatted currency string
 */
export function formatCurrency(value: number): string {
  return new Intl.NumberFormat(APP_LOCALE, {
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
  return new Intl.NumberFormat(APP_LOCALE).format(value);
}
